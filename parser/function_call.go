package parser

import lex "tkom/lexer"

type FunctionCall struct {
	Identifier Identifier
	Arguments  []Variable
	Position   lex.Position
}

func newFunctionCall(identifier Identifier, arguments []Variable, position lex.Position) FunctionCall {
	return FunctionCall{
		Identifier: identifier,
		Arguments:  arguments,
        Position: position,
	}
}
