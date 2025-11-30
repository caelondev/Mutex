package runtime

import (
	"fmt"

	"github.com/caelondev/mutex/src/errors"
)

type Environment interface {
	Environment()
	DeclareVariable(variableName string, value RuntimeValue) RuntimeValue
	AssignVariable(variableName string, value RuntimeValue) RuntimeValue
	ResolveVariable(variableName string) Environment
	GetVariable(variableName string) RuntimeValue  // Capitalized
	LookupVariable(variableName string) RuntimeValue
}

type EnvironmentStruct struct {
	parent    Environment
	variables map[string]RuntimeValue
}

func (e *EnvironmentStruct) Environment() {}

func NewEnvironment(parentEnv Environment) *EnvironmentStruct {
	variables := map[string]RuntimeValue{}
	env := &EnvironmentStruct{
		variables: variables,
		parent:    parentEnv,
	}

	if parentEnv == nil { // This means this is the global environment
		declareGlobalVariables(env)
	}
	return env
}

func declareGlobalVariables(env Environment) {
	env.DeclareVariable("nil", NIL())
	env.DeclareVariable("true", BOOLEAN(true))
	env.DeclareVariable("false", BOOLEAN(false))
}

func (e *EnvironmentStruct) DeclareVariable(variableName string, value RuntimeValue) RuntimeValue {
	// Check if variable already exists in THIS scope (not parent scopes)
	if _, exists := e.variables[variableName]; exists {
		errors.ReportInterpreter(fmt.Sprintf("Cannot declare variable \"%s\" as it is already defined", variableName), 65)
	}

	e.variables[variableName] = value
	return value
}

func (e *EnvironmentStruct) AssignVariable(variableName string, value RuntimeValue) RuntimeValue {
	env := e.ResolveVariable(variableName)
	
	// Type assert to access variables map
	if envStruct, ok := env.(*EnvironmentStruct); ok {
		envStruct.variables[variableName] = value
		return value
	}
	
	errors.ReportInterpreter(fmt.Sprintf("Cannot re-assign variable \"%s\" - invalid environment", variableName), 65)
	return nil
}

func (e *EnvironmentStruct) ResolveVariable(variableName string) Environment {
	// Check if variable exists in current scope
	if _, exists := e.variables[variableName]; exists {
		return e
	}

	// If no parent, variable doesn't exist anywhere
	if e.parent == nil {
		errors.ReportInterpreter(fmt.Sprintf("Cannot resolve variable \"%s\" as it does not exist in the current/outer scopes", variableName), 65)
	}

	// Recursively check parent scopes
	return e.parent.ResolveVariable(variableName)
}

func (e *EnvironmentStruct) GetVariable(variableName string) RuntimeValue {
	return e.variables[variableName]
}

func (e *EnvironmentStruct) LookupVariable(variableName string) RuntimeValue {
	env := e.ResolveVariable(variableName)
	return env.GetVariable(variableName)
}
