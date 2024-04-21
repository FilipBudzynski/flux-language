package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"tkom/lexer"
)

const (
	identifierLimit = 500
	stringLimit     = 1000
	intLimit        = math.MaxInt
)

func main() {
	// scannerTest()
	lexerTest()
}

func scannerTest() {
	file := strings.NewReader("int a = 5\n")
	scanner, _ := lexer.NewScanner(file)
	for {
		scanner.NextRune()
		if scanner.Character() == lexer.EOF {
			break
		}
		position := scanner.Position()
		fmt.Printf("Line: %d, Char: %d, Value: %c\n", position.Line, position.Column, scanner.Current)
	}
}

func lexerTest() {
	file, err := os.Open("example.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	source, _ := lexer.NewScanner(file)
	lex := lexer.NewLexer(source, identifierLimit, stringLimit, intLimit)

	for {
		token := lex.GetNextToken()

		if token == nil || token.Type == lexer.ETX {
			break
		}

		fmt.Printf("%-2v %-12v %-5v\n", token.Pos, token.Type.TypeName(), token.Value)

		if token.GetType() == lexer.ETX {
			fmt.Println("Koniec pliku")
			break
		}
	}
}
