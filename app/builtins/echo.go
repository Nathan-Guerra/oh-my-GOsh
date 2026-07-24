package builtins

import (
	"strings"
)

type Echo struct{}

func (e *Echo) Exec(args []string) *Response {
	out := strings.Join(args, " ") + "\r\n"

	return &Response{Out: out}
}

func init() {
	Register("echo", &Echo{})
}
