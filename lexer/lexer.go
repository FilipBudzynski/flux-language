package lexer

import (
	"fmt"
	"io"
	"math"
	"strings"
	"unicode"
)

type Lexer struct {
	source Scanner
	pos    Position
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		source: *NewScanner(reader),
	}
}

func (l *Lexer) GetNextToken() (t Token, err error) {
	l.skipWhiteChar()
	pos := l.pos
	t = l.tryMatch()
	if t != nil {
		t.SetPosition(pos)
		return t, nil
	}
	return nil, fmt.Errorf("None token match found for the source")
}

func (l *Lexer) tryMatch() (t Token) {
	if l.source.Current == EOF {
		t = NewBaseToken(ETX)
		l.Consume()
		return t
	}

	if l.source.Current == '\n' {
		t = NewBaseToken(EOL)
		l.Consume()
		return t
	}

	t = l.createString()
	if t != nil {
		return t
	}

	t = l.createOperator()
	if t != nil {
		return t
	}

	t = l.createNumber()
	if t != nil {
		return t
	}

	t = l.createIdentifier()
	if t != nil {
		return t
	}

	return nil
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

func (l *Lexer) createOperator() Token {
	buff := l.source.Current
	if buff == '<' || buff == '>' || buff == '=' || buff == '!' || buff == '-' || buff == ':' {
		char := l.Consume()
		if token_type, ok := DoubleOperators[string([]rune{buff, char})]; ok {
			l.Consume()
			return NewBaseToken(token_type)
		}
	}
	if t_type, ok := SingleChar[buff]; ok {
		l.Consume()
		return NewBaseToken(t_type)
	}

	return nil
}

func (l *Lexer) createNumber() Token {
	if !unicode.IsDigit(l.source.Current) {
		return nil
	}
	if l.source.Current == '0' {
		l.Consume()
		return NewIntToken(0)
	}

	var value int
	for unicode.IsDigit(l.source.Current) {
		digit := int(l.source.Current - '0')
		if value > (math.MaxInt-digit)/10 {
			fmt.Errorf("Error [%d, %d], Int value limit Exceeded", l.pos.Line, l.pos.Column)
			return nil
		}
		value = value*10 + digit
		l.Consume()
	}

	if l.source.Current != '.' {
		return NewIntToken(value)
	}
	l.Consume()

	decimals := 0
	var decValue int
	for unicode.IsDigit(l.source.Current) {
		digit := int(l.source.Current - '0')
		decValue = decValue*10 + digit
		decimals += 1
		l.Consume()
	}

	floatValue := float64(value) + float64(decValue)/math.Pow(10, float64(decimals))
	return NewFloatToken(floatValue)
}

func (l *Lexer) createIdentifier() Token {
	var identifier []rune

	if !unicode.IsLetter(l.source.Current) {
		fmt.Errorf("error [%d, %d] identifier should start with a letter", l.pos.Line, l.pos.Column)
		return nil
	}

	identifier = append(identifier, l.source.Current)
	l.Consume()

	for {
		if unicode.IsLetter(l.source.Current) || unicode.IsDigit(l.source.Current) || l.source.Current == '_' {
			identifier = append(identifier, l.source.Current)
			l.Consume()
		} else {
			break
		}
	}

	if tokenType, ok := KeyWords[string(identifier)]; ok {
		return NewBaseToken(tokenType)
	}

	if len(string(identifier)) > 64*1024 {
		fmt.Errorf("error [%d, %d] Identifier is too long", l.pos.Line, l.pos.Column)
		return nil
	}

	return NewIdentifierToken(string(identifier))
}

func (l *Lexer) createString() Token {
	if l.source.Current != '"' {
		return nil
	}

	var strBuilder strings.Builder
	l.Consume()
	for l.source.Current != '"' && l.source.Current != EOF {
		if l.source.Current == '\\' {
			l.Consume()
			if l.source.Current == EOF {
				return nil
			}
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
				return nil
			}
		} else {
			strBuilder.WriteRune(l.source.Current)
		}
		l.Consume()
	}

	if l.source.Current != '"' {
		return nil
	}

	l.Consume()
	return NewStringToken(strBuilder.String())
}
