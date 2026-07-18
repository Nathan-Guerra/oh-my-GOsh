package lexer

type Lexer struct {
	Source   string
	Position int
}

func (l *Lexer) IsEOL() bool {
	return l.Position >= len(l.Source)
}

func (l *Lexer) IsNumeric(b byte) bool {
	return b >= '0' && b <= '9'
}

func (l *Lexer) IsUppercaseLetter(b byte) bool {
	return b >= 'A' && b <= 'Z'
}

func (l *Lexer) ByteIs(b byte) bool {
	return l.Source[l.Position] == b
}

func (l *Lexer) GetStringSlice(start, end int) string {
	return l.Source[start:end]
}

func (l *Lexer) GetCurByte() byte {
	return l.Source[l.Position]
}

func (l *Lexer) NextByte() {
	l.Position++
}

func (l *Lexer) Peek() *byte {
	if l.Position+1 >= len(l.Source) {
		return nil
	}

	p := l.Source[l.Position+1]

	return &p
}

func (l *Lexer) MatchTokenKind(r byte, lookup *byte) TokenKind {
	switch {
	case l.ByteIs(' '):
		return Whitespace
	case l.ByteIs('\\'):
		return Escape
	case l.ByteIs('$'):
		return Expand
	case l.ByteIs('\''):
		return StringLiteral
	case l.ByteIs('"'):
		return StringExpand
	case l.ByteIs('>'):
		return RedirectOut
	case l.IsNumeric(l.GetCurByte()):
		if l.Peek() != nil && *l.Peek() == '>' {
			switch {
			case l.ByteIs('1'):
				return RedirectOut
			case l.ByteIs('2'):
				return RedirectErr
			}
		}
		return Numeric
	default:
		return Literal
	}
}

func (l *Lexer) CreateToken(k TokenKind) Token {
	start := l.Position
	switch k {
	case Whitespace:
		for !l.IsEOL() && l.ByteIs(' ') {
			l.NextByte()
		}

		return Token{Whitespace, " ", nil}
	case Numeric:
		for !l.IsEOL() && l.IsNumeric(l.GetCurByte()) {
			l.NextByte()
		}
		return Token{Numeric, l.GetStringSlice(start, l.Position), nil}
	case StringLiteral:
		for !l.IsEOL() {
			l.NextByte()
			if !l.IsEOL() && l.ByteIs('\'') {
				l.NextByte()
				break
			}
		}
		// skip opening and closing (')
		return Token{StringLiteral, l.GetStringSlice(start+1, l.Position-1), nil}
	case StringExpand:
		l.NextByte()
		innerTokens := make([]Token, 0)

	innerLoop:
		for !l.IsEOL() {
			if l.IsEOL() {
				panic("==Error== Unmatched double quotes!")
			}

			inner_kind := l.MatchTokenKind(l.GetCurByte(), l.Peek())

			switch inner_kind {
			case StringExpand:
				l.NextByte()
				break innerLoop
			case Whitespace:
				innerStart := l.Position
				for !l.IsEOL() && l.ByteIs(' ') {
					l.NextByte()
				}

				innerTokens = append(innerTokens, Token{Whitespace, l.GetStringSlice(innerStart, l.Position), nil})
			case Numeric:
				innerStart := l.Position
				for !l.IsEOL() && l.IsNumeric(l.GetCurByte()) {
					l.NextByte()
				}
				innerTokens = append(innerTokens, Token{Numeric, l.GetStringSlice(innerStart, l.Position), nil})
			case Expand:
				if l.Peek() == nil || *l.Peek() == ' ' {
					l.NextByte()
					innerTokens = append(innerTokens, Token{Literal, "$", nil})
					break
				} else if l.Peek() != nil && *l.Peek() == '$' {
					l.NextByte() // $<$
					l.NextByte() // $$<
					innerTokens = append(innerTokens, Token{Expand, "$", nil})
					break
				}

				innerStart := l.Position
				l.NextByte()
				for !l.IsEOL() && l.IsUppercaseLetter(l.GetCurByte()) {
					l.NextByte()
				}

				v := l.GetStringSlice(innerStart+1, l.Position)
				if len(v) == 0 {
					innerTokens = append(innerTokens, Token{Literal, "$", nil})
				} else {
					innerTokens = append(innerTokens, Token{Expand, v, nil})
				}

			case Escape:
				if l.Peek() == nil {
					l.NextByte()
					innerTokens = append(innerTokens, Token{Literal, "\\", nil})
					break
				}

				l.NextByte()
				if l.ByteIs('\\') || l.ByteIs('"') || l.ByteIs('$') {
					innerTokens = append(innerTokens, Token{Escape, string(l.GetCurByte()), nil})
					l.NextByte()
				} else {
					// don't skip byte here because it is at an unknown token
					innerTokens = append(innerTokens, Token{Literal, "\\", nil})
				}
			default:
				innerStart := l.Position
				for !l.IsEOL() &&
					!l.ByteIs('\\') &&
					!l.ByteIs('"') &&
					!l.ByteIs('$') &&
					!l.ByteIs(' ') {
					l.NextByte()
				}

				innerTokens = append(innerTokens, Token{Literal, l.GetStringSlice(innerStart, l.Position), nil})
			}
		}

		innerString := l.GetStringSlice(start+1, l.Position-1)

		return Token{StringExpand, innerString, &innerTokens}
	case Expand:
		if l.Peek() == nil || *l.Peek() == ' ' {
			l.NextByte()
			return Token{Literal, "$", nil}
		} else if *l.Peek() == '$' {
			l.NextByte()
			l.NextByte()
			return Token{Expand, "$", nil}
		}

		l.NextByte()
		for !l.IsEOL() && l.IsUppercaseLetter(l.GetCurByte()) {
			l.NextByte()
		}

		v := l.GetStringSlice(start+1, l.Position)
		if len(v) == 0 {
			return Token{Literal, "$", nil}
		} else {
			return Token{Expand, v, nil}
		}
	case Escape:
		l.NextByte()
		l.NextByte()

		return Token{Escape, l.GetStringSlice(start+1, l.Position), nil}
	case RedirectOut:
		if l.ByteIs('1') {
			l.NextByte()
			l.NextByte()
		} else {
			l.NextByte()
		}
		return Token{RedirectOut, l.GetStringSlice(start, l.Position), nil}

	case RedirectErr:
		l.NextByte()
		l.NextByte()

		return Token{RedirectErr, l.GetStringSlice(start, l.Position), nil}

	default: // Literal
		for !l.IsEOL() &&
			!l.ByteIs(' ') &&
			!l.ByteIs('\\') &&
			!l.ByteIs('$') &&
			!l.ByteIs('>') {
			l.NextByte()
		}
		return Token{Literal, l.GetStringSlice(start, l.Position), nil}
	}
}

