package builtins

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/autocomplete"
)

type Type struct{}

func (t *Type) Exec(args []string) *Response {
	autocompleter := autocomplete.GetCommandAutocompleter()
	loadedCommands := autocompleter.GetLoadedCommands()

	var str strings.Builder
	for _, command_name := range args {
		if command_name == "" {
			continue
		}

		path, exists := loadedCommands[command_name]
		if exists {
			if path == "builtin" {
				str.WriteString(fmt.Sprintf("%s is a shell builtin\n", command_name))
			} else {
				str.WriteString(fmt.Sprintf("%s is %s\n", command_name, path))
			}
		} else {
			str.WriteString(fmt.Sprintf("%s: not found\n", command_name))
		}
	}

	return &Response{Out: str.String()}

}

func init() {
	Register("type", &Type{})
}
