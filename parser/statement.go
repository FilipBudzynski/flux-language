package parser

import lex "tkom/lexer"

// zwracanie interfejsów nie jest najelpszym pomysłem w go, ale jak mus to mus
type Statement interface{}

type Parameter struct {
	Name     string
	Type     lex.TokenType
	Position lex.Position
}

func NewParameter(name string, position lex.Position, parameterType lex.TokenType) Parameter {
	return Parameter{
		Name:     name,
		Position: position,
		Type:     parameterType,
	}
}
