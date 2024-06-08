package ast

import (
	"reflect"
	"tkom/shared"
)

type Block struct {
	Statements []Statement
}

func NewBlock(statements []Statement) *Block {
	return &Block{
		Statements: statements,
	}
}

func (b *Block) Accept(v Visitor) {
	v.VisitBlock(b)
}

func (b *Block) Equals(other Expression) bool {
	otherBlock, ok := other.(*Block)
	if !ok {
		return false
	}
	for i, s := range b.Statements {
		if !reflect.DeepEqual(s, otherBlock.Statements[i]) {
			return false
		}
	}
	return true
}

func (b *Block) GetPosition() shared.Position {
	return shared.Position{}
}
