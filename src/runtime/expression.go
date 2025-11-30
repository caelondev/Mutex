package runtime

import (
	"fmt"
	"math"

	"github.com/caelondev/mutex/src/errors"
	"github.com/caelondev/mutex/src/frontend/ast"
	"github.com/caelondev/mutex/src/frontend/lexer"
)

func evaluateNumberExpression(expr *ast.NumberExpression) RuntimeValue {
	return &NumberValue{Value: expr.Value}
}

func evaluateStringExpression(expr *ast.StringExpression) RuntimeValue {
	return &StringValue{Value: expr.Value}
}

func evaluateBinaryExpression(expr *ast.BinaryExpression, env Environment) RuntimeValue {
	// Handle short-circuit evaluation for logical operators
	if expr.Operator.TokenType == lexer.AND {
		left := EvaluateExpression(expr.Left, env)
		if !isTruthy(left) {
			return BOOLEAN(false)
		}
		right := EvaluateExpression(expr.Right, env)
		return BOOLEAN(isTruthy(right))
	}

	if expr.Operator.TokenType == lexer.OR {
		left := EvaluateExpression(expr.Left, env)
		if isTruthy(left) {
			return BOOLEAN(true)
		}
		right := EvaluateExpression(expr.Right, env)
		return BOOLEAN(isTruthy(right))
	}

	// Evaluate both operands for all other operators
	left := EvaluateExpression(expr.Left, env)
	right := EvaluateExpression(expr.Right, env)

	// Handle string operations
	leftStr, leftIsStr := left.(*StringValue)
	rightStr, rightIsStr := right.(*StringValue)

	if leftIsStr && rightIsStr {
		return evaluateStringBinaryExpression(leftStr, rightStr, expr.Operator)
	}

	// Handle numeric operations
	leftNum, leftIsNum := left.(*NumberValue)
	rightNum, rightIsNum := right.(*NumberValue)

	if leftIsNum && rightIsNum {
		return evaluateNumericBinaryExpression(leftNum, rightNum, expr.Operator)
	}

	// Type mismatch
	errors.ReportInterpreter(fmt.Sprintf("Cannot perform operation %s on incompatible types", expr.Operator.Lexeme), 65)
	return NIL()
}

func evaluateStringBinaryExpression(left *StringValue, right *StringValue, operator lexer.Token) RuntimeValue {
	lhs := left.Value
	rhs := right.Value

	switch operator.TokenType {
	case lexer.PLUS:
		return &StringValue{Value: lhs + rhs}
	case lexer.EQUAL_TO:
		return BOOLEAN(lhs == rhs)
	case lexer.NOT_EQUAL:
		return BOOLEAN(lhs != rhs)
	default:
		errors.ReportInterpreter(fmt.Sprintf("Unsupported string operator: %s", operator.Lexeme), 65)
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
			errors.ReportInterpreter("Division by zero", 65)
		}
		result = lhs / rhs
	case lexer.MODULO:
		if rhs == 0 {
			errors.ReportInterpreter("Modulo by zero", 65)
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
		errors.ReportInterpreter(fmt.Sprintf("Unsupported binary operator: %s", operator.Lexeme), 65)
	}

	return &NumberValue{Value: result}
}

func evaluateSymbolExpression(expr *ast.SymbolExpression, env Environment) RuntimeValue {
	return env.LookupVariable(expr.Value)
}

func evaluateAssignmentExpression(expr *ast.AssignmentExpression, env Environment) RuntimeValue {
	symbol := expr.Assignee.(*ast.SymbolExpression)
	value := EvaluateExpression(expr.NewValue, env)

	return env.AssignVariable(symbol.Value, value)
}

func evaluateUnaryExpression(expr *ast.UnaryExpression, env Environment) RuntimeValue {
	operand := EvaluateExpression(expr.Operand, env)

	switch expr.Operator.TokenType {
	case lexer.NOT:
		return BOOLEAN(!isTruthy(operand))

	case lexer.MINUS:
		numValue, ok := operand.(*NumberValue)
		if !ok {
			errors.ReportInterpreter("Unary minus requires numeric operand", 65)
		}
		return &NumberValue{Value: -numValue.Value}

	default:
		errors.ReportInterpreter(fmt.Sprintf("Unknown unary operator: %s", expr.Operator.Lexeme), 65)
	}

	return NIL()
}

