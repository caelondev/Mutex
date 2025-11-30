package ast

type BlockStatement struct {
	Body []Statement
}

func (node *BlockStatement) Statement() {}

type ExpressionStatement struct {
	Expression Expression
}

func (node *ExpressionStatement) Statement() {}

type VariableDeclarationStatement struct {
	IsMutable bool
	Identifier string
	Value Expression
}

func (node *VariableDeclarationStatement) Statement() {}
