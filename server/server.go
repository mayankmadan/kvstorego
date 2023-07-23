package server

import (
	"fmt"
	"kvstore/cmd"
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
		request, err := proto.ParseRequest(conn)
		if err != nil {
			result := &cmd.Result{Data: nil, Err: err}
			conn.Write([]byte(proto.PrepareResponseFromResult(result)))
			continue
		}
		if request == nil {
			break
		}

		result := request.Command.Exec(request.Operands)
		out := proto.PrepareResponseFromResult(result)
		conn.Write([]byte(out))
		if result.Close {
			break
		}

	}
}

func CreateServer(port int, proto Protocol) Server {
	return Server{port, proto}
}
