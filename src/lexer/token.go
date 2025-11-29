package lexer

import "fmt"

type Token struct {
	TokenType TokenType
	Lexeme    string
	Literal   any
	Line      int
}

func NewToken(tokenType TokenType, lexeme string, literal any, line int) Token {
	return Token{
		tokenType, lexeme, literal, line,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("Token Type: %s\nValue: %s\nLiteral: %v\n", TokenTypeString(t.TokenType), t.Lexeme, t.Literal)
}
