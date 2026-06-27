package builtins

import (
	"fmt"
)

func echo(args []string) error {
	var out string
	for i := 0; i < len(args); i++ {
		if args[i] == "" {
			continue
		}
		if i > 0 {
			out += " "
		}
		out += args[i]
	}
	fmt.Println(out)
	return nil
}

func init() {
	Register("echo", echo)
}
