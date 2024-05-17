package ast

import "reflect"

type Case interface {
    Node
}

type SwitchStatement struct {
	Variables  []*Variable
	Expression Expression
	Cases      []Case
}

func NewSwitchStatement(variables []*Variable, expression Expression, cases []Case) *SwitchStatement {
	return &SwitchStatement{
		Variables:  variables,
		Expression: expression,
		Cases:      cases,
	}
}

func (s *SwitchStatement) Equals(other SwitchStatement) bool {
	if !reflect.DeepEqual(s.Expression, other.Expression) {
		return false
	}
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

func (s *SwitchStatement) Accept(v Visitor) {
	v.VisitSwitchStatement(s)
}

type SwitchCase struct {
	Condition        Expression
	OutputExpression Expression // block lub expression
}

func (s *SwitchCase) Accept(v Visitor) {
	v.VisitSwitchCase(s)
}

func (s *SwitchCase) Equals(o Case) bool {
	if other, ok := o.(*SwitchCase); ok {
		if !s.Condition.Equals(other.Condition) {
			return false
		}

		if !s.OutputExpression.Equals(other.OutputExpression) {
			return false
		}
		return true
	}
	return false
}

func NewSwitchCase(condition Expression, outputExpression Expression) *SwitchCase {
	return &SwitchCase{
		Condition:        condition,
		OutputExpression: outputExpression, // expression or *Block sturcture
	}
}

type DefaultSwitchCase struct {
	OutputExpression Expression
}

func NewDefaultCase(outputExpression Expression) *DefaultSwitchCase {
	return &DefaultSwitchCase{
		OutputExpression: outputExpression,
	}
}

func (d *DefaultSwitchCase) Accept(v Visitor) {
	v.VisitDefaultSwitchCase(d)
}

func (s *DefaultSwitchCase) Equals(o Case) bool {
	if other, ok := o.(*SwitchCase); ok {
		return s.OutputExpression.Equals(other.OutputExpression)
	}
	return false
}
