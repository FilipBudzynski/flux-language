package ast

import (
	lex "tkom/lexer"
)

type FunDef struct {
	Name       string
	Block      *Block
	Parameters []*Variable
	Type       TypeAnnotation
	Position   lex.Position
}

func NewFunctionDefinition(name string, parameters []*Variable, funType TypeAnnotation, block *Block, position lex.Position) *FunDef {
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
