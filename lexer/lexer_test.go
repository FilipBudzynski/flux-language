package lexer

import (
	"reflect"
	"strings"
	"testing"
)

func TestSingleTokens(t *testing.T) {
	testCases := []struct {
		expect *Token
		name   string
		input  string
	}{
		{
			name:   "identifierToken",
			input:  "myInt",
			expect: NewToken(IDENTIFIER, NewPosition(1, 1), "myInt"), // NewIdentifierToken("myInt", NewPosition(1, 1)),
		},
		{
			name:   "StringToken",
			input:  "\"String token test\"",
			expect: NewToken(CONST_STRING, NewPosition(1, 1), "String token test"), // NewStringToken("String token test", NewPosition(1, 1)),
		},
		{
			name:   "ConstIntToken",
			input:  "123",
			expect: NewToken(CONST_INT, NewPosition(1, 1), 123), // NewIntToken(123, NewPosition(1, 1)),
		},
		{
			name:   "ConstFloatToken",
			input:  "123.456",
			expect: NewToken(CONST_FLOAT, NewPosition(1, 1), 123.456), // NewFloatToken(123.456, NewPosition(1, 1)),
		},
		{
			name:   "OperatorToken",
			input:  ">=",
			expect: NewToken(GREATER_OR_EQUAL, NewPosition(1, 1), nil), // NewBaseToken(GREATER_OR_EQUAL, NewPosition(1, 1)),
		},
		{
			name:   "OperatorToken",
			input:  "=>",
			expect: NewToken(CASE_ARROW, NewPosition(1, 1), nil), // NewBaseToken(CASE_ARROW, NewPosition(1, 1)),
		},
		{
			name:   "KeyWordToken",
			input:  "return",
			expect: NewToken(RETURN, NewPosition(1, 1), nil), // NewBaseToken(RETURN, NewPosition(1, 1)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.input)
			source, _ := NewScanner(reader)
			lexer := NewLexer(*source)

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
		tokens []*Token
	}{
		{
			input: "int a = 5\nif a == 6\nwhile a >c",
			/*tokens: []Token{
				NewBaseToken(INT, NewPosition(1, 1)),
				NewIdentifierToken("a", NewPosition(1, 5)),
				NewBaseToken(ASSIGN, NewPosition(1, 7)),
				NewIntToken(5, NewPosition(1, 9)),
				NewBaseToken(EOL, NewPosition(1, 10)),
				NewBaseToken(IF, NewPosition(2, 1)),
				NewIdentifierToken("a", NewPosition(2, 4)),
				NewBaseToken(EQUALS, NewPosition(2, 6)),
				NewIntToken(6, NewPosition(2, 9)),
				NewBaseToken(EOL, NewPosition(2, 10)),
				NewBaseToken(WHILE, NewPosition(3, 1)),
				NewIdentifierToken("a", NewPosition(3, 7)),
				NewBaseToken(GREATER_THAN, NewPosition(3, 9)),
				NewIdentifierToken("c", NewPosition(3, 10)),
				NewBaseToken(ETX, NewPosition(3, 11)),
			},*/
			tokens: []*Token{
				NewToken(INT, NewPosition(1, 1), nil),
				NewToken(IDENTIFIER, NewPosition(1, 5), "a"),
				NewToken(ASSIGN, NewPosition(1, 7), nil),
				NewToken(CONST_INT, NewPosition(1, 9), 5),
				NewToken(EOL, NewPosition(1, 10), nil),
				NewToken(IF, NewPosition(2, 1), nil),
				NewToken(IDENTIFIER, NewPosition(2, 4), "a"),
				NewToken(EQUALS, NewPosition(2, 6), nil),
				NewToken(CONST_INT, NewPosition(2, 9), 6),
				NewToken(EOL, NewPosition(2, 10), nil),
				NewToken(WHILE, NewPosition(3, 1), nil),
				NewToken(IDENTIFIER, NewPosition(3, 7), "a"),
				NewToken(GREATER_THAN, NewPosition(3, 9), nil),
				NewToken(IDENTIFIER, NewPosition(3, 10), "c"),
				NewToken(ETX, NewPosition(3, 11), nil),
			},
		},
	}

	for _, tc := range testCases {
		reader := strings.NewReader(tc.input)
		source, _ := NewScanner(reader)
		lexer := NewLexer(*source)

		var tokens []*Token
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
