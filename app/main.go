package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/builtins"
)

func main() {
	const (
		SPACE_RUNE       = ' '
		S_QUOTE_RUNE     = '\''
		D_QUOTE_RUNE     = '"'
		DOLLAR_SIGN_RUNE = '$'
		BACKSLASH_RUNE   = '\\'
	)

	scanner := bufio.NewScanner(os.Stdin)
	regexpLetters, err := regexp.Compile("[a-zA-Z]")
	if err != nil {
		fmt.Println(err)
	}
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
			var arg string
		loop:
			for i < len(input) {
				switch input[i] {
				case BACKSLASH_RUNE:
					i++
					arg += string(input[i])
					i++
					break
				case DOLLAR_SIGN_RUNE:
					i++
					if i >= len(input) || input[i] == SPACE_RUNE {
						// ...$ <- at the end
						arg += string(DOLLAR_SIGN_RUNE)
						break
					}

					if input[i] == DOLLAR_SIGN_RUNE {
						// $$ <- pid
						arg += strconv.Itoa(os.Getpid())
						i++
						break
					}

					var env_var string
					for i < len(input) {
						s := string(input[i])
						if regexpLetters.MatchString(s) {
							env_var += s
							i++
						} else {
							break
						}
					}

					v := os.Getenv(env_var)
					if len(v) > 0 {
						arg += v
					}
				case SPACE_RUNE:
					i++
					break loop
				case S_QUOTE_RUNE:
					i++
					for i < len(input) {
						if input[i] == S_QUOTE_RUNE {
							i++
							break
						}
						arg += string(input[i])
						i++
					}
				case D_QUOTE_RUNE:
					i++
					escapable_runes := []rune{D_QUOTE_RUNE, DOLLAR_SIGN_RUNE, BACKSLASH_RUNE}
					for i < len(input) {
						if input[i] == D_QUOTE_RUNE {
							i++
							break
						}

						if input[i] == BACKSLASH_RUNE {
							i++
							if slices.Contains(escapable_runes[:], rune(input[i])) {
								arg += string(input[i])
							} else {
								arg += string(BACKSLASH_RUNE) + string(input[i])
							}
							i++
							continue
						}

						if input[i] == DOLLAR_SIGN_RUNE {
							i++
							if input[i] == D_QUOTE_RUNE || input[i] == SPACE_RUNE {
								// ...$ <- at the end of string or word
								arg += string(DOLLAR_SIGN_RUNE)
								// do not increment to let outer loop handle the rune
								continue
							}

							if input[i] == DOLLAR_SIGN_RUNE {
								// $$ <- pid
								arg += strconv.Itoa(os.Getpid())
								i++
								continue
							}

							var env_var string
							for i < len(input) {
								s := string(input[i])
								if regexpLetters.MatchString(s) {
									env_var += s
									i++
								} else {
									break
								}
							}

							v := os.Getenv(env_var)
							if len(v) > 0 {
								arg += v
							}
						} else {
							arg += string(input[i])
							i++
						}
					}

				default:
					arg += string(input[i])
					i++
				}
			}

			if !found_command {
				found_command = true
				command_name = arg
			} else if len(arg) > 0 {
				args = append(args, arg)
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
