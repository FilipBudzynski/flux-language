package ast

type ReturnStatement struct {
	Expression Expression
}

func NewReturnStatement(expression Expression) *ReturnStatement {
	return &ReturnStatement{
		Expression: expression,
	}
}

func (r *ReturnStatement) Accept(v Visitor) {
	v.VisitReturnStatement(r)
}
