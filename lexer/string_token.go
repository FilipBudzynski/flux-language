package lexer

import (
	"fmt"
)

type stringToken struct {
	Value string
	Type  TokenTypes
	Pos   Position
}

func NewStringToken(name string) *stringToken {
	return &stringToken{
		Type:  CONST_STRING,
		Value: name,
	}
}

func (s *stringToken) IsEqual(token Token) bool {
	if other, ok := token.(*stringToken); ok {
		return s.Type == other.Type && s.Pos == other.Pos && s.Value == other.Value
	}
	return false
}

func (s *stringToken) ShowDetails() {
	fmt.Printf("%v, %v, %v\n", s.Pos, s.Type.GetName(), s.Value)
}

func (s *stringToken) GetType() TokenTypes {
	return s.Type
}

func (s *stringToken) GetName() string {
	return s.Type.GetName()
}

func (s *stringToken) SetPosition(position Position) {
	s.Pos = position
}
