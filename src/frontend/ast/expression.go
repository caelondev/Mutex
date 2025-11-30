package ast

import "github.com/caelondev/mutex/src/frontend/lexer"


type NumberExpression struct {
	Value float64
}

func (node *NumberExpression) Expression() {}

type StringExpression struct {
	Value string
}

func (node *StringExpression) Expression() {}

type SymbolExpression struct {
	Value string
}

func (node *SymbolExpression) Expression() {}



type BinaryExpression struct {
	Left Expression
	Right Expression
	Operator lexer.Token
}

func (node *BinaryExpression) Expression() {}

type AssignmentExpression struct {
	Assignee Expression
	NewValue Expression
}

func (node *AssignmentExpression) Expression() {}

type UnaryExpression struct {
	Operator lexer.Token
	Operand  Expression
}

func (node *UnaryExpression) Expression() {}

type PostfixExpression struct {
	Operator lexer.Token
	Operand  Expression
}

func (node *PostfixExpression) Expression() {}

type ArrayExpression struct {
	Elements []Expression
}

func (node *ArrayExpression) Expression() {}

type ArrayIndexExpression struct {
	Object Expression
	Index Expression
}

func (node *ArrayIndexExpression) Expression() {}

type ArrayIndexAssignmentExpression struct {
	Object Expression
	Index Expression
	NewValue Expression
}

func (node *ArrayIndexAssignmentExpression) Expression() {}
