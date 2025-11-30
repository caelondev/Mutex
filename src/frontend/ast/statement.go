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

type IfStatement struct {
	Condition   Expression
	Consequent  Statement
	Alternate   Statement
}

func (node *IfStatement) Statement() {}

type WhileStatement struct {
	Condition   Expression
	Body        Statement
}

func (node *WhileStatement) Statement() {}

type ForStatement struct {
	Initializer Statement
	Condition   Expression
	Increment   Expression
	Body        Statement
}

func (node *ForStatement) Statement() {}
