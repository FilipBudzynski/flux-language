package parser


type FunctionCall struct {
	Arguments  []Expression
	Identifier Identifier
}

func newFunctionCall(identifier Identifier, arguments []Expression) FunctionCall {
	return FunctionCall{
		Identifier: identifier,
		Arguments:  arguments,
	}
}
