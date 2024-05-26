package ast

import (
	"tkom/shared"
)

type FunDef struct {
	Name       string
	Block      *Block
	Parameters []*Variable
	Type       shared.TypeAnnotation
	Position   shared.Position
}

func NewFunctionDefinition(name string, parameters []*Variable, funType shared.TypeAnnotation, block *Block, position shared.Position) *FunDef {
	return &FunDef{
		Name:       name,
		Type:       funType,
		Parameters: parameters,
		Block:      block,
		Position:   position,
	}
}

func (f *FunDef) Equals(other *FunDef) bool {
	if f.Type != other.Type {
		return false
	}
	if f.Name != other.Name {
		return false
	}
	if len(f.Parameters) != len(other.Parameters) {
		return false
	}
	for i, param := range f.Parameters {
		if !param.Equals(*other.Parameters[i]) {
			return false
		}
	}
	return f.Block.Equals(other.Block)
}

func (f *FunDef) Accept(v Visitor) {
	v.VisitFunDef(f)
}
