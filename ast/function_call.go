package ast

import (
	"reflect"
	"tkom/lexer"
)

type FunctionCall struct {
	Name      string
	Arguments []Expression
	Position  lexer.Position
}

func NewFunctionCall(name string, position lexer.Position, arguments []Expression) *FunctionCall {
	return &FunctionCall{
		Name:      name,
		Position:  position,
		Arguments: arguments,
	}
}

func (f *FunctionCall) Equals(other Expression) bool {
	if o, ok := other.(*FunctionCall); ok {
		if f.Name != o.Name {
			return false
		}
		if !reflect.DeepEqual(f.Position, o.Position) {
			return false
		}
		for i, e := range f.Arguments {
			if !f.Arguments[i].Equals(e) {
				return false
			}
		}
	}
	return true
}

func (f *FunctionCall) Accept(v Visitor) {
	v.VisitFunctionCall(f)
}
