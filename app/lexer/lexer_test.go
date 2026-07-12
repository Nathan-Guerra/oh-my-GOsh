package lexer

import "testing"

func TestCanCreateWhitespaceToken(t *testing.T) {
	result := Tokenize("   ")
	t.Logf("Output: %v", result)
	if len(result) != 1 {
		t.Errorf("Expected 1 token, %d tokens received.", len(result))
	}

	if result[0].Kind != WHITE_SPACE {
		t.Errorf(
			"Expected token kind to be %s, got %s.",
			get_token_kind(WHITE_SPACE),
			get_token_kind(result[0].Kind),
		)
	}
}
