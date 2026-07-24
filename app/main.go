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
)

func buffer() string {
	var lastKey byte
	var line []byte
	buffer := make([]byte, 1)

	autocompleter := autocomplete.GetCommandAutocompleter()
	autocompleter.SetBuiltins(slices.Collect(maps.Keys(builtins.Builtins)))
	autocompleter.SetPATH(os.Getenv("PATH"))
	autocompleter.EagerLoad()

	// reader := bufio.NewReader(os.Stdin)
	// oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	// if err != nil {
	// 	panic(err)
	// }
	// defer term.Restore(int(os.Stdin.Fd()), oldState)

loop:
	for {
		n, err := os.Stdin.Read(buffer)
		if err != nil {
			panic(err)
		}

		switch buffer[0] {
		case keyboard.EnterCR, keyboard.Enter:
			fmt.Println("")
			break loop
		case keyboard.Tab:
			old := autocompleter.Retrieve()
			if len(old) >= 2 && lastKey == keyboard.Tab {
				fmt.Printf("\n%s", strings.Join(old, "  "))
				fmt.Printf("\n$ %s", buffer.String())
				continue loop
			} else {
				matches := autocompleter.Match(buffer.String())
				if len(matches) == 1 {
					size := len(buffer.String())
					suffix := matches[0][size:] + " "
					buffer.Write([]byte(suffix))
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
			if len(buffer.String()) > 0 {
				v := buffer.String()
				buffer.Reset()
				buffer.WriteString(v[:len(v)-1])
				// back one char, erase (prints " ")
				// go back again (cursor is over the space char, looking like it deleted the content)
				fmt.Print("\b \b")
			}
		default:
			// printable characters
			if input >= keyboard.Space && input <= keyboard.Tilde {
				buffer.WriteByte(input)
				fmt.Printf("%c", input)
			}

		}

		lastKey = input
	}

	return buffer.String()
}

func main() {
	var code int

	for {
		fmt.Print("$ ")
		buffer := buffer()

		cmd := parser.CreateCommand(lexer.Tokenize(buffer))
		if cmd.CommandName == "" {
			continue
		}

		command, exists := builtins.Builtins[cmd.CommandName]
		if exists {
			response := command.Exec(cmd.Arguments)
			if len(response.Out) > 0 {
				// cmd.Stdout.Write([]byte(response.Out))
				n, err := cmd.Stdout.Write([]byte(response.Out))
				fmt.Fprintf(os.Stderr, "WRITE: n=%d err=%v out=%q\n", n, err, response.Out)
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
			fmt.Printf("%s: command not found\n", cmd.CommandName)
		}
	}

	os.Exit(code)
}
