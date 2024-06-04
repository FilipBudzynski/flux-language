package interpreter

import (
	"fmt"
	"reflect"
	"strconv"
	"tkom/ast"
	"tkom/shared"
)

// implementation of interpreter
// based on the visitor pattern
//
// https://en.wikipedia.org/wiki/Visitor_pattern
type CodeVisitor struct {
	LastResult        any
	FunctionsMap      map[string]ast.Function
	CurrentScope      *Scope
	CallStack         CallStack
	ScopeStack        Stack
	ReturnFlag        bool
	SwitchEndFlag     bool
	CurrentReturnType shared.TypeAnnotation
	MaxRecursionDepth int
}

func NewCodeVisitor(maxRecursionDepth int) *CodeVisitor {
	return &CodeVisitor{
		FunctionsMap:      embeddedFunctions,
		ScopeStack:        Stack{elem: []*Scope{}},
		CallStack:         CallStack{elem: map[string]int{}},
		LastResult:        nil,
		CurrentScope:      nil,
		ReturnFlag:        false,
		SwitchEndFlag:     false,
		CurrentReturnType: shared.VOID,
		MaxRecursionDepth: maxRecursionDepth,
	}
}

// helper function for visiting CastExpression
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

// helper function for visiting CastExpression
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

// helper function for visiting CastExpression
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

// helper function for visiting CastExpression
func (v *CodeVisitor) tryCastToString(value any) (string, error) {
	return fmt.Sprintf("%v", value), nil
}

func (v *CodeVisitor) checkType(value any, expectedType shared.TypeAnnotation, pos shared.Position) error {
	switch expectedType {
	case shared.INT:
		if _, ok := value.(int); !ok {
			return NewSemanticError(fmt.Sprintf(TYPE_MISMATCH, expectedType, reflect.TypeOf(value)), pos)
		}
	case shared.FLOAT:
		if _, ok := value.(float64); !ok {
			return NewSemanticError(fmt.Sprintf(TYPE_MISMATCH, expectedType, reflect.TypeOf(value)), pos)
		}
	case shared.BOOL:
		if _, ok := value.(bool); !ok {
			return NewSemanticError(fmt.Sprintf(TYPE_MISMATCH, expectedType, reflect.TypeOf(value)), pos)
		}
	case shared.STRING:
		if _, ok := value.(string); !ok {
			return NewSemanticError(fmt.Sprintf(TYPE_MISMATCH, expectedType, reflect.TypeOf(value)), pos)
		}
	default:
		return fmt.Errorf("unknown type: %v", expectedType)
	}
	return nil
}

// helper function for getting the return type of the current function
func (v *CodeVisitor) getCurrentFunctionReturnType() shared.TypeAnnotation {
	for i := len(v.ScopeStack.elem) - 1; i >= 0; i-- {
		if v.ScopeStack.elem[i].ReturnType != nil {
			return *v.ScopeStack.elem[i].ReturnType
		}
	}
	return shared.VOID
}

func (v *CodeVisitor) VisitIntExpression(intExp *ast.IntExpression) {
	v.LastResult = intExp.Value
}

func (v *CodeVisitor) VisitFloatExpression(floatExp *ast.FloatExpression) {
	v.LastResult = floatExp.Value
}

func (v *CodeVisitor) VisitStringExpression(strExp *ast.StringExpression) {
	v.LastResult = strExp.Value
}

func (v *CodeVisitor) VisitBoolExpression(boolExp *ast.BoolExpression) {
	v.LastResult = boolExp.Value
}

func (v *CodeVisitor) VisitIdentifier(idExp *ast.Identifier) {
	sc, err := v.CurrentScope.GetVariable(idExp.Name)
	if err != nil {
		panic(NewSemanticError(err.Error(), idExp.Position))
	}
	v.LastResult = sc.Value
}

func (v *CodeVisitor) VisitNegateExpression(negateExp *ast.NegateExpression) {
	negateExp.Expression.Accept(v)
	ne := v.LastResult

	switch value := ne.(type) {
	case int:
		v.LastResult = -value
	case float64:
		v.LastResult = -value
	case bool:
		v.LastResult = !value
	default:
		panic(NewSemanticError(fmt.Sprintf(INVALID_NEGATE_EXPRESSION, ne, "string"), negateExp.Position))
	}
}

