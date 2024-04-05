package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"tkom/lexer"
)

func main() {
	// scannerTest()
	lexerTest()
}

func scannerTest() {
	file := strings.NewReader("int a = 5\n")
	scanner := lexer.NewScanner(file)
	for {
		err := scanner.NextRune()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error scanning file:", err)
			return
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

	// text := "int a = 5\nif a == 6\nwhile a > c" + string(0xFFFF)
	// reader := strings.NewReader(text)

	lex := lexer.NewLexer(file)
	lex.Consume()

	for {
		token, err := lex.GetNextToken()
		if err != nil {
			fmt.Println("Błąd podczas parsowania:", err)
			break
		}

		if token == nil {
			break
		}

		token.ShowDetails()
		if token.GetType() == lexer.ETX {
			fmt.Println("Koniec pliku")
			break
		}
	}
}
