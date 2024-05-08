package ast


type FunctionCall struct {
	Arguments  []Expression
	Identifier Identifier
}

func NewFunctionCall(identifier Identifier, arguments []Expression) *FunctionCall {
	return &FunctionCall{
		Identifier: identifier,
		Arguments:  arguments,
	}
}

func (f *FunctionCall) Equals(other Expression) bool {
	if o, ok := other.(*FunctionCall); ok {
        if !f.Identifier.Equals(o.Identifier){
            return false
        }
        for i, e := range f.Arguments {
            if !f.Arguments[i].Equals(e){
                return false
            }
        }
    }
    return true
}

func (f *FunctionCall) Accept(v Visitor) {
    v.VisitFunctionCall(f)
}
