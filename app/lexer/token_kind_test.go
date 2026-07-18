package lexer

import "testing"

func TestCanPrintToken(t *testing.T) {
	var tkn TokenKind
	tkn = 1 // whitespace

	if tkn.String() != "WHITESPACE" {
		t.Errorf("Cannot stringify token name, (%s) given, (%s) expected.", tkn.String(), "WHITESPACE")
	}

	var unknown TokenKind = 999
	if unknown.String() != "999" {
		t.Errorf("Cannot stringify token name, (%s) given, (%s) expected.", tkn.String(), "999")
	}
}
