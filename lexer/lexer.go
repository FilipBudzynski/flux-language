package lexer

import (
	"fmt"
	"math"
	"strings"
	"unicode"
)

type Lexer struct {
	source *Scanner
	pos    Position
}

func NewLexer(source *Scanner) *Lexer {
	l := &Lexer{source: source}
	l.pos = l.source.Position()
	return l
}

func (l *Lexer) GetNextToken() (t *Token, err error) {
	l.skipWhiteChar()
	pos := l.pos

	if l.source.Current == EOF {
		return NewToken(ETX, pos, nil), nil
	}

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

	return nil, fmt.Errorf("none token match found for the source")
}

func (l *Lexer) Consume() rune {
	l.source.NextRune()
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
	if l.source.Current == '0' {
		l.Consume()
		return NewToken(CONST_INT, position, 0), nil
	}

	var value int
	for l.isDigit(l.source.Current) {
		digit := int(l.source.Current - '0')
		if value >= (math.MaxInt-digit)/10 {
			return nil, fmt.Errorf("error [%d, %d], Int value limit Exceeded", position.Line, position.Column)
		}
		value = value*10 + digit
		l.Consume()
	}

	if l.source.Current != '.' {
		return NewToken(CONST_INT, position, value), nil
	}
	l.Consume()

	decimals := 0
	var decValue int
	for l.isDigit(l.source.Current) {
		digit := int(l.source.Current - '0')
		if decValue >= (math.MaxInt-digit)/10 {
			return nil, fmt.Errorf("error [%d, %d], decimal value limit Exceeded", l.pos.Line, l.pos.Column)
		}
		decValue = decValue*10 + digit
		decimals += 1
		l.Consume()
	}

	floatValue := float64(value) + float64(decValue)/math.Pow(10, float64(decimals))
	return NewToken(CONST_FLOAT, position, floatValue), nil
}

func (l *Lexer) createIdentifier(position Position) (*Token, error) {
	var strBuilder strings.Builder

	if !unicode.IsLetter(l.source.Current) {
		return nil, nil
	}

	strBuilder.WriteRune(l.source.Current)
	l.Consume()

	for unicode.IsLetter(l.source.Current) || unicode.IsDigit(l.source.Current) || l.source.Current == '_' {
		if strBuilder.Len() >= 64*1024 {
			return nil, fmt.Errorf("error [%d, %d] Identifier capacity exceeded", l.pos.Line, l.pos.Column)
		}
		strBuilder.WriteRune(l.source.Current)
		l.Consume()
	}

	if tokenType, ok := KeyWords[strBuilder.String()]; ok {
		if tokenType == CONST_BOOL {
			if strBuilder.String() == "true" {
				return NewToken(CONST_BOOL, position, true), nil
			} else {
				return NewToken(CONST_BOOL, position, false), nil
			}
		}
		return NewToken(tokenType, position, nil), nil
	}

	return NewToken(IDENTIFIER, position, strBuilder.String()), nil
}

func (l *Lexer) createString(position Position) (*Token, error) {
	if l.source.Current != '"' {
		return nil, nil
	}
	var strBuilder strings.Builder
	l.Consume()
	for l.source.Current != '"' && l.source.Current != EOF {
		if strBuilder.Len() >= 64*1024 {
			return nil, fmt.Errorf("error [%d, %d] Identifier capacity exceeded", l.pos.Line, l.pos.Column)
		}
		if l.source.Current == '\\' {
			l.Consume()
			switch l.source.Current {
			case 'n':
				strBuilder.WriteRune('\n')
			case 't':
				strBuilder.WriteRune('\t')
			case '"':
				strBuilder.WriteRune('"')
			case '\\':
				strBuilder.WriteRune('\\')
			default:
				return nil, fmt.Errorf("error [%d, %d] Invalid syntax escaping", l.pos.Line, l.pos.Column)
			}
		} else {
			strBuilder.WriteRune(l.source.Current)
		}
		l.Consume()
	}

	if l.source.Current != '"' {
		return nil, fmt.Errorf("error [%d, %d] String not closed, perhaps you forgot \"", l.pos.Line, l.pos.Column)
	}

	l.Consume()
	return NewToken(CONST_STRING, position, strBuilder.String()), nil
}
