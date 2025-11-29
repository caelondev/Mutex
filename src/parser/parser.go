package parser

import (
	"fmt"

	"github.com/caelondev/mutex/src/ast"
	"github.com/caelondev/mutex/src/errors"
	"github.com/caelondev/mutex/src/lexer"
)

type parser struct {
	tokens   []*lexer.Token
	position int
}

func ProduceAST(tokens []*lexer.Token) ast.BlockStatement {
	body := make([]ast.Statement, 0)
	p := instantiateParser(tokens)

	for !p.isEOF() {
		body = append(body, parseStatement(p))
	}

	return ast.BlockStatement{
		Body: body,
	}
}

func instantiateParser(tokens []*lexer.Token) *parser {
	createTokenLookups()

	return &parser{
		tokens:   tokens,
		position: 0,
	}
}

func (p *parser) previousToken() *lexer.Token {
	return p.tokens[p.position-1]
}

func (p *parser) previousTokenType() lexer.TokenType {
	return p.previousToken().TokenType
}

func (p *parser) currentToken() *lexer.Token {
	return p.tokens[p.position]
}

func (p *parser) currentTokenType() lexer.TokenType {
	return p.currentToken().TokenType
}

func (p *parser) advance() *lexer.Token {
	token := p.currentToken()
	p.position++
	return token
}

func (p *parser) isEOF() bool {
	return p.position >= len(p.tokens) || p.currentTokenType() == lexer.EOF
}

func (p *parser) expect(expectedType lexer.TokenType) *lexer.Token {
	return p.expectError(expectedType, "")
}

func (p *parser) expectError(expectedType lexer.TokenType, err string) *lexer.Token {
	token := p.currentToken()
	tokenType := token.TokenType

	if tokenType != expectedType {
		if err == "" {
			err = fmt.Sprintf("Expected %s but got %s instead", lexer.TokenTypeString(expectedType), lexer.TokenTypeString(tokenType))
		}

		errors.ReportParser(err, 65)
	}

	return p.advance()
}
