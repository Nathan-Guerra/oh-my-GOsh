package parser

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/lexer"
)

type Command struct {
	CommandName string
	Arguments   []string
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
	for _, token := range tokens {
		switch token.Kind {
		case lexer.WHITE_SPACE:
			cmd.push(arg.String())
			arg.Reset()
		case lexer.LITERAL,
			lexer.NUMERIC,
			lexer.STRING_LITERAL,
			lexer.ESCAPE:
			arg.WriteString(token.Value)
		case lexer.EXPAND:
			if token.Value == "$" {
				arg.WriteString(strconv.Itoa(os.Getpid()))
			} else {
				arg.WriteString(os.Getenv(token.Value))
			}
		case lexer.STRING_EXPAND:
			// already have a command, treat as a string
			var s strings.Builder
			for _, token := range *token.TokenizedValue {
				switch token.Kind {
				case lexer.WHITE_SPACE:
					if cmd.CommandName != "" {
						s.WriteString(token.Value)
					} else {
						cmd.push(s.String())
						s.Reset()
					}
				case lexer.LITERAL,
					lexer.NUMERIC,
					lexer.ESCAPE:
					s.WriteString(token.Value)
				case lexer.EXPAND:
					if token.Value == "$" {
						s.WriteString(strconv.Itoa(os.Getpid()))
					} else {
						s.WriteString(os.Getenv(token.Value))
					}
				default:
					panic(fmt.Sprintf("==Error== Token kind not identified {%s}.", token.Kind))
				}
			}

			if len(s.String()) != 0 {
				cmd.push(s.String())
			}
		default:
			panic(fmt.Sprintf("==Error== Token kind not identified {%s}.", token.Kind))
		}
	}

	if len(arg.String()) != 0 {
		cmd.push(arg.String())
	}
	return cmd
}
