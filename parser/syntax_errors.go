package parser

const (
	SYNTAX_ERROR_FUNC_DEF_NO_PARENTHASIS = "error [%v, %v]: no parenthasis after identifier in function definition, perhaps you forgot '(' or to close the function definition with ')'"
	SYNTAX_ERROR_FUNCTION_REDEFINITION   = "error [%v, %v]: redefinition of function that already exsists at: %v, %v"
	SYNTAX_ERROR_NO_BLOCK                = "error [%v, %v]: no block defined for definiton"
	SYNTAX_ERROR_NO_IDENTIFIER           = "error [%v, %v]: identifier requierd here but was ommited"
	SYNTAX_ERROR_NO_VARIABLE_IDETIFIER   = "error [%v, %v]: no identifier in variable declaration"
	SYNTAX_ERROR_NO_TYPE                 = "error [%v, %v]: no type for parameter group"
	NO_ETX_TOKEN                         = "program parsed but no ETX was found"
	SYNTAX_ERROR_EXPECTED_RIGHT_BRACE    = "error [%v, %v]: expected right brace"
	SYNTAX_ERROR_UNKNOWN_STATEMENT       = "error [%v, %v]: unknown statement"
	SYNTAX_ERROR_MISSING_COLON_ASSIGN    = "error [%v, %v]: missing ':' after identifier in variable declaration"
	SYNTAX_ERROR_FUNC_CALL_NOT_CLOSED    = "error [%v, %v]: function call not closed, perhaps you forgot '('"
	ERROR_ASIGNMENT_TO_FUNCTION_CALL     = "error [%v, %v]: cannot assign value to function call"
	ERROR_MISSING_EXPRESSION             = "error [%v, %v]: missing expression"
)