func (v *CodeVisitor) VisitCastExpression(castExp *ast.CastExpression) {
	castExp.LeftExpression.Accept(v)
	leftExpValue := v.LastResult

	var result any
	var err error

	switch castExp.TypeAnnotation {
	case shared.INT:
		result, err = v.tryCastToInt(leftExpValue)
	case shared.FLOAT:
		result, err = v.tryCastToFloat(leftExpValue)
	case shared.BOOL:
		result, err = v.tryCastToBool(leftExpValue)
	case shared.STRING:
		result, err = v.tryCastToString(leftExpValue)
	default:
		panic(NewSemanticError(fmt.Sprintf(INVALID_TYPE_ANNOTATION, castExp.TypeAnnotation), castExp.Position))
	}

	if err != nil {
		panic(NewSemanticError(err.Error(), castExp.Position))
	}

	v.LastResult = result
}

func (v *CodeVisitor) VisitMultiplyExpression(mulExp *ast.MultiplyExpression) {
	mulExp.LeftExpression.Accept(v)
	leftResult := v.LastResult

	mulExp.RightExpression.Accept(v)
	rightResult := v.LastResult

	switch leftVal := leftResult.(type) {
	case int:
		switch rigthVal := rightResult.(type) {
		case int:
			v.LastResult = leftVal * rigthVal

		case float64:
			v.LastResult = float64(leftVal) * rigthVal
		}
	case float64:
		switch rightVal := rightResult.(type) {
		case float64:
			v.LastResult = leftVal * rightVal
		case int:
			v.LastResult = leftVal * float64(rightVal)
		}
	default:
		panic(NewSemanticError(fmt.Sprintf(INVALID_MULTIPLY_EXPRESSION, leftResult, rightResult), mulExp.Position))
	}
}

func (v *CodeVisitor) VisitDivideExpression(divExp *ast.DivideExpression) {
	divExp.LeftExpression.Accept(v)
	leftResult := v.LastResult

	divExp.RightExpression.Accept(v)
	rightResult := v.LastResult

	// check for zero as
	if right, ok := rightResult.(int); ok && right == 0 {
		panic(NewSemanticError("Division by zero", divExp.Position))
	} else if right, ok := rightResult.(float64); ok && right == 0.0 {
		panic(NewSemanticError("Division by zero", divExp.Position))
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
			v.LastResult = leftResult.(float64) / rightResult.(float64)
		case int:
			v.LastResult = leftResult.(float64) / float64(rightResult.(int))
		}
	default:
		panic(NewSemanticError(fmt.Sprintf(INVALID_DIVISION_EXPRESSION, leftResult, rightResult), divExp.Position))
	}
}

func (v *CodeVisitor) VisitSumExpression(sumExp *ast.SumExpression) {
	sumExp.LeftExpression.Accept(v)
	leftResult := v.LastResult

	sumExp.RightExpression.Accept(v)
	rightResult := v.LastResult

	var result any
	var err error

	switch left := leftResult.(type) {
	case int:
		result, err = sumInt(left, rightResult)
	case float64:
		result, err = sumFloat64(left, rightResult)
	case string:
		result, err = sumString(left, rightResult)
	default:
		err = fmt.Errorf("invalid left operand type: %T", leftResult)
	}

	if err != nil {
		panic(NewSemanticError(fmt.Sprintf(INVALID_SUM_EXPRESSION, leftResult, rightResult), sumExp.Position))
	}

	v.LastResult = result
}

// helper function for the sum expression function
func sumInt(left int, right any) (any, error) {
	switch right := right.(type) {
	case int:
		return left + right, nil
	case float64:
		return float64(left) + right, nil
	case string:
		return fmt.Sprintf("%d%s", left, right), nil
	default:
		return nil, fmt.Errorf("invalid right operand type: %T", right)
	}
}

// helper function for the sum expression function
func sumFloat64(left float64, right any) (any, error) {
	switch right := right.(type) {
	case int:
		return left + float64(right), nil
	case float64:
		return left + right, nil
	case string:
		return fmt.Sprintf("%f%s", left, right), nil
	default:
		return nil, fmt.Errorf("invalid right operand type: %T", right)
	}
}

// helper function for the sum expression function
func sumString(left string, right any) (any, error) {
	switch right := right.(type) {
	case int:
		return left + fmt.Sprintf("%d", right), nil
	case float64:
		return left + fmt.Sprintf("%f", right), nil
	case string:
		return left + right, nil
	default:
		return nil, fmt.Errorf("invalid right operand type: %T", right)
	}
}

