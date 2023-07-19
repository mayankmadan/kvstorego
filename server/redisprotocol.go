package server

import (
	"bufio"
	"errors"
	"io"
	"kvstore/cmd"
	"strconv"
	"strings"
)

type RedisProtocol struct {
	name string
}

func (proto *RedisProtocol) GetName() string {
	return proto.name
}

func getTypeFromInitialChar(ch byte) cmd.ResultType {
	switch ch {
	case '+':
		return cmd.String
	case ':':
		return cmd.Number
	case '*':
		return cmd.Array
	case '$':
		return cmd.BulkString
	case '-':
		return cmd.Error
	}
	return cmd.Error
}

func readLine(reader *bufio.Reader) ([]byte, error) {
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	if len(line) > 2 && line[len(line)-2] == '\r' {
		return line[:len(line)-2], nil
	}
	return line[:len(line)-1], nil
}

func readArray(reader *bufio.Reader, elementcount int) ([]string, error) {
	elements := make([]string, 0, elementcount)
	for i := 0; i < elementcount; i++ {
		line, err := readLine(reader)
		if err != nil {
			return nil, err
		}
		switch initialChar := getTypeFromInitialChar(line[0]); initialChar {
		case cmd.BulkString:
			length, err := strconv.Atoi(string(line[1:]))
			if err != nil {
				return nil, err
			}
			blkString, err := readBulkString(reader, length)
			if err != nil {
				return nil, err
			}
			elements = append(elements, string(blkString))
		case cmd.String:
			elements = append(elements, string(readString(line)))
		case cmd.Number:
			elements = append(elements, string(readString(line))) // Since we are treating numbers as strings for now
		}
	}
	return elements, nil
}

func readString(line []byte) []byte {
	return line[1:]
}

func readNumber(line []byte) (int, error) {
	val, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return 0, nil
	}
	return val, nil
}

func readBulkString(reader *bufio.Reader, length int) ([]byte, error) {
	val, err := readLine(reader)
	if err != nil {
		return nil, err
	}
	if len(val) < length {
		return nil, errors.New("input length does not match")
	}
	return val[:length], nil
}

func (proto *RedisProtocol) ParseRequest(reader io.Reader) (*Request, error) {
	bufreader := bufio.NewReader(reader)
	lineBytes, err := readLine(bufreader)
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}
	line := string(lineBytes)
	inputType := getTypeFromInitialChar(lineBytes[0])

	if inputType == cmd.Array {
		elements, err := strconv.Atoi(string(line[1:]))
		if err != nil {
			return nil, err
		}
		inputData, err := readArray(bufreader, elements)
		if err != nil {
			return nil, err
		}
		if len(inputData) == 0 {
			return nil, errors.New("empty input")
		}
		commandString := inputData[0]
		operands := inputData[1:]

		command := cmd.GetCommand(commandString)
		if command == nil {
			return nil, errors.New("unknown command")
		}

		if len(operands) != command.GetNumberOfOperands() {
			return nil, errors.New("invalid operands")
		}

		return &Request{Command: command, Operands: operands}, nil
	}
	return nil, errors.New("unparsable input")
}

func (proto *RedisProtocol) PrepareResponseFromResult(res *cmd.Result) string {
	var sb strings.Builder
	if res.Err != nil {
		sb.WriteString("-" + res.Err.Error() + "\r\n")
		return sb.String()
	}
	data := res.Data
	numberOfElements := len(data)
	if numberOfElements > 1 {
		sb.WriteString("*" + strconv.Itoa(numberOfElements) + "\r\n")
	}
	for _, val := range data {
		initialChars := val.ElementType
		sb.WriteByte(byte(initialChars))
		sb.WriteString(val.Value)
		sb.WriteString("\r\n")
	}

	return sb.String()
}
