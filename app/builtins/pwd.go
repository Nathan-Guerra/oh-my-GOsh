package builtins

import (
	"fmt"
	"os"
)

func pwd(args []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	fmt.Println(wd)

	return nil
}

func init() {
	Register("pwd", pwd)
}