func (v *CodeVisitor) VisitSubstractExpression(subExp *ast.SubstractExpression) {
	subExp.LeftExpression.Accept(v)
	leftResult := v.LastResult

	subExp.RightExpression.Accept(v)
	rightResult := v.LastResult

	switch leftVal := leftResult.(type) {
	case int:
		switch rightVal := rightResult.(type) {
		case int:
			v.LastResult = leftVal - rightVal
		}
	case float64:
		switch rightVal := rightResult.(type) {
		case float64:
			v.LastResult = leftVal - rightVal
		case int:
			v.LastResult = leftVal - float64(rightVal)
		}
	default:
		panic(NewSemanticError(fmt.Sprintf(INVALID_SUBSTRACT_EXPRESSION, leftResult, rightResult), subExp.Position))
	}
}

func (v *CodeVisitor) VisitEqualsExpression(eqExp *ast.EqualsExpression) {
	eqExp.LeftExpression.Accept(v)
	leftResult := v.LastResult

	eqExp.RightExpression.Accept(v)
	rightResult := v.LastResult

	if reflect.TypeOf(leftResult) != reflect.TypeOf(rightResult) {
		panic(NewSemanticError(fmt.Sprintf(INVALID_EQUALS_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), eqExp.Position))
	}

	v.LastResult = leftResult == rightResult
}

func (v *CodeVisitor) VisitNotEqualsExpression(neExp *ast.NotEqualsExpression) {
	neExp.LeftExpression.Accept(v)
	leftResult := v.LastResult

	neExp.RightExpression.Accept(v)
	rightResult := v.LastResult

	if reflect.TypeOf(leftResult) != reflect.TypeOf(rightResult) {
		panic(NewSemanticError(fmt.Sprintf(INVALID_NOT_EQUALS_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), neExp.Position))
	}

	v.LastResult = leftResult != rightResult
}

func (v *CodeVisitor) VisitGreaterThanExpression(gtExp *ast.GreaterThanExpression) {
	gtExp.LeftExpression.Accept(v)
	leftResult := v.LastResult

	gtExp.RightExpression.Accept(v)
	rightResult := v.LastResult

	if reflect.TypeOf(leftResult) != reflect.TypeOf(rightResult) {
		panic(NewSemanticError(fmt.Sprintf(INVALID_GREATER_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), gtExp.Position))
	}

	switch left := leftResult.(type) {
	case int:
		right, ok := rightResult.(int)
		if !ok {
			panic(NewSemanticError(fmt.Sprintf(INVALID_GREATER_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), gtExp.Position))
		}
		v.LastResult = left > right
	case float64:
		right, ok := rightResult.(float64)
		if !ok {
			panic(NewSemanticError(fmt.Sprintf(INVALID_GREATER_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), gtExp.Position))
		}
		v.LastResult = left > right
	default:
		panic(NewSemanticError(fmt.Sprintf(INVALID_GREATER_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), gtExp.Position))
	}
}

func (v *CodeVisitor) VisitGreaterOrEqualExpression(geExp *ast.GreaterOrEqualExpression) {
	geExp.LeftExpression.Accept(v)
	leftResult := v.LastResult

	geExp.RightExpression.Accept(v)
	rightResult := v.LastResult

	if reflect.TypeOf(leftResult) != reflect.TypeOf(rightResult) {
		panic(NewSemanticError(fmt.Sprintf(INVALID_GREATER_OR_EQUALS_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), geExp.Position))
	}

	switch left := leftResult.(type) {
	case int:
		right, ok := rightResult.(int)
		if !ok {
			panic(NewSemanticError(fmt.Sprintf(INVALID_GREATER_OR_EQUALS_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), geExp.Position))
		}
		v.LastResult = left >= right
	case float64:
		right, ok := rightResult.(float64)
		if !ok {
			panic(NewSemanticError(fmt.Sprintf(INVALID_GREATER_OR_EQUALS_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), geExp.Position))
		}
		v.LastResult = left >= right
	default:
		panic(NewSemanticError(fmt.Sprintf(INVALID_GREATER_OR_EQUALS_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), geExp.Position))
	}
}

