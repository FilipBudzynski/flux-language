package parser

type ReturnStatmenet struct {
    Expression Expression
}

func NewReturnStatement(expression Expression) *ReturnStatmenet {
    return &ReturnStatmenet{
        Expression: expression,
    }
}
