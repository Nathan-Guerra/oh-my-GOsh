package lexer

import "slices"

type TokenKind int

func (k TokenKind) String() string {
	return get_token_kind(k)
}

type Token struct {
	Kind           TokenKind
	Value          string
	TokenizedValue *[]Token
}

const (
	WHITE_SPACE TokenKind = iota
	LITERAL
	NUMERIC
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
		if kind == STRING_EXPAND {
			tokenized_value := tokenize_for_double_quotes(value)
			tokens = append(tokens, Token{kind, value, tokenized_value})
		} else {
			tokens = append(tokens, Token{kind, value, nil})
		}
	}

	return tokens
}

func tokenize_for_double_quotes(input string) *[]Token {
	i := 0
	tokens := make([]Token, 0)
	for i < len(input) {
		var lookup byte
		var has_lookup bool
		if i+1 < len(input) {
			lookup = input[i+1]
			has_lookup = true
		}
		kind := which_kind_for_double_quotes(input[i], lookup, has_lookup)
		value := get_value_from(kind, &i, input)
		tokens = append(tokens, Token{kind, value, nil})
	}

	return &tokens
}

func which_kind_for_double_quotes(r byte, lookup byte, has_lookup bool) TokenKind {
	switch {
	case r == ' ':
		return WHITE_SPACE
	case r == '\\':
		escapable_runes := []byte{'"', '$', '\\'}
		if has_lookup && slices.Contains(escapable_runes[:], lookup) {
			return ESCAPE
		}
		return LITERAL
	case r == '$':
		if !has_lookup || lookup == ' ' {
			return LITERAL
		} else {
			return EXPAND
		}
	case is_numeric(r):
		return NUMERIC
	default:
		return LITERAL
	}
}

func which_kind(r byte, lookup byte, has_lookup bool) TokenKind {
	switch {
	case r == ' ':
		return WHITE_SPACE
	case r == '\\':
		return ESCAPE
	case r == '$':
		if !has_lookup {
			return LITERAL
		} else if lookup == ' ' {
			return LITERAL
		} else {
			return EXPAND
		}
	case r == '\'':
		return STRING_LITERAL
	case r == '"':
		return STRING_EXPAND
	case r == '>':
		return REDIRECT_OUT
	case is_numeric(r):
		if has_lookup && lookup == '>' {
			switch r {
			case '1':
				return REDIRECT_OUT
			case '2':
				return REDIRECT_ERR
			}
		}
		return NUMERIC
	default:
		return LITERAL
	}
}

func get_value_from(k TokenKind, i *int, input string) string {
	start := *i
	switch k {
	case WHITE_SPACE:
		for (*i) < len(input) && input[*i] == ' ' {
			(*i)++
		}
		return input[start:min(*i, len(input))]
	case LITERAL:
		for (*i) < len(input) &&
			input[*i] != ' ' &&
			input[*i] != '\\' {
			if input[*i] == '$' &&
				*i+1 < len(input) &&
				input[*i+1] != ' ' {
				break
			}
			(*i)++
		}
		return input[start:min(*i, len(input))]
	case NUMERIC:
		for (*i) < len(input) && is_numeric(input[*i]) {
			(*i)++
		}
		return input[start:min(*i, len(input))]
	case STRING_LITERAL:
		start++
		(*i)++
		for (*i) < len(input) {
			(*i)++
			if (*i) < len(input) && input[*i] == '\'' {
				(*i)++
				break
			}
		}
		return input[start:min((*i)-1, len(input))]
	case STRING_EXPAND:
		start++
		(*i)++
		for (*i) < len(input) {
			(*i)++
			if (*i) < len(input) && input[*i] == '"' {
				(*i)++
				break
			}
		}
		return input[start:min((*i)-1, len(input))]
	case EXPAND:
		start += 1
		(*i) += 2

		if (*i) < len(input) {
			char := input[*i]
			if char == '$' {
				(*i)++
			} else {
				for (*i) < len(input) && input[*i] >= 'A' && input[*i] <= 'Z' {
					(*i)++
				}
			}
		}
		return input[start:min(*i, len(input))]
	case ESCAPE:
		start += 1
		(*i) += 2
		return input[start:min(*i, len(input))]
	case REDIRECT_OUT:
		char := input[*i]
		if char == '1' {
			(*i) += 2
		} else {
			(*i)++
		}
		return input[start:min(*i, len(input))]
	case REDIRECT_ERR:
		(*i) += 2
	}

	return input[start:min(*i, len(input))]
}

func is_numeric(c byte) bool {
	return c >= '0' && c <= '9'
}

func get_token_kind(k TokenKind) string {
	switch k {
	case WHITE_SPACE:
		return "WHITE_SPACE"
	case ESCAPE:
		return "ESCAPE"
	case EXPAND:
		return "EXPAND"
	case STRING_LITERAL:
		return "STRING_LITERAL"
	case STRING_EXPAND:
		return "STRING_EXPAND"
	case REDIRECT_OUT:
		return "REDIRECT_OUT"
	case REDIRECT_ERR:
		return "REDIRECT_ERR"
	case NUMERIC:
		return "NUMBER"
	case LITERAL:
		return "LITERAL"
	default:
		panic("Unknown token kind")
	}
}
