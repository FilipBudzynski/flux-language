package ast

import (
	"reflect"
	lex "tkom/lexer"
)

type FunDef struct {
	Type       *TypeAnnotation
	Parameters []*Variable
	Name       string
	Statements []Statement
	Position   lex.Position
}

func NewFunctionDefinition(name string, parameters []*Variable, funType *TypeAnnotation, statements []Statement, position lex.Position) *FunDef {
	return &FunDef{
		Name:       name,
		Type:       funType,
		Parameters: parameters,
		Statements: statements,
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

	for i, statement := range f.Statements {
		if !reflect.DeepEqual(statement, other.Statements[i]) {
			return false
		}
	}
	return true
}

func (f *FunDef) Accept(v Visitor) {
    v.VisitFunDef(f)
}
