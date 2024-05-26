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
		return "INT"
	case FLOAT:
		return "FLOAT"
	case BOOL:
		return "BOOL"
	case STRING:
		return "STRING"
	case VOID:
		return "VOID"
	default:
		return fmt.Sprintf("Unknown TypeAnnotation: %d", t)
	}
}
