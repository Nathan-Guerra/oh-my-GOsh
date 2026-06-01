package builtins

import (
	"os"
	"strconv"
)

func exit(args []string) error {
	if len(args) == 0 {
		os.Exit(0)
	} else {
		code, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		os.Exit(code)
	}
	os.Exit(0)

	return nil
}

func init() {
	Register("exit", exit)
}
