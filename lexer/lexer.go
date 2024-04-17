package lexer

import (
	"math"
	"strings"
	"unicode"
)

type Lexer struct {
	source          *Scanner
	errorHandler    func()
	pos             Position
	identifierLimit int
	stringLimit     int
	intLimit        int
}

func NewLexer(source *Scanner, identifierLimit, stringLimit, intLimit int) *Lexer {
	l := &Lexer{
		source:          source,
		pos:             source.Position(),
		identifierLimit: identifierLimit,
		stringLimit:     stringLimit,
		intLimit:        intLimit,
	}
	return l
}

func (l *Lexer) GetNextToken() (t *Token, err error) {
	l.skipWhiteChar()
	pos := l.pos

	if l.source.Current == EOF {
		return NewToken(ETX, pos, nil), nil
	}

	// EOL nie będzie potrzebny
	if l.source.Current == '\n' {
		l.Consume()
		return NewToken(EOL, pos, nil), nil
	}

	t = l.createComment(pos)
	if t != nil {
		return t, nil
	}

	t, err = l.createString(pos)
	if err != nil {
		return nil, err
	}
	if t != nil {
		return t, nil
	}

	t = l.createOperator(pos)
	if t != nil {
		return t, nil
	}

	t, err = l.createNumber(pos)
	if err != nil {
		return nil, err
	}
	if t != nil {
		return t, nil
	}

	t, err = l.createIdentifier(pos)
	if err != nil {
		return nil, err
	}
	if t != nil {
		return t, nil
	}

	return nil, NewLexerError(NONE_TOKEN_MATCH, l.pos)
}

func (l *Lexer) Consume() rune {
	_ = l.source.NextRune()
	l.pos = l.source.Position()
	return l.source.Current
}

func (l *Lexer) skipWhiteChar() {
	for unicode.IsSpace(l.source.Current) && l.source.Current != '\n' {
		l.Consume()
	}
}

func (l *Lexer) createComment(position Position) *Token {
	if l.source.Current != '#' {
		return nil
	}

	for l.source.Current != '\n' && l.source.Current != EOF {
		l.Consume()
	}

	return NewToken(COMMENT, position, nil)
}

func (l *Lexer) createOperator(position Position) *Token {
	buff := l.source.Current
	if buff == '<' || buff == '>' || buff == '=' || buff == '!' || buff == '-' || buff == ':' {
		char := l.Consume()
		if token_type, ok := DoubleOperators[string([]rune{buff, char})]; ok {
			l.Consume()
			return NewToken(token_type, position, nil)
		} else {
			if t_type, ok := Operators[buff]; ok {
				return NewToken(t_type, position, nil)
			}
		}
	} else if t_type, ok := Operators[buff]; ok {
		l.Consume()
		return NewToken(t_type, position, nil)
	}
	return nil
}

func (l *Lexer) isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func (l *Lexer) createNumber(position Position) (*Token, error) {
	if !l.isDigit(l.source.Current) {
		return nil, nil
	}

	value := int(l.source.Current - '0')
	l.Consume()

	if value != 0 {
		for l.isDigit(l.source.Current) {
			digit := int(l.source.Current - '0')
			if value > (l.intLimit-digit)/10 {
				return nil, NewLexerError(INT_CAPACITY_EXCEEDED, position)
			}
			value = value*10 + digit
			l.Consume()
		}
	}

	if l.source.Current != '.' {
		return NewToken(CONST_INT, position, value), nil
	}
	l.Consume()

	decimals := 0
	var decValue int
	for l.isDigit(l.source.Current) {
		digit := int(l.source.Current - '0')
		if decValue > (l.intLimit-digit)/10 {
			return nil, NewLexerError(FLOAT_CAPACITY_EXCEEDED, l.pos)
		}
		decValue = decValue*10 + digit
		decimals += 1
		l.Consume()
	}

	floatValue := float64(value) + float64(decValue)/math.Pow(10, float64(decimals))
	return NewToken(CONST_FLOAT, position, floatValue), nil
}

func (l *Lexer) createIdentifier(position Position) (*Token, error) {
	if !unicode.IsLetter(l.source.Current) {
		return nil, nil
	}

	var strBuilder strings.Builder

	strBuilder.WriteRune(l.source.Current)
	l.Consume()

	for unicode.IsLetter(l.source.Current) || unicode.IsDigit(l.source.Current) || l.source.Current == '_' {
		// consts for limits
		if strBuilder.Len() == l.identifierLimit {
			return nil, NewLexerError(IDENTIFIER_CAPACITY_EXCEEDED, l.pos)
		}
		strBuilder.WriteRune(l.source.Current)
		l.Consume()
	}

	builtString := strBuilder.String()
	if tokenType, ok := KeyWords[builtString]; ok {
		return NewToken(tokenType, position, nil), nil
	}

	return NewToken(IDENTIFIER, position, builtString), nil
}

func (l *Lexer) handleEscaping() rune {
	if l.source.Current != '\\' {
		return l.source.Current
	}
	l.Consume()
	switch l.source.Current {
	case 'n':
		return '\n'
	case 't':
		return '\t'
	case '"':
		return '"'
	case '\\':
		return '\\'
	}
	return EOF
}

func (l *Lexer) createString(position Position) (*Token, error) {
	if l.source.Current != '"' {
		return nil, nil
	}
	var strBuilder strings.Builder
	l.Consume()
	for l.source.Current != '"' && l.source.Current != EOF {
		if strBuilder.Len() == l.stringLimit {
			return nil, NewLexerError(STRING_CAPACITY_EXCEEDED, l.pos)
		}
		charToAppend := l.handleEscaping()
		strBuilder.WriteRune(charToAppend)
		l.Consume()
	}

	if l.source.Current != '"' {
		return nil, NewLexerError(STRING_NOT_CLOSED, l.pos)
	}

	l.Consume()
	return NewToken(CONST_STRING, position, strBuilder.String()), nil
}
