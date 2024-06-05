package ast

type Program struct {
	Functions map[string]*FunctionDefinition
}

func NewProgram(functions map[string]*FunctionDefinition) *Program {
	return &Program{Functions: functions}
}

func (p *Program) Accept(v Visitor) {
	v.VisitProgram(p)
}
