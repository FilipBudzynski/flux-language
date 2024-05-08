package parser

const (
	SYNTAX_ERROR_FUNC_DEF_NO_PARENTHASIS                   = "error [%v, %v]: no parenthasis after identifier in function definition, perhaps you forgot '(' or to close the function definition with ')'"
	SYNTAX_ERROR_FUNCTION_REDEFINITION                     = "error [%v, %v]: redefinition of function that already exsists at: %v, %v"
	SYNTAX_ERROR_NO_BLOCK                                  = "error [%v, %v]: no block defined"
	SYNTAX_ERROR_NO_IDENTIFIER                             = "error [%v, %v]: identifier requierd here but was ommited"
	SYNTAX_ERROR_NO_VARIABLE_IDETIFIER                     = "error [%v, %v]: no identifier in variable declaration"
	SYNTAX_ERROR_NO_TYPE                                   = "error [%v, %v]: no type for parameter group"
	SYNTAX_ERROR_NO_PARAMETERS_AFTER_COMMA                 = "error [%v, %v]: no parameters defined after comma"
	SYNTAX_ERROR_NO_TYPE_IN_CAST                           = "error [%v, %v]: no type in casted expression"
    ERROR_NO_ETX_TOKEN                                           = "error [%v, %v]: program parsed but no ETX was found"
	SYNTAX_ERROR_EXPECTED_RIGHT_BRACE                      = "error [%v, %v]: expected right brace"
	SYNTAX_ERROR_UNKNOWN_STATEMENT                         = "error [%v, %v]: unknown statement"
	SYNTAX_ERROR_MISSING_COLON_ASSIGN                      = "error [%v, %v]: missing ':' after identifier in variable declaration"
	SYNTAX_ERROR_FUNC_CALL_NOT_CLOSED                      = "error [%v, %v]: function call not closed, perhaps you forgot '('"
	ERROR_ASIGNMENT_TO_FUNCTION_CALL                       = "error [%v, %v]: cannot assign value to function call"
	ERROR_MISSING_EXPRESSION                               = "error [%v, %v]: missing expression after: %v"
	SYNTAX_ERROR_NO_TERM                                   = "error [%v, %v]: no term defined for expression"
	SYNTAX_ERROR_NO_EXPRESSION_IN_VARIABLE_DECLARATION     = "error [%v, %v]: no expression defined for variable declaration"
	SYNTAX_ERROR_NO_TYPE_IN_DECLARATION                    = "error [%v, %v]: no type defined for variable declaration"
	SYNTAX_ERROR_NO_RIGHT_PARENTHESIS_IN_NESTED_EXPRESSION = "error [%v, %v]: no right parenthesis in nested expression"
	SYNTAX_ERROR_NO_RETURN                                 = "error [%v, %v]: no return statement defined"
	SYNTAX_ERROR_NO_LEFT_CURLY_BRACKET_IN_SWITCH           = "error [%v, %v]: no left curly bracket in switch statement"
	SYNTAX_ERROR_NO_ARROW                                  = "error [%v, %v]: no arrow in case condition, perhaps you have ',' after the last case"
	SYNTAX_ERROR_NOT_CLOSED_SWITCH                         = "error [%v, %v]: switch statement not closed"
	SYNTAX_ERROR_NO_SWITCH_CASES                           = "error [%v, %v]: no switch cases defined"
	ERROR_MISSING_SWITCH_CASE                              = "error [%v, %v]: missing switch case, perhaps you have ',' after the last case"
	SYNTAX_ERROR_UNDEFIND_RELATION_FOR_SWITCH_CASE         = "error [%v, %v]: undefined relation for switch case"
	SYNTA_ERROR_NO_RELATION_FOR_SWITCH_CASE                = "error [%v, %v]: no relation oporator for switch case"
	SYNTAX_ERROR_NOT_VALID_TYPE_IN_FUNC                    = "error [%v, %v]: not valid type in function declaration"
	SYNTA_ERROR_NO_BLOCK_DEFINED                           = "error [%v, %v]: no block defined for the function declaration"
	SYNTAX_ERROR_NO_VARIABLE_AFTER_COMMA                   = "error [%v, %v]: no variable after comma in switch case"
	SYNTAX_ERROR_BAD_VARIABLE_DECLARATION                  = "error [%v, %v]: bad variable declaration in switch statement"
	SYNTAX_ERROR_EMPTY_BLOCK_IN_IF_STATEMENT               = "error [%v, %v]: empty block in if statement"
	SYNTAX_ERROR_EMPTY_BLOCK_IN_WHILE_STATEMENT            = "error [%v, %v]: empty block in while statement"
)

type ParserError struct {
	Message string
}

func NewParserError(message string) *ParserError {
	return &ParserError{
		Message: message,
	}
}

func (e *ParserError) Error() string {
	return e.Message
}
