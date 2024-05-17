package ast

type WhileStatement struct {
	Condition    Expression
	Instructions *Block
}

func NewWhileStatement(condition Expression, instructions *Block) *WhileStatement {
	return &WhileStatement{
		Condition:    condition,
		Instructions: instructions,
	}
}

func (w *WhileStatement) Accept(v Visitor) {
	v.VisitWhileStatement(w)
}

func (w *WhileStatement) Equals(other *WhileStatement) bool {
	if !w.Condition.Equals(other.Condition) {
		return false
	}
	if !w.Instructions.Equals(other.Instructions) {
		return false
	}
	return true
}
