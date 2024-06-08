package lexer

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
	"tkom/shared"
)

const (
	identifierLimit = 500
	stringLimit     = 1000
	intLimit        = math.MaxInt
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
			expect: NewToken(IDENTIFIER, shared.NewPosition(1, 1), "myInt"),
		},
		{
			name:   "StringToken",
			input:  "\"String token test\"",
			expect: NewToken(CONST_STRING, shared.NewPosition(1, 1), "String token test"),
		},
		{
			name:   "ConstIntToken",
			input:  "123",
			expect: NewToken(CONST_INT, shared.NewPosition(1, 1), 123),
		},
		{
			name:   "ConstFloatToken",
			input:  "123.456",
			expect: NewToken(CONST_FLOAT, shared.NewPosition(1, 1), 123.456),
		},
		{
			name:   "OperatorToken",
			input:  ">=",
			expect: NewToken(GREATER_OR_EQUAL, shared.NewPosition(1, 1), nil),
		},
		{
			name:   "OperatorToken",
			input:  "=>",
			expect: NewToken(CASE_ARROW, shared.NewPosition(1, 1), nil),
		},
		{
			name:   "KeyWordToken",
			input:  "return",
			expect: NewToken(RETURN, shared.NewPosition(1, 1), nil),
		},
		{
			name:   "CommentToken",
			input:  "# This is just a comment",
			expect: NewToken(COMMENT, shared.NewPosition(1, 1), "# This is just a comment"),
		},
		{
			name:   "ConstBoolTokenFalse",
			input:  "false",
			expect: NewToken(CONST_FALSE, shared.NewPosition(1, 1), nil),
		},
		{
			name:   "ConstBoolTokenTrue",
			input:  "true",
			expect: NewToken(CONST_TRUE, shared.NewPosition(1, 1), nil),
		},
		{
			name:   "DivideToken",
			input:  "/",
			expect: NewToken(DIVIDE, shared.NewPosition(1, 1), nil),
		},
		{
			name:   "Float",
			input:  "0.24",
			expect: NewToken(CONST_FLOAT, shared.NewPosition(1, 1), 0.24),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.input)
			source, _ := NewScanner(reader)
			lexer := NewLexer(source, identifierLimit, stringLimit, intLimit)

			token := lexer.GetNextToken()

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
				NewToken(INT, shared.NewPosition(1, 1), nil),
				NewToken(IDENTIFIER, shared.NewPosition(1, 5), "a"),
				NewToken(ASSIGN, shared.NewPosition(1, 7), nil),
				NewToken(CONST_INT, shared.NewPosition(1, 9), 5),
				NewToken(IF, shared.NewPosition(2, 1), nil),
				NewToken(IDENTIFIER, shared.NewPosition(2, 4), "a"),
				NewToken(EQUALS, shared.NewPosition(2, 6), nil),
				NewToken(CONST_INT, shared.NewPosition(2, 9), 6),
				NewToken(WHILE, shared.NewPosition(3, 1), nil),
				NewToken(IDENTIFIER, shared.NewPosition(3, 7), "a"),
				NewToken(GREATER_THAN, shared.NewPosition(3, 9), nil),
				NewToken(IDENTIFIER, shared.NewPosition(3, 10), "c"),
				NewToken(ETX, shared.NewPosition(3, 11), nil),
			},
		},
		{
			input: "param1 int, param2 string, param3 bool",
			tokens: []*Token{
				NewToken(IDENTIFIER, shared.NewPosition(1, 1), "param1"),
				NewToken(INT, shared.NewPosition(1, 8), nil),
				NewToken(COMMA, shared.NewPosition(1, 11), nil),
				NewToken(IDENTIFIER, shared.NewPosition(1, 13), "param2"),
				NewToken(STRING, shared.NewPosition(1, 20), nil),
				NewToken(COMMA, shared.NewPosition(1, 26), nil),
				NewToken(IDENTIFIER, shared.NewPosition(1, 28), "param3"),
				NewToken(BOOL, shared.NewPosition(1, 35), nil),
				NewToken(ETX, shared.NewPosition(1, 39), nil),
			},
		},
	}

	for _, tc := range testCases {
		reader := strings.NewReader(tc.input)
		source, _ := NewScanner(reader)
		lexer := NewLexer(source, identifierLimit, stringLimit, intLimit)

		var tokens []*Token
		for {
			token := lexer.GetNextToken()
			tokens = append(tokens, token)
			if token.GetType() == ETX {
				break
			}
		}

		if !reflect.DeepEqual(tokens, tc.tokens) {
			fmt.Println("Expected tokens:")
			for _, token := range tc.tokens {
				fmt.Printf("%v %v %v\n", token.Position, token.Type.TypeName(), token.Value)
			}

			fmt.Println("Got tokens:")
			for _, token := range tokens {
				fmt.Printf("%v %v %v\n", token.Position, token.Type.TypeName(), token.Value)
			}

			t.Errorf("Input: %s\nExpected: %+v\nGot: %+v\n", tc.input, tc.tokens, tokens)
		}
	}
}

