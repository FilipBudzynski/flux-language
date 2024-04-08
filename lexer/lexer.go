package lexer

import (
	"fmt"
	"math"
	"strings"
	"unicode"
)

type Lexer struct {
	source Scanner
	pos    Position
}

func NewLexer(source Scanner) *Lexer { return &Lexer{source: source} }

func (l *Lexer) GetNextToken() (t *Token, err error) {
	l.skipWhiteChar()
	l.pos = l.source.Position()
	t = l.tryMatch()
	if t != nil {
		return t, nil
	}

	return nil, fmt.Errorf("None token match found for the source")
}

func (l *Lexer) tryMatch() (t *Token) {
	pos := l.pos

	if l.source.Current == EOF {
		return NewToken(ETX, pos, nil)
	}

	if l.source.Current == '\n' {
		t = NewToken(EOL, pos, nil)
		l.Consume()
		return t
	}

	// TODO: dodac bledy
	t = l.createString(pos)
	if t != nil {
		return t
	}

	// err != nil{
	//   fsdjakfklas (er)
	// }

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

func (l *Lexer) createOperator(position Position) *Token {
	buff := l.source.Current
	if buff == '<' || buff == '>' || buff == '=' || buff == '!' || buff == '-' || buff == ':' {
		char := l.Consume()
		if token_type, ok := DoubleOperators[string([]rune{buff, char})]; ok {
			l.Consume()
			return NewToken(token_type, position, nil)
		} else {
			if t_type, ok := SingleChar[buff]; ok {
				return NewToken(t_type, position, nil)
			}
		}
	} else if t_type, ok := SingleChar[buff]; ok {
		l.Consume()
		return NewToken(t_type, position, nil)
	}

	return nil
}

func (l *Lexer) isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func (l *Lexer) createNumber(position Position) *Token {
	// TODO: change to one of 0..9
	if !l.isDigit(l.source.Current) {
		return nil
	}
	if l.source.Current == '0' {
		l.Consume()
		return NewToken(CONST_INT, position, 0)
	}

	var value int
	for l.isDigit(l.source.Current) {
		digit := int(l.source.Current - '0')
		if value > (math.MaxInt-digit)/10 {
			// TODO: zglosic error
			fmt.Errorf("Error [%d, %d], Int value limit Exceeded", l.pos.Line, l.pos.Column)
			return nil //, err
		}
		value = value*10 + digit
		l.Consume()
	}

	if l.source.Current != '.' {
		return NewToken(CONST_INT, position, value)
	}
	l.Consume()

	decimals := 0
	var decValue int
	for l.isDigit(l.source.Current) {
		digit := int(l.source.Current - '0')
		if decValue > (math.MaxInt-digit)/10 {
			// TODO: zglosic error
			fmt.Errorf("Error [%d, %d], Int value limit Exceeded", l.pos.Line, l.pos.Column)
			return nil //, err
		}
		decValue = decValue*10 + digit
		decimals += 1
		l.Consume()
	}

	floatValue := float64(value) + float64(decValue)/math.Pow(10, float64(decimals))
	return NewToken(CONST_FLOAT, position, floatValue)
}

func (l *Lexer) createIdentifier(position Position) *Token {
	// TODO: zamienic na string buildera

	var strBuilder strings.Builder

	if !unicode.IsLetter(l.source.Current) {
		fmt.Errorf("error [%d, %d] identifier should start with a letter", l.pos.Line, l.pos.Column)
		return nil
	}

	strBuilder.WriteRune(l.source.Current)
	l.Consume()

	for {
		if unicode.IsLetter(l.source.Current) || unicode.IsDigit(l.source.Current) || l.source.Current == '_' {
			if strBuilder.Len() > 64*1024 {
				fmt.Errorf("error [%d, %d] Identifier is too long", l.pos.Line, l.pos.Column)
				return nil
			}
			strBuilder.WriteRune(l.source.Current)
			l.Consume()
		} else {
			break
		}
	}

	if tokenType, ok := KeyWords[strBuilder.String()]; ok {
		return NewToken(tokenType, position, nil)
	}

	return NewToken(IDENTIFIER, position, strBuilder.String())
}

func (l *Lexer) createString(position Position) *Token {
	if l.source.Current != '"' {
		return nil
	}

	var strBuilder strings.Builder
	l.Consume()
	for l.source.Current != '"' && l.source.Current != EOF {
		strBuilder.WriteRune(l.source.Current)
		l.Consume()
	}

	if l.source.Current != '"' {
		return nil
	}

	l.Consume()
	return NewToken(CONST_STRING, position, strBuilder.String())
}
