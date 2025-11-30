package parser

import (
	"fmt"

	"github.com/caelondev/mutex/src/errors"
	"github.com/caelondev/mutex/src/frontend/ast"
	"github.com/caelondev/mutex/src/frontend/lexer"
)

func parseStatement(p *parser) ast.Statement {
	if p.isEOF() {
		return nil
	}
	statementFunction, exists := statementLU[p.currentTokenType()]

	if exists {
		return statementFunction(p)
	}

	expression := parseExpression(p, DEFAULT_BP)

	p.expect(lexer.SEMICOLON)

	return  &ast.ExpressionStatement{
		Expression:  expression,
	}
}

func parseVariableDeclaration(p *parser) ast.Statement {
	//  SYNTAX ---
	//
	//  var (mut | imm) variableName = value
	//
	
	var identifier string
	var value ast.Expression
	var isMutable bool
	
	p.advance() // eat var keyword ---

	// check (imm/mut)
	if(p.currentTokenType() != lexer.IMMUTABLE &&
		 p.currentTokenType() != lexer.MUTABLE) {
		errors.ReportParser(fmt.Sprintf("Expected token (MUTABLE/IMMUTABLE) but gor %s instead", lexer.TokenTypeString(p.currentTokenType())), 65)
	}

	isMutable = p.currentTokenType() == lexer.MUTABLE 
	p.advance() // eat imm/mut keyword ---

	identifier = p.expect(lexer.IDENTIFIER).Lexeme

	if p.currentTokenType() != lexer.SEMICOLON {
		p.expect(lexer.ASSIGNMENT, lexer.SEMICOLON) 
		// NOTE: EXPECTING A SEMICOLON HERE IS USELESS... I JUST USED IT FOR ERROR MESSAGE ---

		value = parseExpression(p, DEFAULT_BP)
	}

	p.expect(lexer.SEMICOLON)

	return &ast.VariableDeclarationStatement{
		Identifier: identifier,
		IsMutable: isMutable,
		Value: value,
	}
}
