package builtins

import (
	"os"
)

func pwd(args []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	println(wd)

	return nil
}

func init() {
	Register("pwd", pwd)
}
