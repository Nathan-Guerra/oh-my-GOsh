package builtins

import (
	"fmt"
	"os"
	"strings"
)

func cd(args []string) error {
	cur_dir, err := os.Getwd()
	if len(args) == 0 {
		home_dir, err := os.UserHomeDir()
		if err == nil {
			return err
		}

		err = os.Chdir(home_dir)
		if err == nil {
			os.Setenv("OLDPWD", cur_dir)
		}

		return err
	}

	fmt.Printf("'%s'\n", args[0])
	if args[0] == "-" {
		args[0] = os.Getenv("OLDPWD")
	} else if strings.HasPrefix(args[0], "~") {
		home_dir, err := os.UserHomeDir()
		if err != nil {
			return err
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
		return fmt.Errorf("cd: %s: No such file or directory", args[0])
	}

	os.Setenv("OLDPWD", cur_dir)
	return nil
}

func init() {
	Register("cd", cd)
}
