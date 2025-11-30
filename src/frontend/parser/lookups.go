package parser

import (
	"github.com/caelondev/mutex/src/frontend/ast"
	"github.com/caelondev/mutex/src/frontend/lexer"
)

type BindingPower int

const (
	DEFAULT_BP BindingPower = iota
	COMMA
	ASSIGNMENT
	LOGICAL
	RELATIONAL
	ADDITIVE
	MULTIPLICATIVE
	UNARY
	POSTFIX
	CALL
	MEMBER
	PRIMARY
)

type StatementHandler func(p *parser) ast.Statement
type NudHandler func(p *parser) ast.Expression
type LedHandler func(p *parser, left ast.Expression, bp BindingPower) ast.Expression

type StatementLookup map[lexer.TokenType]StatementHandler
type NudLookup map[lexer.TokenType]NudHandler
type LedLookup map[lexer.TokenType]LedHandler
type BPLookup map[lexer.TokenType]BindingPower

var bindingPowerLU = BPLookup{}
var nudLU = NudLookup{}
var ledLU = LedLookup{}
var statementLU = StatementLookup{}

func led(tokenType lexer.TokenType, bp BindingPower, ledFunction LedHandler) {
	bindingPowerLU[tokenType] = bp
	ledLU[tokenType] = ledFunction
}

func nud(tokenType lexer.TokenType, nudFunction NudHandler) {
	bindingPowerLU[tokenType] = PRIMARY
	nudLU[tokenType] = nudFunction
}

func statement(tokenType lexer.TokenType, statementFunction StatementHandler) {
	bindingPowerLU[tokenType] = DEFAULT_BP
	statementLU[tokenType] = statementFunction
}

func createTokenLookups() {

	// LITERALS AND SYMBOLS ---
	nud(lexer.NUMBER, parsePrimaryExpression)
	nud(lexer.IDENTIFIER, parsePrimaryExpression)
	nud(lexer.STRING, parsePrimaryExpression)
	nud(lexer.LEFT_PARENTHESIS, parsePrimaryExpression)

	// ARRAYS ---
	nud(lexer.LEFT_BRACKET, parseArrayExpression)
	led(lexer.LEFT_BRACKET, CALL, parseIndexExpression)

	// PREFIX
	nud(lexer.NOT, parseUnaryExpression)
	nud(lexer.MINUS, parseUnaryExpression)

	// ASSIGNMENT & COMPOUND ASSIGNMENT ---
	led(lexer.ASSIGNMENT, ASSIGNMENT, parseVariableAssignmentExpression)
	led(lexer.PLUS_EQUALS, ASSIGNMENT, parseVariableAssignmentExpression)
	led(lexer.MINUS_EQUALS, ASSIGNMENT, parseVariableAssignmentExpression)
	led(lexer.STAR_EQUALS, ASSIGNMENT, parseVariableAssignmentExpression)
	led(lexer.SLASH_EQUALS, ASSIGNMENT, parseVariableAssignmentExpression)
	led(lexer.MODULO_EQUALS, ASSIGNMENT, parseVariableAssignmentExpression)

	// INCREMENT/DECREMENT (Postfix) ---
	led(lexer.PLUS_PLUS, POSTFIX, parsePostfixExpression)
	led(lexer.MINUS_MINUS, POSTFIX, parsePostfixExpression)

	// RELATIONAL ---
	led(lexer.LESS, RELATIONAL, parseBinaryExpression)
	led(lexer.LESS_EQUAL, RELATIONAL, parseBinaryExpression)
	led(lexer.GREATER, RELATIONAL, parseBinaryExpression)
	led(lexer.GREATER_EQUAL, RELATIONAL, parseBinaryExpression)
	led(lexer.EQUAL_TO, RELATIONAL, parseBinaryExpression)
	led(lexer.NOT_EQUAL, RELATIONAL, parseBinaryExpression)

	// ADDITIVE & MULTIPLICATIVE ---
	led(lexer.PLUS, ADDITIVE, parseBinaryExpression)
	led(lexer.MINUS, ADDITIVE, parseBinaryExpression)
	led(lexer.SLASH, MULTIPLICATIVE, parseBinaryExpression)
	led(lexer.MODULO, MULTIPLICATIVE, parseBinaryExpression)
	led(lexer.STAR, MULTIPLICATIVE, parseBinaryExpression)

	// LOGICAL ---
	led(lexer.AND, LOGICAL, parseBinaryExpression)
	led(lexer.OR, LOGICAL, parseBinaryExpression)

	// Statements ---
	statement(lexer.VAR, parseVariableDeclaration)
	statement(lexer.IF, parseIfStatement)
	statement(lexer.WHILE, parseWhileStatement)
	statement(lexer.FOR, parseForStatement)
}
