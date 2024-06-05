package ast

import (
	"tkom/shared"
)

type FunctionDefinition struct {
	Name       string
	Block      *Block
	Parameters []*Variable
	Type       shared.TypeAnnotation
	Position   shared.Position
}

func NewFunctionDefinition(name string, parameters []*Variable, funType shared.TypeAnnotation, block *Block, position shared.Position) *FunctionDefinition {
	return &FunctionDefinition{
		Name:       name,
		Type:       funType,
		Parameters: parameters,
		Block:      block,
		Position:   position,
	}
}

func (f *FunctionDefinition) GetParametersLen() int {
	return len(f.Parameters)
}

// every user based function is not variadic
// only embeded functions such as 'print' are variadic
func (f *FunctionDefinition) IsVariadic() bool {
	return false
}

func (f *FunctionDefinition) Accept(v Visitor) {
	v.VisitFunctionDefinition(f)
}
