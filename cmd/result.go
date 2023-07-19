package cmd

type ResultType byte

const (
	String     ResultType = '+'
	Number     ResultType = ':'
	Array      ResultType = '*'
	Error      ResultType = '-'
	BulkString ResultType = '$'
)

type Element struct {
	ElementType ResultType
	Value       string
}

type Result struct {
	Err   error
	Data  []Element
	Close bool
}
