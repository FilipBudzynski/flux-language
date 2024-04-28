package parser

import (
	"reflect"
	"strings"
	"testing"
	"tkom/lexer"
)

// Helper function to test equality of slices
func isSliceEqual(t *testing.T, slice1, slice2 []*Variable) bool {
	if len(slice1) != len(slice2) {
		t.Errorf("slice length mismatch: %d != %d", len(slice1), len(slice2))
		return false
	}

	// Iterate over the elements of the slices and compare them
	for i := range slice1 {
		if !reflect.DeepEqual(slice1[i], slice2[i]) {
			t.Errorf("element mismatch at index %d: %v != %v", i, slice1[i], slice2[i])
			return false
		}
	}

	return true
}

func TestParseParameterGroup(t *testing.T) {
	input := "param1, param2, param3 string"
	lex := createLexer(input)
	parser := NewParser(lex, func(err error) {
		t.Errorf("ParseParameterGroup error: %v", err)
	})

	params := parser.parseParameterGroup()

	if len(params) != 3 {
		t.Errorf("Expected 3 parameters, got %d", len(params))
	}

	expected := []*Variable{
		NewVariable(lexer.STRING, NewIdentifier("param1", lexer.NewPosition(1, 1)), nil),
		NewVariable(lexer.STRING, NewIdentifier("param2", lexer.NewPosition(1, 9)), nil),
		NewVariable(lexer.STRING, NewIdentifier("param3", lexer.NewPosition(1, 17)), nil),
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
	lex := createLexer(input)
	errorHandler := func(err error) {
		t.Errorf("ParseParameterGroup error: %v", err)
	}
	parser := NewParser(lex, errorHandler)

	params := parser.parseParameters()

	for _, param := range params {
		t.Log(param)
	}

	expected := []*Variable{
		NewVariable(lexer.INT, NewIdentifier("param1", lexer.NewPosition(1, 1)), nil),
		NewVariable(lexer.STRING, NewIdentifier("param2", lexer.NewPosition(1, 13)), nil),
		NewVariable(lexer.BOOL, NewIdentifier("param3", lexer.NewPosition(1, 28)), nil),
	}

	if len(params) != len(expected) {
		t.Errorf("Expected %d parameters, got %d", len(expected), len(params))
		return
	}

	for i, param := range params {
		if param.Idetifier != expected[i].Idetifier {
			t.Errorf("Expected parameter name %s, got %s", expected[i].Idetifier.Name, param.Idetifier.Name)
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
	variable1 := NewVariable(lexer.INT, NewIdentifier("a", lexer.NewPosition(1, 12)), nil)
	variable2 := NewVariable(lexer.STRING, NewIdentifier("b", lexer.NewPosition(1, 19)), nil)
	parameters := []*Variable{variable1, variable2}
	expected := NewFunctionDefinition("myFunction", parameters, nil, Block{Statements: []Statement{}}, lexer.NewPosition(1, 1))

	lex := createLexer(input)
	errorHandler := func(err error) {
		t.Errorf("Parse Identifier error: %v", err)
	}
	parser := NewParser(lex, errorHandler)

	functionDefinition := parser.parseFunDef()

	if !isFunctionDefinitionEqual(t, *functionDefinition, *expected) {
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

	if expression != expected {
		t.Errorf("expressions are not equal, expected: %v, got: %v", expected, expression)
	}
}

// func TestParseExpressionGreater(t *testing.T) {
// 	input := "a >= 2"
// 	expected := NewExpression(NewIdentifier("a", lexer.NewPosition(1, 1)), GREATER_OR_EQUAL, 2)
//
// 	lex := createLexer(input)
// 	errorHandler := func(err error) {
// 		t.Errorf("Parse Identifier error: %v", err)
// 	}
// 	parser := NewParser(lex, errorHandler)
//
// 	expression := parser.parseExpression()
//
// 	expressionOpExpr, ok := expression.(OperationExpression)
// 	if !ok {
// 		t.Errorf("Parsed expression is not of type OperationExpression")
// 		return
// 	}
//
// 	if !reflect.DeepEqual(expected, expressionOpExpr) {
// 		t.Errorf("Expressions are not equal, expected: %v, got: %v", expected, expressionOpExpr)
// 		t.Errorf("Actual type: %v", reflect.TypeOf(expressionOpExpr))
// 		t.Errorf("Expected type: %v", reflect.TypeOf(expected))
// 	}
// }

func TestParseExpressionGreater(t *testing.T) {
	input := "a >= 2"
	expected := NewGreaterOrEqualExpression(NewIdentifier("a", lexer.NewPosition(1, 1)), 2)

	lex := createLexer(input)
	errorHandler := func(err error) {
		t.Errorf("Parse Identifier error: %v", err)
	}
	parser := NewParser(lex, errorHandler)

	expression := parser.parseExpression()

    // type assertion
	if _, ok := expression.(GreaterOrEqualExpression); !ok {
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
	input := "int a := 2 - 3"

	identifier := NewIdentifier("a", lexer.NewPosition(1, 5))
	expected := NewVariable(lexer.INT, identifier, NewSubstractExpression(2, 3))
	parser := craeateParser(t, input)

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

func TestParseVariableDeclarationNester(t *testing.T) {
	input := "int a := 2 + 2 + 2"

	identifier := NewIdentifier("a", lexer.NewPosition(1, 5))
	expected := NewVariable(lexer.INT, identifier, NewSumExpression(NewSumExpression(2, 2), 2))
	parser := craeateParser(t, input)

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

// Helper function to create a lexer from a given input string
func createLexer(input string) *lexer.Lexer {
	source, _ := lexer.NewScanner(strings.NewReader(input))
	return lexer.NewLexer(source, 1000, 1000, 1000)
}

func craeateParser(t *testing.T, input string) *Parser {
	lex := createLexer(input)
	errorHandler := func(err error) {
		t.Errorf("Parse Identifier error: %v", err)
	}
	return NewParser(lex, errorHandler)
}

func isFunctionDefinitionEqual(t *testing.T, functionDefinition, expected FunDef) bool {
	if functionDefinition.Name != expected.Name {
		t.Errorf("expected Name: %s, got: %s", expected.Name, functionDefinition.Name)
		return false
	}

	if !isSliceEqual(t, functionDefinition.Parameters, expected.Parameters) {
		t.Errorf("expected Parameters: %v, got: %v", expected.Parameters, functionDefinition.Parameters)
		return false
	}

	if !reflect.DeepEqual(functionDefinition.Block, expected.Block) {
		t.Errorf("expected Block: %v, got: %v", expected.Block, functionDefinition.Block)
		return false
	}

	if functionDefinition.Type != nil && expected.Type != nil && *functionDefinition.Type != *expected.Type {
		t.Errorf("expected Type: %v, got: %v", *expected.Type, *functionDefinition.Type)
		return false
	}

	if functionDefinition.Position != expected.Position {
		t.Errorf("expected Position: %v, got: %v", expected.Position, functionDefinition.Position)
		return false
	}

	if functionDefinition.Position != expected.Position {
		t.Errorf("expected Position: %v, got: %v", expected.Position, functionDefinition.Position)
		return false
	}
	return true
}
