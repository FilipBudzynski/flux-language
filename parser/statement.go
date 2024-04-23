package parser

import lex "tkom/lexer"

// zwracanie interfejsów nie jest najelpszym pomysłem w go, ale jak mus to mus
type Statement interface {
	Accept()
}

type Parameter struct {
	Name     string
	Type     lex.TokenTypes
	Position lex.Position
}

func NewParameter(name string, position lex.Position) Parameter {
	return Parameter{
		Name:     name,
		Position: position,
	}
}
