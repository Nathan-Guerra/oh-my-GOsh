package parser

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/lexer"
)

type Command struct {
	CommandName string
	Arguments   []string
	Stdout      io.Writer
	Stderr      io.Writer
}

func (c *Command) push(s string) {
	if len(c.CommandName) == 0 {
		if strings.Trim(s, " ") != "" {
			c.CommandName = s
		}
	} else {
		c.Arguments = append(c.Arguments, s)
	}
}

func CreateCommand(tokens []lexer.Token) *Command {
	var arg strings.Builder
	cmd := &Command{}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

parsingLoop:
	for i, token := range tokens {
		switch token.Kind {
		case lexer.Whitespace:
			if len(arg.String()) > 0 {
				cmd.push(arg.String())
				arg.Reset()
			}
		case lexer.Literal,
			lexer.Numeric,
			lexer.StringLiteral,
			lexer.Escape:
			arg.WriteString(token.Value)
		case lexer.Expand:
			if token.Value == "$" {
				arg.WriteString(strconv.Itoa(os.Getpid()))
			} else {
				arg.WriteString(os.Getenv(token.Value))
			}
		case lexer.StringExpand:
			// already have a command, treat as a string
			for _, token := range *token.TokenizedValue {
				switch token.Kind {
				case lexer.Whitespace:
					arg.WriteString(token.Value)
				case lexer.Literal,
					lexer.Numeric,
					lexer.Escape:
					arg.WriteString(token.Value)
				case lexer.Expand:
					if token.Value == "$" {
						arg.WriteString(strconv.Itoa(os.Getpid()))
					} else {
						arg.WriteString(os.Getenv(token.Value))
					}
				default:
					panic(fmt.Sprintf("==Error== Inner token kind not identified {%s}.", token.Kind))
				}
			}
		case lexer.RedirectOut:
			value := findFileName(tokens[i+1:])
			fi, err := os.Stat(value)

			if err != nil { // error, try to create file
				newFile, err := os.Create(value)
				if err != nil {
					panic(err)
				}

				cmd.Stdout = newFile
			} else {
				if fi.IsDir() {
					panic("==Error== Cannot write to a directory.")
				}

				file, err := os.OpenFile(value, os.O_WRONLY|os.O_TRUNC, 0666)
				if err != nil {
					panic(err)
				}
				cmd.Stdout = file
			}

			break parsingLoop
		case lexer.RedirectErr:
			value := findFileName(tokens[i+1:])
			fi, err := os.Stat(value)

			if err != nil { // error, try to create file
				newFile, err := os.Create(value)
				if err != nil {
					panic(err)
				}

				cmd.Stderr = newFile
			} else {
				if fi.IsDir() {
					panic("==Error== Cannot write to a directory.")
				}

				file, err := os.OpenFile(value, os.O_WRONLY, 0666)
				if err != nil {
					panic(err)
				}
				cmd.Stderr = file
			}

			break parsingLoop
		case lexer.RedirectOutAppend:
			value := findFileName(tokens[i+1:])
			fi, err := os.Stat(value)

			if err != nil { // error, try to create file
				newFile, err := os.Create(value)
				if err != nil {
					panic(err)
				}

				cmd.Stdout = newFile
			} else {
				if fi.IsDir() {
					panic("==Error== Cannot write to a directory.")
				}

				file, err := os.OpenFile(value, os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					panic(err)
				}
				cmd.Stdout = file
			}

			break parsingLoop
		default:
			panic(fmt.Sprintf("==Error== Token kind not identified {%s}.", token.Kind))
		}
	}

	if len(arg.String()) != 0 {
		cmd.push(arg.String())
	}

	return cmd
}

func findFileName(tokens []lexer.Token) string {
	var value strings.Builder
loop:
	for _, subToken := range tokens {
		switch subToken.Kind {
		case lexer.Whitespace:
			if len(value.String()) > 0 {
				break loop
			}

			continue loop
		case lexer.Expand:
			if subToken.Value == "$" {
				value.WriteString(strconv.Itoa(os.Getpid()))
			} else {
				value.WriteString(os.Getenv(subToken.Value))
			}
		case lexer.Literal, lexer.StringLiteral,
			lexer.Numeric,
			lexer.Escape:
			value.WriteString(subToken.Value)
		}
	}

	return value.String()
}
