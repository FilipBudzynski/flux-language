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

func (w *WhileStatement) Accept(v Visitor) {
	v.VisitWhileStatement(w)
}
