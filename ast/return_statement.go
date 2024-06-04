package ast

type ReturnStatement struct {
	Value Expression

}

func NewReturnStatement(expression Expression) *ReturnStatement {
	return &ReturnStatement{
		Value: expression,
	}
}

func (r *ReturnStatement) Accept(v Visitor) {
	v.VisitReturnStatement(r)
}
