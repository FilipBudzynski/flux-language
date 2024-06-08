package ast

import (
	"tkom/shared"
)

type Case interface {
	Node
	GetPosition() shared.Position
}

type SwitchStatement struct {
	Variables []*Variable
	Cases     []Case
	Position  shared.Position
}

func NewSwitchStatement(variables []*Variable, cases []Case, position shared.Position) *SwitchStatement {
	return &SwitchStatement{
		Variables: variables,
		Cases:     cases,
		Position:  position,
	}
}

func (s *SwitchStatement) Equals(other Expression) bool {
	if other, ok := other.(*SwitchStatement); ok {
		if len(s.Cases) != len(other.Cases) {
			return false
		}
		for i, c := range s.Cases {
			if c == other.Cases[i] {
				return false
			}
		}
		return true
	}
	return false
}

func (s *SwitchStatement) Accept(v Visitor) {
	v.VisitSwitchStatement(s)
}

func (s *SwitchStatement) GetPosition() shared.Position {
	return s.Position
}

type SwitchCase struct {
	Condition        Expression
	OutputExpression Expression
	Position         shared.Position
}

func (s *SwitchCase) Accept(v Visitor) {
	v.VisitSwitchCase(s)
}

func (s *SwitchCase) GetPosition() shared.Position {
	return s.Position
}

func NewSwitchCase(condition Expression, outputExpression Expression, position shared.Position) *SwitchCase {
	return &SwitchCase{
		Condition:        condition,
		OutputExpression: outputExpression,
		Position:         position,
	}
}

type DefaultSwitchCase struct {
	OutputExpression Expression
	Position         shared.Position
}

func NewDefaultCase(outputExpression Expression, position shared.Position) *DefaultSwitchCase {
	return &DefaultSwitchCase{
		OutputExpression: outputExpression,
		Position:         position,
	}
}

func (d *DefaultSwitchCase) Accept(v Visitor) {
	v.VisitDefaultSwitchCase(d)
}

func (d *DefaultSwitchCase) GetPosition() shared.Position {
	return d.Position
}
