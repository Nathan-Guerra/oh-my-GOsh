package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
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

		params := strings.Split(input, " ")
		args := params[1:]
		command, exists := builtins.Builtins[params[0]]

		if !exists {
			fmt.Printf("%s: command not found\n", params[0])
			continue
		}

		err := command(args)
		if err != nil {
			fmt.Println(err)
		}
	}
}
