package builtins

import (
	"fmt"
	"os/exec"
)

func typeCommand(args []string) error {
	for _, command_name := range args {
		if command_name == "" {
			continue
		}

		_, exists := Builtins[command_name]
		if exists {
			fmt.Printf("%s is a shell builtin\n", command_name)
		} else {
			command_path, err := exec.LookPath(command_name)
			if err == nil {
				fmt.Printf("%s is %s\n", command_name, command_path)
			} else {
				fmt.Printf("%s: not found\n", command_name)
			}

		}
	}

	return nil
}

func init() {
	Register("type", typeCommand)
}
