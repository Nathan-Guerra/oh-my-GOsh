package builtins

import "fmt"

func echo(args []string) error {
	var out string
	for i := 0; i < len(args); i++ {
		if args[i] == "" {
			continue
		}
		out += args[i] + " "
	}

	if out != "" {
		out = out[:len(out)-1]
	}

	fmt.Printf("%s\n", out)

	return nil
}

func init() {
	Register("echo", echo)
}
