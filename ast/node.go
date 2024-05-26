package ast

type Node interface {
	Accept(Visitor)
}
