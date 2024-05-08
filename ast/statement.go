package ast

type Statement interface{
    Accept(Visitor)
}
