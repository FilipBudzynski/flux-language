package interpreter

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
	"tkom/ast"
	"tkom/shared"
)

const MAX_RECURSION_DEPTH = 5

func TestVisitIntExpression(t *testing.T) {
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	visitor.VisitIntExpression(&ast.IntExpression{Value: 42})
	if visitor.LastResult != 42 {
		t.Errorf("expected LastResult to be 42, got %v", visitor.LastResult)
	}
}

func TestVisitFloatExpression(t *testing.T) {
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	visitor.VisitFloatExpression(&ast.FloatExpression{Value: 42.0})
	if visitor.LastResult != 42.0 {
		t.Errorf("expected LastResult to be 42.0, got %v", visitor.LastResult)
	}
}

func TestVisitStringExpression(t *testing.T) {
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	visitor.VisitStringExpression(&ast.StringExpression{Value: "42"})
	if visitor.LastResult != "42" {
		t.Errorf("expected LastResult to be 42, got %v", visitor.LastResult)
	}
}

func TestVisitBoolExpression(t *testing.T) {
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	visitor.VisitBoolExpression(&ast.BoolExpression{Value: true})
	if visitor.LastResult != true {
		t.Errorf("expected LastResult to be true, got %v", visitor.LastResult)
	}
}

// testing the visit negate expression for all types
func TestVisitNegateExpression(t *testing.T) {
	tests := []struct {
		expression     ast.Expression
		expected       interface{}
		name           string
		expectingPanic bool
	}{
		{
			name:           "NegateInt",
			expression:     &ast.NegateExpression{Expression: &ast.IntExpression{Value: 42}},
			expected:       -42,
			expectingPanic: false,
		},
		{
			name:           "NegateFloat",
			expression:     &ast.NegateExpression{Expression: &ast.FloatExpression{Value: 3.14}},
			expected:       -3.14,
			expectingPanic: false,
		},
		{
			name:           "NegateBool",
			expression:     &ast.NegateExpression{Expression: &ast.BoolExpression{Value: true}},
			expected:       false,
			expectingPanic: false,
		},
		{
			name:           "NegateString",
			expression:     &ast.NegateExpression{Expression: &ast.StringExpression{Value: "test"}},
			expected:       nil,
			expectingPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
			defer func() {
				if r := recover(); r != nil {
					if !tt.expectingPanic {
						t.Errorf("expected no panic, but got %v", r)
					}
				} else {
					if tt.expectingPanic {
						t.Errorf("expected panic, but got none")
					}
				}
			}()
			visitor.VisitNegateExpression(tt.expression.(*ast.NegateExpression))
			if visitor.LastResult != tt.expected {
				t.Errorf("expected LastResult to be %v, got %v", tt.expected, visitor.LastResult)
			}
		})
	}
}

func TestVisitAndExpression(t *testing.T) {
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	visitor.VisitAndExpression(&ast.AndExpression{
		LeftExpression: &ast.EqualsExpression{
			LeftExpression:  &ast.IntExpression{Value: 42},
			RightExpression: &ast.IntExpression{Value: 42},
		},
		RightExpression: &ast.EqualsExpression{
			LeftExpression:  &ast.IntExpression{Value: 10},
			RightExpression: &ast.IntExpression{Value: 10},
		},
	})
	if visitor.LastResult != true {
		t.Errorf("expected LastResult to be true, got %v", visitor.LastResult)
	}
}

func TestVisitOrExpression(t *testing.T) {
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	visitor.VisitOrExpression(&ast.OrExpression{
		LeftExpression: &ast.EqualsExpression{
			LeftExpression:  &ast.IntExpression{Value: 50},
			RightExpression: &ast.IntExpression{Value: 42},
		},
		RightExpression: &ast.BoolExpression{Value: true},
	})
	if visitor.LastResult != true {
		t.Errorf("expected LastResult to be true, got %v", visitor.LastResult)
	}
}

func TestVisitSumExpression(t *testing.T) {
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	visitor.VisitSumExpression(
		&ast.SumExpression{
			LeftExpression: &ast.IntExpression{Value: 42},
			RightExpression: &ast.NegateExpression{
				Expression: &ast.IntExpression{Value: 20},
			},
		},
	)
	if visitor.LastResult != 22 {
		t.Errorf("expected LastResult to be 84, got %v", visitor.LastResult)
	}
}

func TestVisitSumExpressionString(t *testing.T) {
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	visitor.VisitSumExpression(
		&ast.SumExpression{
			LeftExpression:  &ast.StringExpression{Value: "even "},
			RightExpression: &ast.StringExpression{Value: "2"},
		},
	)
	if visitor.LastResult != "even 2" {
		t.Errorf("expected LastResult to be 'even 2', got %v", visitor.LastResult)
	}
}

func TestVisitSubstrackExpressionInt(t *testing.T) {
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	visitor.VisitSubstractExpression(
		&ast.SubstractExpression{
			LeftExpression:  &ast.IntExpression{Value: 42},
			RightExpression: &ast.IntExpression{Value: 42},
		},
	)
	if visitor.LastResult != 0 {
		t.Errorf("expected LastResult to be 0, got %v", visitor.LastResult)
	}
}

func TestVisitSubstrackExpressionFloat(t *testing.T) {
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	visitor.VisitSubstractExpression(
		&ast.SubstractExpression{
			LeftExpression:  &ast.FloatExpression{Value: 3.14},
			RightExpression: &ast.FloatExpression{Value: 3.14},
		},
	)
	if visitor.LastResult != 0.0 {
		t.Errorf("expected LastResult to be 0.0, got %v", visitor.LastResult)
	}
}

func TestVisitSubstrackExpressionFloatMinusInt(t *testing.T) {
	expectedOutput := 0.0
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	visitor.VisitSubstractExpression(
		&ast.SubstractExpression{
			LeftExpression:  &ast.FloatExpression{Value: 3.0},
			RightExpression: &ast.IntExpression{Value: 3},
		},
	)
	if visitor.LastResult != expectedOutput {
		t.Errorf("expected LastResult to be %v, got %v", expectedOutput, visitor.LastResult)
	}
}

func TestVisitSubstrackExpressionIntMinusFloat(t *testing.T) {
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("expected panic, but didn't get one")
		}
	}()

	visitor.VisitSubstractExpression(
		&ast.SubstractExpression{
			LeftExpression:  &ast.IntExpression{Value: 3},
			RightExpression: &ast.FloatExpression{Value: 3.0},
		},
	)
}

