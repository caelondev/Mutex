package runtime

import (
	"fmt"
	"math"
	"os"

	"github.com/caelondev/mutex/src/frontend/ast"
	"github.com/caelondev/mutex/src/frontend/lexer"
	"github.com/sanity-io/litter"
)

func EvaluateStatement(node ast.Statement, env Environment) RuntimeValue {
	switch n := node.(type) {
	case *ast.BlockStatement:
		return evaluateBlockStatement(n, env)
	case *ast.ExpressionStatement:
		return EvaluateExpression(n.Expression, env)
	case *ast.VariableDeclarationStatement:
		return evaluateVariableDeclarationStatement(n, env)
	case *ast.IfStatement: // Add this case
		return evaluateIfStatement(n, env)
	case *ast.WhileStatement:
		return evaluateWhileStatement(n, env)
	case *ast.ForStatement:
		return evaluateForStatement(n, env)

	default:
		litter.Dump(fmt.Sprintf("Unsupported statement node type: %T", node))
		litter.Dump(node)
		os.Exit(65)
	}

	return nil
}

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

	default:
		litter.Dump(fmt.Sprintf("Unsupported expression node type: %T", node))
		litter.Dump(node)
	}

	return nil
}

func evaluateBlockStatement(block *ast.BlockStatement, env Environment) RuntimeValue {
	// Create a new child environment for block scope
	blockEnv := NewEnvironment(env)

	var lastEvaluated RuntimeValue = &NilValue{}

	for _, statement := range block.Body {
		lastEvaluated = EvaluateStatement(statement, blockEnv) // Use blockEnv instead of env
	}

	return lastEvaluated
}

func evaluateNumberExpression(expr *ast.NumberExpression) RuntimeValue {
	return &NumberValue{Value: expr.Value}
}

func evaluateStringExpression(expr *ast.StringExpression) RuntimeValue {
	return &StringValue{Value: expr.Value}
}

func evaluateBinaryExpression(expr *ast.BinaryExpression, env Environment) RuntimeValue {
	left := EvaluateExpression(expr.Left, env)
	right := EvaluateExpression(expr.Right, env)

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

	case lexer.LESS:
		return BOOLEAN(lhs < rhs)
	case lexer.LESS_EQUAL:
		return BOOLEAN(lhs <= rhs)
	case lexer.GREATER:
		return BOOLEAN(lhs > rhs)
	case lexer.GREATER_EQUAL:
		return BOOLEAN(lhs >= rhs)
	case lexer.EQUAL_TO:
		return BOOLEAN(lhs == rhs)
	case lexer.NOT_EQUAL:
		return BOOLEAN(lhs != rhs)

	default:
		panic(fmt.Sprintf("Unsupported binary operator: %s", operator.Lexeme))
	}

	return &NumberValue{Value: result}
}

func evaluateVariableDeclarationStatement(stmt *ast.VariableDeclarationStatement, env Environment) RuntimeValue {
	var value RuntimeValue

	if stmt.Value != nil {
		value = EvaluateExpression(stmt.Value, env)
	} else {
		value = &NilValue{}
	}

	return env.DeclareVariable(stmt.Identifier, value, !stmt.IsMutable)
}

func evaluateSymbolExpression(expr *ast.SymbolExpression, env Environment) RuntimeValue {
	return env.LookupVariable(expr.Value)
}

func evaluateAssignmentExpression(expr *ast.AssignmentExpression, env Environment) RuntimeValue {
	symbol := expr.Assignee.(*ast.SymbolExpression)
	value := EvaluateExpression(expr.NewValue, env)

	return env.AssignVariable(symbol.Value, value)
}

func evaluateIfStatement(stmt *ast.IfStatement, env Environment) RuntimeValue {
	condition := EvaluateExpression(stmt.Condition, env)

	// Check if condition is truthy
	if isTruthy(condition) {
		return EvaluateStatement(stmt.Consequent, env)
	} else if stmt.Alternate != nil {
		return EvaluateStatement(stmt.Alternate, env)
	}

	return NIL()
}

func isTruthy(value RuntimeValue) bool {
	switch v := value.(type) {
	case *NilValue:
		return false
	case *BooleanValue:
		return v.Value
	case *NumberValue:
		return v.Value != 0
	case *StringValue:
		return v.Value != ""
	default:
		return true
	}
}

func evaluateWhileStatement(stmt *ast.WhileStatement, env Environment) RuntimeValue {
	for {
		condition := EvaluateExpression(stmt.Condition, env)

		if !isTruthy(condition) {
			break
		}

		EvaluateStatement(stmt.Body, env)
	}

	return NIL()
}

func evaluateForStatement(stmt *ast.ForStatement, env Environment) RuntimeValue {
	loopEnv := NewEnvironment(env)

	EvaluateStatement(stmt.Initializer, loopEnv)

	for {
		condition := EvaluateExpression(stmt.Condition, loopEnv)

		if !isTruthy(condition) {
			break
		}

		EvaluateStatement(stmt.Body, loopEnv)
		EvaluateExpression(stmt.Increment, loopEnv) // Increment after body
	}

	return NIL()
}
