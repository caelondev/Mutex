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
		errors.ReportParser(fmt.Sprintf("Unrecognized token found in the begining of an expression: %s", lexer.TokenTypeString(tokenType)), 0)
	}

	left := nudFunction(p)

	for !p.isEOF() && bindingPowerLU[p.currentTokenType()] > bp {
		tokenType = p.currentTokenType()
		ledFunction, exists := ledLU[tokenType]

		if !exists {
			errors.ReportParser(fmt.Sprintf("Unrecognized token found in the middle of an expression: %s (LED)", lexer.TokenTypeString(tokenType)), 0)
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
		p.expect(lexer.RIGHT_PARENTHESIS)
		return value

	default:
		panic(fmt.Sprintf("Unrecognized primary token found: '%s'", lexer.TokenTypeString(p.currentTokenType())))
	}
}

func parseBinaryExpression(p *parser, left ast.Expression, bp BindingPower) ast.Expression {
	operatorToken := p.advance()
	right := parseExpression(p, bp)

	return &ast.BinaryExpression{
		Left:     left,
		Right:    right,
		Operator: *operatorToken,
	}
}

func parseVariableAssignmentExpression(p *parser, left ast.Expression, bp BindingPower) ast.Expression {
	operatorToken := p.advance() // Get the operator (=, +=, -=, etc.)
	value := parseExpression(p, ASSIGNMENT)

	// Check if left side is an index expression (arr[0] = ...)
	if indexExpr, ok := left.(*ast.ArrayIndexExpression); ok {
		// Handle compound assignment for arrays
		if operatorToken.TokenType != lexer.ASSIGNMENT {
			var binaryOp lexer.TokenType
			switch operatorToken.TokenType {
			case lexer.PLUS_EQUALS:
				binaryOp = lexer.PLUS
			case lexer.MINUS_EQUALS:
				binaryOp = lexer.MINUS
			case lexer.STAR_EQUALS:
				binaryOp = lexer.STAR
			case lexer.SLASH_EQUALS:
				binaryOp = lexer.SLASH
			case lexer.MODULO_EQUALS:
				binaryOp = lexer.MODULO
			default:
				errors.ReportParser(fmt.Sprintf("Unrecognized compound assignment operator: %s", lexer.TokenTypeString(operatorToken.TokenType)), 65)
			}

			// arr[i] += 5  becomes  arr[i] = arr[i] + 5
			value = &ast.BinaryExpression{
				Left:  indexExpr, // arr[i]
				Right: value,
				Operator: lexer.Token{
					TokenType: binaryOp,
					Lexeme:    lexer.TokenTypeString(binaryOp),
					Literal:   nil,
					Line:      operatorToken.Line,
				},
			}
		}

		return &ast.ArrayIndexAssignmentExpression{
			Object:   indexExpr.Object,
			Index:    indexExpr.Index,
			NewValue: value,
		}
	}

	// Handle regular variable assignment
	if operatorToken.TokenType != lexer.ASSIGNMENT {
		var binaryOp lexer.TokenType
		switch operatorToken.TokenType {
		case lexer.PLUS_EQUALS:
			binaryOp = lexer.PLUS
		case lexer.MINUS_EQUALS:
			binaryOp = lexer.MINUS
		case lexer.STAR_EQUALS:
			binaryOp = lexer.STAR
		case lexer.SLASH_EQUALS:
			binaryOp = lexer.SLASH
		case lexer.MODULO_EQUALS:
			binaryOp = lexer.MODULO
		default:
			errors.ReportParser(fmt.Sprintf("Unrecognized compound assignment operator: %s", lexer.TokenTypeString(operatorToken.TokenType)), 65)
		}

		value = &ast.BinaryExpression{
			Left:  left,
			Right: value,
			Operator: lexer.Token{
				TokenType: binaryOp,
				Lexeme:    lexer.TokenTypeString(binaryOp),
				Literal:   nil,
				Line:      operatorToken.Line,
			},
		}
	}

	return &ast.AssignmentExpression{
		Assignee: left,
		NewValue: value,
	}
}

func parseUnaryExpression(p *parser) ast.Expression {
	operatorToken := p.advance()
	operand := parseExpression(p, UNARY)

	return &ast.UnaryExpression{
		Operator: *operatorToken,
		Operand:  operand,
	}
}

func parsePostfixExpression(p *parser, left ast.Expression, bp BindingPower) ast.Expression {
	operatorToken := p.advance()

	return &ast.PostfixExpression{
		Operator: *operatorToken,
		Operand:  left,
	}
}

func parseArrayExpression(p *parser) ast.Expression {
	// Parse: [1, 2, 3]
	p.advance() // eat '['
	
	var elements []ast.Expression
	
	// Handle empty array: []
	if p.currentTokenType() == lexer.RIGHT_BRACKET {
		p.advance() // eat ']'
		return &ast.ArrayExpression{
			Elements: elements,
		}
	}
	
	// Parse first element
	elements = append(elements, parseExpression(p, DEFAULT_BP))
	
	// Parse remaining elements
	for p.currentTokenType() == lexer.COMMA {
		p.advance() // eat ','
		elements = append(elements, parseExpression(p, DEFAULT_BP))
	}
	
	p.expect(lexer.RIGHT_BRACKET) // eat ']'
	
	return &ast.ArrayExpression{
		Elements: elements,
	}
}
func parseIndexExpression(p *parser, left ast.Expression, bp BindingPower) ast.Expression {
	// Parse: arr[0] or arr[i+1]
	p.advance() // eat '['
	
	index := parseExpression(p, DEFAULT_BP)
	
	p.expect(lexer.RIGHT_BRACKET) // eat ']'
	
	return &ast.ArrayIndexExpression{
		Object: left,
		Index:  index,
	}
}
