package lexer

import (
	"math"
	"strings"
	"tkom/shared"
	"unicode"
)

type Lexer struct {
	scanner         *Scanner
	ErrorHandler    func(err error)
	pos             shared.Position
	identifierLimit int
	stringLimit     int
	intLimit        int
}

func NewLexer(scanner *Scanner, identifierLimit, stringLimit, intLimit int) *Lexer {
	return &Lexer{
		scanner:         scanner,
		pos:             scanner.Position(),
		identifierLimit: identifierLimit,
		stringLimit:     stringLimit,
		intLimit:        intLimit,
	}
}

func (l *Lexer) GetNextToken() (t *Token) {
	defer func() {
		if err := recover(); err != nil {
			l.ErrorHandler(err.(error))
			t = NewToken(ETX, l.pos, nil)
		}
	}()

	l.skipWhiteChar()
	pos := l.pos

	if l.scanner.Character() == EOF {
		return NewToken(ETX, pos, nil)
	}

	t = l.createComment(pos)
	if t != nil {
		return t
	}

	t = l.createString(pos)
	if t != nil {
		return t
	}

	t = l.createOperator(pos)
	if t != nil {
		return t
	}

	t = l.createNumber(pos)
	if t != nil {
		return t
	}

	t = l.createIdentifier(pos)
	if t != nil {
		return t
	}

	return nil
}

func (l *Lexer) consume() rune {
	l.scanner.NextRune()
	l.pos = l.scanner.Position()
	return l.scanner.Character()
}

func (l *Lexer) skipWhiteChar() {
	for unicode.IsSpace(l.scanner.Character()) {
		l.consume()
	}
}

func (l *Lexer) createComment(position shared.Position) *Token {
	if l.scanner.Character() != '#' {
		return nil
	}

	var strBuilder strings.Builder
	for l.scanner.Character() != '\n' && l.scanner.Character() != EOF {
		strBuilder.WriteRune(l.scanner.Character())
		l.consume()
	}

	return NewToken(COMMENT, position, strBuilder.String())
}

func (l *Lexer) createOperator(position shared.Position) *Token {
	buff := l.scanner.Character()
	if buff == '<' || buff == '>' || buff == '=' || buff == '!' || buff == '-' || buff == ':' {
		char := l.consume()
		if token_type, ok := DoubleOperators[string([]rune{buff, char})]; ok {
			l.consume()
			return NewToken(token_type, position, nil)
		} else {
			if t_type, ok := Operators[buff]; ok {
				return NewToken(t_type, position, nil)
			}
		}
	} else if t_type, ok := Operators[buff]; ok {
		l.consume()
		return NewToken(t_type, position, nil)
	}
	return nil
}

func (l *Lexer) isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func (l *Lexer) createNumber(position shared.Position) *Token {
	if !l.isDigit(l.scanner.Character()) {
		return nil
	}

	value := int(l.scanner.Character() - '0')
	l.consume()

	if value != 0 {
		for l.isDigit(l.scanner.Character()) {
			digit := int(l.scanner.Character() - '0')
			if value > (l.intLimit-digit)/10 {
				panic(NewLexerError(INT_CAPACITY_EXCEEDED, position))
			}
			value = value*10 + digit
			l.consume()
		}
	}

	if l.scanner.Character() != '.' {
		return NewToken(CONST_INT, position, value)
	}
	l.consume()

	decimals := 0
	var decValue int
	for l.isDigit(l.scanner.Character()) {
		digit := int(l.scanner.Character() - '0')
		if decValue > (l.intLimit-digit)/10 {
			panic(NewLexerError(FLOAT_CAPACITY_EXCEEDED, l.pos))
		}
		decValue = decValue*10 + digit
		decimals += 1
		l.consume()
	}

	floatValue := float64(value) + float64(decValue)/math.Pow(10, float64(decimals))
	return NewToken(CONST_FLOAT, position, floatValue)
}

func (l *Lexer) createIdentifier(position shared.Position) *Token {
	if !unicode.IsLetter(l.scanner.Character()) {
		return nil
	}

	var strBuilder strings.Builder

	strBuilder.WriteRune(l.scanner.Character())
	l.consume()

	for unicode.IsLetter(l.scanner.Character()) || unicode.IsDigit(l.scanner.Character()) || l.scanner.Character() == '_' {
		if strBuilder.Len() == l.identifierLimit {
			panic(NewLexerError(IDENTIFIER_CAPACITY_EXCEEDED, l.pos))
		}
		strBuilder.WriteRune(l.scanner.Character())
		l.consume()
	}

	builtString := strBuilder.String()
	if tokenType, ok := KeyWords[builtString]; ok {
		return NewToken(tokenType, position, nil)
	}

	return NewToken(IDENTIFIER, position, builtString)
}

func (l *Lexer) handleEscaping() rune {
	if l.scanner.Character() != '\\' {
		return l.scanner.Character()
	}
	l.consume()
	switch l.scanner.Character() {
	case 'n':
		return '\n'
	case 't':
		return '\t'
	case '"':
		return '"'
	case '\\':
		return '\\'
	default:
		panic(NewLexerError(INVALID_ESCAPING, l.pos))
	}
}

func (l *Lexer) createString(position shared.Position) *Token {
	if l.scanner.Character() != '"' {
		return nil
	}
	var strBuilder strings.Builder
	l.consume()
	for l.scanner.Character() != '"' && l.scanner.Character() != EOF {
		if strBuilder.Len() == l.stringLimit {
			panic(NewLexerError(STRING_CAPACITY_EXCEEDED, l.pos))
		}
		charToAppend := l.handleEscaping()
		strBuilder.WriteRune(charToAppend)
		l.consume()
	}

	if l.scanner.Character() != '"' {
		panic(NewLexerError(STRING_NOT_CLOSED, l.pos))
	}

	l.consume()
	return NewToken(CONST_STRING, position, strBuilder.String())
}
