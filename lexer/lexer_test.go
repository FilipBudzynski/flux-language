package lexer

import (
	"fmt"
	"math"
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
			expect: NewToken(IDENTIFIER, NewPosition(1, 1), "myInt"),
		},
		{
			name:   "StringToken",
			input:  "\"String token test\"",
			expect: NewToken(CONST_STRING, NewPosition(1, 1), "String token test"),
		},
		{
			name:   "ConstIntToken",
			input:  "123",
			expect: NewToken(CONST_INT, NewPosition(1, 1), 123),
		},
		{
			name:   "ConstFloatToken",
			input:  "123.456",
			expect: NewToken(CONST_FLOAT, NewPosition(1, 1), 123.456),
		},
		{
			name:   "OperatorToken",
			input:  ">=",
			expect: NewToken(GREATER_OR_EQUAL, NewPosition(1, 1), nil),
		},
		{
			name:   "OperatorToken",
			input:  "=>",
			expect: NewToken(CASE_ARROW, NewPosition(1, 1), nil),
		},
		{
			name:   "KeyWordToken",
			input:  "return",
			expect: NewToken(RETURN, NewPosition(1, 1), nil),
		},
		{
			name:   "CommentToken",
			input:  "# This is just a comment",
			expect: NewToken(COMMENT, NewPosition(1, 1), nil),
		},
		{
			name:   "ConstBoolTokenFalse",
			input:  "false",
			expect: NewToken(CONST_BOOL, NewPosition(1, 1), false),
		},
		{
			name:   "ConstBoolTokenTrue",
			input:  "true",
			expect: NewToken(CONST_BOOL, NewPosition(1, 1), true),
		},
		{
			name:   "DivideToken",
			input:  "/",
			expect: NewToken(DIVIDE, NewPosition(1, 1), nil),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.input)
			source, _ := NewScanner(reader)
			lexer := NewLexer(source)

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

func TestLexerCodeExample(t *testing.T) {
	testCases := []struct {
		input  string
		tokens []*Token
	}{
		{
			input: "int a = 5\nif a == 6\nwhile a >c",
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
		lexer := NewLexer(source)

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

func TestStringNotClosed(t *testing.T) {
	input := `"unclosed string`
	expectedError := NewLexerError(STRING_NOT_CLOSED, NewPosition(1, 17))

	reader := strings.NewReader(input)
	scanner, _ := NewScanner(reader)
	lexer := NewLexer(scanner)

	_, err := lexer.GetNextToken()

	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("Expected error: %v, but got: %v", expectedError, err)
	}
}

func TestIntValueLimitExceeded(t *testing.T) {
	input := fmt.Sprintf("int a := %d0", math.MaxInt)
	expectedError := NewLexerError(INT_CAPACITY_EXCEEDED, NewPosition(1, 10))

	reader := strings.NewReader(input)
	scanner, _ := NewScanner(reader)
	lexer := NewLexer(scanner)

	var err error
	for err == nil {
		_, err = lexer.GetNextToken()
	}

	if err.Error() != expectedError.Error() {
		t.Errorf("Expected error: %v, but got: %v", expectedError, err)
	}
}

func TestLexerStringTokenEscaping(t *testing.T) {
	input := `"Hello\nWorld\t!\"\\"`
	expected := "Hello\nWorld\t!\"\\"

	source, _ := NewScanner(strings.NewReader(input))
	lexer := NewLexer(source)

	token, err := lexer.GetNextToken()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if token.Type != CONST_STRING {
		t.Errorf("Expected CONST_STRING token type, got %v", token.Type)
	}

	if token.Value != expected {
		t.Errorf("Expected token value: %s, got: %s", expected, token.Value)
	}
}

func TestLexerInvalidStringTokenEscaping(t *testing.T) {
	input := `"Hello\nWorld\!\"\\"`
	expectedError := NewLexerError(INVALID_ESCAPING, NewPosition(1, 15))

	source, _ := NewScanner(strings.NewReader(input))
	lexer := NewLexer(source)

	_, err := lexer.GetNextToken()
	if err == nil {
		t.Error("Expected error for invalid escaping, got nil")
	} else {
		if err.Error() != expectedError.Error() {
			t.Errorf("Expected error message: %s, got: %s", expectedError, err.Error())
		}
	}
}
