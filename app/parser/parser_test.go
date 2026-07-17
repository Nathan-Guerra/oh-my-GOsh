package parser

import (
	"testing"

	"github.com/codecrafters-io/shell-starter-go/app/lexer"
)

func TestParseSimpleCommand(t *testing.T) {
	tkns := lexer.Tokenize("echo foo")
	t.Logf("Output: %v", tkns)
	cmd := CreateCommand(tkns)
	t.Logf("Output: %v", cmd)
	if cmd.CommandName != "echo" {
		t.Errorf("Expected name 'echo', '%s' given.", cmd.CommandName)
	}
	if len(cmd.Arguments) != 1 {
		t.Errorf("Expected exatcly 1 argument, (%d) given.", len(cmd.Arguments))
	}
}

func TestParseNotSoSimpleCommand(t *testing.T) {
	tkns := lexer.Tokenize("\"echo \\ \\\"$PWD\\\"\"")
	t.Logf("Output: %v", tkns)
	cmd := CreateCommand(tkns)
	t.Logf("Output: %v", cmd)
	if cmd.CommandName != "echo" {
		t.Errorf("Expected name 'echo', '%s' given.", cmd.CommandName)
	}
	if len(cmd.Arguments) != 1 {
		t.Errorf("Expected exatcly 1 argument, (%d) given.", len(cmd.Arguments))
	}
}
