package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"tkom/ast"
	"tkom/interpreter"
	"tkom/lexer"
	"tkom/parser"
)

const (
	IDENTIFIERLIMIT     = 500
	STRING_LIMIT        = 1000
	INT_LIMIT           = math.MaxInt
	MAX_RECURSION_DEPTH = 200
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "%v\n", r)
			//os.Exit(1)
		}
	}()

	if len(os.Args) < 2 {
		fmt.Println("Missing parameter, provide file name or use piped input!")
		return
	}

	var program *ast.Program
	var err error

	if len(os.Args) == 2 && os.Args[1] == "-" {
		program, err = parseProgramFromStdin()
	} else {
		fileName := os.Args[1]
		program, err = parseProgramFromFile(fileName)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	visitor := interpreter.NewCodeVisitor(MAX_RECURSION_DEPTH)
	visitor.MaxRecursionDepth = MAX_RECURSION_DEPTH
	for _, fd := range program.Functions {
		visitor.FunctionsMap[fd.Name] = fd
	}

	arguments := os.Args[2:]
	functionCallArgs := make([]ast.Expression, len(arguments))
	for i, arg := range arguments {
		if intValue, err := strconv.Atoi(arg); err == nil {
			functionCallArgs[i] = &ast.IntExpression{Value: intValue}
		} else {
			functionCallArgs[i] = &ast.StringExpression{Value: arg}
		}
	}

	functionCall := &ast.FunctionCall{
		Name:      "main",
		Arguments: functionCallArgs,
	}
	visitor.VisitFunctionCall(functionCall)
}

func parseProgramFromFile(fileName string) (*ast.Program, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	return parseProgram(reader)
}

func parseProgramFromStdin() (*ast.Program, error) {
	reader := bufio.NewReader(os.Stdin)
	return parseProgram(reader)
}

func parseProgram(reader *bufio.Reader) (*ast.Program, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "%v\n", r)
			os.Exit(1)
		}
	}()
	source, _ := lexer.NewScanner(reader)
	lex := lexer.NewLexer(source, IDENTIFIERLIMIT, STRING_LIMIT, INT_LIMIT)
	errorHandler := func(err error) {
		panic(err)
	}
    lex.ErrorHandler = errorHandler
	parser := parser.NewParser(lex, errorHandler)

	program := parser.ParseProgram()
	return program, nil
}
