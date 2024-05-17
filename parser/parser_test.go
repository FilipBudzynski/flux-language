package parser

import (
	"reflect"
	"strings"
	"testing"
	. "tkom/ast"
	"tkom/lexer"
)

// Helper function to create a lexer from a given input string
func createLexer(input string) *lexer.Lexer {
	source, _ := lexer.NewScanner(strings.NewReader(input))
	return lexer.NewLexer(source, 1000, 1000, 1000)
}

// Helper function to create parser
func createParser(t *testing.T, input string) *Parser {
	lex := createLexer(input)
	errorHandler := func(err error) {
		t.Errorf("Parse Identifier error: %v", err)
	}
	return NewParser(lex, errorHandler)
}

func TestParseParameterGroup(t *testing.T) {
	input := "param1, param2, param3 string"
	parser := createParser(t, input)

	params := parser.parseParameterGroup()

	if len(params) != 3 {
		t.Errorf("Expected 3 parameters, got %d", len(params))
	}

	expected := []*Variable{
		NewVariable(STRING, "param1", nil, lexer.NewPosition(1, 1)),
		NewVariable(STRING, "param2", nil, lexer.NewPosition(1, 9)),
		NewVariable(STRING, "param3", nil, lexer.NewPosition(1, 17)),
	}

	for i, param := range params {
		if !reflect.DeepEqual(param, expected[i]) {
			t.Errorf("Expected %d parameters, got %d", len(expected), len(params))
			t.Errorf("expected variable: %v, got: %v", expected[i], param)
		}
		if param.Type != expected[i].Type {
			t.Errorf("Expected parameter type %v, got %v", expected[i].Type, param.Type)
		}
	}
}

func TestParseParameters(t *testing.T) {
	input := "param1 int, param2 string, param3 bool"
	expected := []*Variable{
		NewVariable(INT, "param1", nil, lexer.NewPosition(1, 1)),
		NewVariable(STRING, "param2", nil, lexer.NewPosition(1, 13)),
		NewVariable(BOOL, "param3", nil, lexer.NewPosition(1, 28)),
	}
	parser := createParser(t, input)

	params := parser.parseParameters()

	for _, param := range params {
		t.Log(param)
	}

	if len(params) != len(expected) {
		t.Errorf("Expected %d parameters, got %d", len(expected), len(params))
		return
	}

	for i, param := range params {
		if !param.Equals(*expected[i]) {
			t.Errorf("Expected parameter name %s, got %s", expected[i].Name, param.Name)
		}
		if param.Type != expected[i].Type {
			t.Errorf("Expected parameter type %v, got %v", expected[i].Type, param.Type)
		}
	}
}

func TestParseIdentifier(t *testing.T) {
	input := "identifier1"
	expected := NewIdentifier("identifier1", lexer.NewPosition(1, 1))
	lex := createLexer(input)
	errorHandler := func(err error) {
		t.Errorf("Parse Identifier error: %v", err)
	}
	parser := NewParser(lex, errorHandler)
	identifier := parser.parseAssignment()

	t.Log(identifier)
	if identifier != expected {
		t.Errorf("expected: %v, got: %v ", identifier, expected)
	}
}

func TestParseFunctionDefinitions(t *testing.T) {
	tests := []struct {
		expected *FunDef
		input    string
	}{
		{
			input: "myFunction(a int, b string) { }",
			expected: NewFunctionDefinition(
				"myFunction",
				[]*Variable{
					NewVariable(INT, "a", nil, lexer.NewPosition(1, 12)),
					NewVariable(STRING, "b", nil, lexer.NewPosition(1, 19)),
				},
				VOID,
				NewBlock([]Statement{}),
				lexer.NewPosition(1, 1),
			),
		},
		{
			input: "secondFunc(a, b string) string { return a }",
			expected: NewFunctionDefinition(
				"secondFunc",
				[]*Variable{
					NewVariable(STRING, "a", nil, lexer.NewPosition(1, 12)),
					NewVariable(STRING, "b", nil, lexer.NewPosition(1, 15)),
				},
				STRING,
				NewBlock([]Statement{NewReturnStatement(NewIdentifier("a", lexer.NewPosition(1, 41)))}),
				lexer.NewPosition(1, 1),
			),
		},
	}

	for _, tt := range tests {
		parser := createParser(t, tt.input)
		functionDefinition := parser.parseFunDef()

		if !functionDefinition.Equals(tt.expected) {
			t.Errorf("function definitions are not equal, expected: %v, got: %v", tt.expected, functionDefinition)
		}
	}
}

