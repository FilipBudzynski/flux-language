package parser

import lex "tkom/lexer"

type Block struct {
	Statements []*Statement
}

type Statement struct{}

type Parameter struct {
	Name     string
	Type     lex.TokenTypes
	Position lex.Position
}

func newParameter(name string, position lex.Position) Parameter {
	return Parameter{
		Name:     name,
		Position: position,
	}
}
