package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
)

func main() {
	const (
		SPACE_CHAR   = ' '
		S_QUOTE_CHAR = '\''
	)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("$ ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
			continue
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		var command_name string
		var args []string
		var found_command bool

		i := 0
		for {
			if !found_command {
				for i < len(input) {
					if input[i] == SPACE_CHAR {
						found_command = true
						break
					}

					command_name += string(input[i])
					i++
				}
			} else {
				var arg string
			loop:
				for i < len(input) {
					switch input[i] {
					case SPACE_CHAR:
						i++
						break loop
					case S_QUOTE_CHAR:
						i++
						for i < len(input) {
							if input[i] == S_QUOTE_CHAR {
								i++
								break
							}
							arg += string(input[i])
							i++
						}

					default:
						arg += string(input[i])
						i++
					}

				}

				if len(arg) > 0 {
					args = append(args, arg)
				}
			}

			if i == len(input) {
				break
			}
		}

		command, exists := builtins.Builtins[command_name]

		if exists {
			err := command(args)
			if err != nil {
				fmt.Println(err)
			}
		} else if _, err := exec.LookPath(command_name); err == nil {
			cmd := exec.Command(command_name, args...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			cmd.Run()
		} else {
			fmt.Printf("%s: command not found\n", command_name)
		}
	}
}
