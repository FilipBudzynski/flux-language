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

func (a *Assignment) Accept(v Visitor) {
	v.VisitAssignement(a)
}
