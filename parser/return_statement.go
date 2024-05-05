package parser

type ReturnStatement struct {
    Expression Expression
}

func NewReturnStatement(expression Expression) *ReturnStatement {
    return &ReturnStatement{
        Expression: expression,
    }
}
