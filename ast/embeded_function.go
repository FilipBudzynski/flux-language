package ast

type EmbeddedFunction struct {
	Func       func(...any) any
	Name       string
	Parameters []any
	Variadic   bool
}

func (ef *EmbeddedFunction) Accept(v Visitor) {
	v.VisitEmbeddedFunction(ef)
}

func (ef *EmbeddedFunction) GetParametersLen() int {
	if ef.Variadic {
		return -1
	}
	return len(ef.Parameters)
}

func (ef *EmbeddedFunction) IsVariadic() bool {
	return ef.Variadic
}
