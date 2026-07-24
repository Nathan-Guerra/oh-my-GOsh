package builtins

import (
	"fmt"

	"github.com/codecrafters-io/shell-starter-go/app/autocomplete"
)

type Type struct{}

func (t *Type) Exec(args []string) *Response {
	autocompleter := autocomplete.GetCommandAutocompleter()
	loadedCommands := autocompleter.GetLoadedCommands()

	for _, command_name := range args {
		if command_name == "" {
			continue
		}

		path, exists := loadedCommands[command_name]
		if exists {
			if path == "builtin" {
				return &Response{Out: fmt.Sprintf("%s is a shell builtin\n", command_name)}
			} else {
				return &Response{Out: fmt.Sprintf("%s is %s\n", command_name, path)}

			}
		}

		return &Response{Out: fmt.Sprintf("%s: not found\n", command_name)}
	}

	return &Response{}

}

func init() {
	Register("type", &Type{})
}
