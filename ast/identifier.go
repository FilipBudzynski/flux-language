package ast

import "tkom/shared"

type Identifier struct {
	Name     string
	Position shared.Position
}

func NewIdentifier(name string, position shared.Position) *Identifier {
	return &Identifier{
		Name:     name,
		Position: position,
	}
}

func (i *Identifier) Equals(other Expression) bool {
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

func (i *Identifier) Accept(v Visitor) {
	v.VisitIdentifier(i)
}

func (i *Identifier) GetPosition() shared.Position {
    return i.Position
}
