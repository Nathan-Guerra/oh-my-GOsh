package builtins

import (
	"fmt"
	"os"
	"strings"
)

func cd(args []string) (out string, err error) {
	cur_dir, err := os.Getwd()
	if len(args) == 0 {
		home_dir, e := os.UserHomeDir()
		if e == nil {
			return
		}

		e = os.Chdir(home_dir)
		if e == nil {
			os.Setenv("OLDPWD", cur_dir)
		}

		return
	}

	if args[0] == "-" {
		args[0] = os.Getenv("OLDPWD")
	} else if strings.HasPrefix(args[0], "~") {
		home_dir, e := os.UserHomeDir()
		if e != nil {
			err = e
		}

		ps := string(os.PathSeparator)
		if strings.HasSuffix(home_dir, ps) {
			args[0] = home_dir + args[0][1:]
		} else {
			args[0] = home_dir + ps + args[0][1:]
		}
	}

	err = os.Chdir(args[0])
	if err != nil {
		err = fmt.Errorf("cd: %s: No such file or directory", args[0])
	}

	os.Setenv("OLDPWD", cur_dir)
	return
}

func init() {
	Register("cd", cd)
}
