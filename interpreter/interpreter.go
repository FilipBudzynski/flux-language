package interpreter

import (
	"fmt"
	"reflect"
	"strconv"
	"tkom/ast"
	"tkom/shared"
)

type Stack struct {
	elem []*Scope
}

func (s *Stack) Push(value *Scope) {
	s.elem = append(s.elem, value)
}

func (s *Stack) Pop() (*Scope, error) {
	if len(s.elem) == 0 {
		return nil, fmt.Errorf("stack is empty")
	}
	top := s.elem[len(s.elem)-1]
	s.elem = s.elem[:len(s.elem)-1]
	return top, nil
}

func (s *Stack) Peek() (*Scope, error) {
	if len(s.elem) == 0 {
		return nil, fmt.Errorf("stack is empty")
	}
	return s.elem[len(s.elem)-1], nil
}

func (s *Stack) IsEmpty() bool {
	return len(s.elem) == 0
}

func (s *Stack) Size() int {
	return len(s.elem)
}

type CodeVisitor struct {
	LastResult        any
	Functions         map[string]*ast.FunDef
	CurrentScope      *Scope
	ScopeStack        Stack
	ReturnFlag        bool
	CurrentReturnType shared.TypeAnnotation
}

func NewCodeVisitor(functions map[string]*ast.FunDef) *CodeVisitor {
	return &CodeVisitor{
		Functions:         functions,
		ScopeStack:        Stack{elem: []*Scope{}},
		LastResult:        nil,
		CurrentScope:      nil,
		ReturnFlag:        false,
		CurrentReturnType: shared.VOID,
	}
}

func (c *CodeVisitor) tryCastToInt(value any) (int, error) {
	switch val := value.(type) {
	case int:
		return val, nil
	case float64:
		return int(val), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case string:
		return strconv.Atoi(val)
	default:
		return 0, fmt.Errorf("invalid cast expression: %v to int", value)
	}
}

func (v *CodeVisitor) tryCastToFloat(value any) (float64, error) {
	switch val := value.(type) {
	case int:
		return float64(val), nil
	case float64:
		return val, nil
	case bool:
		if val {
			return 1.0, nil
		}
		return 0.0, nil
	case string:
		return strconv.ParseFloat(val, 64)
	default:
		return 0, fmt.Errorf("invalid cast expression: %v to float", value)
	}
}

func (v *CodeVisitor) tryCastToBool(value any) (bool, error) {
	switch val := value.(type) {
	case int:
		return val != 0, nil
	case float64:
		return val != 0.0, nil
	case bool:
		return val, nil
	case string:
		return strconv.ParseBool(val)
	default:
		return false, fmt.Errorf("invalid cast expression: %v to bool", value)
	}
}

func (v *CodeVisitor) tryCastToString(value any) (string, error) {
	return fmt.Sprintf("%v", value), nil
}

func (v *CodeVisitor) createVariable(name string, variableType shared.TypeAnnotation) {
	// variable := &ast.Variable{
	// 	Name: name,
	// 	Type: variableType,
	// }
	// err := v.CurrentScope.AddVariable(name, variable)
	// if err != nil {
	// 	// bo inaczej kazda funkcja visitor'a musi zwracac error - mozliwe ze to dobrze
	// 	panic(err)
	// }

	// not to propagate the value further
	// v.LastResult = nil
}

func (v *CodeVisitor) returnToFunctionDefScope() {
	for v.CurrentScope.ReturnType == nil {
		parentScope, err := v.ScopeStack.Pop()
		if err != nil {
			panic(err)
		}
		v.CurrentScope = parentScope
	}
	// parentScope, err := v.ScopeStack.Pop()
	// if err != nil {
	// 	panic(err)
	// }
	// v.CurrentScope = parentScope
}

func (v *CodeVisitor) getCurrentFunctionReturnType() shared.TypeAnnotation {
	for i := len(v.ScopeStack.elem) - 1; i >= 0; i-- {
		if v.ScopeStack.elem[i].ReturnType != nil {
			return *v.ScopeStack.elem[i].ReturnType
		}
	}
	return shared.VOID
}

func (v *CodeVisitor) VisitIntExpression(e *ast.IntExpression) {
	v.LastResult = e.Value
}

func (v *CodeVisitor) VisitFloatExpression(e *ast.FloatExpression) {
	v.LastResult = e.Value
}

