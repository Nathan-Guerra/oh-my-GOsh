package builtins

import (
	"fmt"
	"os"
	"strings"
)

type Cd struct{}

func (c *Cd) Exec(args []string) *Response {
	cur_dir, err := os.Getwd()
	if len(args) == 0 {
		home_dir, e := os.UserHomeDir()
		if e == nil {
			return &Response{}
		}

		e = os.Chdir(home_dir)
		if e == nil {
			os.Setenv("OLDPWD", cur_dir)
		}

		return &Response{}
	}

	if args[0] == "-" {
		args[0] = os.Getenv("OLDPWD")
	} else if strings.HasPrefix(args[0], "~") {
		home_dir, e := os.UserHomeDir()
		if e == nil {
			ps := string(os.PathSeparator)
			if strings.HasSuffix(home_dir, ps) {
				args[0] = home_dir + args[0][1:]
			} else {
				args[0] = home_dir + ps + args[0][1:]
			}
		}
	}

	err = os.Chdir(args[0])
	if err != nil {
		return &Response{Err: fmt.Sprintf("cd: %s: No such file or directory\n", args[0])}
	}

	os.Setenv("OLDPWD", cur_dir)
	return &Response{}
}

func init() {
	Register("cd", &Cd{})
}
