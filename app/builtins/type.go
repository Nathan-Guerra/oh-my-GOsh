package builtins

import (
	"fmt"
	"os/exec"
)

type Type struct{}

func (t *Type) Exec(args []string) *Response {
	for _, command_name := range args {
		if command_name == "" {
			continue
		}

		_, exists := Builtins[command_name]
		if exists {
			return &Response{Out: fmt.Sprintf("%s is a shell builtin\n", command_name)}
		}

		command_path, err := exec.LookPath(command_name)
		if err == nil {
			return &Response{Out: fmt.Sprintf("%s is %s\n", command_name, command_path)}
		}

		return &Response{Out: fmt.Sprintf("%s: not found\n", command_name)}
	}

	return &Response{}

}

func init() {
	Register("type", &Type{})
}
