package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
	"github.com/codecrafters-io/shell-starter-go/app/lexer"
	"github.com/codecrafters-io/shell-starter-go/app/parser"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("$ ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
			continue
		}

		input := strings.TrimSpace(scanner.Text())
		cmd := parser.CreateCommand(lexer.Tokenize(input))
		fmt.Printf("Output: %v\n", cmd)
		fmt.Printf("Output: %v\n", cmd.CommandName)
		fmt.Printf("Output: %v\n", cmd.Arguments)
		// continue
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
