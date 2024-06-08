package shared

type Position struct {
	Line   int
	Column int
}

func NewPosition(line, column int) Position {
	return Position{Line: line, Column: column}
}
