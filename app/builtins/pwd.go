package builtins

import (
	"os"
)

func pwd(args []string) (out string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}

	out = wd
	return
}

func init() {
	Register("pwd", pwd)
}
