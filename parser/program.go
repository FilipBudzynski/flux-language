package parser

type Program struct {
	functions map[string]FunDef
}

func NewProgram(functions map[string]FunDef) *Program {
	return &Program{functions: functions}
}
