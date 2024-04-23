package parser

import (
	"strings"
	"testing"
	"tkom/lexer"
)

// Helper function to create a lexer from a given input string
func createLexer(input string) *lexer.Lexer {
	source, _ := lexer.NewScanner(strings.NewReader(input))
	return lexer.NewLexer(source, 1000, 1000, 1000)
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

	expected := []struct {
		Name string
		Type lexer.TokenTypes
	}{
		{"param1", lexer.STRING},
		{"param2", lexer.STRING},
		{"param3", lexer.STRING},
	}

	for i, param := range params {
		if param.Name != expected[i].Name {
			t.Errorf("Expected parameter name %s, got %s", expected[i].Name, param.Name)
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

	expected := []struct {
		Name string
		Type lexer.TokenTypes
	}{
		{"param1", lexer.INT},
		{"param2", lexer.STRING},
		{"param3", lexer.BOOL},
	}

	if len(params) != len(expected) {
		t.Errorf("Expected %d parameters, got %d", len(expected), len(params))
		return
	}

	for i, param := range params {
		if param.Name != expected[i].Name {
			t.Errorf("Expected parameter name %s, got %s", expected[i].Name, param.Name)
		}
		if param.Type != expected[i].Type {
			t.Errorf("Expected parameter type %v, got %v", expected[i].Type, param.Type)
		}
	}
}