func (v *CodeVisitor) VisitStringExpression(e *ast.StringExpression) {
	v.LastResult = e.Value
}

func (v *CodeVisitor) VisitBoolExpression(e *ast.BoolExpression) {
	v.LastResult = e.Value
}

func (v *CodeVisitor) VisitIdentifier(e *ast.Identifier) {
	v.LastResult = e.Name
}

func (v *CodeVisitor) VisitNegateExpression(e *ast.NegateExpression) {
	e.Expression.Accept(v)
	ne := v.LastResult

	switch value := ne.(type) {
	case int:
		v.LastResult = -value
	case float64:
		v.LastResult = -value
	case bool:
		v.LastResult = !value
	default:
		panic(NewSemanticError(fmt.Sprintf(INVALID_NEGATE_EXPRESSION, ne, "string"), e.Position))
	}
}

func (v *CodeVisitor) VisitCastExpression(e *ast.CastExpression) {
	e.LeftExpression.Accept(v)
	ce := v.LastResult

	var result any
	var err error

	switch e.TypeAnnotation {
	case shared.INT:
		result, err = v.tryCastToInt(ce)
	case shared.FLOAT:
		result, err = v.tryCastToFloat(ce)
	case shared.BOOL:
		result, err = v.tryCastToBool(ce)
	case shared.STRING:
		result, err = v.tryCastToString(ce)
	default:
		panic(NewSemanticError(fmt.Sprintf(INVALID_TYPE_ANNOTATION, e.TypeAnnotation), e.Position))
	}

	if err != nil {
		panic(NewSemanticError(err.Error(), e.Position))
	}

	v.LastResult = result
}

func (v *CodeVisitor) VisitMultiplyExpression(e *ast.MultiplyExpression) {
	e.LeftExpression.Accept(v)
	leftResult := v.LastResult

	e.RightExpression.Accept(v)
	rightResult := v.LastResult

	switch leftResult.(type) {
	case int:
		switch rightResult.(type) {
		case int:
			v.LastResult = leftResult.(int) * rightResult.(int)
		}
	case float64:
		switch rightResult.(type) {
		case float64:
			v.LastResult = leftResult.(float64) * rightResult.(float64)
		}
	default:
		panic(NewSemanticError(fmt.Sprintf(INVALID_MULTIPLY_EXPRESSION, leftResult, rightResult), e.Position))
	}
}

func (v *CodeVisitor) VisitDivideExpression(e *ast.DivideExpression) {
	e.LeftExpression.Accept(v)
	leftResult := v.LastResult

	e.RightExpression.Accept(v)
	rightResult := v.LastResult

	// check for zero as
	if right, ok := rightResult.(int); ok && right == 0 {
		panic(NewSemanticError("Division by zero", e.Position))
	} else if right, ok := rightResult.(float64); ok && right == 0.0 {
		panic(NewSemanticError("Division by zero", e.Position))
	}

	switch leftResult.(type) {
	case int:
		switch rightResult.(type) {
		case int:
			v.LastResult = leftResult.(int) / rightResult.(int)
		}
	case float64:
		switch rightResult.(type) {
		case float64:
			v.LastResult = leftResult.(float64) * rightResult.(float64)
		}
	default:
		panic(NewSemanticError(fmt.Sprintf(INVALID_DIVISION_EXPRESSION, leftResult, rightResult), e.Position))
	}
}

func (v *CodeVisitor) VisitSumExpression(e *ast.SumExpression) {
	e.LeftExpression.Accept(v)
	leftResult := v.LastResult

	e.RightExpression.Accept(v)
	rightResult := v.LastResult

	switch leftResult.(type) {
	case int:
		switch rightResult.(type) {
		case int:
			v.LastResult = leftResult.(int) + rightResult.(int)
		case float64:
			switch rightResult.(type) {
			case float64:
				v.LastResult = leftResult.(float64) + rightResult.(float64)
			}
		case string:
			switch rightResult.(type) {
			case string:
				v.LastResult = leftResult.(string) + rightResult.(string)
			}
		default:
			panic(NewSemanticError(fmt.Sprintf(INVALID_SUM_EXPRESSION, leftResult, rightResult), e.Position))
		}
	}
}

