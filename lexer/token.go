package lexer

import "fmt"

type Token interface {
	IsEqual(token Token) bool
	GetType() TokenTypes
	GetName() string
	ShowDetails()
}

type baseToken struct {
	Type TokenTypes
	Pos  Position
}

func NewBaseToken(token_type TokenTypes, pos Position) *baseToken {
	return &baseToken{
		Type: token_type,
		Pos:  pos,
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
	NO_EQUALS
	GREATER_THAN
	LESS_THAN
	GREATER_THAN_OR_EQUAL
	LESS_THAN_OREQUAL
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
	"NO_EQUALS",
	"GREATER_THAN",
	"LESS_THAN",
	"GREATER_THAN_OR_EQUAL",
	"LESS_THAN_OREQUAL",
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
