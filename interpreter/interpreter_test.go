package interpreter

import (
	"fmt"
	"testing"
	"tkom/ast"
	"tkom/shared"
)

func TestVisitIntExpression(t *testing.T) {
	visitor := NewCodeVisitor(nil)
	visitor.VisitIntExpression(&ast.IntExpression{Value: 42})
	if visitor.LastResult != 42 {
		t.Errorf("expected LastResult to be 42, got %v", visitor.LastResult)
	}
}

func TestVisitFloatExpression(t *testing.T) {
	visitor := NewCodeVisitor(nil)
	visitor.VisitFloatExpression(&ast.FloatExpression{Value: 42.0})
	if visitor.LastResult != 42.0 {
		t.Errorf("expected LastResult to be 42.0, got %v", visitor.LastResult)
	}
}

func TestVisitStringExpression(t *testing.T) {
	visitor := NewCodeVisitor(nil)
	visitor.VisitStringExpression(&ast.StringExpression{Value: "42"})
	if visitor.LastResult != "42" {
		t.Errorf("expected LastResult to be 42, got %v", visitor.LastResult)
	}
}

func TestVisitBoolExpression(t *testing.T) {
	visitor := NewCodeVisitor(nil)
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
			visitor := NewCodeVisitor(nil)
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
	visitor := NewCodeVisitor(nil)
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
	visitor := NewCodeVisitor(nil)
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
	visitor := NewCodeVisitor(nil)
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

func TestVisitSubstrackExpression(t *testing.T) {
	visitor := NewCodeVisitor(nil)
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
		{name: "StringToBoolFalse", initialValue: "false", initialType: shared.STRING, targetType: shared.BOOL, expectedResult: false, expectingPanic: false},
		{name: "StringToBoolInvalid", initialValue: "abc", initialType: shared.STRING, targetType: shared.BOOL, expectedResult: nil, expectingPanic: true},
		{name: "StringToString", initialValue: "hello", initialType: shared.STRING, targetType: shared.STRING, expectedResult: "hello", expectingPanic: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			visitor := NewCodeVisitor(nil)
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
	visitor := NewCodeVisitor(nil)
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
	visitor := NewCodeVisitor(nil)
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
	variable := visitor.CurrentScope.InScope("a")

	if variable.Type != shared.STRING {
		t.Errorf("expected variable type to be %v, got %v", shared.STRING, variable.Type)
	}
	if variable.Value != expected {
		t.Errorf("expected variable value to be %v, got %v", expected, variable.Value)
	}
}

func TestGettingValueFromIdentifier(t *testing.T) {
	expected := 22
	visitor := NewCodeVisitor(nil)
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
	visitor := NewCodeVisitor(nil)
	scope := NewScope(nil, nil)
	visitor.CurrentScope = scope
	expectedError := "undefined: a"

	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok || err.Error() != expectedError {
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
	sumAandBfunction := &ast.FunDef{
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
								Expression: &ast.Identifier{
									Name: "e",
								},
							},
						},
					},
				},
				&ast.ReturnStatement{
					Expression: &ast.IntExpression{
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
						&ast.Assignemnt{
							Identifier: ast.Identifier{Name: "c"},
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
	functionMap := map[string]*ast.FunDef{
		"sum_a_b": sumAandBfunction,
	}
	visitor := NewCodeVisitor(functionMap)
	scopeReturnType := shared.VOID
	scope := NewScope(nil, &scopeReturnType)
	visitor.ScopeStack.Push(scope)
	visitor.CurrentScope = scope

	expectedError := fmt.Sprintf(UNDEFINED_VARIABLE, "e")
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok || err.Error() != expectedError {
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
							Expression: ast.NewIntExpression(42, shared.NewPosition(1, 1)),
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
										Expression: ast.NewIntExpression(82, shared.NewPosition(1, 1)),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	funMap := map[string]*ast.FunDef{"main": ast.NewFunctionDefinition("main", []*ast.Variable{}, shared.STRING, block, shared.NewPosition(1, 1))}
	visitor := NewCodeVisitor(funMap)
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
func TestVisitIfStatement_ConditionTrue(t *testing.T) {
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

	visitor := NewCodeVisitor(map[string]*ast.FunDef{})
	visitor.ScopeStack.Push(NewScope(nil, nil))
	visitor.LastResult = 99 // initial value to test the clearing
	visitor.VisitIfStatement(ifStmt)

	if visitor.LastResult != nil {
		t.Errorf("Expected LastResult to be 42, but got %v", visitor.LastResult)
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
func TestVisitIfStatement_ElseBlock(t *testing.T) {
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
	visitor := NewCodeVisitor(nil)
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

	if variable.Value != 42 {
		t.Errorf("expected variable y's value to be 42, got %v", variable.Value)
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
	sumAandBfunction := &ast.FunDef{
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
					Expression: &ast.SumExpression{
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
	functionMap := map[string]*ast.FunDef{
		"sum_a_b": sumAandBfunction,
	}
	visitor := NewCodeVisitor(functionMap)
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

	if visitor.LastResult != 3 {
		t.Errorf("expected lastResult to be %v, got %v", 3, visitor.LastResult)
	}
}

func TestVisitAsignmentWithFunctionCall(t *testing.T) {
	sumAandBfunction := &ast.FunDef{
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
					Expression: &ast.SumExpression{
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
	functionMap := map[string]*ast.FunDef{
		"sum_a_b": sumAandBfunction,
	}
	visitor := NewCodeVisitor(functionMap)
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

	if variable := visitor.CurrentScope.InScope("c"); variable == nil {
		t.Errorf("variable not in scope but should be, got: %v", variable)
	}
	if variable := visitor.CurrentScope.InScope("c"); variable.Value != 3 {
		t.Errorf("expected variable value to be %v, got %v", 3, variable.Value)
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
	sumAandBfunction := &ast.FunDef{
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
								Expression: &ast.SumExpression{
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
					Expression: &ast.IntExpression{
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
						&ast.Assignemnt{
							Identifier: ast.Identifier{Name: "c"},
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
	functionMap := map[string]*ast.FunDef{
		"sum_a_b": sumAandBfunction,
	}
	visitor := NewCodeVisitor(functionMap)
	scopeReturnType := shared.VOID
	scope := NewScope(nil, &scopeReturnType)
	visitor.ScopeStack.Push(scope)
	visitor.CurrentScope = scope
	visitor.VisitBlock(mainBlock)

	if visitor.LastResult != nil {
		t.Errorf("expected lastResult to be %v, got %v", nil, visitor.LastResult)
	}
	if variable := visitor.CurrentScope.InScope("c"); variable == nil {
		t.Errorf("variable not in scope but should be, got: %v", variable)
	}
	if variable := visitor.CurrentScope.InScope("c"); variable.Value != 3 {
		t.Errorf("expected variable value to be %v, got %v", 3, variable.Value)
	}
}

func TestParametersAndArguments(t *testing.T) {
	sumAandBfunction := &ast.FunDef{
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
								Expression: &ast.SumExpression{
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
		arguments     []ast.Expression
		expectedError string
	}{
		{
			name:         "String type mismatch",
			functionName: "sum_a_b",
			arguments: []ast.Expression{
				&ast.StringExpression{Value: "one"},
				&ast.IntExpression{Value: 2},
			},
			expectedError: "Type mismatch: expected INT, got string",
		},
		{
			name:         "float type mismatch",
			functionName: "sum_a_b",
			arguments: []ast.Expression{
				&ast.FloatExpression{Value: 4.20},
				&ast.IntExpression{Value: 2},
			},
			expectedError: "Type mismatch: expected INT, got float64",
		},
		{
			name:         "bool type mismatch",
			functionName: "sum_a_b",
			arguments: []ast.Expression{
				&ast.BoolExpression{Value: true},
				&ast.IntExpression{Value: 2},
			},
			expectedError: "Type mismatch: expected INT, got bool",
		},
		{
			name:         "Wrong number of arguments",
			functionName: "sum_a_b",
			arguments: []ast.Expression{
				&ast.IntExpression{Value: 1},
			},
			expectedError: "function sum_a_b expects 2 arguments but got 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			functionMap := map[string]*ast.FunDef{
				"sum_a_b": sumAandBfunction,
			}
			visitor := NewCodeVisitor(functionMap)
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
					LeftExpression:  ast.Identifier{Name: "a"},
					RightExpression: &ast.IntExpression{Value: 2},
				},
				RightExpression: &ast.LessOrEqualExpression{
					LeftExpression:  ast.Identifier{Name: "a"},
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
				LeftExpression:  ast.Identifier{Name: "a"},
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
					LeftExpression:  ast.Identifier{Name: "a"},
					RightExpression: &ast.IntExpression{Value: 5},
				},
				RightExpression: &ast.LessThanExpression{
					LeftExpression:  ast.Identifier{Name: "a"},
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
				LeftExpression:  ast.Identifier{Name: "a"},
				RightExpression: &ast.IntExpression{Value: 15},
			},
			expectedResult: "Whole bottle",
			conditionMet:   true,
		},
		{
			name:            "a is undefined",
			initialVariable: 3,
			condition: &ast.EqualsExpression{
				LeftExpression:  ast.Identifier{Name: "b"},
				RightExpression: &ast.IntExpression{Value: 5},
			},
			expectedPanicMsg: "undefined: b",
			conditionMet:     false,
		},
		{
			name:            "Condition not met",
			initialVariable: 1,
			condition: &ast.AndExpression{
				LeftExpression: &ast.GreaterThanExpression{
					LeftExpression:  ast.Identifier{Name: "a"},
					RightExpression: &ast.IntExpression{Value: 2},
				},
				RightExpression: &ast.LessOrEqualExpression{
					LeftExpression:  ast.Identifier{Name: "a"},
					RightExpression: &ast.IntExpression{Value: 4},
				},
			},
			expectedResult: "",
			conditionMet:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			visitor := NewCodeVisitor(make(map[string]*ast.FunDef))
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
	visitor := NewCodeVisitor(make(map[string]*ast.FunDef))
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

func TestWhileStatement(t *testing.T) {
	// TODO
}

func TestFunctionCall(t *testing.T) {
	// TODO
}

func TestMultipleFunctionDefinition(t *testing.T) {
	// TODO
}
