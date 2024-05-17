package ast

type Node interface {
    Accept(Visitor)
}

type Visitor interface {
	VisitIntExpression(*IntExpression)
	VisitFloatExpression(*FloatExpression)
	VisitStringExpression(*StringExpression)
	VisitBoolExpression(*BoolExpression)
	VisitIdentifier(*Identifier)
	VisitFunctionCall(*FunctionCall)
	VisitVariable(*Variable)
	VisitAssignement(*Assignemnt)
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
	VisitFunDef(*FunDef)
	VisitProgram(*Program)
}