func TestStringNotClosed(t *testing.T) {
	input := `"unclosed string`
	expectedError := NewLexerError(STRING_NOT_CLOSED, shared.NewPosition(1, 17))

	reader := strings.NewReader(input)
	scanner, _ := NewScanner(reader)
	lexer := NewLexer(scanner, identifierLimit, stringLimit, intLimit)
	errors := []error{}
	lexer.ErrorHandler = func(err error) { errors = append(errors, err) }

	_ = lexer.GetNextToken()

	err := errors[0]
	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("Expected error: %v, but got: %v", expectedError, err)
	}
}

func TestIntValueLimitExceeded(t *testing.T) {
	input := fmt.Sprintf("int a := %d0", math.MaxInt)
	expectedError := NewLexerError(INT_CAPACITY_EXCEEDED, shared.NewPosition(1, 10))

	reader := strings.NewReader(input)
	scanner, _ := NewScanner(reader)
	lexer := NewLexer(scanner, identifierLimit, stringLimit, intLimit)
	errors := []error{}
	lexer.ErrorHandler = func(err error) { errors = append(errors, err) }

	for len(errors) == 0 {
		_ = lexer.GetNextToken()
	}

	err := errors[0]
	if err.Error() != expectedError.Error() {
		t.Errorf("Expected error: %v, but got: %v", expectedError, err)
	}
}

func TestLexerStringTokenEscaping(t *testing.T) {
	input := `"Hello\nWorld\t!\"\\"`
	expected := "Hello\nWorld\t!\"\\"

	source, _ := NewScanner(strings.NewReader(input))
	lexer := NewLexer(source, identifierLimit, stringLimit, intLimit)

	token := lexer.GetNextToken()

	if token.Type != CONST_STRING {
		t.Errorf("Expected CONST_STRING token type, got %v", token.Type)
	}

	if token.Value != expected {
		t.Errorf("Expected token value: %s, got: %s", expected, token.Value)
	}
}

func TestLexerInvalidStringTokenEscaping(t *testing.T) {
	input := `"Hello\nWorld\!\"\\"`
	expectedError := NewLexerError(INVALID_ESCAPING, shared.NewPosition(1, 15))

	source, _ := NewScanner(strings.NewReader(input))
	lexer := NewLexer(source, identifierLimit, stringLimit, intLimit)
	errors := []error{}
	lexer.ErrorHandler = func(err error) { errors = append(errors, err) }

	_ = lexer.GetNextToken()

	err := errors[0]
	if err == nil {
		t.Error("Expected error for invalid escaping, got nil")
	} else {
		if err.Error() != expectedError.Error() {
			t.Errorf("Expected error message: %s, got: %s", expectedError, err.Error())
		}
	}
}

func TestStringValueLimitExceeded(t *testing.T) {
	stringLimit := 5
	input := `"WAAAAA"`
	expectedError := NewLexerError(STRING_CAPACITY_EXCEEDED, shared.NewPosition(1, 7))

	reader := strings.NewReader(input)
	scanner, _ := NewScanner(reader)
	lexer := NewLexer(scanner, identifierLimit, stringLimit, intLimit)
	var testError error
	lexer.ErrorHandler = func(err error) { testError = err }

	for testError == nil {
		_ = lexer.GetNextToken()
	}

	if testError.Error() != expectedError.Error() {
		t.Errorf("Expected error: %v, but got: %v", expectedError, testError)
	}
}

func TestLexerErrorHandling(t *testing.T) {
	input := "abc 34 karma"
	source, _ := NewScanner(strings.NewReader(input))

	expectedError := NewLexerError(INT_CAPACITY_EXCEEDED, shared.NewPosition(1, 5))
	expectedTokens := []*Token{
		NewToken(IDENTIFIER, shared.NewPosition(1, 1), "abc"),
		NewToken(ETX, shared.NewPosition(1, 6), nil),
	}

	lexer := NewLexer(source, 10, 10, 10)
	var externalErrors []error
	lexer.ErrorHandler = func(err error) {
		externalErrors = append(externalErrors, err)
	}
	tokens := []*Token{}

	for {
		token := lexer.GetNextToken()
		tokens = append(tokens, token)
		if token.Type == ETX {
			break
		}
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Expected: %+v\nGot: %+v\n", tokens, expectedTokens)
	}

	if len(externalErrors) == 0 {
		t.Error("Expected lexer to stop after encountering an error, but it continued parsing")
	}

	if len(externalErrors) == 1 && externalErrors[0].Error() != expectedError.Error() {
		t.Errorf("Expected error: %s, but got: %s", expectedError, externalErrors[0].Error())
	}
}
