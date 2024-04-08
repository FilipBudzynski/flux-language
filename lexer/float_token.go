package lexer

// import (
// 	"fmt"
// )
//
// type floatToken struct {
// 	Type  TokenTypes
// 	Pos   Position
// 	Value float64
// }
//
// func NewFloatToken(value float64, position Position) *floatToken {
// 	return &floatToken{
// 		Type:  CONST_FLOAT,
// 		Value: value,
// 		Pos:   position,
// 	}
// }
//
// func (f *floatToken) IsEqual(token Token) bool {
// 	if other, ok := token.(*floatToken); ok {
// 		return f.Type == other.Type && f.Pos == other.Pos && f.Value == other.Value
// 	}
// 	return false
// }
//
// func (f *floatToken) ShowDetails() {
// 	fmt.Printf("%v, %v, %v\n", f.Pos, f.Type.GetName(), f.Value)
// }
//
// func (f *floatToken) GetType() TokenTypes {
// 	return f.Type
// }
//
// func (f *floatToken) GetName() string {
// 	return f.Type.GetName()
// }
//
// func (f *floatToken) SetPosition(position Position) {
// 	f.Pos = position
// }
