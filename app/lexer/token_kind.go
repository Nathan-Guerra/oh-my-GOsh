package lexer

import "strconv"

type TokenKind int

func (k TokenKind) String() string {
	switch k {
	case Whitespace:
		return "WHITESPACE"
	case Literal:
		return "LITERAL"
	case Numeric:
		return "NUMBER"
	case StringLiteral:
		return "STRING_LITERAL"
	case StringExpand:
		return "STRING_EXPAND"
	case Expand:
		return "EXPAND"
	case Escape:
		return "ESCAPE"
	case RedirectOut:
		return "REDIRECT_OUT"
	case RedirectErr:
		return "REDIRECT_ERR"
	case RedirectOutAppend:
		return "REDIRECT_OUT_APPEND"
	case RedirectErrAppend:
		return "REDIRECT_ERR_APPEND"
	case ParsingStringLiteral:
		return "PARSING_STRING_LITERAL"
	case ParsingStringExpand:
		return "PARSING_STRING_EXPAND"
	default:
		return strconv.Itoa(int(k))
	}
}

const (
	EOL TokenKind = iota
	Whitespace
	Literal
	Numeric
	StringLiteral
	StringExpand
	Expand
	Escape
	RedirectOut
	RedirectErr
	RedirectOutAppend
	RedirectErrAppend

	ParsingStringLiteral
	ParsingStringExpand
)
