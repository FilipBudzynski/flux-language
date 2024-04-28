package parser

import lex "tkom/lexer"

type FunctionCall struct {
	Arguments  []Expression
	Identifier Identifier
	Position   lex.Position
}

func newFunctionCall(identifier Identifier, arguments []Expression, position lex.Position) FunctionCall {
	return FunctionCall{
		Identifier: identifier,
		Arguments:  arguments,
		Position:   position,
	}
}
