package shared

import "fmt"

type TypeAnnotation int

const (
	INT TypeAnnotation = iota
	FLOAT
	BOOL
	STRING
	VOID
)

func (t TypeAnnotation) String() string {
	switch t {
	case INT:
		return "int"
	case FLOAT:
		return "float"
	case BOOL:
		return "bool"
	case STRING:
		return "string"
	case VOID:
		return "void"
	default:
		return fmt.Sprintf("Unknown TypeAnnotation: %d", t)
	}
}