func TestParseExpressionIdentifierOnly(t *testing.T) {
	input := "identifier1"
	expected := NewIdentifier("identifier1", lexer.NewPosition(1, 1))

	lex := createLexer(input)
	errorHandler := func(err error) {
		t.Errorf("Parse Identifier error: %v", err)
	}
	parser := NewParser(lex, errorHandler)

	expression := parser.parseExpression()

	if expr, ok := expression.(Identifier); !ok {
		t.Errorf("expressions are not equal, expected: %v, got: %v", expected, expr)
	}
}

func TestParseExpressions(t *testing.T) {
	tests := []struct {
		expected Expression
		name     string
		input    string
	}{
		{
			name:  "GreaterOrEqual",
			input: "a >= 2",
			expected: NewGreaterOrEqualExpression(
				NewIdentifier("a", lexer.NewPosition(1, 1)),
				NewIntExpression(2, lexer.NewPosition(1, 6)),
				lexer.NewPosition(1, 3),
			),
		},
		{
			name:  "LessThan",
			input: "a < 2",
			expected: NewLessThanExpression(
				NewIdentifier("a", lexer.NewPosition(1, 1)),
				NewIntExpression(2, lexer.NewPosition(1, 5)),
				lexer.NewPosition(1, 3),
			),
		},
		{
			name:  "LessOrEqual",
			input: "a <= 2",
			expected: NewLessOrEqualExpression(
				NewIdentifier("a", lexer.NewPosition(1, 1)),
				NewIntExpression(2, lexer.NewPosition(1, 6)),
				lexer.NewPosition(1, 3),
			),
		},
		{
			name:  "GreaterThan",
			input: "a > 2",
			expected: NewGreaterThanExpression(
				NewIdentifier("a", lexer.NewPosition(1, 1)),
				NewIntExpression(2, lexer.NewPosition(1, 5)),
				lexer.NewPosition(1, 3),
			),
		},
		{
			name:  "Equal",
			input: "a == 2",
			expected: NewEqualsExpression(
				NewIdentifier("a", lexer.NewPosition(1, 1)),
				NewIntExpression(2, lexer.NewPosition(1, 6)),
				lexer.NewPosition(1, 3),
			),
		},
		{
			name:  "NotEqual",
			input: "a != 2",
			expected: NewNotEqualsExpression(
				NewIdentifier("a", lexer.NewPosition(1, 1)),
				NewIntExpression(2, lexer.NewPosition(1, 6)),
				lexer.NewPosition(1, 3),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex := createLexer(tt.input)
			errorHandler := func(err error) {
				t.Errorf("Parse Identifier error: %v", err)
			}
			parser := NewParser(lex, errorHandler)

			expression := parser.parseExpression()

			// type assertion
			if reflect.TypeOf(expression) != reflect.TypeOf(tt.expected) {
				t.Errorf("Parsed expression is not of expected type. got: %v, want: %v",
					reflect.TypeOf(expression), reflect.TypeOf(tt.expected))
				return
			}

			if !expression.Equals(tt.expected) {
				t.Errorf("Expressions are not equal, expected: %v, got: %v", tt.expected, expression)
				t.Errorf("Actual type: %v", reflect.TypeOf(expression))
				t.Errorf("Expected type: %v", reflect.TypeOf(tt.expected))
			}
		})
	}
}

func TestParseVariableDeclaration(t *testing.T) {
	tests := []struct {
		expected *Variable
		input    string
	}{
		{
			input: "int a := 2 - 3",
			expected: &Variable{
				Type: INT,
				Name: "a",
				Value: NewSubstractExpression(
					NewIntExpression(2, lexer.NewPosition(1, 10)),
					NewIntExpression(3, lexer.NewPosition(1, 14)),
					lexer.NewPosition(1, 12),
				),
				Position: lexer.NewPosition(1, 5),
			},
		},
		{
			input: "int a := 2 + 2 + 2",
			expected: &Variable{
				Type: INT,
				Name: "a",
				Value: NewSumExpression(
					NewSumExpression(
						NewIntExpression(2, lexer.NewPosition(1, 10)),
						NewIntExpression(2, lexer.NewPosition(1, 14)),
						lexer.NewPosition(1, 12)),
					NewIntExpression(2, lexer.NewPosition(1, 18)),
					lexer.NewPosition(1, 16),
				),
				Position: lexer.NewPosition(1, 5),
			},
		},
	}

	for _, test := range tests {
		parser := createParser(t, test.input)
		statement := parser.parseVariableDeclaration()

		if !reflect.DeepEqual(test.expected, statement) {
			t.Errorf("Input: %s\nExpressions are not equal, expected: %v, got: %v", test.input, test.expected, statement)
		}
	}
}

