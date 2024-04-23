package parser

import lex "tkom/lexer"

type FunDef struct {
	Name       string
	parameters []Parameter
	block      Block
	Type       lex.TokenTypes
	Position   lex.Position
}

func NewFunctionDefinition(name string, parameters []Parameter, funType *lex.TokenTypes, block *Block, position lex.Position) *FunDef {
	return &FunDef{
		Name:       name,
		Type:       *funType,
		parameters: parameters,
		block:      *block,
		Position:   position,
	}
}
