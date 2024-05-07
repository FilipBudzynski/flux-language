package visitor

import (
	"tkom/ast"
)

type Visitor interface {
	visitIntExpression(*ast.IntExpression)
	visitFloatExpression(*ast.FloatExpression)
	visitStringExpression(*ast.StringExpression)
	visitBoolExpression(*ast.BoolExpression)
	visitIdentifier(*ast.Identifier)
	visitFunctionCall(*ast.FunctionCall)
	visitVariable(*ast.Variable)
	visitNegateExpression(*ast.NegateExpression)
	visitCastExpression(*ast.CastExpression)
	visitMultiplyExpression(*ast.MultiplyExpression)
	visitDivideExpression(*ast.DivideExpression)
	visitSumExpression(*ast.SumExpression)
	visitSubstractExpression(*ast.SubstractExpression)
	visitEqualsExpression(*ast.EqualsExpression)
	visitNotEqualsExpression(*ast.NotEqualsExpression)
	visitGreaterThenExpression(*ast.GreaterThanExpression)
	visitLessThenExpression(*ast.LessThanExpression)
	visitGreaterOrEqualExpression(*ast.GreaterOrEqualExpression)
	visitLessOrEqualExpression(*ast.LessOrEqualExpression)
	visitAndExpression(*ast.AndExpression)
	visitOrExpression(*ast.OrExpression)
	visitIfStatement(*ast.IfStatement)
	visitReturnStatement(*ast.ReturnStatement)
	visitSwitchStatement(*ast.SwitchStatement)
	visitSwitchCase(*ast.SwitchCase)
	visitWhileStatement(*ast.WhileStatement)
	visitFunDef(*ast.FunDef)
	visitProgram(*ast.Program)
}
