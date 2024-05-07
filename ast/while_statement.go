package ast

type WhileStatement struct {
	Condition    Expression
	Instructions []Statement
}

func NewWhileStatement(condition Expression, instructions []Statement) *WhileStatement {
	return &WhileStatement{
		Condition:    condition,
		Instructions: instructions,
	}
}