func TestParseNegateExpression(t *testing.T) {
	tests := []struct {
		expected Expression
		input    string
	}{
		{
			input:    "-22",
			expected: NewNegateExpression(NewIntExpression(22, lexer.NewPosition(1, 2)), lexer.NewPosition(1, 1)),
		},
		{
			input:    "!a",
			expected: NewNegateExpression(NewIdentifier("a", lexer.NewPosition(1, 2)), lexer.NewPosition(1, 1)),
		},
	}

	for _, test := range tests {
		parser := createParser(t, test.input)
		statement := parser.parseExpression()

		if statement, ok := statement.(*NegateExpression); !ok {
			t.Errorf("Parsed statement is not of type Variable")
			t.Errorf("Actual type: %v", reflect.TypeOf(statement))
			t.Errorf("Expected type: %v", reflect.TypeOf(test.expected))
			return
		}

		if !reflect.DeepEqual(test.expected, statement) {
			t.Errorf("Expressions are not equal, expected: %v, got: %v", test.expected, statement)
		}
	}
}

func TestParseVariableDeclarations(t *testing.T) {
	testCases := []struct {
		expected *Variable
		input    string
	}{
		{
			input:    "bool c := true",
			expected: NewVariable(BOOL, "c", NewBoolExpression(true, lexer.NewPosition(1, 11)), lexer.NewPosition(1, 6)),
		},
		{
			input:    "int a := 2",
			expected: NewVariable(INT, "a", NewIntExpression(2, lexer.NewPosition(1, 10)), lexer.NewPosition(1, 5)),
		},
		{
			input:    "float b := 5.14",
			expected: NewVariable(FLOAT, "b", NewFloatExpression(5.14, lexer.NewPosition(1, 12)), lexer.NewPosition(1, 7)),
		},
	}

	for _, tc := range testCases {
		parser := createParser(t, tc.input)
		statement := parser.parseVariableDeclaration()

		if !reflect.DeepEqual(tc.expected, statement) {
			t.Errorf("Expressions are not equal, expected: %v, got: %v", tc.expected, statement)
		}
	}
}

func TestParseFloatExpression(t *testing.T) {
	input := "3.14"

	expected := NewFloatExpression(3.14, lexer.NewPosition(1, 1))
	parser := createParser(t, input)

	statement := parser.parseExpression()

	if statement, ok := statement.(*FloatExpression); !ok {
		t.Errorf("Parsed statement is not of type Float")
		t.Errorf("Actual type: %v", reflect.TypeOf(statement))
		t.Errorf("Expected type: %v", reflect.TypeOf(expected))
		return
	}

	if !reflect.DeepEqual(expected, statement) {
		t.Errorf("Expressions are not equal, expected: %v, got: %v", expected, statement)
	}
}

func TestParseStringExpression(t *testing.T) {
	input := `"This is test a string"`

	expected := NewStringExpression("This is test a string", lexer.NewPosition(1, 1))
	parser := createParser(t, input)

	statement := parser.parseExpression()

	if statement, ok := statement.(*StringExpression); !ok {
		t.Errorf("Parsed statement is not of type Float")
		t.Errorf("Actual type: %v", reflect.TypeOf(statement))
		t.Errorf("Expected type: %v", reflect.TypeOf(expected))
		return
	}

	if !reflect.DeepEqual(expected, statement) {
		t.Errorf("Expressions are not equal, expected: %v, got: %v", expected, statement)
	}
}

