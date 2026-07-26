package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/autocomplete"
	"github.com/codecrafters-io/shell-starter-go/app/keyboard"
	"github.com/codecrafters-io/shell-starter-go/app/lexer"
	"github.com/domonda/go-pretty"
	"golang.org/x/term"
)

const Prompt string = "\r$ "

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
	autocomplete.GetCommandAutocompleter().Clear()
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
			var matches []string
			replacements := 0
			oldMatches := autocomplete.GetCommandAutocompleter().Retrieve()
			// fmt.Printf("\r\n%v\r\n", oldMatches)
			if len(oldMatches) > 0 {
				matches = oldMatches
			} else {
				matches, replacements = autocomplete.GetCommandAutocompleter().Match(string(line))
			}
			switch len(matches) {
			case 0:
				fmt.Printf("%c", keyboard.Bell)
			case 1:
				line = slices.Delete(line, len(line)-replacements, len(line))
				line = append(line, matches[0]...)
				for i := 0; i < replacements; i++ {
					fmt.Print("\b")
				}
				fmt.Print(matches[0])

				autocomplete.GetCommandAutocompleter().Clear()
				oldMatches = make([]string, 0)
				buffer[0] = keyboard.Null
			default: // >= 2
				largestPrefix := autocomplete.GetCommandAutocompleter().LargestCommonPrefix(string(line))
				if len(largestPrefix) > 0 {
					line = append(line, largestPrefix...)
					fmt.Print(string(largestPrefix))
					// lastKey = keyboard.Null // erase last byte to prevent next iteration to always show the list of elements
				} else if lastKey != keyboard.Tab {
					fmt.Printf("%c", keyboard.Bell)
				} else {
					fmt.Println("\r")
					fmt.Println(strings.Join(oldMatches, " "))
					fmt.Printf("%s%s", Prompt, string(line))
				}
			}
		case keyboard.Backspace:
			if len(line) > 0 {
				line = line[0 : len(line)-1]
				// back one char, erase (prints " ")
				// go back again (cursor is over the space char, looking like it deleted the content)
				fmt.Print("\b \b")
				autocomplete.GetCommandAutocompleter().Clear()
			}
		default:
			// printable characters
			if buffer[0] >= keyboard.Space && buffer[0] <= keyboard.Tilde {
				line = append(line, buffer[0])
				fmt.Printf("%c", buffer[0])
				autocomplete.GetCommandAutocompleter().Clear()
			}
		}

		lastKey = buffer[0]

		writter, err := os.OpenFile("logs/debug.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Println("cannot debug")
		} else {
			lexer := lexer.Tokenize(line)
			pretty.Fprint(writter, lexer, "    ")
		}
	}

	return line
}
