package builtins

import "os"

func cd(args []string) error {
	if len(args) == 0 {
		home_dir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		return os.Chdir(home_dir)
	}
	return os.Chdir(args[0])
}

func init() {
	Register("cd", cd)
}
