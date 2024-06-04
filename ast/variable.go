package ast

import (
	"tkom/shared"
)

const ERROR_WRONG_VALUE_IN_DECLARATION = "cannot use \"%s\", as %s value in variable declaration"

type Variable struct {
	Value    Expression
	Name     string
	Type     shared.TypeAnnotation
	Position shared.Position
}

func NewVariable(variableType shared.TypeAnnotation, name string, value Expression, position shared.Position) *Variable {
	return &Variable{
		Type:     variableType,
		Name:     name,
		Value:    value,
		Position: position,
	}
}

// func (v *Variable) Equals(other Variable) bool {
// 	if v.Type != other.Type {
// 		return false
// 	}
// 	if v.Value != other.Value {
// 		return false
// 	}
// 	if v.Name != other.Name {
// 		return false
// 	}
// 	if !reflect.DeepEqual(v.Position, other.Position) {
// 		return false
// 	}
// 	return true
// }

func (a *Variable) Accept(v Visitor) {
	v.VisitVariable(a)
}

type Assignment struct {
	Value      Expression
	Identifier *Identifier
}

func NewAssignment(identifier *Identifier, value Expression) *Assignment {
	return &Assignment{
		Identifier: identifier,
		Value:      value,
	}
}

// func (v *Assignemnt) Equals(other Assignemnt) bool {
// 	if v.Value != other.Value {
// 		return false
// 	}
// 	if !v.Identifier.Equals(other.Identifier) {
// 		return false
// 	}
// 	return false
// }

func (a *Assignment) Accept(v Visitor) {
	v.VisitAssignement(a)
}
