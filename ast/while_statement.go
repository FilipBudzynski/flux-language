package ast

type WhileStatement struct {
	Condition    Expression
	InstructionsBlock *Block
}

func NewWhileStatement(condition Expression, instructions *Block) *WhileStatement {
	return &WhileStatement{
		Condition:    condition,
		InstructionsBlock: instructions,
	}
}

func (w *WhileStatement) Accept(v Visitor) {
	v.VisitWhileStatement(w)
}

