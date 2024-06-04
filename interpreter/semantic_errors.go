package interpreter

import (
	"fmt"
	"tkom/shared"
)

type SemantciError struct {
	Message  string
	Position shared.Position
}

func NewSemanticError(message string, position shared.Position) *SemantciError {
	msg := fmt.Sprintf("error [%v, %v]: %s", position.Line, position.Column, message)
	return &SemantciError{
		Message:  msg,
		Position: position,
	}
}

func (err *SemantciError) Error() string {
	return err.Message
}

const (
	UNDEFINED_VARIABLE                       = "undefined: %s"
	UNDEFINED_FUNCTION                       = "undefined function: %s"
	REDECLARED_VARIABLE                      = "redeclared variable: %s, variable with that name already exists"
	TYPE_MISMATCH                            = "type mismatch: expected %v, got %v"
	WRONG_NUMBER_OF_ARGUMENTS                = "function %s expects %d arguments but got: %d"
	INVALID_NEGATE_EXPRESSION                = "cannot negate %v of type %v"
	INVALID_MULTIPLY_EXPRESSION              = "cannot evaluate '*' operation with instances of %v and %v"
	INVALID_DIVISION_EXPRESSION              = "cannot evaluate '/' operation with instances of %v and %v"
	INVALID_SUM_EXPRESSION                   = "cannot evaluate '+' operation with instances of %v and %v"
	INVALID_SUBSTRACT_EXPRESSION             = "cannot evaluate '-' operation with instances of %v and %v"
	INVALID_EQUALS_MISSMATCH                 = "cannot evaluate '==' operation with instances, mismatched types of %v and %v"
	INVALID_NOT_EQUALS_MISSMATCH             = "cannot evaluate '!=' operation with instances, mismatched types of %v and %v"
	INVALID_GREATER_THAN_MISSMATCH           = "cannot evaluate '>' operation with instances, mismatched types of %v and %v"
	INVALID_GREATER_OR_EQUALS_THAN_MISSMATCH = "cannot evaluate '>=' operation with instances, mismatched types of %v and %v"
	INVALID_LESS_THAN_MISSMATCH              = "cannot evaluate '<' operation with instances, mismatched types of %v and %v"
	INVALID_LESS_OR_EQUALS_THAN_MISSMATCH    = "cannot evaluate '<=' operation with instances, mismatched types of %v and %v"
	INVALID_ASSIGNMENT_TYPES                 = "cannot assign value of type %v to variable of type: %v"
	INVALID_TYPE_ANNOTATION                  = "invalid type annotation: %s"
	INVALID_RETURN_TYPE                      = "invalid return type: %s, expected: %s"
	MISSING_RETURN                           = "missing return, function should return type: %v"
	FUNCTION_REDEFINITION                    = "function redefinition, function with name: '%s' already defined here: %v"
	MULTIPLE_DEFAULT_CASES                   = "multiple default cases in switch statement"
	INVALID_CASE_TYPE                        = "unknown case type in switch statement"
	INVALID_WHILE_CONDITION                  = "expected boolean expression as condition but got: %v"
	ERROR_ARGUMENTS_NOT_FOUND                = "Last result is not an array of arguments: %v"
	INVALID_ARGUMENTS_TYPE                   = "invalid arguments type: %v"
	MAX_RECURSION_DEPTH_EXCEEDED             = "maximum recursion depth exceeded for function: %s"
	EXPECTED_BOOLEAN_EXPRESSION              = "expected boolean expression but got: %v"
)
