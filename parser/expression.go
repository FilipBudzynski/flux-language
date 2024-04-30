package parser

import "tkom/lexer"

// import lex "tkom/lexer"

type (
	Operation      interface{}
	OperationType  int
	TypeAnnotation int
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

const (
	INT TypeAnnotation = iota
	FLOAT
	BOOL
	STRING
)

var validTypes = map[lexer.TokenType]TypeAnnotation{
	lexer.INT:    INT,
	lexer.FLOAT:  FLOAT,
	lexer.BOOL:   BOOL,
	lexer.STRING: STRING,
}

type Expression interface{}

type OperationExpression struct {
	LeftExpression  Expression
	RightExpression Expression
	Operation       Operation
}

func NewOperationExpression(leftExpression Expression, operation Operation, rightExpression Expression) OperationExpression {
	return OperationExpression{
		LeftExpression:  leftExpression,
		Operation:       operation,
		RightExpression: rightExpression,
	}
}

type OrExpression struct {
	LeftExpression  Expression
	RightExpression Expression
}

func NewOrExpression(leftExpression Expression, rightExpression Expression) OrExpression {
	return OrExpression{
		LeftExpression:  leftExpression,
		RightExpression: rightExpression,
	}
}

type AndExpression struct {
	LeftExpression  Expression
	RightExpression Expression
}

func NewAndExpression(leftExpression Expression, rightExpression Expression) AndExpression {
	return AndExpression{
		LeftExpression:  leftExpression,
		RightExpression: rightExpression,
	}
}


type EqualsExpression struct {
	LeftExpression  Expression
	RightExpression Expression
}

func NewEqualsExpression(leftExpression Expression, rightExpression Expression) EqualsExpression {
	return EqualsExpression{
		LeftExpression:  leftExpression,
		RightExpression: rightExpression,
	}
}

type NotEqualsExpression struct {
	LeftExpression  Expression
	RightExpression Expression
}

func NewNotEqualsExpression(leftExpression Expression, rightExpression Expression) NotEqualsExpression {
	return NotEqualsExpression{
		LeftExpression:  leftExpression,
		RightExpression: rightExpression,
	}
}

type GreaterThanExpression struct {
	LeftExpression  Expression
	RightExpression Expression
}

func NewGreaterThanExpression(leftExpression Expression, rightExpression Expression) GreaterThanExpression {
	return GreaterThanExpression{
		LeftExpression:  leftExpression,
		RightExpression: rightExpression,
	}
}

type LessThanExpression struct {
	LeftExpression  Expression
	RightExpression Expression
}

func NewLessThanExpression(leftExpression Expression, rightExpression Expression) LessThanExpression {
	return LessThanExpression{
		LeftExpression:  leftExpression,
		RightExpression: rightExpression,
	}
}

type GreaterOrEqualExpression struct {
	LeftExpression  Expression
	RightExpression Expression
}

func NewGreaterOrEqualExpression(leftExpression Expression, rightExpression Expression) GreaterOrEqualExpression {
	return GreaterOrEqualExpression{
		LeftExpression:  leftExpression,
		RightExpression: rightExpression,
	}
}

type LessOrEqualExpression struct {
	LeftExpression  Expression
	RightExpression Expression
}

func NewLessOrEqualExpression(leftExpression Expression, rightExpression Expression) LessOrEqualExpression {
	return LessOrEqualExpression{
		LeftExpression:  leftExpression,
		RightExpression: rightExpression,
	}
}

type SumExpression struct {
	LeftExpression  Expression
	RightExpression Expression
}

func NewSumExpression(leftExpression Expression, rightExpression Expression) SumExpression {
	return SumExpression{
		LeftExpression:  leftExpression,
		RightExpression: rightExpression,
	}
}

type SubstractExpression struct {
	LeftExpression  Expression
	RightExpression Expression
}

func NewSubstractExpression(leftExpression Expression, rightExpression Expression) SubstractExpression {
	return SubstractExpression{
		LeftExpression:  leftExpression,
		RightExpression: rightExpression,
	}
}

type MultiplyExpression struct {
	LeftExpression  Expression
	RightExpression Expression
}

func NewMultiplyExpression(leftExpression Expression, rightExpression Expression) MultiplyExpression {
	return MultiplyExpression{
		LeftExpression:  leftExpression,
		RightExpression: rightExpression,
	}
}

type DivideExpression struct {
	LeftExpression  Expression
	RightExpression Expression
}

func NewDivideExpression(leftExpression Expression, rightExpression Expression) DivideExpression {
	return DivideExpression{
		LeftExpression:  leftExpression,
		RightExpression: rightExpression,
	}
}

type CastExpression struct {
	LeftExpression Expression
	TypeAnnotation Operation
}

func NewCastExpression(leftExpression Expression, typeAnnotation Operation) CastExpression {
	return CastExpression{
		LeftExpression: leftExpression,
		TypeAnnotation: typeAnnotation,
	}
}

type NegateExpression struct {
	Expression Expression
}

func NewNegateExpression(expression Expression) NegateExpression {
	return NegateExpression{
		Expression: expression,
	}
}
