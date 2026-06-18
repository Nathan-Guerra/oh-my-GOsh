package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
					if input[i] == ' ' {
						found_command = true
						break
					}

					command_name += string(input[i])
					i++
				}
			} else {
				var arg string
				for i < len(input) {
					if input[i] == ' ' {
						i++
						break
					}

					arg += string(input[i])
					i++
				}

				if len(arg) > 0 {
					args = append(args, arg)
				}
			}

			if i == len(input) {
				break
			}
		}

		fmt.Println(command_name)
		fmt.Printf("%v\n", args)

		// params := strings.Split(input, " ")
		// // args := params[1:]
		// command, exists := builtins.Builtins[params[0]]

		// if exists {
		// 	err := command(args)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}
		// } else if _, err := exec.LookPath(params[0]); err == nil {
		// 	cmd := exec.Command(params[0], args...)
		// 	cmd.Stdin = os.Stdin
		// 	cmd.Stdout = os.Stdout
		// 	cmd.Stderr = os.Stderr

		// 	cmd.Run()
		// } else {
		// 	fmt.Printf("%s: command not found\n", params[0])
		// }
	}
}
