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

		if command == "" {
			continue
		} else if command == "exit" {
			break
		} else if command == "echo" {
			var out string
			for i := 1; i < len(inputArr); i++ {
				if inputArr[i] == "" {
					continue
				}
				out += inputArr[i] + " "
			}

			if out != "" {
				out = out[:len(out)-1]
			}
			fmt.Printf("%s\n", out)
			continue
		}

		fmt.Printf("%s: command not found\n", command)
	}
}
