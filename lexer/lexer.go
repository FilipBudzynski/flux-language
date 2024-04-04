package lexer

import (
	"fmt"
	"io"
	"math"
	"unicode"
)

type Lexer struct {
	singleChar map[rune]TokenTypes
	keyWords   map[string]TokenTypes
	scanner    Scanner
	pos        Position
	buffer     []rune
}

func NewLexer(reader io.Reader) *Lexer {
	singleChar := map[rune]TokenTypes{
		'+': PLUS,
		'-': MINUS,
		'*': MULTIPLY,
		'/': DIVIDE,
		'{': LEFT_BRACE,
		'}': RIGHT_BRACE,
		'(': LEFT_PARENTHESIS,
		')': RIGHT_PARENTHESIS,
		',': COMMA,
		'>': GREATER_THAN,
	}

	keyWords := map[string]TokenTypes{
		"int":     INT,
		"string":  STRING,
		"bool":    BOOL,
		"float":   FLOAT,
		"switch":  SWITCH,
		"while":   WHILE,
		"if":      IF,
		"else":    ELSE,
		"default": DEFAULT,
		"return":  RETURN,
	}

	return &Lexer{
		singleChar: singleChar,
		keyWords:   keyWords,
		scanner:    *NewScanner(reader),
	}
}

func (l *Lexer) GetNextToken() (t Token, err error) {
	l.skipWhiteChar()

	if l.scanner.Current == 0xFFFF {
		t = NewBaseToken(ETX, l.pos)
		l.scanner.NextRune()
		return t, nil
	}

	// end of line
	if l.scanner.Current == '\n' {
		t = NewBaseToken(EOL, l.scanner.Position())
		l.scanner.NextRune()
		return t, nil
	}

	t = l.createOperator()
	if t != nil {
		return t, nil
	}

	// t = l.createKeyWord()
	// if t != nil {
	// 	return t, nil
	// }

	t = l.createInt()
	if t != nil {
		return t, nil
	}

	t = l.createIdentifier()
	if t != nil {
		return t, nil
	}

	fmt.Println("niedopasowanie")
	return nil, err
}

func (l *Lexer) Consume() {
	l.scanner.NextRune()
	l.pos = l.scanner.Position()
}

func (l *Lexer) skipWhiteChar() {
	for unicode.IsSpace(l.scanner.Current) {
		l.Consume()
	}
}

func (l *Lexer) createOperator() Token {
	pos := l.pos
	if l.scanner.Current == '=' {
		l.Consume()
		if l.scanner.Current == '=' {
			l.Consume()
			return NewBaseToken(EQUALS, pos)
		} else {
			return NewBaseToken(ASSIGN, pos)
		}
	}

	// special one char tokens
	if token_type, ok := l.singleChar[l.scanner.Current]; ok {
		l.Consume()
		return NewBaseToken(token_type, l.scanner.Position())
	}

	return nil
}

func (l *Lexer) createKeyWord() Token {
	pos := l.pos
	keyword := ""
	for unicode.IsLetter(l.scanner.Current) && len(keyword) <= 64*1024 {
		keyword += string(l.scanner.Current)
		l.Consume()
	}

	if tokenType, ok := l.keyWords[keyword]; ok {
		return NewBaseToken(tokenType, pos)
	}

	return nil
}

func (l *Lexer) createInt() Token {
	pos := l.pos
	if !unicode.IsDigit(l.scanner.Current) {
		return nil
	}
	if l.scanner.Current == '0' {
		l.Consume()
		return NewIntToken(pos, 0)
	}

	var value int
	for unicode.IsDigit(l.scanner.Current) {
		if value >= math.MaxInt {
			fmt.Errorf("Error [%d, %d], Int value limit Exceeded", l.pos.Line, l.pos.Column)
			return nil
		}

		value += int(l.scanner.Current - '0')
		l.Consume()
	}

	return NewIntToken(pos, value)
}

func (l *Lexer) createIdentifier() Token {
	pos := l.pos
	var identifier []rune

	if !unicode.IsLetter(l.scanner.Current) {
		fmt.Errorf("error [%d, %d] identifier should start with a letter", l.pos.Line, l.pos.Column)
		return nil
	}

	identifier = append(identifier, l.scanner.Current)
	l.Consume()

	for {
		if unicode.IsLetter(l.scanner.Current) || unicode.IsDigit(l.scanner.Current) || l.scanner.Current == '_' {
			identifier = append(identifier, l.scanner.Current)
			l.Consume()
		} else {
			break
		}
	}

	if tokenType, ok := l.keyWords[string(identifier)]; ok {
		return NewBaseToken(tokenType, pos)
	}

	if len(string(identifier)) > 64*1024 {
		fmt.Errorf("error [%d, %d] Identifier is too long", l.pos.Line, l.pos.Column)
		return nil
	}

	return NewIdentifierToken(pos, string(identifier))
}
