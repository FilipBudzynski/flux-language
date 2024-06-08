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

func (i *IfStatement) Accept(v Visitor) {
	v.VisitIfStatement(i)
}
