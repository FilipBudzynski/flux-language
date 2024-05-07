package ast

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

var ValidTypeAnnotation = map[lexer.TokenType]TypeAnnotation{
	lexer.INT:    INT,
	lexer.FLOAT:  FLOAT,
	lexer.BOOL:   BOOL,
	lexer.STRING: STRING,
}

type Expression interface{}

type OrExpression struct {
	LeftExpression  Expression
	RightExpression Expression
	Position        lexer.Position
}

func NewOrExpression(leftExpression Expression, rightExpression Expression, position lexer.Position) OrExpression {
	return OrExpression{
		LeftExpression:  leftExpression,
		RightExpression: rightExpression,
		Position:        position,
	}
}

func (e *OrExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
}

type AndExpression struct {
	LeftExpression  Expression
	RightExpression Expression
	Position        lexer.Position
}

func NewAndExpression(leftExpression Expression, rightExpression Expression, position lexer.Position) AndExpression {
	return AndExpression{
		LeftExpression:  leftExpression,
		RightExpression: rightExpression,
		Position:        position,
	}
}

func (e *AndExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
}

type EqualsExpression struct {
	LeftExpression  Expression
	RightExpression Expression
	Position        lexer.Position
}

func NewEqualsExpression(leftExpression Expression, rightExpression Expression, position lexer.Position) Expression {
	return EqualsExpression{
		LeftExpression:  leftExpression,
		RightExpression: rightExpression,
		Position:        position,
	}
}

func (e *EqualsExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
}

type NotEqualsExpression struct {
	LeftExpression  Expression
	RightExpression Expression
	Position        lexer.Position
}

func NewNotEqualsExpression(leftExpression Expression, rightExpression Expression, position lexer.Position) Expression {
	return NotEqualsExpression{
		LeftExpression:  leftExpression,
		RightExpression: rightExpression,
		Position:        position,
	}
}

func (e *NotEqualsExpression) Equals(other EqualsExpression) bool {
	return reflect.DeepEqual(e, other)
}

type GreaterThanExpression struct {
	LeftExpression  Expression
	RightExpression Expression
	Position        lexer.Position
}

func NewGreaterThanExpression(leftExpression Expression, rightExpression Expression, position lexer.Position) Expression {
	return GreaterThanExpression{
		LeftExpression:  leftExpression,
		RightExpression: rightExpression,
		Position:        position,
	}
}

func (e *GreaterThanExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
}

type LessThanExpression struct {
	LeftExpression  Expression
	RightExpression Expression
	Position        lexer.Position
}

func NewLessThanExpression(leftExpression Expression, rightExpression Expression, position lexer.Position) Expression {
	return LessThanExpression{
		LeftExpression:  leftExpression,
		RightExpression: rightExpression,
		Position:        position,
	}
}

func (e *LessThanExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
}

type GreaterOrEqualExpression struct {
	LeftExpression  Expression
	RightExpression Expression
	Position        lexer.Position
}

func NewGreaterOrEqualExpression(leftExpression Expression, rightExpression Expression, position lexer.Position) Expression {
	return GreaterOrEqualExpression{
		LeftExpression:  leftExpression,
		RightExpression: rightExpression,
		Position:        position,
	}
}

func (e *GreaterOrEqualExpression) Equals(other GreaterOrEqualExpression) bool {
	return reflect.DeepEqual(e, other)
}

type LessOrEqualExpression struct {
	LeftExpression  Expression
	RightExpression Expression
	Position        lexer.Position
}

func NewLessOrEqualExpression(leftExpression Expression, rightExpression Expression, position lexer.Position) Expression {
	return LessOrEqualExpression{
		LeftExpression:  leftExpression,
		RightExpression: rightExpression,
		Position:        position,
	}
}

func (e *LessOrEqualExpression) Equals(other LessOrEqualExpression) bool {
	return reflect.DeepEqual(e, other)
}

type SumExpression struct {
	LeftExpression  Expression
	RightExpression Expression
	Position        lexer.Position
}

func NewSumExpression(leftExpression Expression, rightExpression Expression, position lexer.Position) Expression {
	return SumExpression{
		LeftExpression:  leftExpression,
		RightExpression: rightExpression,
		Position:        position,
	}
}

func (e *SumExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
}

type SubstractExpression struct {
	LeftExpression  Expression
	RightExpression Expression
	Position        lexer.Position
}

func NewSubstractExpression(leftExpression Expression, rightExpression Expression, position lexer.Position) Expression {
	return SubstractExpression{
		LeftExpression:  leftExpression,
		RightExpression: rightExpression,
		Position:        position,
	}
}

func (e *SubstractExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
}

type MultiplyExpression struct {
	LeftExpression  Expression
	RightExpression Expression
	Position        lexer.Position
}

func NewMultiplyExpression(leftExpression Expression, rightExpression Expression, position lexer.Position) Expression {
	return MultiplyExpression{
		LeftExpression:  leftExpression,
		RightExpression: rightExpression,
		Position:        position,
	}
}

func (e *MultiplyExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
}

type DivideExpression struct {
	LeftExpression  Expression
	RightExpression Expression
	Position        lexer.Position
}

func NewDivideExpression(leftExpression Expression, rightExpression Expression, position lexer.Position) Expression {
	return DivideExpression{
		LeftExpression:  leftExpression,
		RightExpression: rightExpression,
		Position:        position,
	}
}

func (e *DivideExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
}

type CastExpression struct {
	LeftExpression Expression
	TypeAnnotation Operation
	Position       lexer.Position
}

func NewCastExpression(leftExpression Expression, typeAnnotation Operation, position lexer.Position) CastExpression {
	return CastExpression{
		LeftExpression: leftExpression,
		TypeAnnotation: typeAnnotation,
		Position:       position,
	}
}

func (e *CastExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
}

type NegateExpression struct {
	Expression Expression
	Position   lexer.Position
}

func NewNegateExpression(expression Expression, position lexer.Position) NegateExpression {
	return NegateExpression{
		Expression: expression,
		Position:   position,
	}
}

func (e *NegateExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
}

type IntExpression struct {
	Value    int
	Position lexer.Position
}

func NewIntExpression(value int, position lexer.Position) IntExpression {
	return IntExpression{
		Value:    value,
		Position: position,
	}
}

func (e *IntExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
}

type FloatExpression struct {
	Value    float64
	Position lexer.Position
}

func NewFloatExpression(value float64, position lexer.Position) FloatExpression {
	return FloatExpression{
		Value:    value,
		Position: position,
	}
}

func (e *FloatExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
}

type BoolExpression struct {
	Value    bool
	Position lexer.Position
}

func NewBoolExpression(value bool, position lexer.Position) BoolExpression {
	return BoolExpression{
		Value:    value,
		Position: position,
	}
}

func (e *BoolExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
}

type StringExpression struct {
	Value    string
	Position lexer.Position
}

func NewStringExpression(value string, position lexer.Position) StringExpression {
	return StringExpression{
		Value:    value,
		Position: position,
	}
}

func (e *StringExpression) Equals(other Expression) bool {
	return reflect.DeepEqual(e, other)
}
