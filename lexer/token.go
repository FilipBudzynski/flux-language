package lexer

var Operators = map[rune]TokenTypes{
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

var KeyWords = map[string]TokenTypes{
	"int":     INT,
	"string":  STRING,
	"float":   FLOAT,
	"bool":    BOOL,
	"true":    CONST_BOOL,
	"false":   CONST_BOOL,
	"switch":  SWITCH,
	"while":   WHILE,
	"if":      IF,
	"else":    ELSE,
	"default": DEFAULT,
	"return":  RETURN,
	"print":   PRINT,
	"and":     AND,
	"or":      OR,
}

var DoubleOperators = map[string]TokenTypes{
	"<=": LESS_OR_EQUAL,
	">=": GREATER_OR_EQUAL,
	"==": EQUALS,
	"!=": NOT_EQUALS,
	"=>": CASE_ARROW,
	":=": DECLARE,
}

type Token struct {
	Value any
	Type  TokenTypes
	Pos   Position
}

func NewToken(token_type TokenTypes, position Position, value any) *Token {
	return &Token{
		Type:  token_type,
		Pos:   position,
		Value: value,
	}
}

func (b *Token) GetType() TokenTypes {
	return b.Type
}

type TokenTypes int

const (
	IDENTIFIER TokenTypes = iota
	CONST_INT
	CONST_FLOAT
	CONST_STRING
	CONST_BOOL
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
	PRINT
	RETURN
	DECLARE
	ASSIGN
	CASE_ARROW
	LEFT_BRACE
	RIGHT_BRACE
	LEFT_PARENTHESIS
	RIGHT_PARENTHESIS
	COMMA
	PARAMETER
	ARGUMENT
	ERROR
	WARNING
	STX
	ETX
	EOL
	COMMENT
	UNDEFINED
)

var tokenTypeNames = [...]string{
	"IDENTIFIER",
	"CONST_INT",
	"CONST_FLOAT",
	"CONST_STRING",
	"CONST_BOOL",
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
	"PRINT",
	"RETURN",
	"DECLARE",
	"ASSIGN",
	"CASE_ARROW",
	"LEFT_BRACE",
	"RIGHT_BRACE",
	"LEFT_PARENTHESIS",
	"RIGHT_PARENTHESIS",
	"COMMA",
	"PARAMETER",
	"ARGUMENT",
	"ERROR",
	"WARNING",
	"STX",
	"ETX",
	"EOL",
	"COMMENT",
	"UNDEFINED",
}

func (t TokenTypes) TypeName() string {
	if t < 0 || int(t) >= len(tokenTypeNames) {
		return "UNKNOWN"
	}
	return tokenTypeNames[t]
}
