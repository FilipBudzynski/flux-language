package main

import (
	"fmt"
	"io"
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
	text := "int a = 5\nif a == 6\nwhile" + string(0xFFFF)
	reader := strings.NewReader(text)
	lex := lexer.NewLexer(reader)
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

		if token.GetType() == lexer.ETX {
			fmt.Println("Koniec pliku")
			break // Koniec pliku
		}
		token.ShowDetails()
	}
}
