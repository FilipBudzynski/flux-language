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
	Value    any
	Type     shared.TypeAnnotation
	Position shared.Position
}

type Scope struct {
	Parent     *Scope
	variables  map[string]*ScopeVariable
	ReturnType *shared.TypeAnnotation
}

func NewScope(parent *Scope, returnType *shared.TypeAnnotation) *Scope {
	return &Scope{
		Parent:     parent,
		variables:  map[string]*ScopeVariable{},
		ReturnType: returnType,
	}
}

func (s *Scope) InScope(name string) *ScopeVariable {
	currentScope := s
	for currentScope != nil {
		if v, ok := currentScope.variables[name]; ok {
			return v
		}
		if currentScope.ReturnType != nil {
			break
		}
		currentScope = currentScope.Parent
	}
	return nil
}

func (s *Scope) AddVariable(name string, value any, variableType shared.TypeAnnotation, position shared.Position) error {
	if v, ok := s.variables[name]; ok {
		return NewSemanticError(fmt.Sprintf(REDECLARED_VARIABLE, name), v.Position)
	}
	scopeVariable := &ScopeVariable{
		Value:    value,
		Type:     variableType,
		Position: position,
	}

	s.variables[name] = scopeVariable
	return nil
}

// Sets a value to variable in the scope
//
// If no such value is foud returns a UNDEFINED_VARIABLE error
//
// If variable type doesn't match with value type, returns a TYPE_MISMATCH error
func (s *Scope) SetValue(name string, value any) error {
	v := s.InScope(name)
	if v == nil {
		return fmt.Errorf(UNDEFINED_VARIABLE, name)
	}

	err := s.CheckVariableType(v, value)
	if err != nil {
		return err
	}
	v.Value = value

	return nil
}

func (s *Scope) CheckVariableType(variable *ScopeVariable, value any) error {
	typeCheckers := map[shared.TypeAnnotation]func(any) bool{
		shared.INT:    func(v any) bool { _, ok := v.(int); return ok },
		shared.BOOL:   func(v any) bool { _, ok := v.(bool); return ok },
		shared.FLOAT:  func(v any) bool { _, ok := v.(float64); return ok },
		shared.STRING: func(v any) bool { _, ok := v.(string); return ok },
	}

	if checkFunc, found := typeCheckers[variable.Type]; found {
		if !checkFunc(value) {
			return fmt.Errorf(TYPE_MISMATCH, variable.Type, reflect.TypeOf(value))
		}
	}

	return nil
}

func (s *Scope) GetValue(name string) (any, error) {
	v := s.InScope(name)
	if v == nil {
		return nil, NewSemanticError(fmt.Sprintf(UNDEFINED_VARIABLE, name), shared.NewPosition(999, 999))
	}
	return v.Value, nil
}

func (s *Scope) GetVariable(name string) (*ScopeVariable, error) {
	v := s.InScope(name)
	if v == nil {
		return nil, fmt.Errorf(fmt.Sprintf(UNDEFINED_VARIABLE, name))
	}
	return v, nil
}