func (v *CodeVisitor) VisitSubstractExpression(e *ast.SubstractExpression) {
	e.LeftExpression.Accept(v)
	leftResult := v.LastResult

	e.RightExpression.Accept(v)
	rightResult := v.LastResult

	switch leftResult.(type) {
	case int:
		switch rightResult.(type) {
		case int:
			v.LastResult = leftResult.(int) - rightResult.(int)
		}
	case float64:
		switch rightResult.(type) {
		case float64:
			v.LastResult = leftResult.(float64) - rightResult.(float64)
		}
	default:
		panic(NewSemanticError(fmt.Sprintf(INVALID_SUBSTRACT_EXPRESSION, leftResult, rightResult), e.Position))
	}
}

func (v *CodeVisitor) VisitEqualsExpression(e *ast.EqualsExpression) {
	e.LeftExpression.Accept(v)
	leftResult := v.LastResult
	e.RightExpression.Accept(v)
	rightResult := v.LastResult

	if reflect.TypeOf(leftResult) != reflect.TypeOf(rightResult) {
		panic(NewSemanticError(fmt.Sprintf(INVALID_EQUALS_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), e.Position))
	}

	v.LastResult = leftResult == rightResult
}

func (v *CodeVisitor) VisitNotEqualsExpression(e *ast.NotEqualsExpression) {
	e.LeftExpression.Accept(v)
	leftResult := v.LastResult
	e.RightExpression.Accept(v)
	rightResult := v.LastResult

	if reflect.TypeOf(leftResult) != reflect.TypeOf(rightResult) {
		panic(NewSemanticError(fmt.Sprintf(INVALID_NOT_EQUALS_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), e.Position))
	}

	v.LastResult = leftResult != rightResult
}

func (v *CodeVisitor) VisitGreaterThanExpression(e *ast.GreaterThanExpression) {
	e.LeftExpression.Accept(v)
	leftResult := v.LastResult
	e.RightExpression.Accept(v)
	rightResult := v.LastResult

	if reflect.TypeOf(leftResult) != reflect.TypeOf(rightResult) {
		panic(NewSemanticError(fmt.Sprintf(INVALID_GREATER_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), e.Position))
	}

	switch left := leftResult.(type) {
	case int:
		right, ok := rightResult.(int)
		if !ok {
			panic(NewSemanticError(fmt.Sprintf(INVALID_GREATER_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), e.Position))
		}
		v.LastResult = left > right
	case float64:
		right, ok := rightResult.(float64)
		if !ok {
			panic(NewSemanticError(fmt.Sprintf(INVALID_GREATER_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), e.Position))
		}
		v.LastResult = left > right
	default:
		panic(NewSemanticError(fmt.Sprintf("Unsupported types for > operator: %v", reflect.TypeOf(leftResult)), e.Position))
	}
}

func (v *CodeVisitor) VisitGreaterOrEqualExpression(e *ast.GreaterOrEqualExpression) {
	e.LeftExpression.Accept(v)
	leftResult := v.LastResult
	e.RightExpression.Accept(v)
	rightResult := v.LastResult

	if reflect.TypeOf(leftResult) != reflect.TypeOf(rightResult) {
		panic(NewSemanticError(fmt.Sprintf(INVALID_GREATER_OR_EQUALS_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), e.Position))
	}

	switch left := leftResult.(type) {
	case int:
		right, ok := rightResult.(int)
		if !ok {
			panic(NewSemanticError(fmt.Sprintf(INVALID_GREATER_OR_EQUALS_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), e.Position))
		}
		v.LastResult = left >= right
	case float64:
		right, ok := rightResult.(float64)
		if !ok {
			panic(NewSemanticError(fmt.Sprintf(INVALID_GREATER_OR_EQUALS_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), e.Position))
		}
		v.LastResult = left >= right
	default:
		panic(NewSemanticError(fmt.Sprintf("Unsupported types for > operator: %v", reflect.TypeOf(leftResult)), e.Position))
	}
}

func (v *CodeVisitor) VisitLessThanExpression(e *ast.LessThanExpression) {
	e.LeftExpression.Accept(v)
	leftResult := v.LastResult
	e.RightExpression.Accept(v)
	rightResult := v.LastResult

	if reflect.TypeOf(leftResult) != reflect.TypeOf(rightResult) {
		panic(NewSemanticError(fmt.Sprintf(INVALID_GREATER_OR_EQUALS_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), e.Position))
	}

	switch left := leftResult.(type) {
	case int:
		right, ok := rightResult.(int)
		if !ok {
			panic(NewSemanticError(fmt.Sprintf(INVALID_GREATER_OR_EQUALS_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), e.Position))
		}
		v.LastResult = left < right
	case float64:
		right, ok := rightResult.(float64)
		if !ok {
			panic(NewSemanticError(fmt.Sprintf(INVALID_GREATER_OR_EQUALS_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), e.Position))
		}
		v.LastResult = left < right
	default:
		panic(NewSemanticError(fmt.Sprintf("Unsupported types for > operator: %v", reflect.TypeOf(leftResult)), e.Position))
	}
}

