package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/app/builtins"
	"github.com/codecrafters-io/shell-starter-go/app/lexer"
	"github.com/codecrafters-io/shell-starter-go/app/parser"
)

func main() {
	var prefixCompleter readline.PrefixCompleter = *readline.NewPrefixCompleter(
		readline.PcItem("echo"),
		readline.PcItem("exit"),
	)

	var cfg readline.Config = readline.Config{
		Prompt:       "$ ",
		AutoComplete: &prefixCompleter,
	}
	rl, err := readline.NewEx(&cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "starting new term:", err)
	}
	defer rl.Close()
	for {
		input, err := rl.Readline()
		if err != nil {
			fmt.Println(err)
			break
		}

		cmd := parser.CreateCommand(lexer.Tokenize(input))
		if cmd.CommandName == "" {
			continue
		}

		command, exists := builtins.Builtins[cmd.CommandName]
		if exists {
			out, err := command(cmd.Arguments)
			if len(out) > 0 {
				cmd.Stdout.Write([]byte(out))
			}
			if err != nil {
				cmd.Stderr.Write([]byte(err.Error()))

			}
		} else if _, err := exec.LookPath(cmd.CommandName); err == nil {
			externalCommand := exec.Command(cmd.CommandName, cmd.Arguments...)
			externalCommand.Stdin = os.Stdin
			externalCommand.Stdout = cmd.Stdout
			externalCommand.Stderr = cmd.Stderr

			externalCommand.Run()
		} else {
			fmt.Printf("%s: command not found\n", cmd.CommandName)
		}
	}
}