func Tokenize(input string) []Token {
	lexer := Lexer{input, 0}
	tokens := make([]Token, 0)
	for !lexer.IsEOL() {
		kind := lexer.MatchTokenKind(lexer.GetCurByte(), lexer.Peek())
		tokens = append(tokens, lexer.CreateToken(kind))
		// value := get_value_from(kind, &i, input)
		// if kind == StringExpand {
		// 	tokenized_value := tokenize_for_double_quotes(value)
		// 	tokens = append(tokens, Token{kind, value, tokenized_value})
		// } else {
		// 	tokens = append(tokens, Token{kind, value, nil})
		// }
	}

	return tokens
}

// func tokenize_for_double_quotes(input string) *[]Token {
// 	i := 0
// 	tokens := make([]Token, 0)
// 	for i < len(input) {
// 		var lookup byte
// 		var has_lookup bool
// 		if i+1 < len(input) {
// 			lookup = input[i+1]
// 			has_lookup = true
// 		}
// 		kind := which_kind_for_double_quotes(input[i], lookup, has_lookup)
// 		value := get_value_from_for_double_quotes(kind, &i, input)
// 		tokens = append(tokens, Token{kind, value, nil})
// 	}

// 	return &tokens
// }

// func which_kind_for_double_quotes(r byte, lookup byte, has_lookup bool) TokenKind {
// 	switch {
// 	case r == ' ':
// 		return Whitespace
// 	case r == '\\':
// 		escapable_runes := []byte{'"', '$', '\\'}
// 		if has_lookup && slices.Contains(escapable_runes[:], lookup) {
// 			return Escape
// 		}
// 		return Literal
// 	case r == '$':
// 		if !has_lookup || lookup == ' ' {
// 			return Literal
// 		} else {
// 			return Expand
// 		}
// 	case is_numeric(r):
// 		return Numeric
// 	default:
// 		return Literal
// 	}
// }

// func get_value_from_for_double_quotes(k TokenKind, i *int, input string) string {
// 	start := *i
// 	switch k {
// 	case Whitespace:
// 		(*i)++
// 		return input[start:min(*i, len(input))]
// 	case Escape:
// 		start += 1
// 		(*i) += 2
// 		return input[start:min(*i, len(input))]
// 	case Expand:
// 		start += 1
// 		(*i) += 2

