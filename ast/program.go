package ast

type Program struct {
	Functions map[string]*FunctionDefinition
}

func NewProgram(functions map[string]*FunctionDefinition) *Program {
	return &Program{Functions: functions}
}

// func (p *Program) Equals(other *Program) bool {
// 	if len(p.functions) != len(other.functions) {
// 		return false
// 	}
//
// 	for k, v := range p.functions {
// 		if !v.Equals(other.functions[k]) {
// 			return false
// 		}
// 	}
//
// 	return true
// }

func (p *Program) Accept(v Visitor) {
	v.VisitProgram(p)
}
