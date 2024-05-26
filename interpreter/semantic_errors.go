package interpreter

import "tkom/shared"

type SemantciError struct {
	Message  string
	Position shared.Position
}

func NewSemanticError(message string, position shared.Position) *SemantciError {
	return &SemantciError{
		Message:  message,
		Position: position,
	}
}

func (err *SemantciError) Error() string {
	return err.Message
}

const (
	UNDEFINED_VARIABLE                       = "Undefined: %s"
	UNDEFINED_FUNCTION                       = "Undefined function: %s"
	REDECLARED_VARIABLE                      = "Redeclared variable: %s, variable with that name already exists"
	TYPE_MISMATCH                            = "cannot assign value of %v, to variable of type: %v"
	WRONG_NUMBER_OF_ARGUMENTS                = "Function %s expects %d arguments, got %d"
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
	INVALID_ASSIGNMENT_TYPES                 = "cannot assign value of type %v to variable of type %v"
	INVALID_TYPE_ANNOTATION                  = "Invalid type annotation: %s"
	INVALID_RETURN_TYPE                      = "Invalid return type: %s, expected %s"
)
