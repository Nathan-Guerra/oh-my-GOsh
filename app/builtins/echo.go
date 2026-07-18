package builtins

import (
	"strings"
)

func echo(args []string) (out string, err error) {
	var s strings.Builder
	for i := 0; i < len(args); i++ {
		if args[i] == "" {
			continue
		}
		if i > 0 {
			s.WriteRune(' ')
		}
		s.WriteString(args[i])
	}
	out = s.String()
	return
}

func init() {
	Register("echo", echo)
}