// testing cast expression for every type to every type
func TestCastExpression(t *testing.T) {
	tests := []struct {
		initialValue   any
		expectedResult any
		name           string
		expectingPanic bool
		initialType    shared.TypeAnnotation
		targetType     shared.TypeAnnotation
	}{
		// Int to other types
		{name: "IntToInt", initialValue: 10, initialType: shared.INT, targetType: shared.INT, expectedResult: 10, expectingPanic: false},
		{name: "IntToFloat", initialValue: 10, initialType: shared.INT, targetType: shared.FLOAT, expectedResult: 10.0, expectingPanic: false},
		{name: "IntToBool", initialValue: 10, initialType: shared.INT, targetType: shared.BOOL, expectedResult: true, expectingPanic: false},
		{name: "IntToBoolZero", initialValue: 0, initialType: shared.INT, targetType: shared.BOOL, expectedResult: false, expectingPanic: false},
		{name: "IntToString", initialValue: 10, initialType: shared.INT, targetType: shared.STRING, expectedResult: "10", expectingPanic: false},

		// Float to other types
		{name: "FloatToInt", initialValue: 10.5, initialType: shared.FLOAT, targetType: shared.INT, expectedResult: 10, expectingPanic: false},
		{name: "FloatToFloat", initialValue: 10.5, initialType: shared.FLOAT, targetType: shared.FLOAT, expectedResult: 10.5, expectingPanic: false},
		{name: "FloatToBool", initialValue: 10.5, initialType: shared.FLOAT, targetType: shared.BOOL, expectedResult: true, expectingPanic: false},
		{name: "FloatToBoolZero", initialValue: 0.0, initialType: shared.FLOAT, targetType: shared.BOOL, expectedResult: false, expectingPanic: false},
		{name: "FloatToString", initialValue: 10.5, initialType: shared.FLOAT, targetType: shared.STRING, expectedResult: "10.5", expectingPanic: false},

		// Bool to other types
		{name: "BoolToIntTrue", initialValue: true, initialType: shared.BOOL, targetType: shared.INT, expectedResult: 1, expectingPanic: false},
		{name: "BoolToIntFalse", initialValue: false, initialType: shared.BOOL, targetType: shared.INT, expectedResult: 0, expectingPanic: false},
		{name: "BoolToFloatTrue", initialValue: true, initialType: shared.BOOL, targetType: shared.FLOAT, expectedResult: 1.0, expectingPanic: false},
		{name: "BoolToFloatFalse", initialValue: false, initialType: shared.BOOL, targetType: shared.FLOAT, expectedResult: 0.0, expectingPanic: false},
		{name: "BoolToBoolTrue", initialValue: true, initialType: shared.BOOL, targetType: shared.BOOL, expectedResult: true, expectingPanic: false},
		{name: "BoolToBoolFalse", initialValue: false, initialType: shared.BOOL, targetType: shared.BOOL, expectedResult: false, expectingPanic: false},
		{name: "BoolToStringTrue", initialValue: true, initialType: shared.BOOL, targetType: shared.STRING, expectedResult: "true", expectingPanic: false},
		{name: "BoolToStringFalse", initialValue: false, initialType: shared.BOOL, targetType: shared.STRING, expectedResult: "false", expectingPanic: false},

		// String to other types
		{name: "StringToInt", initialValue: "10", initialType: shared.STRING, targetType: shared.INT, expectedResult: 10, expectingPanic: false},
		{name: "StringToIntInvalid", initialValue: "abc", initialType: shared.STRING, targetType: shared.INT, expectedResult: nil, expectingPanic: true},
		{name: "StringToFloat", initialValue: "10.5", initialType: shared.STRING, targetType: shared.FLOAT, expectedResult: 10.5, expectingPanic: false},
		{name: "StringToFloatInvalid", initialValue: "abc", initialType: shared.STRING, targetType: shared.FLOAT, expectedResult: nil, expectingPanic: true},
		{name: "StringToBoolTrue", initialValue: "true", initialType: shared.STRING, targetType: shared.BOOL, expectedResult: true, expectingPanic: false},
		{name: "StringToBoolFalse", initialValue: "", initialType: shared.STRING, targetType: shared.BOOL, expectedResult: false, expectingPanic: false},
		{name: "StringToBoolInvalid", initialValue: "abc", initialType: shared.STRING, targetType: shared.BOOL, expectedResult: true, expectingPanic: false},
		{name: "StringToString", initialValue: "hello", initialType: shared.STRING, targetType: shared.STRING, expectedResult: "hello", expectingPanic: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
			defer func() {
				if r := recover(); r != nil {
					if !tt.expectingPanic {
						t.Errorf("expected no panic, but got %v", r)
					}
				} else {
					if tt.expectingPanic {
						t.Errorf("expected panic, but got none")
					}
				}
			}()
			var castExpr *ast.CastExpression
			switch tt.initialType {
			case shared.INT:
				castExpr = &ast.CastExpression{LeftExpression: &ast.IntExpression{Value: tt.initialValue.(int)}, TypeAnnotation: tt.targetType}
			case shared.FLOAT:
				castExpr = &ast.CastExpression{LeftExpression: &ast.FloatExpression{Value: tt.initialValue.(float64)}, TypeAnnotation: tt.targetType}
			case shared.BOOL:
				castExpr = &ast.CastExpression{LeftExpression: &ast.BoolExpression{Value: tt.initialValue.(bool)}, TypeAnnotation: tt.targetType}
			case shared.STRING:
				castExpr = &ast.CastExpression{LeftExpression: &ast.StringExpression{Value: tt.initialValue.(string)}, TypeAnnotation: tt.targetType}
			default:
				t.Fatalf("unsupported initial type %v", tt.initialType)
			}

			visitor.VisitCastExpression(castExpr)

			if !tt.expectingPanic && visitor.LastResult != tt.expectedResult {
				t.Errorf("expected LastResult to be %v, got %v", tt.expectedResult, visitor.LastResult)
			}
		})
	}
}

func TestVisitIdentifier(t *testing.T) {
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	scope := NewScope(nil, nil)
	err := scope.AddVariable("a", 10, shared.INT, shared.NewPosition(1, 1))
	visitor.CurrentScope = scope
	visitor.VisitIdentifier(
		&ast.Identifier{
			Name: "a",
		},
	)
	if err != nil {
		t.Error("unexpected error", err)
	}
	if visitor.LastResult != 10 {
		t.Errorf("expected LastResult to be 0, got %v", visitor.LastResult)
	}
}

func TestVisitVariable(t *testing.T) {
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	scope := NewScope(nil, nil)
	visitor.CurrentScope = scope
	visitor.VisitVariable(
		&ast.Variable{
			Name:  "a",
			Value: ast.NewStringExpression("some string", shared.NewPosition(1, 1)),
			Type:  shared.STRING,
		},
	)
	expected := "some string"
	_, variable := visitor.CurrentScope.InScope("a")
	variableType := visitor.DetermineType(variable)

	if variableType != shared.STRING {
		t.Errorf("expected variable type to be %v, got %v", shared.STRING, variableType)
	}
	if variable != expected {
		t.Errorf("expected variable value to be %v, got %v", expected, variable)
	}
}

func TestGettingValueFromIdentifier(t *testing.T) {
	expected := 22
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	scope := NewScope(nil, nil)
	visitor.CurrentScope = scope
	visitor.VisitVariable(
		&ast.Variable{
			Name:  "a",
			Value: ast.NewIntExpression(22, shared.NewPosition(1, 1)),
			Type:  shared.INT,
		},
	)

	visitor.VisitIdentifier(
		&ast.Identifier{
			Name: "a",
		},
	)

	if visitor.LastResult != expected {
		t.Errorf("expected lastResult to be %v, got %v", 10, visitor.LastResult)
	}
}

