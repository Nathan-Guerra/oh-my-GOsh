package lexer

type Token struct {
	Kind           TokenKind
	Value          string
	TokenizedValue *[]Token
}
