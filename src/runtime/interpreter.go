package runtime

import (
	"fmt"
	"math"
	"os"

	"github.com/caelondev/mutex/src/frontend/ast"
	"github.com/caelondev/mutex/src/frontend/lexer"
	"github.com/sanity-io/litter"
)

// EvaluateStatement evaluates a statement node
func EvaluateStatement(node ast.Statement, env Environment) RuntimeValue {
	switch n := node.(type) {
	case *ast.BlockStatement:
		return evaluateBlockStatement(n, env)
	case *ast.ExpressionStatement:
		return EvaluateExpression(n.Expression, env)
	case *ast.VariableDeclarationStatement:
		return evaluateVariableDeclarationStatement(n, env)
	default:
		litter.Dump(fmt.Sprintf("Unsupported statement node type: %T", node))
		os.Exit(65)
	}

	return nil
}

// EvaluateExpression evaluates an expression node
func EvaluateExpression(node ast.Expression, env Environment) RuntimeValue {
	switch n := node.(type) {
	case *ast.NumberExpression:
		return evaluateNumberExpression(n)
	case *ast.BinaryExpression:
		return evaluateBinaryExpression(n, env)
	case *ast.SymbolExpression:
		return evaluateSymbolExpression(n, env)
	default:
		panic(fmt.Sprintf("Unsupported expression node type: %T\n\nAST: %s", node, node))
	}
}

func evaluateBlockStatement(block *ast.BlockStatement, env Environment) RuntimeValue {
	var lastEvaluated RuntimeValue = &NilValue{}
	
	for _, statement := range block.Body {
		lastEvaluated = EvaluateStatement(statement, env)
	}
	
	return lastEvaluated
}

func evaluateNumberExpression(expr *ast.NumberExpression) RuntimeValue {
	return &NumberValue{Value: expr.Value}
}

func evaluateBinaryExpression(expr *ast.BinaryExpression, env Environment) RuntimeValue {
	left := EvaluateExpression(expr.Left, env)
	right := EvaluateExpression(expr.Right, env)
	
	// Ensure both sides are numbers
	leftNum, leftOk := left.(*NumberValue)
	rhsNum, rightOk := right.(*NumberValue)
	
	if leftOk && rightOk {
		return evaluateNumericBinaryExpression(leftNum, rhsNum, expr.Operator)
	}
	
	return NIL()
}

func evaluateNumericBinaryExpression(left *NumberValue, right *NumberValue, operator lexer.Token) RuntimeValue {
	result := 0.0

	lhs := left.Value
	rhs := right.Value
	
	switch operator.TokenType {
	case lexer.PLUS:
		result = lhs + rhs
	case lexer.MINUS:
		result = lhs - rhs
	case lexer.STAR:
		result = lhs * rhs
	case lexer.SLASH:
		if rhs == 0 {
			panic("Division by zero")
		}
		result = lhs / rhs
	case lexer.MODULO:
		if rhs == 0 {
			panic("Modulo by zero")
		}
		result = math.Mod(lhs, rhs)
	default:
		panic(fmt.Sprintf("Unsupported binary operator: %s", operator.Lexeme))
	}
	
	return &NumberValue{Value: result}
}

func evaluateVariableDeclarationStatement(stmt *ast.VariableDeclarationStatement, env Environment) RuntimeValue {
	var value RuntimeValue
	
	// If there's an initial value, evaluate it
	if stmt.Value != nil {
		value = EvaluateExpression(stmt.Value, env)
	} else {
		// If no initial value, default to nil
		value = &NilValue{}
	}
	
	// Declare the variable in the environment
	return env.DeclareVariable(stmt.Identifier, value)
}

func evaluateSymbolExpression(expr *ast.SymbolExpression, env Environment) RuntimeValue {
	return env.LookupVariable(expr.Value)
}

