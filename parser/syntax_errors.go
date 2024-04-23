package parser

const (
	SYNTAX_ERROR_FUNC_DEF_NO_PARENTHASIS = "error [%v, %v]: no parenthasis after identifier in function definition, perhaps you forgot '(' or to close the function definition with ')'"
	SYNTAX_ERROR_FUNCTION_REDEFINITION   = "error [%v, %v]: redefinition of function that already exsists at: %v, %v"
	SYNTAX_ERROR_NO_BLOCK                = "error [%v, %v]: no block defined for definiton"
	SYNTAX_ERROR_NO_IDENTIFIER           = "error [%v, %v]: no identifier in parameter group"
	SYNTAX_ERROR_NO_TYPE                 = "error [%v, %v]: no type for parameter group"
	NO_ETX_TOKEN                         = "program parsed but no ETX was found"
)
