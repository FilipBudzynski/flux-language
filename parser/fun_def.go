package parser

import lex "tkom/lexer"

type FunDef struct {
	Name       string
	Type       lex.TokenTypes
	parameters []Parameter
	block      Block
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
