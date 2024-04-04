package lexer

import (
	"fmt"

	"github.com/google/uuid"
)

type identifierToken struct {
	Name string
	Id   string
	Type TokenTypes
	Pos  Position
}

func NewIdentifierToken(pos Position, name string) *identifierToken {
	uid := uuid.New().String()
	return &identifierToken{
		Type: IDENTIFIER,
		Name: name,
		Id:   uid,
		Pos:  pos,
	}
}

func (i *identifierToken) IsEqual(token Token) bool {
	if other, ok := token.(*identifierToken); ok {
		return i.Type == other.Type && i.Pos == other.Pos && i.Id == other.Id
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
