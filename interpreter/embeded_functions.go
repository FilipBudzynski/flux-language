package interpreter

import (
	"fmt"
	"math"
	"reflect"
	"tkom/ast"
	"tkom/shared"
)

var PrintFunction = &ast.EmbeddedFunction{
	Name: "println",
	Func: func(args ...any) any {
		var output string
		for _, arg := range args {
			output += fmt.Sprintf("%v", arg)
		}
		fmt.Println(output)
		return nil
	},
	Parameters: []any{},
	Variadic:   true,
}

var PrintlnFunction = &ast.EmbeddedFunction{
	Name: "println",
	Func: func(args ...any) any {
		for _, arg := range args {
			fmt.Println(arg)
		}
		return nil
	},
	Parameters: []any{},
	Variadic:   true,
}

// Define the modulo function
var ModuloFunction = &ast.EmbeddedFunction{
	Name: "modulo",
	Func: func(args ...any) any {
		a := args[0].(int)
		b := args[1].(int)
		return a%b == 0
	},
	Parameters: []any{
		shared.FLOAT,
		shared.FLOAT,
	},
	Variadic: false,
}

// Define the square function
var SquareRootFunction = &ast.EmbeddedFunction{
	Name: "sqrt",
	Func: func(args ...any) any {
		if a, ok := args[0].(float64); ok {
			return math.Sqrt(a)
		} else {
			panic(fmt.Errorf(INVALID_ARGUMENTS_TYPE, reflect.TypeOf(args[0])))
		}
	},
	Parameters: []any{
		shared.FLOAT,
	},
	Variadic: false,
}

// Define the square function
var PowerFunction = &ast.EmbeddedFunction{
	Name: "power",
	Func: func(args ...any) any {
		a := args[0].(float64)
		b := args[1].(float64)

		return math.Pow(a, b)
	},
	Parameters: []any{
		shared.FLOAT,
		shared.FLOAT,
	},
	Variadic: false,
}

var embeddedFunctions = map[string]ast.Function{
	"print":   PrintFunction,
	"println": PrintlnFunction,
	"modulo":  ModuloFunction,
	"sqrt":    SquareRootFunction,
	"power":   PowerFunction,
}
