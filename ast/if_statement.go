package ast

type IfStatement struct {
	Condition             Expression
	InstructionsBlock     *Block
	ElseInstructionsBlock *Block
}

func NewIfStatement(conditions Expression, instructions *Block, elseInstructions *Block) *IfStatement {
	return &IfStatement{
		Condition:             conditions,
		InstructionsBlock:     instructions,
		ElseInstructionsBlock: elseInstructions,
	}
}

// func (i *IfStatement) Equals(other *IfStatement) bool {
// 	if reflect.DeepEqual(i.Condition, other.Condition) {
// 		return false
// 	}
// 	if !i.InstructionsBlock.Equals(other.InstructionsBlock) {
// 		return false
// 	}
// 	if !i.ElseInstructionsBlock.Equals(other.ElseInstructionsBlock) {
// 		return false
// 	}
// 	return true
// }

func (i *IfStatement) Accept(v Visitor) {
	v.VisitIfStatement(i)
}