func TestVariableNotInScope(t *testing.T) {
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	scope := NewScope(nil, nil)
	visitor.CurrentScope = scope
	expectedError := NewSemanticError("undefined: a", shared.NewPosition(0, 0))

	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok || err.Error() != expectedError.Error() {
				t.Errorf("Expected panic with error: %v, but got: %v", expectedError, r)
			}
		} else {
			t.Errorf("Expected panic due to undefined variable 'e', but did not panic")
		}
	}()

	visitor.VisitIdentifier(
		&ast.Identifier{
			Name: "a",
		},
	)
}

// in this scenario we expect the error of undefined "a" because we dont allow to use variables
// from different function scopes
//
// TESTING CASE:
//
//	sum_a_b (a, b int) int {
//	  if a > 0 {
//	    return e
//	  }
//	  return 0
//	}
//
//	{
//	  int e := 22
//	  int c := 0
//	  if true {
//	    c = sum_a_b(1, 2)
//	  }
//	}
func TestSearchVariableInScope(t *testing.T) {
	sumAandBfunction := &ast.FunctionDefinition{
		Name: "sum_a_b",
		Type: shared.INT,
		Parameters: []*ast.Variable{
			{
				Name: "a",
				Type: shared.INT,
			},
			{
				Name: "b",
				Type: shared.INT,
			},
		},
		Block: &ast.Block{
			Statements: []ast.Statement{
				&ast.IfStatement{
					Condition: &ast.GreaterThanExpression{
						LeftExpression: &ast.Identifier{
							Name: "a",
						},
						RightExpression: &ast.IntExpression{
							Value: 0,
						},
					},
					InstructionsBlock: &ast.Block{
						Statements: []ast.Statement{
							&ast.ReturnStatement{
								Value: &ast.Identifier{
									Name: "e",
								},
							},
						},
					},
				},
				&ast.ReturnStatement{
					Value: &ast.IntExpression{
						Value: 0,
					},
				},
			},
		},
	}
	mainBlock := &ast.Block{
		Statements: []ast.Statement{
			&ast.Variable{
				Name:  "e",
				Type:  shared.INT,
				Value: &ast.IntExpression{Value: 22},
			},
			&ast.Variable{
				Name:  "c",
				Type:  shared.INT,
				Value: &ast.IntExpression{Value: 0},
			},
			&ast.IfStatement{
				Condition: &ast.BoolExpression{Value: true},
				InstructionsBlock: &ast.Block{
					Statements: []ast.Statement{
						&ast.Assignment{
							Identifier: &ast.Identifier{Name: "c"},
							Value: &ast.FunctionCall{
								Name: "sum_a_b",
								Arguments: []ast.Expression{
									&ast.IntExpression{Value: 1},
									&ast.IntExpression{Value: 2},
								},
							},
						},
					},
				},
			},
		},
	}
	// functionMap := map[string]*ast.FunctionDefinition{
	functionsMap := map[string]ast.Function{
		"sum_a_b": sumAandBfunction,
	}
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	visitor.FunctionsMap = functionsMap
	scopeReturnType := shared.VOID
	scope := NewScope(nil, &scopeReturnType)
	visitor.ScopeStack.Push(scope)
	visitor.CurrentScope = scope

	expectedError := NewSemanticError(fmt.Sprintf(UNDEFINED_VARIABLE, "e"), shared.NewPosition(0, 0))
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok || err.Error() != expectedError.Error() {
				t.Errorf("Expected panic with error: %v, but got: %v", expectedError, r)
			}
		} else {
			t.Errorf("Expected panic due to undefined variable 'e', but did not panic")
		}
	}()

	visitor.VisitBlock(mainBlock)
}

