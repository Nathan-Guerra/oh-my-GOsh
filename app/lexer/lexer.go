package lexer

type TokenKind int

type Token struct {
	Kind  TokenKind
	Value string
}

const (
	WHITE_SPACE TokenKind = iota
	LITERAL
	NUMBER
	STRING_LITERAL
	STRING_EXPAND
	EXPAND
	ESCAPE
	REDIRECT_OUT
	REDIRECT_ERR
)

func Tokenize(input string) []Token {
	i := 0
	tokens := make([]Token, 0)
	for i < len(input) {
		var lookup byte
		var has_lookup bool

		if i+1 < len(input) {
			lookup = input[i+1]
			has_lookup = true
		}
		kind := which_kind(input[i], lookup, has_lookup)
		value := get_value_from(kind, &i, input)
		tokens = append(tokens, Token{kind, value})
	}

	return tokens
}

func which_kind[char rune | byte](r char, lookup char, has_lookup bool) TokenKind {
	switch r {
	case ' ':
		return WHITE_SPACE
	case '\\':
		return ESCAPE
	case '$':
		return EXPAND
	case '\'':
		return STRING_LITERAL
	case '"':
		return STRING_EXPAND
	case '1':
		if has_lookup && lookup == '>' {
			return REDIRECT_OUT
		}
		return NUMBER
	case '2':
		if has_lookup && lookup == '>' {
			return REDIRECT_ERR
		}
		return NUMBER

	default:
		return LITERAL
	}
}

func get_value_from(k TokenKind, i *int, input string) string {
	start := *i
	switch k {
	case WHITE_SPACE:
		for input[*i] == ' ' {
			(*i)++
		}
	case LITERAL:
		char := input[*i]
		for ('a' >= char && char <= 'z') || ('A' >= char && char <= 'Z') || (char >= '0' && char <= '9') {
			(*i)++
		}
	case NUMBER:
		char := input[*i]
		for char >= '0' && char <= '9' {
			(*i)++
		}
	case STRING_LITERAL:
		for input[*i] != '\'' {
			(*i)++
		}
	case STRING_EXPAND:
		for input[*i] != '"' {
			(*i)++
		}
	case EXPAND:
		char := input[*i]
		if char == '$' {
			(*i)++
		} else {
			for 'A' >= char && char <= 'Z' {
				(*i)++
			}
		}
	case ESCAPE:
		(*i)++
	case REDIRECT_OUT:
		char := input[*i]
		if char == '1' {
			(*i) += 2
		} else {
			(*i)++
		}
	case REDIRECT_ERR:
		(*i) += 2
	}

	return input[start:*i]
}

func is_alphabetic(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func is_numeric(c byte) bool {
	return c >= '0' && c <= '9'
}

func is_alphanum(c byte) bool {
	return is_alphabetic(c) || is_numeric(c)
}