func (v *CodeVisitor) VisitLessOrEqualExpression(e *ast.LessOrEqualExpression) {
	e.LeftExpression.Accept(v)
	leftResult := v.LastResult
	e.RightExpression.Accept(v)
	rightResult := v.LastResult

	if reflect.TypeOf(leftResult) != reflect.TypeOf(rightResult) {
		panic(NewSemanticError(fmt.Sprintf(INVALID_GREATER_OR_EQUALS_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), e.Position))
	}

	switch left := leftResult.(type) {
	case int:
		right, ok := rightResult.(int)
		if !ok {
			panic(NewSemanticError(fmt.Sprintf(INVALID_GREATER_OR_EQUALS_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), e.Position))
		}
		v.LastResult = left <= right
	case float64:
		right, ok := rightResult.(float64)
		if !ok {
			panic(NewSemanticError(fmt.Sprintf(INVALID_GREATER_OR_EQUALS_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), e.Position))
		}
		v.LastResult = left <= right
	default:
		panic(NewSemanticError(fmt.Sprintf("Unsupported types for > operator: %v", reflect.TypeOf(leftResult)), e.Position))
	}
}

func (v *CodeVisitor) VisitOrExpression(e *ast.OrExpression) {
	e.LeftExpression.Accept(v)
	leftResult := v.LastResult

	leftBool, ok := leftResult.(bool)
	if !ok {
		panic(NewSemanticError(fmt.Sprintf("Expected boolean expression but got %v", reflect.TypeOf(leftResult)), e.Position))
	}

	// If the left expression is true, return true
	if leftBool {
		v.LastResult = true
		return
	}

	e.RightExpression.Accept(v)
	rightResult := v.LastResult

	// Check if the right result is a boolean
	rightBool, ok := rightResult.(bool)
	if !ok {
		panic(NewSemanticError(fmt.Sprintf("Expected boolean expression but got %v", reflect.TypeOf(rightResult)), e.Position))
	}

	v.LastResult = rightBool
}

func (v *CodeVisitor) VisitAndExpression(e *ast.AndExpression) {
	e.LeftExpression.Accept(v)
	leftResult := v.LastResult

	leftBool, ok := leftResult.(bool)
	if !ok {
		panic(NewSemanticError(fmt.Sprintf("Expected boolean expression but got %v", reflect.TypeOf(leftResult)), e.Position))
	}

	// If the left expression is false, return false
	if !leftBool {
		v.LastResult = false
		return
	}

	e.RightExpression.Accept(v)
	rightResult := v.LastResult

	rightBool, ok := rightResult.(bool)
	if !ok {
		panic(NewSemanticError(fmt.Sprintf("Expected boolean expression but got %v", reflect.TypeOf(rightResult)), e.Position))
	}

	v.LastResult = rightBool
}

func (v *CodeVisitor) VisitAssignement(e *ast.Assignemnt) {
	e.Value.Accept(v)
	value := v.LastResult
	err := v.CurrentScope.SetVariableValue(e.Identifier.Name, value)
	if err != nil {
		panic(err)
	}
	// not to propagate the value further
	v.LastResult = nil
}

func (v *CodeVisitor) VisitVariable(e *ast.Variable) {
	err := v.CurrentScope.AddVariable(e.Name, e.Value, e.Type, e.Position)
	if err != nil {
		panic(err)
	}
}

func (v *CodeVisitor) VisitBlock(e *ast.Block) {
	for _, statement := range e.Statements {
		statement.Accept(v)
		if v.ReturnFlag {
			break
		}
	}
	// if v.ReturnFlag {
	// 	v.returnToFunctionDefScope()
	// }
	// parentScope, err := v.ScopeStack.Pop()
	// if err != nil {
	// 	panic(err)
	// }
	// v.CurrentScope = parentScope
}

