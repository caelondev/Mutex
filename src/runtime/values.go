package runtime

import "fmt"

type ValueTypes string

const (
	BOOLEAN_VALUE ValueTypes = "boolean"
	NIL_VALUE    ValueTypes = "nil"
	NUMBER_VALUE ValueTypes = "number"
	STRING_VALUE ValueTypes = "string"
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


