package runtime

import (
	"fmt"
	"strings"
)

type ValueTypes string

const (
	BOOLEAN_VALUE ValueTypes = "boolean"
	NIL_VALUE    ValueTypes = "nil"
	NUMBER_VALUE ValueTypes = "number"
	STRING_VALUE ValueTypes = "string"
	ARRAY_VALUE ValueTypes = "array"
)

type RuntimeValue interface {
	Type() ValueTypes
}

type NilValue struct{}

func (n *NilValue) Type() ValueTypes {
	return NIL_VALUE
}

func (n *NilValue) String() string {
	return "nil"
}

type StringValue struct {
	Value string
}

func (n *StringValue) Type() ValueTypes {
	return STRING_VALUE
}

func (n *StringValue) String() string {
	return fmt.Sprintf("%v", n.Value)
}

type NumberValue struct {
	Value float64
}

func (n *NumberValue) Type() ValueTypes {
	return NUMBER_VALUE
}

func (n *NumberValue) String() string {
	return fmt.Sprintf("%v", n.Value)
}

type BooleanValue struct {
	Value bool
}

func (n *BooleanValue) Type() ValueTypes {
	return NUMBER_VALUE
}

func (n *BooleanValue) String() string {
	return fmt.Sprintf("%v", n.Value)
}


func NIL() *NilValue {
	return &NilValue{}
}

func BOOLEAN(value bool) *BooleanValue {
	return &BooleanValue{ Value: value }
}

type ArrayValue struct {
	Elements []RuntimeValue
}

func (a *ArrayValue) Type() ValueTypes {
	return ARRAY_VALUE
}

func (a *ArrayValue) String() string {
	if len(a.Elements) == 0 {
		return "[]"
	}

	var elements []string
	for _, elem := range a.Elements {
		elements = append(elements, fmt.Sprintf("%v", elem))
	}
	return "[" + strings.Join(elements, ", ") + "]"
}

func ARRAY(elements []RuntimeValue) *ArrayValue {
	return &ArrayValue{Elements: elements}
}
