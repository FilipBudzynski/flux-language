package parser

import "reflect"

type Case interface {
    Equals(other Case) bool
}


type SwitchStatement struct {
	Variables  []*Variable
	Expression Expression
	Cases      []Statement
}

func NewSwitchStatement(variables []*Variable, expression Expression, cases []Statement) *SwitchStatement {
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

type SwitchCase struct {
	CaseExpression   Expression
	OutputExpression Expression // block lub expression
}

func NewSwitchCase(caseExpression Expression, outputExpression Expression) *SwitchCase {
	return &SwitchCase{
		CaseExpression:   caseExpression,
		OutputExpression: outputExpression,
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