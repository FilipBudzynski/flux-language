package parser

import lex "tkom/lexer"

type Identifier struct {
	Name     string
	Position lex.Position
}

func newIdentifier(name string, position lex.Position) Identifier {
	return Identifier{
		Name:     name,
		Position: position,
	}
}
