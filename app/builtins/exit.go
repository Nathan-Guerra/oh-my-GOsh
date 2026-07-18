package builtins

import (
	"os"
	"strconv"
)

func exit(args []string) (out string, err error) {
	if len(args) == 0 {
		os.Exit(0)
	} else {
		code, e := strconv.Atoi(args[0])
		if e != nil {
			err = e
			return
		}

		os.Exit(code)
	}
	os.Exit(0)
	return
}

func init() {
	Register("exit", exit)
}