func (v *CodeVisitor) VisitLessThanExpression(ltExp *ast.LessThanExpression) {
	ltExp.LeftExpression.Accept(v)
	leftResult := v.LastResult

	ltExp.RightExpression.Accept(v)
	rightResult := v.LastResult

	if reflect.TypeOf(leftResult) != reflect.TypeOf(rightResult) {
		panic(NewSemanticError(fmt.Sprintf(INVALID_LESS_OR_EQUALS_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), ltExp.Position))
	}

	switch left := leftResult.(type) {
	case int:
		right, ok := rightResult.(int)
		if !ok {
			panic(NewSemanticError(fmt.Sprintf(INVALID_LESS_OR_EQUALS_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), ltExp.Position))
		}
		v.LastResult = left < right
	case float64:
		right, ok := rightResult.(float64)
		if !ok {
			panic(NewSemanticError(fmt.Sprintf(INVALID_LESS_OR_EQUALS_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), ltExp.Position))
		}
		v.LastResult = left < right
	default:
		panic(NewSemanticError(fmt.Sprintf(INVALID_LESS_OR_EQUALS_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), ltExp.Position))
	}
}

func (v *CodeVisitor) VisitLessOrEqualExpression(leExp *ast.LessOrEqualExpression) {
	leExp.LeftExpression.Accept(v)
	leftResult := v.LastResult

	leExp.RightExpression.Accept(v)
	rightResult := v.LastResult

	if reflect.TypeOf(leftResult) != reflect.TypeOf(rightResult) {
		panic(NewSemanticError(fmt.Sprintf(INVALID_GREATER_OR_EQUALS_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), leExp.Position))
	}

	switch left := leftResult.(type) {
	case int:
		right, ok := rightResult.(int)
		if !ok {
			panic(NewSemanticError(fmt.Sprintf(INVALID_GREATER_OR_EQUALS_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), leExp.Position))
		}
		v.LastResult = left <= right
	case float64:
		right, ok := rightResult.(float64)
		if !ok {
			panic(NewSemanticError(fmt.Sprintf(INVALID_GREATER_OR_EQUALS_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), leExp.Position))
		}
		v.LastResult = left <= right
	default:
		panic(NewSemanticError(fmt.Sprintf(INVALID_GREATER_OR_EQUALS_THAN_MISSMATCH, reflect.TypeOf(leftResult), reflect.TypeOf(rightResult)), leExp.Position))
	}
}

func (v *CodeVisitor) VisitOrExpression(orExp *ast.OrExpression) {
	orExp.LeftExpression.Accept(v)
	leftResult := v.LastResult

	leftBool, ok := leftResult.(bool)
	if !ok {
		panic(NewSemanticError(fmt.Sprintf(EXPECTED_BOOLEAN_EXPRESSION, reflect.TypeOf(leftResult)), orExp.Position))
	}

	// If the left expression is true, return true
	if leftBool {
		v.LastResult = true
		return
	}

	orExp.RightExpression.Accept(v)
	rightResult := v.LastResult

	// Check if the right result is a boolean
	rightBool, ok := rightResult.(bool)
	if !ok {
		panic(NewSemanticError(fmt.Sprintf(EXPECTED_BOOLEAN_EXPRESSION, reflect.TypeOf(rightResult)), orExp.Position))
	}

	v.LastResult = rightBool
}

func (v *CodeVisitor) VisitAndExpression(andExp *ast.AndExpression) {
	andExp.LeftExpression.Accept(v)
	leftResult := v.LastResult

	leftBool, ok := leftResult.(bool)
	if !ok {
		panic(NewSemanticError(fmt.Sprintf(EXPECTED_BOOLEAN_EXPRESSION, reflect.TypeOf(leftResult)), andExp.Position))
	}

	// If the left expression is false, return false
	if !leftBool {
		v.LastResult = false
		return
	}

	andExp.RightExpression.Accept(v)
	rightResult := v.LastResult

	rightBool, ok := rightResult.(bool)
	if !ok {
		panic(NewSemanticError(fmt.Sprintf(EXPECTED_BOOLEAN_EXPRESSION, reflect.TypeOf(rightResult)), andExp.Position))
	}

	v.LastResult = rightBool
}

func (v *CodeVisitor) VisitAssignement(assignment *ast.Assignment) {
	assignment.Value.Accept(v)
	value := v.LastResult

	err := v.CurrentScope.SetValue(assignment.Identifier.Name, value)
	if err != nil {
		panic(err)
	}

	v.LastResult = nil
}

