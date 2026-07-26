package builtins

type Complete struct{}

func (e *Complete) Exec(args []string) *Response {
	return &Response{}
}

func init() {
	Register("complete", &Complete{})
}
