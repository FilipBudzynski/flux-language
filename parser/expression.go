package parser

import (
	"reflect"
	"tkom/lexer"
)

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

// type OperationExpression struct {
// 	LeftExpression  Expression
// 	RightExpression Expression
// 	Operation       Operation
// }
//
// func NewOperationExpression(leftExpression Expression, operation Operation, rightExpression Expression) OperationExpression {
// 	return OperationExpression{
// 		LeftExpression:  leftExpression,
// 		Operation:       operation,
// 		RightExpression: rightExpression,
// 	}
// }
//
// func (e *OperationExpression) Equals(other OperationExpression) bool {
// 	return reflect.DeepEqual(e, other)
// }

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

func (e *OrExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
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

func (e *AndExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
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

func (e *EqualsExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
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

func (e *NotEqualsExpression) Equals(other EqualsExpression) bool {
	return reflect.DeepEqual(e, other)
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

func (e *GreaterThanExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
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

func (e *LessThanExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
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

func (e *GreaterOrEqualExpression) Equals(other GreaterOrEqualExpression) bool {
	return reflect.DeepEqual(e, other)
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

func (e *LessOrEqualExpression) Equals(other LessOrEqualExpression) bool {
	return reflect.DeepEqual(e, other)
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

func (e *SumExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
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

func (e *SubstractExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
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

func (e *MultiplyExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
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

func (e *DivideExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
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

func (e *CastExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
}

type NegateExpression struct {
	Expression Expression
}

func NewNegateExpression(expression Expression) NegateExpression {
	return NegateExpression{
		Expression: expression,
	}
}

func (e *NegateExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
}
