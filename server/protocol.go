package server

type Protocol interface {
	ParseRequest(data string) (*Request, error)
	PrepareResponse(res *Response) string
}
