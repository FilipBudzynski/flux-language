package ast

import lex "tkom/lexer"

type Identifier struct {
	Name     string
	Position lex.Position
}

func NewIdentifier(name string, position lex.Position) Identifier {
	return Identifier{
		Name:     name,
		Position: position,
	}
}

func (i Identifier) Equals(other Expression) bool {
	if other, ok := other.(*Identifier); ok {
		if i.Name != other.Name {
			return false
		}
		if i.Position.Line != other.Position.Line {
			return false
		}
		if i.Position.Column != other.Position.Column {
			return false
		}
		return true
	} else {
		return false
	}
}
