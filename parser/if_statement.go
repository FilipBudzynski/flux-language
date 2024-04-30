package parser

type IfStatement struct {
	Conditions       Expression
	Instructions     []Statement
	ElseInstructions []Statement
}

func NewIfStatement(conditions Expression, instructions []Statement, elseInstructions []Statement) *IfStatement {
	return &IfStatement{
		Conditions:       conditions,
		Instructions:     instructions,
		ElseInstructions: elseInstructions,
	}
}
