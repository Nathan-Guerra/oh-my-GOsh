package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	// TODO: Uncomment the code below to pass the first stage
	for {
		var input string
		fmt.Print("$ ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
			continue
		}

		input = strings.TrimSpace(scanner.Text())
		inputArr := strings.Split(input, " ")
		command := inputArr[0]
		args := inputArr[1:]

		if command == "" {
			continue
		} else if command == "exit" {
			break
		} else if command == "echo" {
			var out string
			for i := 0; i < len(args); i++ {
				if args[i] == "" {
					continue
				}
				out += args[i] + " "
			}

			if out != "" {
				out = out[:len(out)-1]
			}
			fmt.Printf("%s\n", out)
			continue
		} else if command == "type" {
			for _, command_name := range args {
				if command_name == "" {
					continue
				}
				var found bool
				for _, v := range []string{"echo", "exit", "type"} {
					if v == command_name {
						found = true
						break
					}
				}

				if found {
					fmt.Printf("%s is a shell builtin\n", command_name)
				} else {
					fmt.Printf("%s: not found\n", command_name)
				}
			}

			continue
		}

		fmt.Printf("%s: command not found\n", command)
	}
}
