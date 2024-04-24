package parser

import (
	"fmt"
	"reflect"
	lex "tkom/lexer"
)

const ERROR_WRONG_VALUE_IN_DECLARATION = "cannot use \"%s\", as %s value in variable declaration"

type Variable struct {
	Value     any
	Idetifier Identifier
	Type      lex.TokenType
}

func newVariable(variableType lex.TokenType, identifier Identifier, value any) Variable {

	// switch variableType {
	// case lex.INT:
	// 	v, err := convertValue(value, reflect.Int)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	value = v
	// case lex.FLOAT:
	// 	v, err := convertValue(value, reflect.Float64)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	value = v
	// case lex.BOOL:
	// 	v, err := convertValue(value, reflect.Bool)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	value = v
	// case lex.STRING:
	// 	v, err := convertValue(value, reflect.String)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	value = v
	// }

	return Variable{
		Type:      variableType,
		Idetifier: identifier,
		Value:     value,
	}
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
	return nil, fmt.Errorf(ERROR_WRONG_VALUE_IN_DECLARATION, expectedType, value)
}
