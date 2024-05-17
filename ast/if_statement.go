package ast

import "reflect"

type IfStatement struct {
	Condition        Expression
	Instructions     *Block
	ElseInstructions *Block
}

func NewIfStatement(conditions Expression, instructions *Block, elseInstructions *Block) *IfStatement {
	return &IfStatement{
		Condition:        conditions,
		Instructions:     instructions,
		ElseInstructions: elseInstructions,
	}
}

func (i *IfStatement) Equals(other *IfStatement) bool {
	if reflect.DeepEqual(i.Condition, other.Condition) {
		return false
	}
    if !i.Instructions.Equals(other.Instructions) {
        return false
    }
    if !i.ElseInstructions.Equals(other.ElseInstructions) {
        return false
    }
	return true
}

func (i *IfStatement) Accept(v Visitor) {
	v.VisitIfStatement(i)
}
