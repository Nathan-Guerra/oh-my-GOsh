package lexer

import "testing"

func TestCanCreateWhitespaceToken(t *testing.T) {
	result := Tokenize("   ")
	t.Logf("Output: %v", result)

	assert_length_is(1, result, t)
	assert_token_kind_is(result[0], WHITE_SPACE, t)
	assert_token_value_is(result[0], " ", t)
}

func TestCanCreateLiteralToken(t *testing.T) {
	result := Tokenize("echo")
	t.Logf("Output: %v", result)
	assert_length_is(1, result, t)
	assert_token_kind_is(result[0], LITERAL, t)
	assert_token_value_is(result[0], "echo", t)
}

func TestCanCreateMultipleTokens(t *testing.T) {
	result := Tokenize("echo foo!   bar ")
	t.Logf("Output: %v", result)

	assert_length_is(6, result, t)
	assert_token_kind_is(result[0], LITERAL, t)
	assert_token_kind_is(result[1], WHITE_SPACE, t)
	assert_token_kind_is(result[2], LITERAL, t)
	assert_token_kind_is(result[3], WHITE_SPACE, t)
	assert_token_kind_is(result[4], LITERAL, t)
	assert_token_kind_is(result[5], WHITE_SPACE, t)

}

func TestCanCreateEscapeToken(t *testing.T) {
	result := Tokenize("\\a")
	t.Logf("Output: %v", result)
	assert_length_is(1, result, t)
	assert_token_kind_is(result[0], ESCAPE, t)
	assert_token_value_is(result[0], "a", t)
}

func TestCanCreateEscapeTokenInsideWord(t *testing.T) {
	result := Tokenize("b\\ac")
	t.Logf("Output: %v", result)
	assert_token_kind_is(result[1], ESCAPE, t)
	assert_token_value_is(result[1], "a", t)
}

func TestCanCreateExpandToken(t *testing.T) {
	resultA := Tokenize("$USER")
	assert_length_is(1, resultA, t)
	assert_token_kind_is(resultA[0], EXPAND, t)
	assert_token_value_is(resultA[0], "USER", t)

	resultB := Tokenize("$USERtest")
	assert_length_is(2, resultB, t)
	assert_token_kind_is(resultB[0], EXPAND, t)
	assert_token_value_is(resultB[0], "USER", t)

	resultC := Tokenize("test$USER")
	assert_length_is(2, resultC, t)
	assert_token_kind_is(resultC[1], EXPAND, t)
	assert_token_value_is(resultC[1], "USER", t)
}

func TestCanIdentifyDollarSignAsLiteralOnStringEnds(t *testing.T) {
	resultA := Tokenize("foo$")
	t.Logf("Output: %v", resultA)

	assert_length_is(1, resultA, t)
	assert_token_kind_is(resultA[0], LITERAL, t)
	assert_token_value_is(resultA[0], "foo$", t)

	resultB := Tokenize("foo$ bar")
	t.Logf("Output: %v", resultB)
	assert_length_is(3, resultB, t)
	assert_token_kind_is(resultB[0], LITERAL, t)
	assert_token_value_is(resultB[0], "foo$", t)
}

func TestCanCreateNumericToken(t *testing.T) {
	resultA := Tokenize("432")
	t.Logf("Output: %v", resultA)
	assert_token_kind_is(resultA[0], NUMERIC, t)
	assert_token_value_is(resultA[0], "432", t)

	resultB := Tokenize("432foo")
	t.Logf("Output: %v", resultB)
	assert_token_kind_is(resultB[0], NUMERIC, t)
	assert_token_value_is(resultB[0], "432", t)

}

func TestCanCreateStringLiteralToken(t *testing.T) {
	result := Tokenize("'foo $BAR  123'")
	t.Logf("Output: %v", result)
	assert_token_kind_is(result[0], STRING_LITERAL, t)
	assert_token_value_is(result[0], "foo $BAR  123", t)
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
			get_token_kind(k),
			get_token_kind(tkn.Kind),
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
