package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"slices"
	"strings"

	"golang.org/x/term"

	"github.com/codecrafters-io/shell-starter-go/app/autocomplete"
	"github.com/codecrafters-io/shell-starter-go/app/builtins"
	"github.com/codecrafters-io/shell-starter-go/app/keyboard"
	"github.com/codecrafters-io/shell-starter-go/app/lexer"
	"github.com/codecrafters-io/shell-starter-go/app/parser"
)

func main() {
	var lastKey byte
	var autocompleter autocomplete.Autocompleter
	autocompleter.SetBuiltins(slices.Collect(maps.Keys(builtins.Builtins)))
	autocompleter.SetPATH(os.Getenv("PATH"))
	autocompleter.EagerLoadPathCommands()
	reader := bufio.NewReader(os.Stdin)
	var buffer strings.Builder
termLoop:
	for {
		buffer.Reset()
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			panic(err)
		}
		defer term.Restore(int(os.Stdin.Fd()), oldState)

		fmt.Print("$ ")
	lineLoop:
		for {
			input, err := reader.ReadByte()
			if err != nil {
				panic(err)
			}

			switch input {
			case keyboard.Tab:
				old := autocompleter.Retrieve()

				if len(old) >= 2 && lastKey == keyboard.Tab {
					fmt.Printf("\r\n%s\r\n", strings.Join(old, "  "))
					fmt.Printf("$ %s", buffer.String())
					goto lineLoop
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

			case keyboard.Enter:
				fmt.Print("\r\n")
				break lineLoop
			case keyboard.CtrlC:
				buffer.Reset()
				os.Stdin.Write([]byte{})
				goto termLoop
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
				}

				fmt.Printf("%c", input)
			}

			lastKey = input
		}

		cmd := parser.CreateCommand(lexer.Tokenize(string(buffer.String())))
		if cmd.CommandName == "" {
			continue
		}

		if cmd.CommandName == "exit" {
			goto exit
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
		} else if _, err := exec.LookPath(cmd.CommandName); err == nil {
			externalCommand := exec.Command(cmd.CommandName, cmd.Arguments...)
			externalCommand.Stdin = os.Stdin
			externalCommand.Stdout = cmd.Stdout
			externalCommand.Stderr = cmd.Stderr

			externalCommand.Run()
		} else {
			fmt.Printf("%s: command not found\r\n", cmd.CommandName)
		}
	}
exit:
}
