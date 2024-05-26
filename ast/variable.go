package ast

import (
	"reflect"
	"tkom/shared"
)

const ERROR_WRONG_VALUE_IN_DECLARATION = "cannot use \"%s\", as %s value in variable declaration"

type Variable struct {
	Value    any
	Name     string
	Type     shared.TypeAnnotation
	Position shared.Position
}

func NewVariable(variableType shared.TypeAnnotation, name string, value any, position shared.Position) *Variable {
	return &Variable{
		Type:     variableType,
		Name:     name,
		Value:    value,
		Position: position,
	}
}

func (v *Variable) Equals(other Variable) bool {
	if v.Type != other.Type {
		return false
	}
	if v.Value != other.Value {
		return false
	}
	if v.Name != other.Name {
		return false
	}
	if !reflect.DeepEqual(v.Position, other.Position) {
		return false
	}
	return true
}

func (a *Variable) Accept(v Visitor) {
	v.VisitVariable(a)
}

type Assignemnt struct {
	Value      Expression
	Identifier Identifier
}

func NewAssignment(identifier Identifier, value Expression) *Assignemnt {
	return &Assignemnt{
		Identifier: identifier,
		Value:      value,
	}
}

func (v *Assignemnt) Equals(other Assignemnt) bool {
	if v.Value != other.Value {
		return false
	}
	if !v.Identifier.Equals(&other.Identifier) {
		return false
	}
	return false
}

func (a *Assignemnt) Accept(v Visitor) {
	v.VisitAssignement(a)
}
