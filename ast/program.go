package ast

type Program struct {
	functions map[string]*FunDef
}

func NewProgram(functions map[string]*FunDef) *Program {
	return &Program{functions: functions}
}

func (p *Program) Equals(other *Program) bool {
	if len(p.functions) != len(other.functions) {
		return false
	}

	for k, v := range p.functions {
	    if !v.Equals(other.functions[k]) {
	        return false
	    }
	}
    
	return true
}
