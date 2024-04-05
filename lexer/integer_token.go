package lexer

import (
	"fmt"
)

type intToken struct {
	Type  TokenTypes
	Pos   Position
	Value int
}

func NewIntToken(value int) *intToken {
	return &intToken{
		Type:  CONSNT_INT,
		Value: value,
	}
}

func (i *intToken) IsEqual(token Token) bool {
	if other, ok := token.(*intToken); ok {
		return i.Type == other.Type && i.Pos == other.Pos && i.Value == other.Value
	}
	return false
}

func (i *intToken) ShowDetails() {
	fmt.Printf("%v, %v, %v\n", i.Pos, i.Type.GetName(), i.Value)
}

func (i *intToken) GetType() TokenTypes {
	return i.Type
}

func (i *intToken) GetName() string {
	return i.Type.GetName()
}

func (i *intToken) SetPosition(position Position) {
	i.Pos = position
}
