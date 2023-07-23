package server

import (
	"io"
	"kvstore/cmd"
)

type Protocol interface {
	RequestParser
	ResponseProcessor
	GetName() string
}

type RequestParser interface {
	ParseRequest(reader io.Reader) (*Request, error)
}

type ResponseProcessor interface {
	PrepareResponseFromResult(res *cmd.Result) string
}
