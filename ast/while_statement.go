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

// func (w *WhileStatement) Equals(other *WhileStatement) bool {
// 	if !w.Condition.Equals(other.Condition) {
// 		return false
// 	}
// 	if !w.InstructionsBlock.Equals(other.InstructionsBlock) {
// 		return false
// 	}
// 	return true
// }
