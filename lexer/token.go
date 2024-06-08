package lexer

import (
	"fmt"
	"reflect"
	"tkom/shared"
)

const WRONG_TYPE_ERROR = "wrong value type to token type match, expected: %s, got: %s"

type Token struct {
	Value    any
	Type     TokenType
	Position shared.Position
}

func convertValue(value any, expectedType reflect.Kind) (any, error) {
	switch expectedType {
	case reflect.Int:
		if v, ok := value.(int); ok {
			return v, nil
		}
	case reflect.Float64:
		if v, ok := value.(float64); ok {
			return v, nil
		}
	case reflect.String:
		if v, ok := value.(string); ok {
			return v, nil
		}
	}
	return nil, fmt.Errorf(WRONG_TYPE_ERROR, expectedType, value)
}

func NewToken(tokenType TokenType, position shared.Position, value any) *Token {
	switch tokenType {
	case CONST_INT:
		v, err := convertValue(value, reflect.Int)
		if err != nil {
			panic(err)
		}
		value = v
	case CONST_FLOAT:
		v, err := convertValue(value, reflect.Float64)
		if err != nil {
			panic(err)
		}
		value = v
	case CONST_STRING, IDENTIFIER:
		v, err := convertValue(value, reflect.String)
		if err != nil {
			panic(err)
		}
		value = v
	}
	return &Token{
		Type:     tokenType,
		Position: position,
		Value:    value,
	}
}

func (b *Token) GetType() TokenType {
	return b.Type
}

var Operators = map[rune]TokenType{
	'+': PLUS,
	'-': MINUS,
	'*': MULTIPLY,
	'/': DIVIDE,
	'{': LEFT_BRACE,
	'}': RIGHT_BRACE,
	'(': LEFT_PARENTHESIS,
	')': RIGHT_PARENTHESIS,
	',': COMMA,
	'>': GREATER_THAN,
	'<': LESS_THAN,
	'=': ASSIGN,
	'!': NEGATE,
}

var KeyWords = map[string]TokenType{
	"int":     INT,
	"string":  STRING,
	"float":   FLOAT,
	"bool":    BOOL,
	"true":    CONST_TRUE,
	"false":   CONST_FALSE,
	"switch":  SWITCH,
	"while":   WHILE,
	"if":      IF,
	"else":    ELSE,
	"default": DEFAULT,
	"return":  RETURN,
	"and":     AND,
	"or":      OR,
	"as":      AS,
}

var DoubleOperators = map[string]TokenType{
	"<=": LESS_OR_EQUAL,
	">=": GREATER_OR_EQUAL,
	"==": EQUALS,
	"!=": NOT_EQUALS,
	"=>": CASE_ARROW,
	":=": DECLARE,
}

type TokenType int

const (
	IDENTIFIER TokenType = iota
	CONST_INT
	CONST_FLOAT
	CONST_STRING
	CONST_TRUE
	CONST_FALSE
	INT
	FLOAT
	STRING
	BOOL
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
	AND
	OR
	NEGATE
	IF
	ELSE
	WHILE
	SWITCH
	DEFAULT
	AS
	RETURN
	DECLARE
	ASSIGN
	CASE_ARROW
	LEFT_BRACE
	RIGHT_BRACE
	LEFT_PARENTHESIS
	RIGHT_PARENTHESIS
	COMMA
	ETX
	COMMENT
	UNDEFINED
)

var tokenTypeNames = [...]string{
	"IDENTIFIER",
	"CONST_INT",
	"CONST_FLOAT",
	"CONST_STRING",
	"CONST_TRUE",
	"CONST_FALSE",
	"INT",
	"FLOAT",
	"STRING",
	"BOOL",
	"PLUS",
	"MINUS",
	"MULTIPLY",
	"DIVIDE",
	"EQUALS",
	"NOT_EQUALS",
	"GREATER_THAN",
	"LESS_THAN",
	"GREATER_OR_EQUAL",
	"LESS_OR_EQUAL",
	"AND",
	"OR",
	"NEGATE",
	"IF",
	"ELSE",
	"WHILE",
	"SWITCH",
	"DEFAULT",
	"AS",
	"RETURN",
	"DECLARE",
	"ASSIGN",
	"CASE_ARROW",
	"LEFT_BRACE",
	"RIGHT_BRACE",
	"LEFT_PARENTHESIS",
	"RIGHT_PARENTHESIS",
	"COMMA",
	"ETX",
	"COMMENT",
	"UNDEFINED",
}

func (t TokenType) TypeName() string {
	if t < 0 || int(t) >= len(tokenTypeNames) {
		return "UNKNOWN"
	}
	return tokenTypeNames[t]
}
