package builtins

import (
	"fmt"
	"os"
	"strings"
)

func echo(args []string) error {
	var out string

	if len(args) == 1 && strings.HasPrefix(args[0], "$") {
		fmt.Println(os.Getenv(args[0][1:]))
		return nil
	}

	for i := 0; i < len(args); i++ {
		if args[i] == "" {
			continue
		}
		out += args[i] + " "
	}

	if out != "" {
		out = out[:len(out)-1]
	}

	fmt.Println(out)

	return nil
}

func init() {
	Register("echo", echo)
}
