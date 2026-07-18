package lexer

import (
	"testing"
)

func TestCanCreateWhitespaceToken(t *testing.T) {
	result := Tokenize("   ")
	t.Logf("Output: %v", result)

	assert_length_is(1, result, t)
	assert_token_kind_is(result[0], Whitespace, t)
	assert_token_value_is(result[0], " ", t)
}

func TestCanCreateLiteralToken(t *testing.T) {
	result := Tokenize("echo")
	t.Logf("Output: %v", result)
	assert_length_is(1, result, t)
	assert_token_kind_is(result[0], Literal, t)
	assert_token_value_is(result[0], "echo", t)
}

func TestCanCreateMultipleTokens(t *testing.T) {
	result := Tokenize("echo foo!   bar ")
	t.Logf("Output: %v", result)

	assert_length_is(6, result, t)
	assert_token_kind_is(result[0], Literal, t)
	assert_token_kind_is(result[1], Whitespace, t)
	assert_token_kind_is(result[2], Literal, t)
	assert_token_kind_is(result[3], Whitespace, t)
	assert_token_kind_is(result[4], Literal, t)
	assert_token_kind_is(result[5], Whitespace, t)

}

func TestCanCreateEscapeToken(t *testing.T) {
	result := Tokenize("\\a")
	t.Logf("Output: %v", result)
	assert_length_is(1, result, t)
	assert_token_kind_is(result[0], Escape, t)
	assert_token_value_is(result[0], "a", t)
}

func TestCanCreateEscapeTokenInsideWord(t *testing.T) {
	result := Tokenize("b\\ac")
	t.Logf("Output: %v", result)
	assert_length_is(3, result, t)
	assert_token_kind_is(result[1], Escape, t)
	assert_token_value_is(result[1], "a", t)
}

func TestCanCreateExpandToken(t *testing.T) {
	resultA := Tokenize("$USER")
	assert_length_is(1, resultA, t)
	assert_token_kind_is(resultA[0], Expand, t)
	assert_token_value_is(resultA[0], "USER", t)

	resultB := Tokenize("$USERtest")
	assert_length_is(2, resultB, t)
	assert_token_kind_is(resultB[0], Expand, t)
	assert_token_value_is(resultB[0], "USER", t)

	resultC := Tokenize("test$USER")
	assert_length_is(2, resultC, t)
	assert_token_kind_is(resultC[1], Expand, t)
	assert_token_value_is(resultC[1], "USER", t)

	resultD := Tokenize("$$")
	t.Logf("Output: %v", resultD)
	assert_length_is(1, resultD, t)
	assert_token_kind_is(resultD[0], Expand, t)
	assert_token_value_is(resultD[0], "$", t)

	resultE := Tokenize("$FOO$BAR")
	t.Logf("Output: %v", resultE)
	assert_length_is(2, resultE, t)
	assert_token_kind_is(resultE[0], Expand, t)
	assert_token_kind_is(resultE[1], Expand, t)
	assert_token_value_is(resultE[0], "FOO", t)
	assert_token_value_is(resultE[1], "BAR", t)
}

func TestCanIdentifyDollarSignAsLiteralOnStringEnds(t *testing.T) {
	resultA := Tokenize("foo$")
	t.Logf("Output: %v", resultA)
	assert_length_is(2, resultA, t)
	assert_token_kind_is(resultA[0], Literal, t)
	assert_token_value_is(resultA[0], "foo", t)
	assert_token_kind_is(resultA[1], Literal, t)
	assert_token_value_is(resultA[1], "$", t)

	resultB := Tokenize("foo$ bar")
	t.Logf("Output: %v", resultB)
	assert_length_is(4, resultB, t)
	assert_token_kind_is(resultB[0], Literal, t)
	assert_token_value_is(resultB[0], "foo", t)
	assert_token_kind_is(resultB[1], Literal, t)
	assert_token_value_is(resultB[1], "$", t)

	resultC := Tokenize("$")
	t.Logf("Output: %v", resultC)
	assert_length_is(1, resultC, t)
	assert_token_kind_is(resultC[0], Literal, t)
	assert_token_value_is(resultC[0], "$", t)
}

func TestCanCreateNumericToken(t *testing.T) {
	resultA := Tokenize("432")
	t.Logf("Output: %v", resultA)
	assert_token_kind_is(resultA[0], Numeric, t)
	assert_token_value_is(resultA[0], "432", t)

	resultB := Tokenize("432foo")
	t.Logf("Output: %v", resultB)
	assert_token_kind_is(resultB[0], Numeric, t)
	assert_token_value_is(resultB[0], "432", t)

	resultC := Tokenize("foo00bar")
	t.Logf("Output: %v", resultC)
	assert_length_is(1, resultC, t)
	assert_token_kind_is(resultC[0], Literal, t)

}

func TestCanCreateStringLiteralToken(t *testing.T) {
	result := Tokenize("'foo $BAR  123'")
	t.Logf("Output: %v", result)
	assert_length_is(1, result, t)
	assert_token_kind_is(result[0], StringLiteral, t)
	assert_token_value_is(result[0], "foo $BAR  123", t)
}

