package parser

// import lex "tkom/lexer"

type (
	Operation     interface{}
	OperationType int
)

const (
	OR OperationType = iota
	AND
	PLUS
	MINUS
	MULTIPLY
	DIVIDE
	EQUALS
	NOT_EQUALS
	GREATER_THAN
	LESS_THAN
	GREATER_OR_EQUAL
	LESS_OR_EQUAL
	NEGATE
	AS
)

type (
	Expression          interface{}
	OperationExpression struct {
		LeftExpression  Expression
		RightExpression Expression
		Operation       Operation
	}
)

func NewExpression(leftExpression Expression, operation Operation, rightExpression Expression) OperationExpression {
	return OperationExpression{
		LeftExpression:  leftExpression,
		Operation:       operation,
		RightExpression: rightExpression,
	}
}

//TODO: is it even needed???
//
// type CastedTerm struct {
// 	Term           *UnaryTerm
// 	TypeAnnotation *lex.TokenType
// }
//
// func NewCastedTerm(term *UnaryTerm, typeAnnotation *lex.TokenType) *CastedTerm {
// 	return &CastedTerm{
// 		Term:           term,
// 		TypeAnnotation: typeAnnotation,
// 	}
// }
//
// type UnaryTerm struct {
// 	Value  any
// 	Negate *lex.TokenType
// }
//
// func NewUnaryTerm(value any, negate *lex.TokenType) *UnaryTerm {
// 	return &UnaryTerm{
// 		Value:  value,
// 		Negate: negate,
// 	}
// }
//
// type Term struct {
// 	Value any
// }
//
// func NewTerm(value any) *Term {
// 	return &Term{
// 		Value: value,
// 	}
// }
