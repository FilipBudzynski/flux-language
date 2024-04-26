package parser

import lex "tkom/lexer"

//TODO: czy idetifier powinien miec typ??
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