func (v *CodeVisitor) VisitVariable(varDecl *ast.Variable) {
	varDecl.Value.Accept(v)
	value := v.LastResult

	err := v.checkType(value, varDecl.Type, varDecl.Position)
	if err != nil {
		panic(err)
	}

	err = v.CurrentScope.AddVariable(varDecl.Name, v.LastResult, varDecl.Type, varDecl.Position)
	if err != nil {
		panic(err)
	}

	v.LastResult = nil
}

func (v *CodeVisitor) VisitBlock(block *ast.Block) {
	for _, statement := range block.Statements {
		statement.Accept(v)
		if v.ReturnFlag {
			break
		}
	}
}

func (v *CodeVisitor) VisitIfStatement(ifStmt *ast.IfStatement) {
	newScope := NewScope(v.CurrentScope, nil)
	v.ScopeStack.Push(newScope)
	v.CurrentScope = newScope

	ifStmt.Condition.Accept(v)
	conditionResult, ok := v.LastResult.(bool)
	if !ok {
		panic(NewSemanticError(fmt.Sprintf(EXPECTED_BOOLEAN_EXPRESSION, reflect.TypeOf(v.LastResult)), ifStmt.Condition.GetPosition()))
	}

	if conditionResult {
		ifStmt.InstructionsBlock.Accept(v)
	} else if ifStmt.ElseInstructionsBlock != nil {
		ifStmt.ElseInstructionsBlock.Accept(v)
	}

	currentScope, err := v.ScopeStack.Pop()
	if err != nil {
		panic(err)
	}
	v.CurrentScope = currentScope.Parent

	if !v.ReturnFlag {
		v.LastResult = nil
	}
}

func (v *CodeVisitor) VisitReturnStatement(returnStmt *ast.ReturnStatement) {
	if returnStmt.Value != nil {
		returnStmt.Value.Accept(v)
	} else {
		v.LastResult = nil
	}

	expectedReturnType := v.getCurrentFunctionReturnType()
	actualReturnType := v.determineType(v.LastResult)

	if expectedReturnType != actualReturnType {
		panic(NewSemanticError(fmt.Sprintf(INVALID_RETURN_TYPE, actualReturnType, expectedReturnType), returnStmt.Value.GetPosition()))
	}
	v.ReturnFlag = true
}