func TestParseParenthesisExpression(t *testing.T) {
	input := "3 * (1 + (4 - 1))"

	expected := NewMultiplyExpression(
		NewIntExpression(3, lexer.NewPosition(1, 1)),
		NewSumExpression(
			NewIntExpression(1, lexer.NewPosition(1, 6)),
			NewSubstractExpression(
				NewIntExpression(4, lexer.NewPosition(1, 11)),
				NewIntExpression(1, lexer.NewPosition(1, 15)),
				lexer.NewPosition(1, 13)),
			lexer.NewPosition(1, 8)),
		lexer.NewPosition(1, 3))
	parser := createParser(t, input)

	statement := parser.parseExpression()

	if statement, ok := statement.(*MultiplyExpression); !ok {
		t.Errorf("Parsed statement is not of type %v", reflect.TypeOf(expected))
		t.Errorf("Actual type: %v", reflect.TypeOf(statement))
		t.Errorf("Expected type: %v", reflect.TypeOf(expected))
		return
	}

	if !reflect.DeepEqual(expected, statement) {
		t.Errorf("Expressions are not equal, expected: %v, got: %v", expected, statement)
	}
}

func TestParseOrAndExpression(t *testing.T) {
	input := "a or b and c"

	idA := NewIdentifier("a", lexer.NewPosition(1, 1))
	idB := NewIdentifier("b", lexer.NewPosition(1, 6))
	idC := NewIdentifier("c", lexer.NewPosition(1, 12))
	expected := NewOrExpression(idA, NewAndExpression(idB, idC, lexer.NewPosition(1, 8)), lexer.NewPosition(1, 3))
	parser := createParser(t, input)

	statement := parser.parseExpression()

	if statement, ok := statement.(*OrExpression); !ok {
		t.Errorf("Parsed statement is not of type OrExpression")
		t.Errorf("Actual type: %v", reflect.TypeOf(statement))
		t.Errorf("Expected type: %v", reflect.TypeOf(expected))
	}

	if st, ok := statement.(*OrExpression); ok {
		if !expected.Equals(st) {
			t.Errorf("And with Or expressions not parsed correctly, expected: %v, got: %v", expected, statement)
		}
	}
}

func TestParseIfStatement(t *testing.T) {
	input := "if x == 10 { y = 20 }"

	expected := NewIfStatement(
		NewEqualsExpression(NewIdentifier("x", lexer.NewPosition(1, 4)), NewIntExpression(10, lexer.NewPosition(1, 9)), lexer.NewPosition(1, 6)),
		NewBlock([]Statement{NewAssignment(NewIdentifier("y", lexer.NewPosition(1, 14)), NewIntExpression(20, lexer.NewPosition(1, 18)))}),
		nil,
	)
	parser := createParser(t, input)

	statement := parser.parseConditionalStatement()

	if !reflect.DeepEqual(expected, statement) {
		t.Errorf("If statement not parsed correctly, expected: %v, got: %v", expected, statement)
	}
}

func TestParseIfStatementWithElse(t *testing.T) {
	input := "if x == 10 { y = 20 } else { y = 15 }"

	expected := NewIfStatement(
		NewEqualsExpression(NewIdentifier("x", lexer.NewPosition(1, 4)), NewIntExpression(10, lexer.NewPosition(1, 9)), lexer.NewPosition(1, 6)),
		NewBlock([]Statement{NewAssignment(NewIdentifier("y", lexer.NewPosition(1, 14)), NewIntExpression(20, lexer.NewPosition(1, 18)))}),
		NewBlock([]Statement{NewAssignment(NewIdentifier("y", lexer.NewPosition(1, 30)), NewIntExpression(15, lexer.NewPosition(1, 34)))}),
	)
	parser := createParser(t, input)
	statement := parser.parseConditionalStatement()

	if !reflect.DeepEqual(expected, statement) {
		t.Errorf("If statement not parsed correctly, expected: %v, got: %v", expected, statement)
	}
}

func TestParseWhileStatement(t *testing.T) {
	input := "while x < 10 { y = y + 1 }"

	idX := NewIdentifier("x", lexer.NewPosition(1, 7))
	idY1 := NewIdentifier("y", lexer.NewPosition(1, 16))
	idY2 := NewIdentifier("y", lexer.NewPosition(1, 20))
	expected := NewWhileStatement(
		NewLessThanExpression(idX, NewIntExpression(10, lexer.NewPosition(1, 11)), lexer.NewPosition(1, 9)),
		NewBlock([]Statement{NewAssignment(idY1, NewSumExpression(idY2, NewIntExpression(1, lexer.NewPosition(1, 24)), lexer.NewPosition(1, 22)))}),
	)
	parser := createParser(t, input)

	statement := parser.parseWhileStatement()

	if !statement.Equals(expected) {
		t.Errorf("While statement not parsed correctly, expected: %v, got: %v", expected, statement)
	}
}

