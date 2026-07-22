package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/app/builtins"
	"github.com/codecrafters-io/shell-starter-go/app/lexer"
	"github.com/codecrafters-io/shell-starter-go/app/parser"
)

func findCommand(search string) []string {
	matches := make([]string, 0)

	for name := range builtins.Builtins {
		if strings.HasPrefix(name, search) {
			matches = append(matches, name)
		}
	}

	path := os.Getenv("PATH")
	if len(path) > 0 {
		dirs := strings.Split(path, string(os.PathListSeparator))

		for _, dir := range dirs {
			entries, err := os.ReadDir(dir)
			if err != nil {
				// fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
				continue
			}

			for _, entry := range entries {
				info, err := entry.Info()
				if err != nil {
					panic(err)
				}

				if entry.Type().IsDir() || !strings.HasPrefix(entry.Name(), search) {
					continue
				}

				if info.Mode()&0111 != 0 {
					matches = append(matches, entry.Name())
				}
			}
		}
	}

	if len(matches) < 2 {
		os.Stdout.Write([]byte{'\a'})
		return make([]string, 0)

	}
	// else if len(matches) >= 20 {
	// 	fmt.Fprintf(os.Stdout, "Display all %d possibilities? (y or n)", len(matches))
	// 	scanner := bufio.NewScanner(os.Stdin)
	// 	if scanner.Scan() {
	// 		if scanner.Text() == "n" {
	// 			return make([]string, 0)
	// 		} else if
	// 	}
	// }

	return matches
}

var commandList = readline.PcItemDynamic(findCommand)
var prefixCompleter = readline.NewPrefixCompleter(commandList)
var cfg readline.Config = readline.Config{
	Prompt:       "$ ",
	AutoComplete: prefixCompleter,
}

func main() {
	rl, err := readline.NewEx(&cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "starting new term:", err)
	}
	defer rl.Close()
	for {
		input, err := rl.Readline()
		if err == readline.ErrInterrupt {
			if len(input) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		} else if err != nil {
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
