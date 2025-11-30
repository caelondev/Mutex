package parser

import (
	"fmt"
	"strconv"

	"github.com/caelondev/mutex/src/errors"
	"github.com/caelondev/mutex/src/frontend/ast"
	"github.com/caelondev/mutex/src/frontend/lexer"
)

func parseExpression(p *parser, bp BindingPower) ast.Expression {
	// Parse NUD ---
	tokenType := p.currentTokenType()
	nudFunction, exists := nudLU[tokenType]

	if p.isEOF() {
		errors.ReportParser("Unexpected end of file expression (EOF)", 0)
	}

	if !exists {
		errors.ReportParser(fmt.Sprintf("Unrecognized token found: %s (NUD)", lexer.TokenTypeString(tokenType)), 0)
	}

	left := nudFunction(p)

	for !p.isEOF() && bindingPowerLU[p.currentTokenType()] > bp {
		tokenType = p.currentTokenType()
		ledFunction, exists := ledLU[tokenType]

		if !exists {
			errors.ReportParser(fmt.Sprintf("Unrecognized token found: %s (LED)", lexer.TokenTypeString(tokenType)), 0)
		}

		left = ledFunction(p, left, bindingPowerLU[p.currentTokenType()])
	}

	return left
}

func parsePrimaryExpression(p *parser) ast.Expression {
	switch p.currentTokenType() {
	case lexer.NUMBER:
		number, _ := strconv.ParseFloat(p.advance().Lexeme, 64)
		return &ast.NumberExpression{
			Value: number,
		}
	case lexer.STRING:
		return &ast.StringExpression{
			Value: p.advance().Lexeme,
		}
	case lexer.IDENTIFIER:
		return &ast.SymbolExpression{
			Value: p.advance().Lexeme,
		}
	case lexer.LEFT_PARENTHESIS:
		p.advance() // eat ( ---
		value := parseExpression(p, DEFAULT_BP)
		p.advance() // eat ) ---
		return value

	default:
		panic(fmt.Sprintf("Unrecognized primary token found: '%s'", lexer.TokenTypeString(p.currentTokenType())))
	}
}

func parseBinaryExpression(p *parser, left ast.Expression, bp BindingPower) ast.Expression {
	operatorToken := p.advance()
	right := parseExpression(p, bp)

	return &ast.BinaryExpression{
		Left: left,
		Right: right,
		Operator: *operatorToken,
	}
}

func parseVariableAssignmentExpression(p *parser, left ast.Expression, bp BindingPower) ast.Expression {
	p.advance() // Eat past = token ---
	return &ast.AssignmentExpression{
		Assignee: left,
		NewValue: parseExpression(p, bp),
	}
}
