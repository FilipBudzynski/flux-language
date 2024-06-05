package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
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
		}
	}()

	if len(os.Args) < 2 {
		fmt.Println("Missing parameter, provide file name and arguments or use '-' to run from stream")
		return
	}

	var program *ast.Program
	var err error

	if len(os.Args) >= 2 && os.Args[1] == "-" {
		program, err = parseProgramFromStdin()
	} else {
		fileName := os.Args[1]
		ext := filepath.Ext(fileName)
		if ext != ".fl" {
			log.Fatal("File must have '.fl' extension")
			os.Exit(1)
		}

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
	defer func() {
		if r := recover(); r != nil {
			file.Close()
			fmt.Fprintf(os.Stderr, "%v\n", r)
			os.Exit(1)
		}
	}()

	reader := bufio.NewReader(file)
	return parseProgram(reader)
}

func parseProgramFromStdin() (*ast.Program, error) {
	reader := bufio.NewReader(os.Stdin)
	return parseProgram(reader)
}

func parseProgram(reader *bufio.Reader) (*ast.Program, error) {
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