// in this scenario we expect the returned value to be "42"
// TESTING CASE:
//
//	{
//	    if ture {
//	        return 42
//	    }
//
//	    while true {
//	        return 82
//	    }
//	}
func TestReturningNestedBlocks(t *testing.T) {
	block := &ast.Block{
		Statements: []ast.Statement{
			&ast.IfStatement{
				Condition: ast.NewBoolExpression(true, shared.NewPosition(1, 1)),
				InstructionsBlock: &ast.Block{
					Statements: []ast.Statement{
						&ast.ReturnStatement{
							Value: ast.NewIntExpression(42, shared.NewPosition(1, 1)),
						},
					},
				},
				ElseInstructionsBlock: &ast.Block{
					Statements: []ast.Statement{
						&ast.WhileStatement{
							Condition: ast.NewBoolExpression(true, shared.NewPosition(1, 1)),
							InstructionsBlock: &ast.Block{
								Statements: []ast.Statement{
									&ast.ReturnStatement{
										Value: ast.NewIntExpression(82, shared.NewPosition(1, 1)),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	//	funMap := map[string]*ast.FunctionDefinition{"main": ast.NewFunctionDefinition("main", []*ast.Variable{}, shared.STRING, block, shared.NewPosition(1, 1))}
	funMap := map[string]ast.Function{"main": ast.NewFunctionDefinition("main", []*ast.Variable{}, shared.STRING, block, shared.NewPosition(1, 1))}
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	visitor.FunctionsMap = funMap
	typeInt := shared.INT
	visitor.ScopeStack.Push(NewScope(nil, &typeInt))
	block.Accept(visitor)

	if visitor.ReturnFlag != true {
		t.Errorf("Expected ReturnFlag to be true, got false")
	}
	if visitor.LastResult != 42 {
		t.Errorf("Expected LastResult to be 42, got %v", visitor.LastResult)
	}
}

// in this scenario lastResult should be cleared from value "99"
// TESTING CASE:
//
//	{
//	    if ture {
//	        a := 42
//	    }
//	}
func TestVisitIfStatementConditionTrue(t *testing.T) {
	condition := ast.NewBoolExpression(true, shared.NewPosition(1, 1))
	block := &ast.Block{
		Statements: []ast.Statement{
			&ast.Variable{
				Value:    ast.NewIntExpression(42, shared.NewPosition(1, 1)),
				Name:     "a",
				Type:     shared.INT,
				Position: shared.NewPosition(1, 1),
			},
		},
	}
	ifStmt := &ast.IfStatement{
		Condition:         condition,
		InstructionsBlock: block,
	}

	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	visitor.ScopeStack.Push(NewScope(nil, nil))
	visitor.LastResult = 99
	visitor.VisitIfStatement(ifStmt)

	if visitor.LastResult != nil {
		t.Errorf("Expected LastResult to be 42, but got %v", visitor.LastResult)
	}
}

// in this scenario lastResult should be cleared from value "99"
// TESTING CASE:
//
//	{
//	    if false {
//	        int b := 42
//	    }
//	    else {
//	        string table := "stół z powyłamywanymi nogami"
//	    }
//	}
func TestVisitIfStatementElseBlock(t *testing.T) {
	condition := &ast.BoolExpression{Value: false}
	block := &ast.Block{
		Statements: []ast.Statement{
			&ast.Variable{
				Value:    ast.NewIntExpression(42, shared.NewPosition(1, 1)),
				Name:     "b",
				Type:     shared.INT,
				Position: shared.NewPosition(1, 1),
			},
		},
	}
	elseBlock := &ast.Block{
		Statements: []ast.Statement{
			&ast.Variable{
				Value:    ast.NewStringExpression("stół z powyłamywanymi nogami", shared.NewPosition(1, 1)),
				Name:     "table",
				Type:     shared.STRING,
				Position: shared.NewPosition(1, 1),
			},
		},
	}
	ifStmt := &ast.IfStatement{
		Condition:             condition,
		InstructionsBlock:     block,
		ElseInstructionsBlock: elseBlock,
	}

	visitor := &CodeVisitor{}
	visitor.LastResult = 99 // initial value to test the cleaning of the LastResult
	visitor.VisitIfStatement(ifStmt)

	if visitor.LastResult != nil {
		t.Errorf("Expected LastResult to be nil, but got %v", visitor.LastResult)
	}
}

// Testing if scope variables are actually stored in the scope
func TestScopeVariables(t *testing.T) {
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	globalScope := NewScope(nil, nil)
	visitor.CurrentScope = globalScope

	block := &ast.Block{
		Statements: []ast.Statement{
			&ast.Variable{Name: "y", Value: ast.NewIntExpression(42, shared.NewPosition(1, 1)), Type: shared.INT},
			&ast.IntExpression{Value: 42},
		},
	}

	newScope := NewScope(visitor.CurrentScope, nil)
	visitor.ScopeStack.Push(newScope)
	visitor.CurrentScope = newScope

	block.Accept(visitor)

	variable, err := newScope.GetVariable("y")
	if err != nil {
		t.Errorf("expected variable y to be defined in current scope, but it was not found")
	}

	if variable != 42 {
		t.Errorf("expected variable y's value to be 42, got %v", variable)
	}

	if visitor.LastResult != 42 {
		t.Errorf("expected LastResult to be 42, got %v", visitor.LastResult)
	}

	poppedScope, err := visitor.ScopeStack.Pop()
	if err != nil {
		t.Errorf("error popping scope: %v", err)
	}
	visitor.CurrentScope = poppedScope
}

// in this scenario we are testing declaring a variable using function call
// expecting result is c variable in scope with value 3
//
// TESTING CASE:
//
//	sum_a_b (a, b int) int {
//	 return a + b
//	}
//
//	{
//	   int c := sum_a_b(1, 2)
//	}
func TestVisitFunctionCall(t *testing.T) {
	sumAandBfunction := &ast.FunctionDefinition{
		Name: "sum_a_b",
		Type: shared.INT,
		Parameters: []*ast.Variable{
			{
				Name: "a",
				Type: shared.INT,
			},
			{
				Name: "b",
				Type: shared.INT,
			},
		},
		Block: &ast.Block{
			Statements: []ast.Statement{
				&ast.ReturnStatement{
					Value: &ast.SumExpression{
						LeftExpression: &ast.Identifier{
							Name: "a",
						},
						RightExpression: &ast.Identifier{
							Name: "b",
						},
					},
				},
			},
		},
	}
	// functionMap := map[string]*ast.FunctionDefinition{
	functionsMap := map[string]ast.Function{
		"sum_a_b": sumAandBfunction,
	}
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	visitor.FunctionsMap = functionsMap
	scope := NewScope(nil, nil)
	visitor.CurrentScope = scope
	visitor.VisitFunctionCall(
		&ast.FunctionCall{
			Name: "sum_a_b",
			Arguments: []ast.Expression{
				&ast.IntExpression{
					Value: 1,
				},
				&ast.IntExpression{
					Value: 2,
				},
			},
		},
	)

	if visitor.ReturnFlag {
		t.Errorf("expected returnFlag to be false but is %v", visitor.ReturnFlag)
	}
	if visitor.LastResult != 3 {
		t.Errorf("expected lastResult to be %v, got %v", 3, visitor.LastResult)
	}
}

func TestVisitFunctionCallWithIdentifier(t *testing.T) {
	sumAandBfunction := &ast.FunctionDefinition{
		Name: "sum_a_b",
		Type: shared.INT,
		Parameters: []*ast.Variable{
			{
				Name: "a",
				Type: shared.INT,
			},
			{
				Name: "b",
				Type: shared.INT,
			},
		},
		Block: &ast.Block{
			Statements: []ast.Statement{
				&ast.ReturnStatement{
					Value: &ast.SumExpression{
						LeftExpression: &ast.Identifier{
							Name: "a",
						},
						RightExpression: &ast.Identifier{
							Name: "b",
						},
					},
				},
			},
		},
	}
	//	functionMap := map[string]*ast.FunctionDefinition{
	functionsMap := map[string]ast.Function{
		"sum_a_b": sumAandBfunction,
	}
	voidType := shared.VOID
	scope := NewScope(nil, &voidType)
	scope.AddVariable("one", 1, shared.INT, shared.Position{Line: 1, Column: 1})
	scope.AddVariable("two", 2, shared.INT, shared.Position{Line: 1, Column: 1})
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	visitor.FunctionsMap = functionsMap
	visitor.CurrentScope = scope
	visitor.VisitFunctionCall(
		&ast.FunctionCall{
			Name: "sum_a_b",
			Arguments: []ast.Expression{
				&ast.Identifier{
					Name: "one",
				},
				&ast.Identifier{
					Name: "two",
				},
			},
		},
	)

	if visitor.LastResult != 3 {
		t.Errorf("expected lastResult to be %v, got %v", 3, visitor.LastResult)
	}
}

func TestVisitAsignmentWithFunctionCall(t *testing.T) {
	sumAandBfunction := &ast.FunctionDefinition{
		Name: "sum_a_b",
		Type: shared.INT,
		Parameters: []*ast.Variable{
			{
				Name: "a",
				Type: shared.INT,
			},
			{
				Name: "b",
				Type: shared.INT,
			},
		},
		Block: &ast.Block{
			Statements: []ast.Statement{
				&ast.ReturnStatement{
					Value: &ast.SumExpression{
						LeftExpression: &ast.Identifier{
							Name: "a",
						},
						RightExpression: &ast.Identifier{
							Name: "b",
						},
					},
				},
			},
		},
	}
	// functionMap := map[string]*ast.FunctionDefinition{
	functionsMap := map[string]ast.Function{
		"sum_a_b": sumAandBfunction,
	}
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	visitor.FunctionsMap = functionsMap
	scope := NewScope(nil, nil)
	visitor.CurrentScope = scope
	visitor.VisitVariable(
		&ast.Variable{
			Name: "c",
			Type: shared.INT,
			Value: &ast.FunctionCall{
				Name: "sum_a_b",
				Arguments: []ast.Expression{
					&ast.IntExpression{
						Value: 1,
					},
					&ast.IntExpression{
						Value: 2,
					},
				},
			},
		},
	)

	if _ ,variable := visitor.CurrentScope.InScope("c"); variable == nil {
		t.Errorf("variable not in scope but should be, got: %v", variable)
	}
	if _, variable := visitor.CurrentScope.InScope("c"); variable != 3 {
		t.Errorf("expected variable value to be %v, got %v", 3, variable)
	}
}

// in this scenario we are testing returning from nested scopes with function calls
//
// TESTING CASE:
//
//	sum_a_b (a, b int) int {
//	  if a > 0 {
//	    return a + b
//	  }
//	  return 0
//	}
//
//	{
//	  int c := 0
//	  if true {
//	    c = sum_a_b(1, 2)
//	  }
//	  string d := "hello"
//	}
func TestVisitNestedFunctionCallWithReturn(t *testing.T) {
	sumAandBfunction := &ast.FunctionDefinition{
		Name: "sum_a_b",
		Type: shared.INT,
		Parameters: []*ast.Variable{
			{
				Name: "a",
				Type: shared.INT,
			},
			{
				Name: "b",
				Type: shared.INT,
			},
		},
		Block: &ast.Block{
			Statements: []ast.Statement{
				&ast.IfStatement{
					Condition: &ast.GreaterThanExpression{
						LeftExpression: &ast.Identifier{
							Name: "a",
						},
						RightExpression: &ast.IntExpression{
							Value: 0,
						},
					},
					InstructionsBlock: &ast.Block{
						Statements: []ast.Statement{
							&ast.ReturnStatement{
								Value: &ast.SumExpression{
									LeftExpression: &ast.Identifier{
										Name: "a",
									},
									RightExpression: &ast.Identifier{
										Name: "b",
									},
								},
							},
						},
					},
				},
				&ast.ReturnStatement{
					Value: &ast.IntExpression{
						Value: 0,
					},
				},
			},
		},
	}
	mainBlock := &ast.Block{
		Statements: []ast.Statement{
			&ast.Variable{
				Name:  "c",
				Type:  shared.INT,
				Value: &ast.IntExpression{Value: 0},
			},
			&ast.IfStatement{
				Condition: &ast.BoolExpression{Value: true},
				InstructionsBlock: &ast.Block{
					Statements: []ast.Statement{
						&ast.Assignment{
							Identifier: &ast.Identifier{Name: "c"},
							Value: &ast.FunctionCall{
								Name: "sum_a_b",
								Arguments: []ast.Expression{
									&ast.IntExpression{Value: 1},
									&ast.IntExpression{Value: 2},
								},
							},
						},
					},
				},
			},
			&ast.Variable{
				Name:  "d",
				Type:  shared.STRING,
				Value: &ast.StringExpression{Value: "hello"},
			},
		},
	}
	// functionMap := map[string]*ast.FunctionDefinition{
	functionsMap := map[string]ast.Function{
		"sum_a_b": sumAandBfunction,
	}
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	visitor.FunctionsMap = functionsMap
	scopeReturnType := shared.VOID
	scope := NewScope(nil, &scopeReturnType)
	visitor.ScopeStack.Push(scope)
	visitor.CurrentScope = scope
	visitor.VisitBlock(mainBlock)

	if visitor.ReturnFlag {
		t.Errorf("expected returnFlag to be false, got %v", visitor.ReturnFlag)
	}
	if visitor.LastResult != nil {
		t.Errorf("expected lastResult to be %v, got %v", nil, visitor.LastResult)
	}
	if _, variable := visitor.CurrentScope.InScope("c"); variable == nil {
		t.Errorf("variable not in scope but should be, got: %v", variable)
	}
	if _, variable := visitor.CurrentScope.InScope("c"); variable != 3 {
		t.Errorf("expected variable value to be %v, got %v", 3, variable)
	}
}

func TestParametersAndArguments(t *testing.T) {
	sumAandBfunction := &ast.FunctionDefinition{
		Name: "sum_a_b",
		Type: shared.INT,
		Parameters: []*ast.Variable{
			{
				Name: "a",
				Type: shared.INT,
			},
			{
				Name: "b",
				Type: shared.INT,
			},
		},
		Block: &ast.Block{
			Statements: []ast.Statement{
				&ast.IfStatement{
					Condition: &ast.GreaterThanExpression{
						LeftExpression: &ast.Identifier{
							Name: "a",
						},
						RightExpression: &ast.IntExpression{
							Value: 0,
						},
					},
					InstructionsBlock: &ast.Block{
						Statements: []ast.Statement{
							&ast.ReturnStatement{
								Value: &ast.SumExpression{
									LeftExpression: &ast.Identifier{
										Name: "a",
									},
									RightExpression: &ast.Identifier{
										Name: "b",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name          string
		functionName  string
		expectedError string
		arguments     []ast.Expression
	}{
		{
			name:         "String type mismatch",
			functionName: "sum_a_b",
			arguments: []ast.Expression{
				&ast.StringExpression{Value: "one"},
				&ast.IntExpression{Value: 2},
			},
			expectedError: NewSemanticError(fmt.Sprintf(WRONG_ARGUMENT_TYPE, shared.STRING, shared.INT), shared.NewPosition(0, 0)).Error(),
		},
		{
			name:         "float type mismatch",
			functionName: "sum_a_b",
			arguments: []ast.Expression{
				&ast.FloatExpression{Value: 4.20},
				&ast.IntExpression{Value: 2},
			},
			expectedError: NewSemanticError(fmt.Sprintf(WRONG_ARGUMENT_TYPE, shared.FLOAT, shared.INT), shared.NewPosition(0, 0)).Error(),
		},
		{
			name:         "bool type mismatch",
			functionName: "sum_a_b",
			arguments: []ast.Expression{
				&ast.BoolExpression{Value: true},
				&ast.IntExpression{Value: 2},
			},
			expectedError: NewSemanticError(fmt.Sprintf(WRONG_ARGUMENT_TYPE, shared.BOOL, shared.INT), shared.NewPosition(0, 0)).Error(),
		},
		{
			name:         "Wrong number of arguments",
			functionName: "sum_a_b",
			arguments: []ast.Expression{
				&ast.IntExpression{Value: 1},
			},
			expectedError: NewSemanticError(fmt.Sprintf(WRONG_NUMBER_OF_ARGUMENTS, "sum_a_b", 2, 1), shared.NewPosition(0, 0)).Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			functionsMap := map[string]ast.Function{
				"sum_a_b": sumAandBfunction,
			}
			visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
			visitor.FunctionsMap = functionsMap
			scopeReturnType := shared.VOID
			scope := NewScope(nil, &scopeReturnType)
			visitor.ScopeStack.Push(scope)
			visitor.CurrentScope = scope

			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok || err.Error() != tt.expectedError {
						t.Errorf("Expected panic with error: %v, but got: %v", tt.expectedError, r)
					}
				} else {
					t.Errorf("Expected panic due to: %v, but did not panic", tt.expectedError)
				}
			}()

			visitor.VisitVariable(
				&ast.Variable{
					Name: "c",
					Type: shared.INT,
					Value: &ast.FunctionCall{
						Name:      tt.functionName,
						Arguments: tt.arguments,
					},
				},
			)
		})
	}
}

// a>2 and a<=4  => "A pint",
func TestVisitSwitchCase(t *testing.T) {
	tests := []struct {
		name             string
		condition        ast.Expression
		expectedResult   string
		expectedPanicMsg string
		initialVariable  int
		conditionMet     bool
	}{
		{
			name:            "a > 2 and a <= 4",
			initialVariable: 3,
			condition: &ast.AndExpression{
				LeftExpression: &ast.GreaterThanExpression{
					LeftExpression:  &ast.Identifier{Name: "a"},
					RightExpression: &ast.IntExpression{Value: 2},
				},
				RightExpression: &ast.LessOrEqualExpression{
					LeftExpression:  &ast.Identifier{Name: "a"},
					RightExpression: &ast.IntExpression{Value: 4},
				},
			},
			expectedResult: "A pint",
			conditionMet:   true,
		},
		{
			name:            "a == 5",
			initialVariable: 5,
			condition: &ast.EqualsExpression{
				LeftExpression:  &ast.Identifier{Name: "a"},
				RightExpression: &ast.IntExpression{Value: 5},
			},
			expectedResult: "Decent beverage",
			conditionMet:   true,
		},
		{
			name:            "a > 5 and a < 15",
			initialVariable: 6,
			condition: &ast.AndExpression{
				LeftExpression: &ast.GreaterThanExpression{
					LeftExpression:  &ast.Identifier{Name: "a"},
					RightExpression: &ast.IntExpression{Value: 5},
				},
				RightExpression: &ast.LessThanExpression{
					LeftExpression:  &ast.Identifier{Name: "a"},
					RightExpression: &ast.IntExpression{Value: 15},
				},
			},
			expectedResult: "A NICE beverage",
			conditionMet:   true,
		},
		{
			name:            "a > 15",
			initialVariable: 16,
			condition: &ast.GreaterThanExpression{
				LeftExpression:  &ast.Identifier{Name: "a"},
				RightExpression: &ast.IntExpression{Value: 15},
			},
			expectedResult: "Whole bottle",
			conditionMet:   true,
		},
		{
			name:            "a is undefined",
			initialVariable: 3,
			condition: &ast.EqualsExpression{
				LeftExpression:  &ast.Identifier{Name: "b"},
				RightExpression: &ast.IntExpression{Value: 5},
			},
			expectedPanicMsg: NewSemanticError("undefined: b", shared.NewPosition(0, 0)).Error(),
			conditionMet:     false,
		},
		{
			name:            "Condition not met",
			initialVariable: 1,
			condition: &ast.AndExpression{
				LeftExpression: &ast.GreaterThanExpression{
					LeftExpression:  &ast.Identifier{Name: "a"},
					RightExpression: &ast.IntExpression{Value: 2},
				},
				RightExpression: &ast.LessOrEqualExpression{
					LeftExpression:  &ast.Identifier{Name: "a"},
					RightExpression: &ast.IntExpression{Value: 4},
				},
			},
			expectedResult: "",
			conditionMet:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
			scopeReturnType := shared.STRING
			baseScope := NewScope(nil, &scopeReturnType)
			visitor.ScopeStack.Push(baseScope)
			scope := NewScope(nil, nil)
			visitor.ScopeStack.Push(scope)
			_ = scope.AddVariable("a", tt.initialVariable, shared.INT, shared.NewPosition(1, 1))
			visitor.CurrentScope = scope

			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok || err.Error() != tt.expectedPanicMsg {
						t.Errorf("Expected panic with error: %v, but got: %v", tt.expectedPanicMsg, r)
					}
				} else if tt.expectedPanicMsg != "" {
					t.Errorf("Expected panic with error: %v, but did not panic", tt.expectedPanicMsg)
				}
			}()

			visitor.VisitSwitchCase(&ast.SwitchCase{
				Condition:        tt.condition,
				OutputExpression: &ast.StringExpression{Value: tt.expectedResult},
			})

			if tt.conditionMet {
				if visitor.LastResult != tt.expectedResult {
					t.Errorf("Expected lastResult to be '%v', got: %v", tt.expectedResult, visitor.LastResult)
				}
			} else {
				if visitor.LastResult != nil {
					t.Errorf("Expected lastResult to be nil, got: %v", visitor.LastResult)
				}
			}
		})
	}
}

// default => "A pint",
func TestVisitDefaultSwitchCase(t *testing.T) {
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	scopeReturnType := shared.STRING
	baseScope := NewScope(nil, &scopeReturnType)
	visitor.ScopeStack.Push(baseScope)
	scope := NewScope(nil, nil)
	visitor.ScopeStack.Push(scope)
	visitor.CurrentScope = scope
	visitor.VisitDefaultSwitchCase(&ast.DefaultSwitchCase{
		OutputExpression: &ast.StringExpression{Value: "A pint"},
	})

	if visitor.LastResult != "A pint" {
		t.Errorf("Expected lastResult to be 'A pint', got: %v", visitor.LastResult)
	}
}

func TestVisitSwitchStatement(t *testing.T) {
	tests := []struct {
		switchStmt     *ast.SwitchStatement
		expectedResult any
		name           string
		expectPanic    bool
	}{
		{
			name: "Single SwitchCase matches",
			switchStmt: &ast.SwitchStatement{
				Variables: []*ast.Variable{
					{Name: "a", Type: shared.INT, Value: &ast.IntExpression{Value: 3}},
				},
				Cases: []ast.Case{
					&ast.SwitchCase{
						Condition: &ast.AndExpression{
							LeftExpression: &ast.GreaterThanExpression{
								LeftExpression:  &ast.Identifier{Name: "a"},
								RightExpression: &ast.IntExpression{Value: 2},
							},
							RightExpression: &ast.LessOrEqualExpression{
								LeftExpression:  &ast.Identifier{Name: "a"},
								RightExpression: &ast.IntExpression{Value: 4},
							},
						},
						OutputExpression: &ast.StringExpression{Value: "A pint"},
					},
				},
			},
			expectedResult: "A pint",
			expectPanic:    false,
		},
		{
			name: "Single SwitchCase does not match, DefaultSwitchCase executed",
			switchStmt: &ast.SwitchStatement{
				Variables: []*ast.Variable{
					{Name: "a", Type: shared.INT, Value: &ast.IntExpression{Value: 5}},
				},
				Cases: []ast.Case{
					&ast.SwitchCase{
						Condition: &ast.AndExpression{
							LeftExpression: &ast.GreaterThanExpression{
								LeftExpression:  &ast.Identifier{Name: "a"},
								RightExpression: &ast.IntExpression{Value: 2},
							},
							RightExpression: &ast.LessOrEqualExpression{
								LeftExpression:  &ast.Identifier{Name: "a"},
								RightExpression: &ast.IntExpression{Value: 4},
							},
						},
						OutputExpression: &ast.StringExpression{Value: "A pint"},
					},
					&ast.DefaultSwitchCase{
						OutputExpression: &ast.StringExpression{Value: "Decent beverage"},
					},
				},
			},
			expectedResult: "Decent beverage",
			expectPanic:    false,
		},
		{
			name: "Multiple DefaultSwitchCase instances",
			switchStmt: &ast.SwitchStatement{
				Variables: []*ast.Variable{
					{Name: "a", Type: shared.INT, Value: &ast.IntExpression{Value: 5}},
				},
				Cases: []ast.Case{
					&ast.DefaultSwitchCase{
						OutputExpression: &ast.StringExpression{Value: "Decent beverage"},
					},
					&ast.DefaultSwitchCase{
						OutputExpression: &ast.StringExpression{Value: "Whole bottle"},
					},
				},
			},
			expectedResult: nil,
			expectPanic:    true,
		},
		{
			name: "No SwitchCase matches, no DefaultSwitchCase",
			switchStmt: &ast.SwitchStatement{
				Variables: []*ast.Variable{
					{Name: "a", Type: shared.INT, Value: &ast.IntExpression{Value: 1}},
				},
				Cases: []ast.Case{
					&ast.SwitchCase{
						Condition: &ast.AndExpression{
							LeftExpression: &ast.GreaterThanExpression{
								LeftExpression:  &ast.Identifier{Name: "a"},
								RightExpression: &ast.IntExpression{Value: 2},
							},
							RightExpression: &ast.LessOrEqualExpression{
								LeftExpression:  &ast.Identifier{Name: "a"},
								RightExpression: &ast.IntExpression{Value: 4},
							},
						},
						OutputExpression: &ast.StringExpression{Value: "A pint"},
					},
				},
			},
			expectedResult: nil,
			expectPanic:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
			scopeReturnType := shared.STRING
			baseScope := NewScope(nil, &scopeReturnType)
			visitor.ScopeStack.Push(baseScope)
			scope := NewScope(nil, nil)
			visitor.ScopeStack.Push(scope)
			visitor.CurrentScope = scope

			defer func() {
				if r := recover(); r != nil {
					if test.expectPanic {
						err, ok := r.(error)
						if !ok || err.Error() != NewSemanticError(MULTIPLE_DEFAULT_CASES, shared.NewPosition(0, 0)).Error() {
							t.Errorf("Expected panic with error '%v', but got: %v", MULTIPLE_DEFAULT_CASES, r)
						}
					} else {
						t.Errorf("Unexpected panic: %v", r)
					}
				} else if test.expectPanic {
					t.Errorf("Expected panic, but did not panic")
				} else if visitor.LastResult != test.expectedResult {
					t.Errorf("Expected lastResult to be %v, got: %v", test.expectedResult, visitor.LastResult)
				}
			}()

			visitor.VisitSwitchStatement(test.switchStmt)
		})
	}
}

//	sum(a, b int) int {
//	    return a + b
//	}
//
//	switch int a := sum(1, 2) {
//	    a>2 and a<=4  => "A pint"
//	}
func TestSwitchStatementWithFunctionCall(t *testing.T) {
	sumAandBfunction := &ast.FunctionDefinition{
		Name: "sum_a_b",
		Type: shared.INT,
		Parameters: []*ast.Variable{
			{
				Name: "a",
				Type: shared.INT,
			},
			{
				Name: "b",
				Type: shared.INT,
			},
		},
		Block: &ast.Block{
			Statements: []ast.Statement{
				&ast.ReturnStatement{
					Value: &ast.SumExpression{
						LeftExpression: &ast.Identifier{
							Name: "a",
						},
						RightExpression: &ast.Identifier{
							Name: "b",
						},
					},
				},
			},
		},
	}
	switchStmt := &ast.SwitchStatement{
		Variables: []*ast.Variable{
			{
				Name: "a",
				Type: shared.INT,
				Value: &ast.FunctionCall{
					Name: "sum_a_b",
					Arguments: []ast.Expression{
						&ast.IntExpression{Value: 1},
						&ast.IntExpression{Value: 2},
					},
				},
			},
		},
		Cases: []ast.Case{
			&ast.SwitchCase{
				Condition: &ast.AndExpression{
					LeftExpression: &ast.GreaterThanExpression{
						LeftExpression:  &ast.Identifier{Name: "a"},
						RightExpression: &ast.IntExpression{Value: 2},
					},
					RightExpression: &ast.LessOrEqualExpression{
						LeftExpression:  &ast.Identifier{Name: "a"},
						RightExpression: &ast.IntExpression{Value: 4},
					},
				},
				OutputExpression: &ast.StringExpression{Value: "A pint"},
			},
		},
	}

	//	functionsMap := map[string]*ast.FunctionDefinition{
	functionsMap := map[string]ast.Function{
		"sum_a_b": sumAandBfunction,
	}
	returnType := shared.STRING
	scope := NewScope(nil, &returnType)
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	visitor.ScopeStack.Push(scope)
	visitor.CurrentScope = scope
	visitor.FunctionsMap = functionsMap
	visitor.VisitSwitchStatement(switchStmt)

	if visitor.LastResult != "A pint" {
		t.Errorf("Expected lastResult to be 'A pint', got: %v", visitor.LastResult)
	}
}

// this switch should not return because it has block, not expression
// it should only change the value of i
//
// int i := 10
//
//	switch {
//	  i > 0   => {
//		   i = i + i
//		 },
//		 default => {
//		   i = 0
//		 }
//	}
func TestSwitchWithBlock(t *testing.T) {
	switchStmt := &ast.SwitchStatement{
		Variables: []*ast.Variable{},
		Cases: []ast.Case{
			&ast.SwitchCase{
				Condition: &ast.GreaterThanExpression{
					LeftExpression:  &ast.Identifier{Name: "i"},
					RightExpression: &ast.IntExpression{Value: 0},
				},
				OutputExpression: &ast.Block{
					Statements: []ast.Statement{
						&ast.Assignment{
							Identifier: &ast.Identifier{Name: "i"},
							Value: &ast.SumExpression{
								LeftExpression:  &ast.Identifier{Name: "i"},
								RightExpression: &ast.Identifier{Name: "i"},
							},
						},
					},
				},
			},
			&ast.DefaultSwitchCase{
				OutputExpression: &ast.Block{
					Statements: []ast.Statement{
						&ast.Assignment{
							Identifier: &ast.Identifier{Name: "i"},
							Value:      &ast.IntExpression{Value: 0},
						},
					},
				},
			},
		},
	}
	outerBlock := &ast.Block{
		Statements: []ast.Statement{
			&ast.Variable{
				Name: "i",
				Type: shared.INT,
				Value: &ast.IntExpression{
					Value: 10,
				},
			},
			switchStmt,
		},
	}

	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	scopeType := shared.VOID
	newScope := NewScope(nil, &scopeType)
	visitor.ScopeStack.Push(newScope)
	visitor.CurrentScope = newScope
	visitor.VisitBlock(outerBlock)

	if visitor.ReturnFlag {
		t.Errorf("Expected returnFlag to be false, got: %v", visitor.ReturnFlag)
	}
    if _, value := visitor.CurrentScope.InScope("i"); value != 20 {
		t.Errorf("Expected lastResult to be 20, got: %v", visitor.LastResult)
	}
}

// here we are testing whether we will get an error due to the lack of
// return in function holding switch statement
func TestFunctionWithSwitch(t *testing.T) {
	switchStatement := &ast.SwitchStatement{
		Variables: []*ast.Variable{},
		Cases: []ast.Case{
			&ast.SwitchCase{
				Condition: &ast.AndExpression{
					LeftExpression: &ast.GreaterThanExpression{
						LeftExpression:  &ast.Identifier{Name: "a"},
						RightExpression: &ast.IntExpression{Value: 2},
					},
					RightExpression: &ast.LessThanExpression{
						LeftExpression:  &ast.Identifier{Name: "a"},
						RightExpression: &ast.IntExpression{Value: 4},
					},
				},
				OutputExpression: &ast.StringExpression{Value: "sample text"},
			},
		},
	}
	functionWithSwitch := &ast.FunctionDefinition{
		Name: "isItThree",
		Type: shared.STRING,
		Parameters: []*ast.Variable{
			{
				Name: "a",
				Type: shared.INT,
			},
		},
		Block: &ast.Block{
			Statements: []ast.Statement{
				switchStatement,
			},
		},
	}
	mainBlock := &ast.Block{
		Statements: []ast.Statement{
			&ast.Variable{
				Name:  "someInt",
				Type:  shared.INT,
				Value: &ast.IntExpression{Value: 22},
			},
			&ast.FunctionCall{
				Name: "isItThree",
				Arguments: []ast.Expression{
					&ast.Identifier{
						Name: "someInt",
					},
				},
			},
		},
	}

	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	visitor.FunctionsMap["isItThree"] = functionWithSwitch
	scopeReturnType := shared.VOID
	baseScope := NewScope(nil, &scopeReturnType)
	visitor.ScopeStack.Push(baseScope)
	visitor.CurrentScope = baseScope

	errorMsg := NewSemanticError(fmt.Sprintf(MISSING_RETURN, shared.STRING), shared.NewPosition(0, 0)).Error()
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok || err.Error() != errorMsg {
				t.Errorf("Expected panic with error '%s', but got: %v", errorMsg, r)
			}
		} else {
			t.Errorf("Expected panic due to missing return, but did not panic")
		}
	}()

	visitor.VisitBlock(mainBlock)
}

// testing if while statement corectly changes variables in parent scope
//
// TESTING CASE:
//
//	{
//		int i := 0
//		while i < 5 {
//			i = i + 1
//		}
//	}
func TestVisitWhileStatement(t *testing.T) {
	block := &ast.Block{
		Statements: []ast.Statement{
			&ast.Variable{
				Value: &ast.IntExpression{Value: 0},
				Name:  "i",
				Type:  shared.INT,
			},
			&ast.WhileStatement{
				Condition: &ast.LessThanExpression{
					LeftExpression:  &ast.Identifier{Name: "i"},
					RightExpression: &ast.IntExpression{Value: 5},
				},
				InstructionsBlock: &ast.Block{
					Statements: []ast.Statement{
						&ast.Assignment{
							Identifier: &ast.Identifier{Name: "i"},
							Value: &ast.SumExpression{
								LeftExpression:  &ast.Identifier{Name: "i"},
								RightExpression: &ast.IntExpression{Value: 1},
							},
						},
					},
				},
			},
		},
	}

	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	voidType := shared.VOID
	newScope := NewScope(nil, &voidType)
	visitor.ScopeStack.Push(newScope)
	visitor.CurrentScope = newScope
	visitor.VisitBlock(block)

	if _, variable := visitor.CurrentScope.InScope("i"); variable != 5 {
		t.Errorf("Expected variable 'i' in parent scope to be 4, but is %v", variable)
	}
}

// testing valid returning from while statement
func TestVisitWhileStatementWithReturn(t *testing.T) {
	block := &ast.Block{
		Statements: []ast.Statement{
			&ast.Variable{
				Value: &ast.IntExpression{Value: 0},
				Name:  "i",
				Type:  shared.INT,
			},
			&ast.WhileStatement{
				Condition: &ast.LessThanExpression{
					LeftExpression:  &ast.Identifier{Name: "i"},
					RightExpression: &ast.IntExpression{Value: 5},
				},
				InstructionsBlock: &ast.Block{
					Statements: []ast.Statement{
						&ast.Assignment{
							Identifier: &ast.Identifier{Name: "i"},
							Value: &ast.SumExpression{
								LeftExpression:  &ast.Identifier{Name: "i"},
								RightExpression: &ast.IntExpression{Value: 1},
							},
						},
						&ast.IfStatement{
							Condition: &ast.EqualsExpression{
								LeftExpression:  &ast.Identifier{Name: "i"},
								RightExpression: &ast.IntExpression{Value: 3},
							},
							InstructionsBlock: &ast.Block{
								Statements: []ast.Statement{
									&ast.ReturnStatement{
										Value: &ast.IntExpression{Value: 22},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	scopeReturnType := shared.INT
	newScope := NewScope(nil, &scopeReturnType)
	visitor.ScopeStack.Push(newScope)
	visitor.CurrentScope = newScope
	visitor.VisitBlock(block)

	if !visitor.ReturnFlag {
		t.Errorf("Expected return flag to be set")
	}
	if _, variable := visitor.CurrentScope.InScope("i"); variable != 3 {
		t.Errorf("Expected variable 'i' in parent scope to be 3, but is %v", variable)
	}
}

func TestEmbededFunction(t *testing.T) {
	block := &ast.Block{
		Statements: []ast.Statement{
			&ast.FunctionCall{
				Name: "println",
				Arguments: []ast.Expression{
					&ast.IntExpression{Value: 22},
				},
			},
			&ast.FunctionCall{
				Name: "println",
				Arguments: []ast.Expression{
					&ast.StringExpression{Value: "halo halo"},
				},
			},
		},
	}
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	scopeReturnType := shared.VOID
	newScope := NewScope(nil, &scopeReturnType)
	visitor.ScopeStack.Push(newScope)
	visitor.CurrentScope = newScope

	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	visitor.VisitBlock(block)

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	expected := "22\nhalo halo\n"
	actual := out
	if actual != expected {
		t.Errorf("Printed value is incorrect. Expected: %s, Got: %s", expected, actual)
	}
}

func TestVariableDeclarationMissmatch(t *testing.T) {
	variableDeclaration := &ast.Variable{
		Value: &ast.StringExpression{Value: "missmatch???"},
		Name:  "i",
		Type:  shared.INT,
	}

	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	scopeReturnType := shared.INT
	newScope := NewScope(nil, &scopeReturnType)
	visitor.ScopeStack.Push(newScope)
	visitor.CurrentScope = newScope
	expectedError := NewSemanticError(fmt.Sprintf(TYPE_MISMATCH, shared.INT, reflect.TypeOf("missmatch???")), shared.NewPosition(0, 0))

	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok || err.Error() != expectedError.Error() {
				t.Errorf("Expected panic with error: %v, but got: %v", expectedError, r)
			}
		} else {
			t.Errorf("Expected panic due to type mismatch, but did not panic")
		}
	}()

	visitor.VisitVariable(variableDeclaration)
}

func TestRecursion(t *testing.T) {
	recursiveFunc := &ast.FunctionDefinition{
		Name: "recursiveFunc",
		Block: &ast.Block{
			Statements: []ast.Statement{
				&ast.FunctionCall{Name: "recursiveFunc", Arguments: []ast.Expression{}},
			},
		},
		Parameters: []*ast.Variable{},
		Type:       shared.VOID,
		Position:   shared.NewPosition(0, 0),
	}
	functioncall := &ast.FunctionCall{
		Name:      "recursiveFunc",
		Arguments: []ast.Expression{},
	}

	functions := map[string]ast.Function{
		recursiveFunc.Name: recursiveFunc,
	}
	visitor := NewCodeVisitor(MAX_RECURSION_DEPTH)
	visitor.FunctionsMap = functions
	scopeReturnType := shared.INT
	newScope := NewScope(nil, &scopeReturnType)
	visitor.ScopeStack.Push(newScope)
	visitor.CurrentScope = newScope
	visitor.MaxRecursionDepth = 2
	expectedError := NewSemanticError(fmt.Sprintf(MAX_RECURSION_DEPTH_EXCEEDED, recursiveFunc.Name), shared.NewPosition(0, 0))

	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok || err.Error() != expectedError.Error() {
				t.Errorf("Expected panic with error: %v, but got: %v", expectedError, r)
			}
		} else {
			t.Errorf("Expected panic due to type mismatch, but did not panic")
		}
	}()

	visitor.VisitFunctionCall(functioncall)
}
