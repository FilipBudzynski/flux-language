package lexer

import (
	"fmt"
)

type ErrorCode int

const (
	INT_CAPACITY_EXCEEDED ErrorCode = iota
	FLOAT_CAPACITY_EXCEEDED
	IDENTIFIER_CAPACITY_EXCEEDED
	STRING_CAPACITY_EXCEEDED
	STRING_NOT_CLOSED
	INVALID_ESCAPING
	NONE_TOKEN_MATCH
)

var errorMessage = map[ErrorCode]string{
	INT_CAPACITY_EXCEEDED:        "error [%d, %d] Int value limit Exceeded",
	FLOAT_CAPACITY_EXCEEDED:      "error [%d, %d], Float decimal value limit Exceeded",
	IDENTIFIER_CAPACITY_EXCEEDED: "error [%d, %d] Identifier capacity exceeded",
	STRING_CAPACITY_EXCEEDED:     "error [%d, %d] String capacity exceeded",
	STRING_NOT_CLOSED:            "error [%d, %d] String not closed, perhaps you forgot \"",
	INVALID_ESCAPING:             "error [%d, %d] Invalid syntax escaping",
	NONE_TOKEN_MATCH:             "error [%d, %d] None token match found for the source",
}

type LexerError struct {
	Code     ErrorCode
	Position Position
}

func (e *LexerError) Error() string {
	msg, ok := errorMessage[e.Code]
	if !ok {
		return "unknown error"
	}
	return fmt.Sprintf(msg, e.Position.Line, e.Position.Column)
}

func NewLexerError(code ErrorCode, position Position) *LexerError {
	return &LexerError{
		Code:     code,
		Position: position,
	}
}
