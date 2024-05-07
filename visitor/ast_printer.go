package visitor

import (
	"fmt"
	"tkom/ast"
)

type ASTPrinter struct{}

func (a *ASTPrinter) visitIntExpression(e *ast.IntExpression) {
	fmt.Print(e.Value)
}

func (a *ASTPrinter) visitFloatExpression(e *ast.FloatExpression) {
	fmt.Print(e.Value)
}

func (a *ASTPrinter) visitStringExpression(e *ast.StringExpression) {
	fmt.Print(e.Value)
}

func (a *ASTPrinter) visitBoolExpression(e *ast.BoolExpression) {
	fmt.Print(e.Value)
}

func (a *ASTPrinter) visitIdentifier(e *ast.Identifier) {
	fmt.Print(e.Name)
}

func (a *ASTPrinter) visitFunctionCall(e *ast.FunctionCall) {
	//     for _, e2 := range e.Arguments {
	//         e2.Accept(v Visitor)
	//     }
	// 	fmt.Printf("%v()", e.Identifier.Name)
}

func (a *ASTPrinter) visitVariable(e *ast.Variable) {
}

func (a *ASTPrinter) visitNegateExpression(e *ast.NegateExpression) {
}

func (a *ASTPrinter) visitCastExpression(e *ast.CastExpression) {
}

func (a *ASTPrinter) visitMultiplyExpression(e *ast.MultiplyExpression) {
}

func (a *ASTPrinter) visitDivideExpression(e *ast.DivideExpression) {
}

func (a *ASTPrinter) visitSumExpression(e *ast.SumExpression) {
}

func (a *ASTPrinter) visitSubstractExpression(e *ast.SubstractExpression) {
}

func (a *ASTPrinter) visitEqualsExpression(e *ast.EqualsExpression) {
}

func (a *ASTPrinter) visitNotEqualsExpression(e *ast.NotEqualsExpression) {
}

func (a *ASTPrinter) visitGreaterThenExpression(e *ast.GreaterThanExpression) {
}

func (a *ASTPrinter) visitLessThenExpression(e *ast.LessThanExpression) {
}

func (a *ASTPrinter) visitGreaterOrEqualExpression(e *ast.GreaterOrEqualExpression) {
}

func (a *ASTPrinter) visitLessOrEqualExpression(e *ast.LessOrEqualExpression) {
}

func (a *ASTPrinter) visitAndExpression(e *ast.AndExpression) {
}

func (a *ASTPrinter) visitOrExpression(e *ast.OrExpression) {
}

func (a *ASTPrinter) visitIfStatement(e *ast.IfStatement) {
}

func (a *ASTPrinter) visitReturnStatement(e *ast.ReturnStatement) {
}

func (a *ASTPrinter) visitSwitchStatement(e *ast.SwitchStatement) {
}

func (a *ASTPrinter) visitSwitchCase(e *ast.SwitchCase) {
}

func (a *ASTPrinter) visitWhileStatement(e *ast.WhileStatement) {
}

func (a *ASTPrinter) visitFunDef(e *ast.FunDef) {
}

func (a *ASTPrinter) visitProgram(e *ast.Program) {
}