func evaluatePostfixExpression(expr *ast.PostfixExpression, env Environment) RuntimeValue {
	symbol, ok := expr.Operand.(*ast.SymbolExpression)
	if !ok {
		errors.ReportInterpreter("Postfix operators can only be applied to variables", 65)
	}

	currentValue := env.LookupVariable(symbol.Value)
	numValue, ok := currentValue.(*NumberValue)
	if !ok {
		errors.ReportInterpreter(fmt.Sprintf("Postfix operator %s requires numeric operand", expr.Operator.Lexeme), 65)
	}

	var newValue *NumberValue
	switch expr.Operator.TokenType {
	case lexer.PLUS_PLUS:
		newValue = &NumberValue{Value: numValue.Value + 1}
	case lexer.MINUS_MINUS:
		newValue = &NumberValue{Value: numValue.Value - 1}
	default:
		errors.ReportInterpreter(fmt.Sprintf("Unknown postfix operator: %s", expr.Operator.Lexeme), 65)
	}

	env.AssignVariable(symbol.Value, newValue)
	return numValue
}

func evaluateArrayExpression(expr *ast.ArrayExpression, env Environment) RuntimeValue {
	elements := make([]RuntimeValue, len(expr.Elements))
	
	for i, elemExpr := range expr.Elements {
		elements[i] = EvaluateExpression(elemExpr, env)
	}
	
	return ARRAY(elements)
}

func evaluateIndexExpression(expr *ast.ArrayIndexExpression, env Environment) RuntimeValue {
	object := EvaluateExpression(expr.Object, env)
	index := EvaluateExpression(expr.Index, env)
	
	// Check if object is an array
	arrayValue, ok := object.(*ArrayValue)
	if !ok {
		errors.ReportInterpreter(fmt.Sprintf("Cannot index into type '%s', expected array", object.Type()), 65)
		return NIL()
	}
	
	// Check if index is a number
	indexNum, ok := index.(*NumberValue)
	if !ok {
		errors.ReportInterpreter(fmt.Sprintf("Array index must be a number, got '%s'", index.Type()), 65)
		return NIL()
	}
	
	// Convert to integer and check bounds
	idx := int(indexNum.Value)
	
	if idx < 0 || idx >= len(arrayValue.Elements) {
		errors.ReportInterpreter(fmt.Sprintf("Array index %d out of bounds (array length: %d)", idx, len(arrayValue.Elements)), 65)
		return NIL()
	}
	
	return arrayValue.Elements[idx]
}

func evaluateIndexAssignmentExpression(expr *ast.ArrayIndexAssignmentExpression, env Environment) RuntimeValue {
	object := EvaluateExpression(expr.Object, env)
	index := EvaluateExpression(expr.Index, env)
	newValue := EvaluateExpression(expr.NewValue, env)
	
	// Check if object is an array
	arrayValue, ok := object.(*ArrayValue)
	if !ok {
		errors.ReportInterpreter(fmt.Sprintf("Cannot index into type '%s', expected array", object.Type()), 65)
		return NIL()
	}
	
	// Check if index is a number
	indexNum, ok := index.(*NumberValue)
	if !ok {
		errors.ReportInterpreter(fmt.Sprintf("Array index must be a number, got '%s'", index.Type()), 65)
		return NIL()
	}
	
	// Convert to integer and check bounds
	idx := int(indexNum.Value)
	
	if idx < 0 || idx >= len(arrayValue.Elements) {
		errors.ReportInterpreter(fmt.Sprintf("Array index %d out of bounds (array length: %d)", idx, len(arrayValue.Elements)), 65)
		return NIL()
	}
	
	// Mutate the array in place
	arrayValue.Elements[idx] = newValue
	
	return NIL()
}