func TestCanCreateStringExpandToken(t *testing.T) {
	result := Tokenize("\"foo $BAR  123\"")
	t.Logf("Output: %v", result)
	assert_length_is(1, result, t)
	assert_token_kind_is(result[0], StringExpand, t)
	assert_token_value_is(result[0], "foo $BAR  123", t)
}

func TestCanCreateComplexStringExpandToken(t *testing.T) {
	result := Tokenize("\"foo $BAR\\\"  123\"")
	t.Logf("Output: %v", result)
	assert_length_is(1, result, t)
	assert_token_kind_is(result[0], StringExpand, t)
	assert_token_value_is(result[0], "foo $BAR\\\"  123", t) // -> literal string (\")
}

func TestStringExpandInnerTokenizer(t *testing.T) {
	result := Tokenize("\"'foo' \\$BAR  123\"")
	assert_length_is(1, result, t)
	assert_token_kind_is(result[0], StringExpand, t)
	assert_token_value_is(result[0], "'foo' \\$BAR  123", t)
	if result[0].TokenizedValue == nil {
		t.Error("Error tokenizing expanded string.")
	}
	t.Logf("Output: %v", (*result[0].TokenizedValue))
	assert_length_is(6, *result[0].TokenizedValue, t)
	assert_token_kind_is((*result[0].TokenizedValue)[0], Literal, t)
	assert_token_kind_is((*result[0].TokenizedValue)[1], Whitespace, t)
	assert_token_kind_is((*result[0].TokenizedValue)[2], Escape, t)
	assert_token_kind_is((*result[0].TokenizedValue)[3], Literal, t)
	assert_token_kind_is((*result[0].TokenizedValue)[4], Whitespace, t)
	assert_token_kind_is((*result[0].TokenizedValue)[5], Numeric, t)

	assert_token_value_is((*result[0].TokenizedValue)[0], "'foo'", t)
	assert_token_value_is((*result[0].TokenizedValue)[1], " ", t)
	assert_token_value_is((*result[0].TokenizedValue)[2], "$", t)
	assert_token_value_is((*result[0].TokenizedValue)[3], "BAR", t)
	assert_token_value_is((*result[0].TokenizedValue)[4], "  ", t)
	assert_token_value_is((*result[0].TokenizedValue)[5], "123", t)

}

func TestCanCreateRedirectOutToken(t *testing.T) {
	resultA := Tokenize(">")
	t.Logf("Output: %v", resultA)
	assert_length_is(1, resultA, t)
	assert_token_kind_is(resultA[0], RedirectOut, t)
	assert_token_value_is(resultA[0], ">", t)

	resultB := Tokenize("1>")
	t.Logf("Output: %v", resultB)
	assert_length_is(1, resultB, t)
	assert_token_kind_is(resultB[0], RedirectOut, t)
	assert_token_value_is(resultB[0], "1>", t)

	resultC := Tokenize("1 >")
	t.Logf("Output: %v", resultC)
	assert_length_is(3, resultC, t)
	assert_token_kind_is(resultC[0], Numeric, t)
	assert_token_kind_is(resultC[1], Whitespace, t)
	assert_token_kind_is(resultC[2], RedirectOut, t)
}

func TestCanCreateRedirectErrToken(t *testing.T) {
	result := Tokenize("2>")
	t.Logf("Output: %v", result)
	assert_length_is(1, result, t)
	assert_token_kind_is(result[0], RedirectErr, t)
	assert_token_value_is(result[0], "2>", t)

	resultB := Tokenize("2 >")
	t.Logf("Output: %v", resultB)
	assert_length_is(3, resultB, t)
	assert_token_kind_is(resultB[0], Numeric, t)
	assert_token_kind_is(resultB[1], Whitespace, t)
	assert_token_kind_is(resultB[2], RedirectOut, t)
}

func TestCanTokenizeComplexInput(t *testing.T) {
	result := Tokenize("echo \"foo   bar\"qux 1000 1>/path/to/file.txt")
	t.Logf("Output: %v", result)
	assert_length_is(9, result, t)
	assert_token_kind_is(result[0], Literal, t)
	assert_token_kind_is(result[1], Whitespace, t)
	assert_token_kind_is(result[2], StringExpand, t)
	assert_token_kind_is(result[3], Literal, t)
	assert_token_kind_is(result[4], Whitespace, t)
	assert_token_kind_is(result[5], Numeric, t)
	assert_token_kind_is(result[6], Whitespace, t)
	assert_token_kind_is(result[7], RedirectOut, t)
	assert_token_kind_is(result[8], Literal, t)
}

func assert_length_is(i int, r []Token, t *testing.T) {
	if len(r) != i {
		t.Errorf("Expected %d token, %d tokens received.", i, len(r))
	}
}

func assert_token_kind_is(tkn Token, k TokenKind, t *testing.T) {
	if tkn.Kind != k {
		t.Errorf(
			"Expected token kind to be \"%s\", got \"%s\".",
			k.String(),
			tkn.Kind.String(),
		)
	}
}

func assert_token_value_is(tkn Token, v string, t *testing.T) {
	if tkn.Value != v {
		t.Errorf(
			"Expected token value to be \"%s\", got \"%s\".",
			v,
			tkn.Value,
		)
	}
}
