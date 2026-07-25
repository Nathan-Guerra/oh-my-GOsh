package main

import (
	"fmt"
	"maps"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/autocomplete"
	"github.com/codecrafters-io/shell-starter-go/app/builtins"
	"github.com/codecrafters-io/shell-starter-go/app/keyboard"
	"github.com/codecrafters-io/shell-starter-go/app/lexer"
	"github.com/codecrafters-io/shell-starter-go/app/parser"
	"golang.org/x/term"
)

var prompt string = "\r$ "

func readByte() ([]byte, error) {
	buffer := make([]byte, 1)
	n, err := os.Stdin.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("cannot read byte: %s", err.Error())
	}
	return buffer[:n], nil
}

var stdinFd int = int(os.Stdin.Fd())

func readLine() []byte {
	old, err := term.MakeRaw(stdinFd)
	if err != nil {
		panic(err)
	}
	defer term.Restore(stdinFd, old)

	var lastKey byte
	var line []byte

loop:
	for {
		buffer, err := readByte()
		if err != nil {
			panic(err)
		}

		switch buffer[0] {
		case keyboard.EnterCR, keyboard.Enter:
			fmt.Println("\r")
			break loop
		case keyboard.Tab:
			old := autocomplete.GetCommandAutocompleter().Retrieve()
			if len(old) >= 2 && lastKey == keyboard.Tab {
				fmt.Println("\r")
				fmt.Println(strings.Join(old, "  "))
				fmt.Printf("%s%s", prompt, string(line))
				continue loop
			} else {
				matches := autocomplete.GetCommandAutocompleter().Match(string(line))
				if len(matches) == 1 {
					size := len(line)
					suffix := matches[0][size:] + " "
					line = append(line, suffix...)
					fmt.Print(suffix)
				} else {
					fmt.Printf("%c", keyboard.Bell)
				}
			}

		// case keyboard.CtrlC:
		// 	fmt.Printf("%c", input)
		// 	fmt.Print("$ ")
		// 	buffer.Reset()
		// 	continue loop
		case keyboard.Backspace:
			if len(line) > 0 {
				line = line[0 : len(line)-1]
				// back one char, erase (prints " ")
				// go back again (cursor is over the space char, looking like it deleted the content)
				fmt.Print("\b \b")
			}
		default:
			// printable characters
			if buffer[0] >= keyboard.Space && buffer[0] <= keyboard.Tilde {
				line = append(line, buffer[0])
				fmt.Printf("%c", buffer[0])
			}
		}

		lastKey = buffer[0]
	}

	return line
}

func REPL() int {
	autocompleter := autocomplete.GetCommandAutocompleter()
	autocompleter.SetBuiltins(slices.Collect(maps.Keys(builtins.Builtins)))
	autocompleter.SetPATH(os.Getenv("PATH"))
	autocompleter.EagerLoad()

	var code int

	for {
		fmt.Print(prompt)
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
