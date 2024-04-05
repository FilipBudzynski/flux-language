package lexer

import "fmt"

var SingleChar = map[rune]TokenTypes{
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
}

var KeyWords = map[string]TokenTypes{
	"int":     INT,
	"string":  STRING,
	"bool":    BOOL,
	"float":   FLOAT,
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
	"->": CASE_ARROW,
	":=": DECLARE,
}

type Token interface {
	IsEqual(token Token) bool
	GetType() TokenTypes
	GetName() string
	ShowDetails()
	SetPosition(Position)
}

type baseToken struct {
	Type TokenTypes
	Pos  Position
}

func NewBaseToken(token_type TokenTypes) *baseToken {
	return &baseToken{
		Type: token_type,
	}
}

func (b *baseToken) IsEqual(token Token) bool {
	if other, ok := token.(*baseToken); ok {
		return b.Type == other.Type && b.Pos == other.Pos
	}
	return false
}

func (b *baseToken) ShowDetails() {
	fmt.Printf("%v, %v\n", b.Pos, b.Type.GetName())
}

func (b *baseToken) GetType() TokenTypes {
	return b.Type
}

func (b *baseToken) GetName() string {
	return b.Type.GetName()
}

func (b *baseToken) SetPosition(position Position) {
	b.Pos = position
}

type TokenTypes int

const (
	// identifier
	IDENTIFIER TokenTypes = iota
	// data types
	CONSNT_INT
	CONST_FLOAT
	CONST_STRING
	CONST_BOOL
	// types annotations
	INT
	FLOAT
	STRING
	BOOL
	// arythmetic operators
	PLUS
	MINUS
	MULTIPLY
	DIVIDE
	// relational operators
	EQUALS
	NOT_EQUALS
	GREATER_THAN
	LESS_THAN
	GREATER_OR_EQUAL
	LESS_OR_EQUAL
	// logic operators
	AND
	OR
	NEGATE
	// keywords
	IF
	ELSE
	WHILE
	SWITCH
	DEFAULT
	AS
	PRINT
	RETURN
	// special symbols
	DECLARE
	ASSIGN
	CASE_ARROW
	LEFT_BRACE
	RIGHT_BRACE
	LEFT_PARENTHESIS
	RIGHT_PARENTHESIS
	COMMA
	// function parameters
	PARAMETER
	ARGUMENT
	// errors and warning
	ERROR
	WARNING
	// other
	STX
	ETX
	EOL
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
	"UNDEFINED",
}

func (t TokenTypes) GetName() string {
	if t < 0 || int(t) >= len(tokenTypeNames) {
		return "UNKNOWN"
	}
	return tokenTypeNames[t]
}
