package parser

const ERROR_WRONG_VALUE_IN_DECLARATION = "cannot use \"%s\", as %s value in variable declaration"

type Variable struct {
	Value      any
	Identifier Identifier
	Type       TypeAnnotation
}

func NewVariable(variableType TypeAnnotation, identifier Identifier, value any) *Variable {
	return &Variable{
		Type:       variableType,
		Identifier: identifier,
		Value:      value,
	}
}

func (v *Variable) Equals(other Variable) bool {
	if v.Type != other.Type {
		return false
	}
	if v.Value != other.Value {
		return false
	}
	if !v.Identifier.Equals(&other.Identifier) {
		return false
	}
	return true
}

// func convertValue(value any, expectedType reflect.Kind) (any, error) {
// 	switch expectedType {
// 	case reflect.Int:
// 		if v, ok := value.(int); ok {
// 			return v, nil
// 		}
// 	case reflect.Float64:
// 		if v, ok := value.(float64); ok {
// 			return v, nil
// 		}
// 	case reflect.String:
// 		if v, ok := value.(string); ok {
// 			return v, nil
// 		}
// 	}
// 	return nil, fmt.Errorf(ERROR_WRONG_VALUE_IN_DECLARATION, expectedType, value)
// }

type Assignemnt struct {
	Value      Expression
	Identifier Identifier
}

func NewAssignment(identifier Identifier, value Expression) *Assignemnt {
	return &Assignemnt{
		Identifier: identifier,
		Value:      value,
	}
}

func (v *Assignemnt) Equals(other Assignemnt) bool {
	if v.Value != other.Value {
		return false
	}
	if !v.Identifier.Equals(&other.Identifier) {
		return false
	}
	return false
}