// 		if (*i) < len(input) {
// 			char := input[*i]
// 			if char == '$' {
// 				(*i)++
// 			} else {
// 				for (*i) < len(input) && input[*i] >= 'A' && input[*i] <= 'Z' {
// 					(*i)++
// 				}
// 			}
// 		}
// 		return input[start:min(*i, len(input))]
// 	case Numeric:
// 		for (*i) < len(input) && is_numeric(input[*i]) {
// 			(*i)++
// 		}
// 		return input[start:min(*i, len(input))]
// 	default: // LITERAL
// 		(*i)++
// 		for (*i) < len(input) &&
// 			input[*i] != ' ' &&
// 			input[*i] != '\\' {
// 			if input[*i] == '$' &&
// 				*i+1 < len(input) &&
// 				input[*i+1] != ' ' {
// 				break
// 			}
// 			(*i)++
// 		}
// 		return input[start:min(*i, len(input))]
// 	}
// }

// func get_value_from(k TokenKind, i *int, input string) string {
// 	start := *i
// 	switch k {
// 	case Whitespace:
// 		for (*i) < len(input) && input[*i] == ' ' {
// 			(*i)++
// 		}
// 		return input[start:min(*i, len(input))]
// 	case Literal:
// 		for (*i) < len(input) &&
// 			input[*i] != ' ' &&
// 			input[*i] != '\\' {
// 			if input[*i] == '$' &&
// 				*i+1 < len(input) &&
// 				input[*i+1] != ' ' {
// 				break
// 			}
// 			(*i)++
// 		}
// 		return input[start:min(*i, len(input))]
// 	case Numeric:
// 		for (*i) < len(input) && is_numeric(input[*i]) {
// 			(*i)++
// 		}
// 		return input[start:min(*i, len(input))]
// 	case StringLiteral:
// 		start++
// 		(*i)++
// 		for (*i) < len(input) {
// 			(*i)++
// 			if (*i) < len(input) && input[*i] == '\'' {
// 				(*i)++
// 				break
// 			}
// 		}
// 		return input[start:min((*i)-1, len(input))]
// 	case StringExpand:
// 		start++
// 		(*i)++
// 		for (*i) < len(input) {
// 			var inner_kind TokenKind
// 			if *i+1 < len(input) {
// 				inner_kind = which_kind(input[*i], input[(*i)+1], true)
// 			} else {
// 				inner_kind = which_kind(input[*i], 'x', false)
// 			}

// 			switch inner_kind {
// 			case Whitespace:
// 				for (*i) < len(input) && input[*i] == ' ' {
// 					(*i)++
// 				}
// 			case Numeric:
// 				for (*i) < len(input) && is_numeric(input[*i]) {
// 					(*i)++
// 				}
// 			case Expand:
// 				(*i)++
// 				if (*i) < len(input) {
// 					char := input[*i]
// 					if char == '$' {
// 						(*i)++
// 					} else {
// 						for (*i) < len(input) && input[*i] >= 'A' && input[*i] <= 'Z' {
// 							(*i)++
// 						}
// 					}
// 				}
// 			case Escape:
// 				(*i) += 2
// 			default:
// 				if input[*i] == '"' {
// 					(*i)++
// 					break
// 				}
// 				(*i)++
// 			}
// 		}
// 		return input[start:min((*i)-1, len(input))]
// 	case Expand:
// 		start += 1
// 		(*i) += 2

// 		if (*i) < len(input) {
// 			char := input[*i]
// 			if char == '$' {
// 				(*i)++
// 			} else {
// 				for (*i) < len(input) && input[*i] >= 'A' && input[*i] <= 'Z' {
// 					(*i)++
// 				}
// 			}
// 		}
// 		return input[start:min(*i, len(input))]
// 	case Escape:
// 		start += 1
// 		(*i) += 2
// 		return input[start:min(*i, len(input))]
// 	case RedirectOut:
// 		char := input[*i]
// 		if char == '1' {
// 			(*i) += 2
// 		} else {
// 			(*i)++
// 		}
// 		return input[start:min(*i, len(input))]
// 	case RedirectErr:
// 		(*i) += 2
// 	}

// 	return input[start:min(*i, len(input))]
// }

// func is_numeric(c byte) bool {
// 	return c >= '0' && c <= '9'
// }
