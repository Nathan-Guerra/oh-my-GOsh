package builtins

import (
	"strconv"
)

type Exit struct{}

func (e *Exit) Exec(args []string) *Response {
	if len(args) == 0 {
		return &Response{ShouldExit: true, ExitSignal: 0}
	}

	code, err := strconv.Atoi(args[0])
	if err != nil {
		return &Response{Err: err.Error()}
	}

	return &Response{ShouldExit: true, ExitSignal: code}
}

func init() {
	Register("exit", &Exit{})
}
