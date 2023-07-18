package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

type Server struct {
	port  int
	proto Protocol
}

func (srv Server) Start() {
	ln, err := net.Listen("tcp", "0.0.0.0:"+fmt.Sprint(srv.port))
	if err != nil {
		fmt.Printf("Error starting up Server: %s\n", err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %s\n", err)
		}
		go handleConnection(conn, srv.proto)
	}
}

func handleConnection(conn net.Conn, proto Protocol) {
	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection Closing")
			} else {
				fmt.Printf("Error Reading: %s\n", err)
			}
			break
		}

		request, err := proto.ParseRequest(string(line))
		if err != nil {
			response := &Response{Data: "", Err: err}
			conn.Write([]byte(proto.PrepareResponse(response) + "\r\n"))
			continue
		}

		resp, close, error := request.Command.Exec(request.Operands)
		response := &Response{Data: resp, Err: error}
		out := proto.PrepareResponse(response)
		conn.Write([]byte(out + "\r\n"))
		if close {
			break
		}

	}
}

func CreateServer(port int, proto Protocol) Server {
	return Server{port, proto}
}
