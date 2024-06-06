package interpreter

import (
	"fmt"
	"reflect"
	"tkom/shared"
)

// implementation of stack for scopes
type Stack struct {
	elem []*Scope
}

func (s *Stack) Push(value *Scope) {
	s.elem = append(s.elem, value)
}

func (s *Stack) Pop() (*Scope, error) {
	if len(s.elem) == 0 {
		return nil, fmt.Errorf("stack is empty")
	}
	top := s.elem[len(s.elem)-1]
	s.elem = s.elem[:len(s.elem)-1]
	return top, nil
}

func (s *Stack) Top() (*Scope, error) {
	if len(s.elem) == 0 {
		return nil, fmt.Errorf("stack is empty")
	}
	return s.elem[len(s.elem)-1], nil
}

func (s *Stack) IsEmpty() bool {
	return len(s.elem) == 0
}

func (s *Stack) Size() int {
	return len(s.elem)
}

type ScopeVariable struct {
	Value any
	Type  shared.TypeAnnotation
}

type Scope struct {
	Parent *Scope
	// variables  map[string]*ScopeVariable
	variables  map[string]any
	ReturnType *shared.TypeAnnotation
}

func NewScope(parent *Scope, returnType *shared.TypeAnnotation) *Scope {
	return &Scope{
		Parent: parent,
		// variables:  map[string]*ScopeVariable{},
		variables:  map[string]any{},
		ReturnType: returnType,
	}
}

func (s *Scope) InScope(name string) map[string]any {
	if _, ok := s.variables[name]; ok {
		return s.variables
	}

	if s.Parent != nil {
		return s.Parent.InScope(name)
	}

	return s.variables
}

func (s *Scope) AddVariable(name string, value any, variableType shared.TypeAnnotation, position shared.Position) error {
	if _, ok := s.variables[name]; ok {
		return NewSemanticError(fmt.Sprintf(REDECLARED_VARIABLE, name), position)
	}

	s.variables[name] = value
	return nil
}

// Sets a value to variable in the scope
//
// If no such value is foud returns a UNDEFINED_VARIABLE error
//
// If variable type doesn't match with value type, returns a TYPE_MISMATCH error
func (s *Scope) SetValue(name string, value any) error {
	variables := s.InScope(name)
	v := variables[name]
	if v == nil {
		return fmt.Errorf(UNDEFINED_VARIABLE, name)
	}

	err := s.CheckVariableType(v, value)
	if err != nil {
		return err
	}
	variables[name] = value

	return nil
}

func (s *Scope) CheckVariableType(variable, value any) error {
	typeCheckers := map[reflect.Type]func(any) bool{
		reflect.TypeOf(1):     func(v any) bool { _, ok := v.(int); return ok },
		reflect.TypeOf(true):  func(v any) bool { _, ok := v.(bool); return ok },
		reflect.TypeOf(0.0):   func(v any) bool { _, ok := v.(float64); return ok },
		reflect.TypeOf("str"): func(v any) bool { _, ok := v.(string); return ok },
	}

	if checkFunc, found := typeCheckers[reflect.TypeOf(variable)]; found {
		if !checkFunc(value) {
			return fmt.Errorf(TYPE_MISMATCH, reflect.TypeOf(variable), reflect.TypeOf(value))
		}
	}

	return nil
}

func (s *Scope) GetVariable(name string) (any, error) {
	variables := s.InScope(name)
	v := variables[name]
	if v == nil {
		return nil, fmt.Errorf(fmt.Sprintf(UNDEFINED_VARIABLE, name))
	}
	return v, nil
}