func TestSwitchStatement(t *testing.T) {
	input := `switch int a := 2 {
        a > 2 and a < 10 => "Kasia",
        a >= 10          => "Asia"
    }`

	idA := NewIdentifier("a", lexer.NewPosition(1, 8))
	variable := NewVariable(INT, "a", NewIntExpression(2, lexer.NewPosition(1, 17)), lexer.NewPosition(1, 8))
	variables := []*Variable{variable}
	cases := []Case{
		NewSwitchCase(NewAndExpression(
			NewGreaterThanExpression(idA, NewIntExpression(2, lexer.NewPosition(2, 12)), lexer.NewPosition(2, 8)),
			NewLessThanExpression(idA, NewIntExpression(10, lexer.NewPosition(2, 22)), lexer.NewPosition(2, 15)), lexer.NewPosition(2, 11)),
			NewStringExpression("Kasia", lexer.NewPosition(2, 29))),
		NewSwitchCase(NewGreaterOrEqualExpression(idA, NewIntExpression(10, lexer.NewPosition(3, 14)), lexer.NewPosition(3, 8)),
			NewStringExpression("Asia", lexer.NewPosition(3, 29))),
	}
	expected := NewSwitchStatement(variables, nil, cases)
	parser := createParser(t, input)

	statement := parser.parseSwitchStatement()

	if !expected.Equals(*statement) {
		t.Errorf("Switch statement not parsed correctly, expected: %v, got: %v", expected, statement)
	}
}

func TestSwitchStatements(t *testing.T) {
	tests := []struct {
		expected      *SwitchStatement
		expectedError *ParserError
		name          string
		input         string
	}{
		{
			name: "SwitchStatement",
			input: `switch int a := 2 {
                a > 2 and a < 10 => "Kasia",
                a >= 10          => "Asia"
            }`,
			expected: NewSwitchStatement(
				[]*Variable{
					NewVariable(INT, "a", NewIntExpression(2, lexer.NewPosition(1, 17)), lexer.NewPosition(1, 8)),
				},
				nil,
				[]Case{
					NewSwitchCase(
						NewAndExpression(
							NewGreaterThanExpression(NewIdentifier("a", lexer.NewPosition(1, 8)), NewIntExpression(2, lexer.NewPosition(2, 12)), lexer.NewPosition(2, 8)),
							NewLessThanExpression(NewIdentifier("a", lexer.NewPosition(1, 8)), NewIntExpression(10, lexer.NewPosition(2, 22)), lexer.NewPosition(2, 15)),
							lexer.NewPosition(2, 11),
						),
						NewStringExpression("Kasia", lexer.NewPosition(2, 29)),
					),
					NewSwitchCase(
						NewGreaterOrEqualExpression(NewIdentifier("a", lexer.NewPosition(1, 8)), NewIntExpression(10, lexer.NewPosition(3, 14)), lexer.NewPosition(3, 8)),
						NewStringExpression("Asia", lexer.NewPosition(3, 29)),
					),
				},
			),
			expectedError: nil,
		},
		{
			name: "SwitchStatementWithDefault",
			input: `switch {
                a > 2   => fun1(),
                a <= 10 => fun2(),
                default => fun3()
            }`,
			expected: NewSwitchStatement(
				nil,
				nil,
				[]Case{
					NewSwitchCase(
						NewGreaterThanExpression(NewIdentifier("a", lexer.NewPosition(2, 9)), NewIntExpression(2, lexer.NewPosition(2, 13)), lexer.NewPosition(2, 11)),
						NewFunctionCall("fun1", lexer.NewPosition(2, 20), nil),
					),
					NewSwitchCase(
						NewLessOrEqualExpression(NewIdentifier("a", lexer.NewPosition(3, 9)), NewIntExpression(2, lexer.NewPosition(3, 14)), lexer.NewPosition(3, 11)),
						NewFunctionCall("fun2", lexer.NewPosition(2, 21), nil),
					),
					NewDefaultCase(
						NewFunctionCall("fun3", lexer.NewPosition(3, 20), nil),
					),
				},
			),
			expectedError: nil,
		},
		{
			name: "SwitchStatementWithBlock",
			input: `switch {
                a > 2   => { return 20 },
            }`,
			expected: NewSwitchStatement(
				nil,
				nil,
				[]Case{
					NewSwitchCase(
						NewGreaterThanExpression(NewIdentifier("a", lexer.NewPosition(2, 9)), NewIntExpression(2, lexer.NewPosition(2, 13)), lexer.NewPosition(2, 11)),
						NewBlock([]Statement{NewReturnStatement(NewIntExpression(20, lexer.NewPosition(2, 29)))}),
					),
				},
			),
			expectedError: nil,
		},
		{
			name: "ParseSwitchError",
			input: `switch {
                >30 => 2
                }`,
			expected: nil,
			expectedError: &ParserError{
				Message: "error [2, 5]: missing or bad switch case condition",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := createParser(t, tt.input)
			errors := []error{}
			errorHandler := func(err error) { errors = append(errors, err) }
			parser.ErrorHandler = errorHandler

			var statement *SwitchStatement

			defer func() {
				if r := recover(); r != nil {
					parser.ErrorHandler(r.(error))
				}
			}()

			statement = parser.parseSwitchStatement()

			if tt.expectedError != nil {
				if len(errors) == 0 {
					t.Errorf("expected error but got none")
				} else if errors[0].Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, but got %v", tt.expectedError.Error(), errors[0].Error())
				}
			} else {
				if len(errors) > 0 {
					t.Errorf("unexpected error: %v", errors[0])
				}
				if !tt.expected.Equals(*statement) {
					t.Errorf("Switch statement not parsed correctly, expected: %v, got: %v", tt.expected, statement)
				}
			}
		})
	}
}

