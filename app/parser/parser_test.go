package parser

import (
	"os"
	"strconv"
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
	tkns := lexer.Tokenize("\"'echo' \\ \\\"$PWD\\\" > /foo\"\\ > /dev/null")
	t.Logf("Output: %v", tkns)
	cmd := CreateCommand(tkns)
	t.Logf("Output: %v", cmd)
	if cmd.CommandName != "'echo'" {
		t.Errorf("Expected name 'echo', '%s' given.", cmd.CommandName)
	}

	if len(cmd.Arguments) != 5 {
		t.Errorf("Expected exatcly 5 arguments, (%d) given.", len(cmd.Arguments))
	}
}

func TestParseCommandWithRedirectOut(t *testing.T) {
	cmd := CreateCommand(lexer.Tokenize("echo foo > /dev/null"))

	if cmd.CommandName != "echo" {
		t.Errorf("Expected name 'echo', '%s' given.", cmd.CommandName)
	}

	if len(cmd.Arguments) != 1 {
		t.Errorf("Expected exatcly 1 argument, (%d) given.", len(cmd.Arguments))
	}

	f, ok := cmd.Stdout.(*os.File)
	t.Logf("Current: '%s' | Stdout: '%s'.", f.Name(), os.Stdout.Name())
	if !ok {
		t.Errorf("Expected redirecting stdout to a file")
	} else if f.Name() == os.Stdout.Name() {
		t.Errorf("Expected output file to be '/dev/null', '%s' given.", os.Stdout.Name())
	}

}

func TestParseCommandWithDirAsRedirectOut(t *testing.T) {
	cmd := CreateCommand(lexer.Tokenize("echo foo > ./$$"))

	if cmd.CommandName != "echo" {
		t.Errorf("Expected name 'echo', '%s' given.", cmd.CommandName)
	}

	if len(cmd.Arguments) != 1 {
		t.Errorf("Expected exatcly 1 argument, (%d) given.", len(cmd.Arguments))
	}

	pid := os.Getpid()
	info, err := os.Stat(strconv.Itoa(pid))
	if err != nil {
		t.Error(err)
	}

	t.Logf("File: %s, created at: %s", info.Name(), info.ModTime().Local())
	os.Remove(strconv.Itoa(pid))
}
