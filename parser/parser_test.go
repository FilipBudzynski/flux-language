package parser

import (
	"log"
	"reflect"
	"strings"
	"testing"
	"tkom/lexer"
)

// Helper function to create a lexer from a given input string
func createLexer(input string) *lexer.Lexer {
	source, _ := lexer.NewScanner(strings.NewReader(input))
	return lexer.NewLexer(source, 1000, 1000, 1000)
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

// Helper function to test equality of slices
func isSliceEqual(t *testing.T, slice1, slice2 []Variable) bool {
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

	expected := []Variable{
		newVariable(lexer.STRING, newIdentifier("param1", lexer.NewPosition(1, 1)), nil),
		newVariable(lexer.STRING, newIdentifier("param2", lexer.NewPosition(1, 9)), nil),
		newVariable(lexer.STRING, newIdentifier("param3", lexer.NewPosition(1, 17)), nil),
	}

	for i, param := range params {
		if param != expected[i] {
			t.Errorf("Expected %d parameters, got %d", len(expected), len(params))
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

	expected := []Variable{
		newVariable(lexer.INT, newIdentifier("param1", lexer.NewPosition(1, 1)), nil),
		newVariable(lexer.STRING, newIdentifier("param2", lexer.NewPosition(1, 13)), nil),
		newVariable(lexer.BOOL, newIdentifier("param3", lexer.NewPosition(1, 28)), nil),
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
	expected := newIdentifier("identifier1", lexer.NewPosition(1, 1))
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
	variable1 := newVariable(lexer.INT, newIdentifier("a", lexer.NewPosition(1, 12)), nil)
	variable2 := newVariable(lexer.STRING, newIdentifier("b", lexer.NewPosition(1, 19)), nil)
	parameters := []Variable{variable1, variable2}
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

func TestParseExpression(t *testing.T) {
	input := "identifier1"
	expected := newIdentifier("identifier1", lexer.NewPosition(1, 1))

	lex := createLexer(input)
	errorHandler := func(err error) {
		t.Errorf("Parse Identifier error: %v", err)
	}
	parser := NewParser(lex, errorHandler)

	expression := parser.parseOrCondition()

	if expression != expected {
		t.Errorf("expressions are not equal, expected: %v, got: %v", expected, expression)
	}
}

func TestParseExpressionGreater(t *testing.T) {
	input := "a > 2"
	expected := NewExpression(newIdentifier("a", lexer.NewPosition(1, 1)), GREATER_THAN, 2)

	lex := createLexer(input)
	errorHandler := func(err error) {
		t.Errorf("Parse Identifier error: %v", err)
	}
	parser := NewParser(lex, errorHandler)

	expression := parser.parseOrCondition()

	expressionOpExpr, ok := expression.(*OperationExpression)
	if !ok {
		t.Errorf("Parsed expression is not of type OperationExpression")
		return
	}
    log.Print(expressionOpExpr)

	if !reflect.DeepEqual(expected, expressionOpExpr) {
		t.Errorf("Expressions are not equal, expected: %v, got: %v", expected, expressionOpExpr)
	}
}
