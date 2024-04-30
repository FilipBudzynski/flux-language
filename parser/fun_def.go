package parser

import lex "tkom/lexer"

type FunDef struct {
	Type       *lex.TokenType
	Parameters []*Variable
	Name       string
	Statements []Statement
	Position   lex.Position
}

func NewFunctionDefinition(name string, parameters []*Variable, funType *lex.TokenType, statements []Statement, position lex.Position) *FunDef {
	return &FunDef{
		Name:       name,
		Type:       funType,
		Parameters: parameters,
		Statements: statements,
		Position:   position,
	}
}
