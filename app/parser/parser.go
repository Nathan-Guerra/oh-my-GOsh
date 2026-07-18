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
					// if cmd.CommandName != "" {
					arg.WriteString(token.Value)
					// } else {
					// 	cmd.push(s.String())
					// 	arg.Reset()
					// }
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

			// if len(arg.String()) != 0 {
			// 	cmd.push(arg.String())
			// 	arg.Reset()
			// }
		case lexer.RedirectOut:
			var value strings.Builder
		redirectLoop:
			for _, subToken := range tokens[i+1:] {
				switch subToken.Kind {
				case lexer.Whitespace:
					if len(value.String()) > 0 {
						break redirectLoop
					}

					continue redirectLoop
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

			fi, err := os.Stat(value.String())

			if err != nil { // error, try to create file
				newFile, err := os.Create(value.String())
				if err != nil {
					panic(err)
				}

				cmd.Stdout = newFile
			} else {
				if fi.IsDir() {
					panic("==Error== Cannot write to a directory.")
				}

				file, err := os.OpenFile(value.String(), os.O_WRONLY, 0666)
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
