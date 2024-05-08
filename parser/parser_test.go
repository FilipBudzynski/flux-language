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

// Helper function to check equality of two while statements
func areWhileStatementsEqual(expected, statement *WhileStatement) bool {
	if !reflect.DeepEqual(expected.Condition, statement.Condition) {
		return false
	}

	if len(expected.Instructions) != len(statement.Instructions) {
		return false
	}

	for i := range expected.Instructions {
		if !reflect.DeepEqual(expected.Instructions[i], statement.Instructions[i]) {
			return false
		}
	}
	return true
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

func TestParseEmptyFunctionDefinition(t *testing.T) {
	input := "myFunction(a int, b string) { }"
	variable1 := NewVariable(INT, "a", nil, lexer.NewPosition(1, 12))
	variable2 := NewVariable(STRING, "b", nil, lexer.NewPosition(1, 19))
	parameters := []*Variable{variable1, variable2}
	expected := NewFunctionDefinition("myFunction", parameters, nil, []Statement{}, lexer.NewPosition(1, 1))

	lex := createLexer(input)
	errorHandler := func(err error) {
		t.Errorf("Parse Identifier error: %v", err)
	}
	parser := NewParser(lex, errorHandler)

	functionDefinition := parser.parseFunDef()

	if !functionDefinition.Equals(expected) {
		t.Errorf("function definitions are not equal, expected: %v, got: %v", functionDefinition, expected)
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

func TestParseExpressionGreater(t *testing.T) {
	input := "a >= 2"
	expected := NewGreaterOrEqualExpression(
		NewIdentifier("a", lexer.NewPosition(1, 1)),
		NewIntExpression(2, lexer.NewPosition(1, 6)),
		lexer.NewPosition(1, 3),
	)

	lex := createLexer(input)
	errorHandler := func(err error) {
		t.Errorf("Parse Identifier error: %v", err)
	}
	parser := NewParser(lex, errorHandler)

	expression := parser.parseExpression()

	// type assertion
	if _, ok := expression.(*GreaterOrEqualExpression); !ok {
		t.Errorf("Parsed expression is not of type OperationExpression")
		return
	}

	if !reflect.DeepEqual(expected, expression) {
		t.Errorf("Expressions are not equal, expected: %v, got: %v", expected, expression)
		t.Errorf("Actual type: %v", reflect.TypeOf(expression))
		t.Errorf("Expected type: %v", reflect.TypeOf(expected))
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
		statement := parser.parseVariableDeclaration().(*Variable)

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

func TestParseBoolExpression(t *testing.T) {
	input := "bool c := true"

	expected := NewVariable(BOOL, "c", NewBoolExpression(true, lexer.NewPosition(1, 11)), lexer.NewPosition(1, 6))
	parser := createParser(t, input)

	statement := parser.parseVariableDeclaration()

	statement, ok := statement.(*Variable)
	if !ok {
		t.Errorf("Parsed statement is not of type Variable")
		t.Errorf("Actual type: %v", reflect.TypeOf(statement))
		t.Errorf("Expected type: %v", reflect.TypeOf(expected))
		return
	}

	if !reflect.DeepEqual(expected, statement) {
		t.Errorf("Expressions are not equal, expected: %v, got: %v", expected, statement)
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
		[]Statement{NewAssignment(NewIdentifier("y", lexer.NewPosition(1, 14)), NewIntExpression(20, lexer.NewPosition(1, 18)))},
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
		[]Statement{NewAssignment(NewIdentifier("y", lexer.NewPosition(1, 14)), NewIntExpression(20, lexer.NewPosition(1, 18)))},
		[]Statement{NewAssignment(NewIdentifier("y", lexer.NewPosition(1, 30)), NewIntExpression(15, lexer.NewPosition(1, 34)))},
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
		[]Statement{NewAssignment(idY1, NewSumExpression(idY2, NewIntExpression(1, lexer.NewPosition(1, 24)), lexer.NewPosition(1, 22)))},
	)
	parser := createParser(t, input)

	statement := parser.parseWhileStatement()

	if !areWhileStatementsEqual(expected, statement) {
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
	cases := []Statement{
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

func TestSwitchStatementWithIdentifier(t *testing.T) {
	// switch opcjonalnie expression
	input := `switch {
        a > 2   => fun1(),
        a <= 10 => fun2(),
        default => fun3()
    }`

	cases := []Statement{
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
	}
	expected := NewSwitchStatement(nil, nil, cases)
	parser := createParser(t, input)

	statement := parser.parseSwitchStatement()

	if !expected.Equals(*statement) {
		t.Errorf("Switch statement not parsed correctly, expected: %v, got: %v", expected, statement)
	}
}

func TestPraseSwitchError(t *testing.T) {
	input := `switch {
    >30 => 2
    }`
	expectedError := ParserError{"error [2, 5]: missing switch case, perhaps you have ',' after the last case"}
	errors := []error{}
	errorHandler := func(err error) { errors = append(errors, err) }

	parser := createParser(t, input)
	parser.ErrorHandler = errorHandler

	defer func() {
		err := recover()
		if err == nil {
			t.Errorf("expected panic but didn't")
			return
		}
		parser.ErrorHandler(err.(error))
	}()

	_ = parser.parseSwitchStatement()

	err := errors[0]
	if err == nil {
		t.Errorf("expected panic but didn't")
	} else if err.Error() != expectedError.Error() {
		t.Errorf("expected %v error message, got: %v", expectedError.Error(), err.Error())
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
			[]Statement{NewAssignment(NewIdentifier("a", lexer.NewPosition(7, 9)), NewIntExpression(0, lexer.NewPosition(7, 13)))},
			[]Statement{NewAssignment(NewIdentifier("b", lexer.NewPosition(10, 9)), NewIntExpression(0, lexer.NewPosition(10, 13)))},
		),
	}

	funDefs := map[string]*FunDef{
		"main": NewFunctionDefinition("main", nil, nil, statements, lexer.NewPosition(1, 1)),
	}

	expected := NewProgram(funDefs)
	parser := createParser(t, input)
	program := parser.ParseProgram()

	// if reflect.DeepEqual(expected, program) {
	if !program.Equals(expected) {
		t.Errorf("Program not parsed correctly, expected: %v, got: %v", expected, program)
	}
}

func TestFunctionsEquals(t *testing.T) {
	funA := NewFunctionDefinition(
		"main",
		nil,
		nil,
		[]Statement{
			NewVariable(INT, "a", 10, lexer.NewPosition(1, 5)),
			NewIfStatement(
				NewGreaterThanExpression(
					NewIdentifier("a", lexer.NewPosition(6, 8)),
					NewIdentifier("b", lexer.NewPosition(6, 12)),
					lexer.NewPosition(5, 10),
				),
				[]Statement{NewAssignment(NewIdentifier("a", lexer.NewPosition(7, 9)), NewIntExpression(0, lexer.NewPosition(7, 13)))},
				[]Statement{NewAssignment(NewIdentifier("b", lexer.NewPosition(10, 9)), NewIntExpression(0, lexer.NewPosition(10, 13)))},
			),
		},
		lexer.NewPosition(1, 1),
	)

	funB := NewFunctionDefinition(
		"main",
		nil,
		nil,
		[]Statement{
			NewVariable(INT, "a", 10, lexer.NewPosition(1, 5)),
			NewIfStatement(
				NewGreaterThanExpression(
					NewIdentifier("a", lexer.NewPosition(6, 8)),
					NewIdentifier("b", lexer.NewPosition(6, 12)),
					lexer.NewPosition(5, 10),
				),
				[]Statement{NewAssignment(NewIdentifier("a", lexer.NewPosition(7, 9)), NewIntExpression(0, lexer.NewPosition(7, 13)))},
				[]Statement{NewAssignment(NewIdentifier("b", lexer.NewPosition(10, 9)), NewIntExpression(0, lexer.NewPosition(10, 13)))},
			),
		},
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
			[]Statement{NewAssignment(NewIdentifier("a", lexer.NewPosition(7, 9)), NewIntExpression(0, lexer.NewPosition(7, 13)))},
			[]Statement{NewAssignment(NewIdentifier("b", lexer.NewPosition(10, 9)), NewIntExpression(0, lexer.NewPosition(10, 13)))},
		),
	}

	funDefsA := map[string]*FunDef{
		"main": NewFunctionDefinition("main", nil, nil, statements, lexer.NewPosition(1, 1)),
	}

	funDefsB := map[string]*FunDef{
		"main": NewFunctionDefinition("main", nil, nil, statements, lexer.NewPosition(1, 1)),
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
