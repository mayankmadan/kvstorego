package server

import (
	"errors"
	"kvstore/cmd"
	"strings"
)

type RedisProtocol struct {
}

func (proto *RedisProtocol) ParseRequest(data string) (*Request, error) {
	str := strings.Split(data, " ")
	inputLength := len(str)
	if inputLength == 0 {
		return nil, errors.New("empty Input")
	}
	commandString := str[0]
	command := cmd.GetCommand(commandString)
	if command == nil {
		return nil, errors.New("unknown Command")
	}
	req := &Request{Command: command, Operands: str[1:]}

	return req, nil

}

func (proto *RedisProtocol) PrepareResponse(res *Response) string {
	var sb strings.Builder
	if res.Err != nil {
		sb.WriteString("-" + res.Err.Error())
	} else {
		sb.WriteByte('+')
		if res.Data == "" {
			sb.WriteString("OK")
		} else {
			sb.WriteString(res.Data)
		}
	}

	return sb.String()
}
