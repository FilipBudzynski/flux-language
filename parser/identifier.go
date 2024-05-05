package parser

import lex "tkom/lexer"

//TODO: czy idetifier powinien miec typ??
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

func (i *Identifier) Equals(other *Identifier) bool {
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
}
