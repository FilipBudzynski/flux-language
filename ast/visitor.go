package ast

type Function interface {
	Node
	GetParametersLen() int
	IsVariadic() bool
}

type Visitor interface {
	VisitIntExpression(*IntExpression)
	VisitFloatExpression(*FloatExpression)
	VisitStringExpression(*StringExpression)
	VisitBoolExpression(*BoolExpression)
	VisitIdentifier(*Identifier)
	VisitFunctionCall(*FunctionCall)
	VisitVariable(*Variable)
	VisitAssignement(*Assignment)
	VisitNegateExpression(*NegateExpression)
	VisitCastExpression(*CastExpression)
	VisitMultiplyExpression(*MultiplyExpression)
	VisitDivideExpression(*DivideExpression)
	VisitSumExpression(*SumExpression)
	VisitSubstractExpression(*SubstractExpression)
	VisitEqualsExpression(*EqualsExpression)
	VisitNotEqualsExpression(*NotEqualsExpression)
	VisitGreaterThanExpression(*GreaterThanExpression)
	VisitLessThanExpression(*LessThanExpression)
	VisitGreaterOrEqualExpression(*GreaterOrEqualExpression)
	VisitLessOrEqualExpression(*LessOrEqualExpression)
	VisitAndExpression(*AndExpression)
	VisitOrExpression(*OrExpression)
	VisitBlock(*Block)
	VisitIfStatement(*IfStatement)
	VisitReturnStatement(*ReturnStatement)
	VisitSwitchStatement(*SwitchStatement)
	VisitSwitchCase(*SwitchCase)
	VisitDefaultSwitchCase(*DefaultSwitchCase)
	VisitWhileStatement(*WhileStatement)
	VisitFunctionDefinition(*FunctionDefinition)
	VisitProgram(*Program)
	VisitEmbeddedFunction(*EmbeddedFunction)
}
