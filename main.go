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
	scanner, _ := lexer.NewScanner(file)
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
	// sprawdzic zeby file byl zamykany przy bledach
	file, err := os.Open("example.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	text := "int a = 5\nif a == 6\nwhile a >c"
	reader := strings.NewReader(text)
	source, _ := lexer.NewScanner(reader)
	lex := lexer.NewLexer(*source)

	for {
		token, err := lex.GetNextToken()
		if err != nil {
			fmt.Println("Błąd podczas parsowania:", err)
			break
		}

		if token == nil {
			break
		}

		fmt.Printf("%-12v %-2v %-5v\n", token.Type.TypeName(), token.Pos, token.Value)

		if token.GetType() == lexer.ETX {
			fmt.Println("Koniec pliku")
			break
		}
	}
}
