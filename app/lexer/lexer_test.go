package lexer

import "testing"

func TestCanCreateWhitespaceToken(t *testing.T) {
	result := Tokenize("   ")
	t.Logf("Output: %v", result)

	assert_length_is(1, result, t)
	assert_token_kind_is(result[0], WHITE_SPACE, t)
}

func TestCanCreateLiteralToken(t *testing.T) {
	result := Tokenize("echo")
	t.Logf("Output: %v", result)
	assert_length_is(1, result, t)
	assert_token_kind_is(result[0], LITERAL, t)

	if result[0].Value != "echo" {
		t.Errorf("Escaped string should be 'echo', '%s' received.", result[0].Value)
	}
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

	if result[0].Value != "a" {
		t.Errorf("Escaped string should be 'a', '%s' received.", result[0].Value)
	}
}

func TestCanCreateEscapeTokenInsideWord(t *testing.T) {
	result := Tokenize("b\\ac")
	t.Logf("Output: %v", result)
	assert_token_kind_is(result[1], ESCAPE, t)

	if result[1].Value != "a" {
		t.Errorf("Escaped string should be 'a', '%s' received.", result[0].Value)
	}
}

func TestCanCreateExpandToken(t *testing.T) {
	resultA := Tokenize("$USER")
	assert_length_is(1, resultA, t)
	assert_token_kind_is(resultA[0], EXPAND, t)
}

func assert_length_is(i int, r []Token, t *testing.T) {
	if len(r) != i {
		t.Errorf("Expected %d token, %d tokens received.", i, len(r))
	}
}

func assert_token_kind_is(tkn Token, k TokenKind, t *testing.T) {
	if tkn.Kind != k {
		t.Errorf(
			"Expected token kind to be %s, got %s.",
			get_token_kind(k),
			get_token_kind(tkn.Kind),
		)
	}
}