func TestParseProgram(t *testing.T) {
	input := `main() {
    int a := 10
    int b := second()
    a = a - b

    if a > b {
        a = 0
    } 
    else {
        b = 0
    }
}`

	statements := []Statement{
		NewVariable(INT, "a", NewIntExpression(10, lexer.NewPosition(2, 14)), lexer.NewPosition(2, 9)),
		NewVariable(
			INT,
			"b",
			NewFunctionCall("second", lexer.NewPosition(3, 14), []Expression{}),
			lexer.NewPosition(3, 9),
		),
		NewAssignment(
			NewIdentifier("a", lexer.NewPosition(4, 5)),
			NewSubstractExpression(
				NewIdentifier("a", lexer.NewPosition(4, 9)),
				NewIdentifier("b", lexer.NewPosition(4, 13)),
				lexer.NewPosition(4, 11),
			),
		),
		NewIfStatement(
			NewGreaterThanExpression(
				NewIdentifier("a", lexer.NewPosition(6, 8)),
				NewIdentifier("b", lexer.NewPosition(6, 12)),
				lexer.NewPosition(6, 10),
			),
			NewBlock([]Statement{NewAssignment(NewIdentifier("a", lexer.NewPosition(7, 9)), NewIntExpression(0, lexer.NewPosition(7, 13)))}),
			NewBlock([]Statement{NewAssignment(NewIdentifier("b", lexer.NewPosition(10, 9)), NewIntExpression(0, lexer.NewPosition(10, 13)))}),
		),
	}

	funDefs := map[string]*FunDef{
		"main": NewFunctionDefinition("main", nil, VOID, NewBlock(statements), lexer.NewPosition(1, 1)),
	}

	expected := NewProgram(funDefs)
	parser := createParser(t, input)
	program := parser.ParseProgram()

	// if reflect.DeepEqual(expected, program) {
	if !program.Equals(expected) {
		t.Errorf("Program not parsed correctly, expected: %v, got: %v", expected, program)
	}
}

func TestParseProgramInt(t *testing.T) {
	input := `main() int {
    int a := 10
}`

	statements := []Statement{
		NewVariable(INT, "a", NewIntExpression(10, lexer.NewPosition(2, 14)), lexer.NewPosition(2, 9)),
		NewVariable(
			INT,
			"b",
			NewFunctionCall("second", lexer.NewPosition(3, 14), []Expression{}),
			lexer.NewPosition(3, 9),
		),
	}

	funDefs := map[string]*FunDef{
		"main": NewFunctionDefinition("main", nil, INT, NewBlock(statements), lexer.NewPosition(1, 1)),
	}

	expected := NewProgram(funDefs)
	parser := createParser(t, input)
	program := parser.ParseProgram()

	if !program.Equals(expected) {
		t.Errorf("Program not parsed correctly, expected: %v, got: %v", expected, program)
	}
}

