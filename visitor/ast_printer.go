package visitor

import (
	"fmt"
	"tkom/ast"
)

type ASTPrinter struct{}

func (a *ASTPrinter) VisitIntExpression(e *ast.IntExpression) {
	fmt.Print(e.Value)
}

func (a *ASTPrinter) VisitFloatExpression(e *ast.FloatExpression) {
	fmt.Print(e.Value)
}

func (a *ASTPrinter) VisitStringExpression(e *ast.StringExpression) {
	fmt.Print(e.Value)
}

func (a *ASTPrinter) VisitBoolExpression(e *ast.BoolExpression) {
	fmt.Print(e.Value)
}

func (a *ASTPrinter) VisitIdentifier(e *ast.Identifier) {
	fmt.Print(e.Name)
}

func (a *ASTPrinter) VisitFunctionCall(e *ast.FunctionCall) {
    e.Accept(a)
	//     for _, e2 := range e.Arguments {
	//         e2.Accept(v Visitor)
	//     }
	// 	fmt.Printf("%v()", e.Identifier.Name)
}

func (a *ASTPrinter) VisitVariable(e *ast.Variable) {
}

func (a *ASTPrinter) VisitNegateExpression(e *ast.NegateExpression) {
}

func (a *ASTPrinter) VisitCastExpression(e *ast.CastExpression) {
}

func (a *ASTPrinter) VisitMultiplyExpression(e *ast.MultiplyExpression) {
}

func (a *ASTPrinter) VisitDivideExpression(e *ast.DivideExpression) {
}

func (a *ASTPrinter) VisitSumExpression(e *ast.SumExpression) {
}

func (a *ASTPrinter) VisitSubstractExpression(e *ast.SubstractExpression) {
}

func (a *ASTPrinter) VisitEqualsExpression(e *ast.EqualsExpression) {
}

func (a *ASTPrinter) VisitNotEqualsExpression(e *ast.NotEqualsExpression) {
}

func (a *ASTPrinter) VisitGreaterThanExpression(e *ast.GreaterThanExpression) {
}

func (a *ASTPrinter) VisitLessThanExpression(e *ast.LessThanExpression) {
}

func (a *ASTPrinter) VisitGreaterOrEqualExpression(e *ast.GreaterOrEqualExpression) {
}

func (a *ASTPrinter) VisitLessOrEqualExpression(e *ast.LessOrEqualExpression) {
}

func (a *ASTPrinter) VisitAndExpression(e *ast.AndExpression) {
}

func (a *ASTPrinter) VisitOrExpression(e *ast.OrExpression) {
}

func (a *ASTPrinter) VisitIfStatement(e *ast.IfStatement) {
}

func (a *ASTPrinter) VisitReturnStatement(e *ast.ReturnStatement) {
}

func (a *ASTPrinter) VisitSwitchStatement(e *ast.SwitchStatement) {
}

func (a *ASTPrinter) VisitSwitchCase(e *ast.SwitchCase) {
}

func (a *ASTPrinter) VisitWhileStatement(e *ast.WhileStatement) {
}

func (a *ASTPrinter) VisitFunDef(e *ast.FunDef) {
}

func (a *ASTPrinter) VisitProgram(e *ast.Program) {
}
