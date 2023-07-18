package server

import "kvstore/cmd"

type Request struct {
	Command  cmd.ICommand
	Operands []string
}