// helper function for determining the type of a value for return
func (v *CodeVisitor) determineType(value any) shared.TypeAnnotation {
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

func (v *CodeVisitor) VisitWhileStatement(whileStmt *ast.WhileStatement) {
	newScope := NewScope(v.CurrentScope, nil)
	v.ScopeStack.Push(newScope)
	v.CurrentScope = newScope

	whileStmt.Condition.Accept(v)

	if ok := v.LastResult.(bool); !ok {
		panic(NewSemanticError(fmt.Sprintf(INVALID_WHILE_CONDITION, reflect.TypeOf(v.LastResult)), shared.NewPosition(0, 0)))
	}

	for v.LastResult.(bool) {
		whileStmt.InstructionsBlock.Accept(v)
		if v.ReturnFlag {
			break
		}
		whileStmt.Condition.Accept(v)
	}

	currentScope, err := v.ScopeStack.Pop()
	if err != nil {
		panic(err)
	}
	v.CurrentScope = currentScope.Parent

	if !v.ReturnFlag {
		// clear the last result if no return flag
		v.LastResult = nil
	}
}

func (v *CodeVisitor) VisitFunctionCall(fc *ast.FunctionCall) {
	functionDef := v.FunctionsMap[fc.Name]
	if functionDef == nil {
		panic(NewSemanticError(fmt.Sprintf(UNDEFINED_FUNCTION, fc.Name), fc.Position))
	}

	if v.CallStack.RecursionDepth(fc.Name) >= v.MaxRecursionDepth {
		panic(NewSemanticError(fmt.Sprintf(MAX_RECURSION_DEPTH_EXCEEDED, fc.Name), fc.Position))
	}

	v.CallStack.Push(fc.Name)
	defer v.CallStack.Pop(fc.Name)

	values := []any{}
	for _, arg := range fc.Arguments {
		arg.Accept(v)
		values = append(values, v.LastResult)
	}
	v.LastResult = values

	functionDef.Accept(v)
}

func (v *CodeVisitor) VisitFunctionDefinition(fd *ast.FunctionDefinition) {
	if _, ok := v.LastResult.([]any); !ok {
		panic(NewSemanticError(fmt.Sprintf(ERROR_ARGUMENTS_NOT_FOUND, reflect.TypeOf(v.LastResult)), fd.Position))
	}
	if len(v.LastResult.([]any)) != len(fd.Parameters) {
		panic(NewSemanticError(fmt.Sprintf(WRONG_NUMBER_OF_ARGUMENTS, fd.Name, len(fd.Parameters), len(v.LastResult.([]any))), fd.Position))
	}
	args := v.LastResult.([]any)

	newScope := NewScope(v.CurrentScope, &fd.Type)
	v.ScopeStack.Push(newScope)
	v.CurrentScope = newScope

	for i, param := range fd.Parameters {

		argumentValue := args[i]
		err := v.checkType(argumentValue, param.Type, param.Position)
		if err != nil {
			panic(err)
		}
		err = v.CurrentScope.AddVariable(param.Name, argumentValue, param.Type, param.Position)
		if err != nil {
			panic(err)
		}
	}

	fd.Block.Accept(v)

	if fd.Type != shared.VOID && !v.ReturnFlag {
		panic(NewSemanticError(fmt.Sprintf(MISSING_RETURN, fd.Type), fd.Position))
	}

	currScope, err := v.ScopeStack.Pop()
	if err != nil {
		panic(err)
	}
	v.CurrentScope = currScope.Parent

	if v.ReturnFlag {
		v.ReturnFlag = false
	} else {
		v.LastResult = nil
	}
}

func (v *CodeVisitor) VisitEmbeddedFunction(ef *ast.EmbeddedFunction) {
	if args, ok := v.LastResult.([]any); ok {
		if !ef.Variadic && len(args) != len(ef.Parameters) {
			panic(fmt.Errorf(WRONG_NUMBER_OF_ARGUMENTS, ef.Name, len(ef.Parameters), len(args)))
		}
		result := ef.Func(args...)
		v.LastResult = result
	} else {
		panic(fmt.Errorf(INVALID_ARGUMENTS_TYPE, reflect.TypeOf(v.LastResult)))
	}
}

func (v *CodeVisitor) VisitSwitchStatement(s *ast.SwitchStatement) {
	newScope := NewScope(v.CurrentScope, nil)
	v.ScopeStack.Push(newScope)
	v.CurrentScope = newScope

	for _, variable := range s.Variables {
		variable.Accept(v)
	}

	var defaultCase *ast.DefaultSwitchCase

	for _, c := range s.Cases {
		switch caseStmt := c.(type) {
		case *ast.SwitchCase:
			caseStmt.Accept(v)
			if v.ReturnFlag || v.SwitchEndFlag {
				break
			}
		case *ast.DefaultSwitchCase:
			if defaultCase != nil {
				panic(NewSemanticError(MULTIPLE_DEFAULT_CASES, defaultCase.GetPosition()))
			}
			defaultCase = caseStmt
		default:
			panic(NewSemanticError(INVALID_CASE_TYPE, caseStmt.GetPosition()))
		}
		if v.ReturnFlag {
			break
		}
	}

	// run default only after cases did not get executed
	if !v.ReturnFlag && defaultCase != nil && !v.SwitchEndFlag {
		defaultCase.Accept(v)
	}

	currentScope, err := v.ScopeStack.Pop()
	if err != nil {
		panic(err)
	}
	v.CurrentScope = currentScope.Parent

	// clear the switch flag
	v.SwitchEndFlag = false

	// clear the lastResult if flag wasnt set
	if !v.ReturnFlag {
		v.LastResult = nil
	}
}

func (v *CodeVisitor) VisitSwitchCase(sc *ast.SwitchCase) {
	sc.Condition.Accept(v)
	condition := v.LastResult

	if condition.(bool) {
		sc.OutputExpression.Accept(v)

		if v.LastResult != nil {
			v.ReturnFlag = true
			return
		} else {
			v.SwitchEndFlag = true
			return
		}
	}

	v.LastResult = nil
}

func (v *CodeVisitor) VisitDefaultSwitchCase(dsc *ast.DefaultSwitchCase) {
	dsc.OutputExpression.Accept(v)
	if v.LastResult != nil {
		v.ReturnFlag = true
		return
	} else {
		v.SwitchEndFlag = true
		return
	}
}

func (v *CodeVisitor) VisitProgram(e *ast.Program) {}
