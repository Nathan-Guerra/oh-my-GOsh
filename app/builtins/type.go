package builtins

import (
	"fmt"
	"os/exec"
)

func typeCommand(args []string) (out string, err error) {
	for _, command_name := range args {
		if command_name == "" {
			continue
		}

		_, exists := Builtins[command_name]
		if exists {
			out = fmt.Sprintf("%s is a shell builtin\n", command_name)
		} else {
			command_path, err := exec.LookPath(command_name)
			if err == nil {
				out = fmt.Sprintf("%s is %s\n", command_name, command_path)
			} else {
				out = fmt.Sprintf("%s: not found\n", command_name)
			}

		}
	}
	return
}

func init() {
	Register("type", typeCommand)
}