func TestFunctionsEquals(t *testing.T) {
	funA := NewFunctionDefinition(
		"main",
		nil,
		VOID,
		NewBlock([]Statement{
			NewVariable(INT, "a", 10, lexer.NewPosition(1, 5)),
			NewIfStatement(
				NewGreaterThanExpression(
					NewIdentifier("a", lexer.NewPosition(6, 8)),
					NewIdentifier("b", lexer.NewPosition(6, 12)),
					lexer.NewPosition(5, 10),
				),
				NewBlock([]Statement{NewAssignment(NewIdentifier("a", lexer.NewPosition(7, 9)), NewIntExpression(0, lexer.NewPosition(7, 13)))}),
				NewBlock([]Statement{NewAssignment(NewIdentifier("b", lexer.NewPosition(10, 9)), NewIntExpression(0, lexer.NewPosition(10, 13)))}),
			),
		}),
		lexer.NewPosition(1, 1),
	)

	funB := NewFunctionDefinition(
		"main",
		nil,
		VOID,
		NewBlock([]Statement{
			NewVariable(INT, "a", 10, lexer.NewPosition(1, 5)),
			NewIfStatement(
				NewGreaterThanExpression(
					NewIdentifier("a", lexer.NewPosition(6, 8)),
					NewIdentifier("b", lexer.NewPosition(6, 12)),
					lexer.NewPosition(5, 10),
				),
				NewBlock([]Statement{NewAssignment(NewIdentifier("a", lexer.NewPosition(7, 9)), NewIntExpression(0, lexer.NewPosition(7, 13)))}),
				NewBlock([]Statement{NewAssignment(NewIdentifier("b", lexer.NewPosition(10, 9)), NewIntExpression(0, lexer.NewPosition(10, 13)))}),
			),
		}),
		lexer.NewPosition(1, 1),
	)

	if !funA.Equals(funB) {
		t.Errorf("Functions are not equal, expected: %v, got: %v", funA, funB)
	}
}

func TestProgramsEquals(t *testing.T) {
	statements := []Statement{
		NewVariable(INT, "a", 10, lexer.NewPosition(1, 5)),
		NewVariable(
			INT,
			"b",
			NewFunctionCall("second", lexer.NewPosition(3, 14), nil),
			lexer.NewPosition(3, 9),
		),
		NewAssignment(
			NewIdentifier("a", lexer.NewPosition(4, 5)),
			NewSubstractExpression(
				NewIdentifier("a", lexer.NewPosition(4, 9)),
				NewIdentifier("b", lexer.NewPosition(4, 13)),
				lexer.NewPosition(4, 11),
			),
		),
		NewIfStatement(
			NewGreaterThanExpression(
				NewIdentifier("a", lexer.NewPosition(6, 8)),
				NewIdentifier("b", lexer.NewPosition(6, 12)),
				lexer.NewPosition(4, 11),
			),
			NewBlock([]Statement{NewAssignment(NewIdentifier("a", lexer.NewPosition(7, 9)), NewIntExpression(0, lexer.NewPosition(7, 13)))}),
			NewBlock([]Statement{NewAssignment(NewIdentifier("b", lexer.NewPosition(10, 9)), NewIntExpression(0, lexer.NewPosition(10, 13)))}),
		),
	}

	funDefsA := map[string]*FunDef{
		"main": NewFunctionDefinition("main", nil, VOID, NewBlock(statements), lexer.NewPosition(1, 1)),
	}

	funDefsB := map[string]*FunDef{
		"main": NewFunctionDefinition("main", nil, VOID, NewBlock(statements), lexer.NewPosition(1, 1)),
	}

	if !funDefsA["main"].Equals(funDefsB["main"]) {
		t.Errorf("Programs are not equal, expected: %v, got: %v", funDefsA, funDefsB)
	}

	programA := NewProgram(funDefsA)
	programB := NewProgram(funDefsB)

	if !programA.Equals(programB) {
		t.Errorf("Programs are not equal, expected: %v, got: %v", programA, programB)
	}
}
