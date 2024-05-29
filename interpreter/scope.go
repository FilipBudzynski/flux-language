package interpreter

import (
	"fmt"
	"tkom/shared"
)

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
func (s *Scope) SetVariableValue(name string, value any) error {
	v := s.InScope(name)
	if v == nil {
		return NewSemanticError(fmt.Sprintf(UNDEFINED_VARIABLE, name), shared.NewPosition(999, 999))
	}

	err := s.CheckVariableType(v, value)
	if err != nil {
		return err
	}
	v.Value = value

	return nil
}

func (s *Scope) CheckVariableType(variable *ScopeVariable, value any) error {
	switch variable.Type {
	case shared.INT:
		if _, ok := value.(int); !ok {
			return NewSemanticError(fmt.Sprintf(TYPE_MISMATCH, value, variable.Type), variable.Position)
		}
	case shared.BOOL:
		if _, ok := value.(bool); !ok {
			return NewSemanticError(fmt.Sprintf(TYPE_MISMATCH, value, variable.Type), variable.Position)
		}
	case shared.FLOAT:
		if _, ok := value.(float64); !ok {
			return NewSemanticError(fmt.Sprintf(TYPE_MISMATCH, value, variable.Type), variable.Position)
		}
	case shared.STRING:
		if _, ok := value.(string); !ok {
			return NewSemanticError(fmt.Sprintf(TYPE_MISMATCH, value, variable.Type), variable.Position)
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
		return nil, NewSemanticError(fmt.Sprintf(UNDEFINED_VARIABLE, name), shared.NewPosition(999, 999))
	}
	return v, nil
}
