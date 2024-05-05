package parser


type FunctionCall struct {
	Arguments  []Expression
	Identifier Identifier
}

func NewFunctionCall(identifier Identifier, arguments []Expression) FunctionCall {
	return FunctionCall{
		Identifier: identifier,
		Arguments:  arguments,
	}
}