func (v *CodeVisitor) VisitIfStatement(e *ast.IfStatement) {
	newScope := NewScope(v.CurrentScope, nil)
	v.ScopeStack.Push(newScope)
	v.CurrentScope = newScope

	e.Condition.Accept(v)
	conditionResult, ok := v.LastResult.(bool)
	if !ok {
		panic(fmt.Sprintf("Expected boolean expression but got %v", reflect.TypeOf(v.LastResult)))
	}

	if conditionResult {
		e.InstructionsBlock.Accept(v)
	} else if e.ElseInstructionsBlock != nil {
		e.ElseInstructionsBlock.Accept(v)
	}

	if v.ReturnFlag {
		v.returnToFunctionDefScope()
	} else {
		v.LastResult = nil
		currentScopre, err := v.ScopeStack.Pop()
		if err != nil {
			panic(err)
		}
		v.CurrentScope = currentScopre
	}
}

func (v *CodeVisitor) determineType(value interface{}) shared.TypeAnnotation {
	switch value.(type) {
	case int:
		return shared.INT
	case float64:
		return shared.FLOAT
	case bool:
		return shared.BOOL
	case string:
		return shared.STRING
	default:
		return shared.VOID
	}
}

func (v *CodeVisitor) VisitReturnStatement(e *ast.ReturnStatement) {
	// TODO: is this enough?
	if e.Expression != nil {
		e.Expression.Accept(v)
	} else {
		v.LastResult = nil
	}

	expectedReturnType := v.getCurrentFunctionReturnType()
	// expectedReturnType := shared.INT
	actualReturnType := v.determineType(v.LastResult)

	if expectedReturnType != actualReturnType {
		panic(fmt.Sprintf(INVALID_RETURN_TYPE, actualReturnType, expectedReturnType))
	}
	v.ReturnFlag = true
}

func (v *CodeVisitor) VisitWhileStatement(e *ast.WhileStatement) {
	newScope := NewScope(v.CurrentScope, nil)
	v.ScopeStack.Push(newScope)
	v.CurrentScope = newScope

	e.Condition.Accept(v)
	for v.LastResult.(bool) {
		e.InstructionsBlock.Accept(v)
		if v.ReturnFlag {
			break
		}
	}

	if v.ReturnFlag {
		v.returnToFunctionDefScope()
	} else {
		v.LastResult = nil
	}
	// parentScope, err := v.ScopeStack.Pop()
	// if err != nil {
	// 	panic(err)
	// }
	// v.CurrentScope = parentScope
}

func (v *CodeVisitor) VisitSwitchStatement(e *ast.SwitchStatement) {
}

func (v *CodeVisitor) VisitSwitchCase(e *ast.SwitchCase) {
}

func (v *CodeVisitor) VisitDefaultSwitchCase(e *ast.DefaultSwitchCase) {
}

func (v *CodeVisitor) VisitFunctionCall(fc *ast.FunctionCall) {
	functionDef := v.Functions[fc.Name]
	if functionDef == nil {
		panic(NewSemanticError(fmt.Sprintf(UNDEFINED_FUNCTION, fc.Name), fc.Position))
	}

	newScope := NewScope(v.CurrentScope, &functionDef.Type)

	if len(fc.Arguments) != len(functionDef.Parameters) {
		panic(NewSemanticError(fmt.Sprintf(WRONG_NUMBER_OF_ARGUMENTS, fc.Name, len(functionDef.Parameters), len(fc.Arguments)), fc.Position))
	}

	values := []ast.Expression{}
	for _, arg := range fc.Arguments {
		arg.Accept(v)
		values = append(values, v.LastResult.(ast.Expression))
	}

	v.ScopeStack.Push(v.CurrentScope)
	v.CurrentScope = newScope

	for i, param := range functionDef.Parameters {
		param.Accept(v)
		v.CurrentScope.SetVariableValue(param.Name, values[i])
	}

	functionDef.Accept(v)
}

func (v *CodeVisitor) VisitFunDef(e *ast.FunDef) {
}

func (v *CodeVisitor) VisitProgram(e *ast.Program) {
}

// VisitFunctionCall(*FunctionCall)
// VisitBlock(*Block)
// VisitIfStatement(*IfStatement)
// VisitReturnStatement(*ReturnStatement)
// VisitSwitchStatement(*SwitchStatement)
// VisitSwitchCase(*SwitchCase)
// VisitDefaultSwitchCase(*DefaultSwitchCase)
// VisitWhileStatement(*WhileStatement)
// VisitFunDef(*FunDef)
// VisitProgram(*Program)
