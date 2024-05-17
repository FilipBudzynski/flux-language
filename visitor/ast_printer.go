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

func (a *ASTPrinter) VisitAssignement(e *ast.Assignemnt) {
}

func (a *ASTPrinter) VisitFunctionCall(e *ast.FunctionCall) {
	for _, arg := range e.Arguments {
		arg.Accept(a)
	}
	fmt.Println(e.Name)
	//     for _, e2 := range e.Arguments {
	//         e2.Accept(v Visitor)
	//     }
	// 	fmt.Printf("%v()", e.Identifier.Name)
}

func (a *ASTPrinter) VisitVariable(e *ast.Variable) {
	fmt.Printf("%v %v\n", e.Type, e.Name)
}

func (a *ASTPrinter) VisitNegateExpression(e *ast.NegateExpression) {
	e.Expression.Accept(a)
	fmt.Printf("Negate")
}

func (a *ASTPrinter) VisitCastExpression(e *ast.CastExpression) {
	e.LeftExpression.Accept(a)
	fmt.Printf("cast to %v\n", e.TypeAnnotation)
}

func (a *ASTPrinter) VisitMultiplyExpression(e *ast.MultiplyExpression) {
	e.LeftExpression.Accept(a)
	fmt.Printf(" * ")
	e.RightExpression.Accept(a)
	fmt.Printf("\n")
}

func (a *ASTPrinter) VisitDivideExpression(e *ast.DivideExpression) {
	e.LeftExpression.Accept(a)
	fmt.Printf(" / ")
	e.RightExpression.Accept(a)
	fmt.Printf("\n")
}

func (a *ASTPrinter) VisitSumExpression(e *ast.SumExpression) {
	e.LeftExpression.Accept(a)
	fmt.Printf(" + ")
	e.RightExpression.Accept(a)
	fmt.Printf("\n")
}

func (a *ASTPrinter) VisitSubstractExpression(e *ast.SubstractExpression) {
	e.LeftExpression.Accept(a)
	fmt.Printf(" - ")
	e.RightExpression.Accept(a)
	fmt.Printf("\n")
}

func (a *ASTPrinter) VisitEqualsExpression(e *ast.EqualsExpression) {
	e.LeftExpression.Accept(a)
	fmt.Printf(" == ")
	e.RightExpression.Accept(a)
	fmt.Printf("\n")
}

func (a *ASTPrinter) VisitNotEqualsExpression(e *ast.NotEqualsExpression) {
	e.LeftExpression.Accept(a)
	fmt.Printf(" != ")
	e.RightExpression.Accept(a)
	fmt.Printf("\n")
}

func (a *ASTPrinter) VisitGreaterThanExpression(e *ast.GreaterThanExpression) {
	e.LeftExpression.Accept(a)
	fmt.Printf(" > ")
	e.RightExpression.Accept(a)
	fmt.Printf("\n")
}

func (a *ASTPrinter) VisitLessThanExpression(e *ast.LessThanExpression) {
	e.LeftExpression.Accept(a)
	fmt.Printf(" < ")
	e.RightExpression.Accept(a)
	fmt.Printf("\n")
}

func (a *ASTPrinter) VisitGreaterOrEqualExpression(e *ast.GreaterOrEqualExpression) {
	e.LeftExpression.Accept(a)
	fmt.Printf(" >= ")
	e.RightExpression.Accept(a)
	fmt.Printf("\n")
}

func (a *ASTPrinter) VisitLessOrEqualExpression(e *ast.LessOrEqualExpression) {
	e.LeftExpression.Accept(a)
	fmt.Printf(" <= ")
	e.RightExpression.Accept(a)
	fmt.Printf("\n")
}

func (a *ASTPrinter) VisitAndExpression(e *ast.AndExpression) {
	e.LeftExpression.Accept(a)
	fmt.Printf(" and ")
	e.RightExpression.Accept(a)
	fmt.Printf("\n")
}

func (a *ASTPrinter) VisitOrExpression(e *ast.OrExpression) {
	e.LeftExpression.Accept(a)
	fmt.Printf(" or ")
	e.RightExpression.Accept(a)
	fmt.Printf("\n")
}

func (a *ASTPrinter) VisitBlock(b *ast.Block) {
}

func (a *ASTPrinter) VisitIfStatement(e *ast.IfStatement) {
	fmt.Printf("if ")
	e.Condition.Accept(a)
	fmt.Printf(" { \n")
	for _, inst := range e.Instructions.Statements {
		inst.Accept(a)
		fmt.Printf("\n")
	}
	fmt.Printf(" } \n")
	fmt.Printf(" else \n")
	fmt.Printf(" { \n")
	for _, inst := range e.ElseInstructions.Statements {
		inst.Accept(a)
		fmt.Printf("\n")
	}
	fmt.Printf(" } \n")
}

func (a *ASTPrinter) VisitReturnStatement(e *ast.ReturnStatement) {
	fmt.Printf("return ")
	e.Accept(a)
	fmt.Printf("\n")
}

func (a *ASTPrinter) VisitSwitchStatement(e *ast.SwitchStatement) {
}

func (a *ASTPrinter) VisitSwitchCase(e *ast.SwitchCase) {
}

func (a *ASTPrinter) VisitDefaultSwitchCase(e *ast.DefaultSwitchCase) {
}

func (a *ASTPrinter) VisitWhileStatement(e *ast.WhileStatement) {
}

func (a *ASTPrinter) VisitFunDef(e *ast.FunDef) {
}

func (a *ASTPrinter) VisitProgram(e *ast.Program) {
}
