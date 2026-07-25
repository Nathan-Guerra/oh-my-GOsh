package main

import (
	"fmt"
	"maps"
	"os"
	"os/exec"
	"slices"

	"github.com/codecrafters-io/shell-starter-go/app/autocomplete"
	"github.com/codecrafters-io/shell-starter-go/app/builtins"
	"github.com/codecrafters-io/shell-starter-go/app/lexer"
	"github.com/codecrafters-io/shell-starter-go/app/parser"
)

func REPL() int {
	autocompleter := autocomplete.GetCommandAutocompleter()
	autocompleter.SetBuiltins(slices.Collect(maps.Keys(builtins.Builtins)))
	autocompleter.SetPATH(os.Getenv("PATH"))
	autocompleter.EagerLoad()

	var code int

	for {
		fmt.Print(Prompt)
		buffer := readLine()

		cmd := parser.CreateCommand(lexer.Tokenize(buffer))
		if cmd.CommandName == "" {
			continue
		}

		command, exists := builtins.Builtins[cmd.CommandName]
		if exists {
			response := command.Exec(cmd.Arguments)
			if len(response.Out) > 0 {
				cmd.Stdout.Write([]byte(response.Out))
			}
			if len(response.Err) > 0 {
				cmd.Stderr.Write([]byte(response.Err))

			}

			if response.ShouldExit {
				code = response.ExitSignal
				break
			}
		} else if _, err := exec.LookPath(cmd.CommandName); err == nil {
			externalCommand := exec.Command(cmd.CommandName, cmd.Arguments...)
			externalCommand.Stdin = os.Stdin
			externalCommand.Stdout = cmd.Stdout
			externalCommand.Stderr = cmd.Stderr

			externalCommand.Run()
		} else {
			fmt.Println(cmd.CommandName + ": command not found")
		}
	}

	return code
}

func main() {
	os.Exit(REPL())
}
