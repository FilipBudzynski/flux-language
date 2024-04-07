package lexer

import (
	"reflect"
	"strings"
	"testing"
)

func TestSingleTokens(t *testing.T) {
	testCases := []struct {
		expect Token
		name   string
		input  string
	}{
		{
			name:   "identifierToken",
			input:  "myInt",
			expect: NewIdentifierToken("myInt", NewPosition(1, 1)),
		},
		{
			name:   "StringToken",
			input:  "\"String token test\"",
			expect: NewStringToken("String token test", NewPosition(1, 1)),
		},
		{
			name:   "ConstIntToken",
			input:  "123",
			expect: NewIntToken(123, NewPosition(1, 1)),
		},
		{
			name:   "ConstFloatToken",
			input:  "123.456",
			expect: NewFloatToken(123.456, NewPosition(1, 1)),
		},
		{
			name:   "OperatorToken",
			input:  ">=",
			expect: NewBaseToken(GREATER_OR_EQUAL, NewPosition(1, 1)),
		},
		{
			name:   "OperatorToken",
			input:  "=>",
			expect: NewBaseToken(CASE_ARROW, NewPosition(1, 1)),
		},
		{
			name:   "KeyWordToken",
			input:  "return",
			expect: NewBaseToken(RETURN, NewPosition(1, 1)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lexer := NewLexer(strings.NewReader(tc.input))
			lexer.Consume()
			token, err := lexer.GetNextToken()
			if err != nil {
				t.Fatalf("Error while getting token: %v", err)
			}

			if !reflect.DeepEqual(token, tc.expect) {
				t.Errorf("Expected token: %v, Got: %v", tc.expect, token)
			}
		})
	}
}

func TestLexer(t *testing.T) {
	testCases := []struct {
		input  string
		tokens []Token
	}{
		{
			input: "int a = 5\nif a == 6\nwhile a >c",
			tokens: []Token{
				NewBaseToken(INT, NewPosition(1, 1)),
				NewIdentifierToken("a", NewPosition(1, 5)),
				NewBaseToken(ASSIGN, NewPosition(1, 7)),
				NewIntToken(5, NewPosition(1, 9)),
				NewBaseToken(EOL, NewPosition(2, 0)),
				NewBaseToken(IF, NewPosition(2, 1)),
				NewIdentifierToken("a", NewPosition(2, 4)),
				NewBaseToken(EQUALS, NewPosition(2, 6)),
				NewIntToken(6, NewPosition(2, 9)),
				NewBaseToken(EOL, NewPosition(3, 0)),
				NewBaseToken(WHILE, NewPosition(3, 1)),
				NewIdentifierToken("a", NewPosition(3, 7)),
				NewBaseToken(GREATER_THAN, NewPosition(3, 9)),
				NewIdentifierToken("c", NewPosition(3, 10)),
				NewBaseToken(ETX, NewPosition(3, 10)),
			},
		},
	}

	for _, tc := range testCases {
		reader := strings.NewReader(tc.input)
		lexer := NewLexer(reader)
		lexer.Consume()

		var tokens []Token
		for {
			token, err := lexer.GetNextToken()
			if err != nil {
				break
			}
			tokens = append(tokens, token)
			if token.GetType() == ETX {
				break
			}
		}

		if !reflect.DeepEqual(tokens, tc.tokens) {
			t.Errorf("Input: %s\nExpected: %+v\nGot: %+v\n", tc.input, tc.tokens, tokens)
		}
	}
}
