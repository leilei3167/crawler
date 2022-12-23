package collect

type Request struct {
	URL       string
	Cookie    string
	ParseFunc func(c []byte, request *Request) ParseResult
}

type ParseResult struct {
	Requests []*Request
	Items    []any
}
