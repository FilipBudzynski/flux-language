package lexer

import (
	"fmt"
)

type identifierToken struct {
	Name string
	Type TokenTypes
	Pos  Position
}

func NewIdentifierToken(name string, position Position) *identifierToken {
	return &identifierToken{
		Type: IDENTIFIER,
		Name: name,
		Pos:  position,
	}
}

func (i *identifierToken) IsEqual(token Token) bool {
	if other, ok := token.(*identifierToken); ok {
		return i.Type == other.Type && i.Pos == other.Pos && i.Name == other.Name
	}
	return false
}

func (i *identifierToken) ShowDetails() {
	fmt.Printf("%v, %v, %v\n", i.Pos, i.Type.GetName(), i.Name)
}

func (i *identifierToken) GetType() TokenTypes {
	return i.Type
}

func (i *identifierToken) GetName() string {
	return i.Type.GetName()
}

func (i *identifierToken) SetPosition(position Position) {
	i.Pos = position
}
