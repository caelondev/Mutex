package runtime

import "github.com/caelondev/mutex/src/frontend/ast"

func EvaluateExpression(node ast.Expression, env Environment) RuntimeValue {
	switch n := node.(type) {
	case *ast.NumberExpression:
		return evaluateNumberExpression(n)
	case *ast.StringExpression:
		return evaluateStringExpression(n)
	case *ast.BinaryExpression:
		return evaluateBinaryExpression(n, env)
	case *ast.SymbolExpression:
		return evaluateSymbolExpression(n, env)
	case *ast.AssignmentExpression:
		return evaluateAssignmentExpression(n, env)
	case *ast.UnaryExpression:
		return evaluateUnaryExpression(n, env)
	case *ast.PostfixExpression:
		return evaluatePostfixExpression(n, env)
	case *ast.ArrayExpression:
		return evaluateArrayExpression(n, env)
	case *ast.ArrayIndexExpression:
		return evaluateIndexExpression(n, env)
	case *ast.ArrayIndexAssignmentExpression:
		return evaluateIndexAssignmentExpression(n, env)

	default:
		panic("Unsupported expression node type")
	}
}

func EvaluateStatement(node ast.Statement, env Environment) RuntimeValue {
	switch n := node.(type) {
	case *ast.BlockStatement:
		return evaluateBlockStatement(n, env)
	case *ast.ExpressionStatement:
		return EvaluateExpression(n.Expression, env)
	case *ast.VariableDeclarationStatement:
		return evaluateVariableDeclarationStatement(n, env)
	case *ast.IfStatement:
		return evaluateIfStatement(n, env)
	case *ast.WhileStatement:
		return evaluateWhileStatement(n, env)
	case *ast.ForStatement:
		return evaluateForStatement(n, env)
	default:
		panic("Unsupported statement node type")
	}
}
