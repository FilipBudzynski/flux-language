package parser

import lex "tkom/lexer"

type FunDef struct {
	Type       *lex.TokenType
	Parameters []*Variable
	Name       string
	Block      Block
	Position   lex.Position
}

func NewFunctionDefinition(name string, parameters []*Variable, funType *lex.TokenType, block Block, position lex.Position) *FunDef {
	return &FunDef{
		Name:       name,
		Type:       funType,
		Parameters: parameters,
		Block:      block,
		Position:   position,
	}
}
