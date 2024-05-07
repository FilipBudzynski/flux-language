package ast

import "reflect"

type IfStatement struct {
	Conditions       Expression
	Instructions     []Statement
	ElseInstructions []Statement
}

func NewIfStatement(conditions Expression, instructions []Statement, elseInstructions []Statement) *IfStatement {
	return &IfStatement{
		Conditions:       conditions,
		Instructions:     instructions,
		ElseInstructions: elseInstructions,
	}
}

func (i *IfStatement) Equals(other *IfStatement) bool {
	if reflect.DeepEqual(i.Conditions, other.Conditions) {
		return false
	}
	if len(i.Instructions) != len(other.Instructions) {
		return false
	}
	for i, instruction := range i.Instructions {
		if !reflect.DeepEqual(instruction, other.Instructions[i]) {
			return false
		}
	}
	if len(i.ElseInstructions) != len(other.ElseInstructions) {
		return false
	}
	for i, instruction := range i.ElseInstructions {
		if !reflect.DeepEqual(instruction, other.ElseInstructions[i]) {
			return false
		}
	}
	return true
}
