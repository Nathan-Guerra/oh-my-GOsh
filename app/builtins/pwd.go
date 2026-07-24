package builtins

import (
	"os"
)

type Pwd struct{}

func (p *Pwd) Exec(args []string) *Response {
	wd, err := os.Getwd()
	if err != nil {
		return &Response{}
	}

	return &Response{Out: wd + string('\n')}
}

func init() {
	Register("pwd", &Pwd{})
}
