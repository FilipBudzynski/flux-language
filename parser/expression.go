package parser

import "tkom/lexer"

// czy to jest aby na pewno dobrze?
type Expression struct {
	Conjunctions []*Conjunction
}

type Conjunction struct {
	RelationTerms []*RelationTerm
	Operators     []string
}

type RelationTerm struct {
	Operator      *RelationOperator // This will be nil if no relation operator is present
	AdditiveTerms []*AdditiveTerm
}

type RelationOperator struct {
	Operator string // This will store one of: ">=", ">", "<=", "<", "==", "!="
}

type AdditiveTerm struct {
	MultiplicativeTerms []*MultiplicativeTerm
	Operators           []string // This will store "+" or "-" operators between terms
}

type MultiplicativeTerm struct {
	UnaryOperators []*UnaryOperator
	Operators      []string // This will store "*" or "/" operators between terms
}

type UnaryOperator struct {
	Term   *Term
	Negate bool // Indicates if the operator is a "-" or "!"
}

type Term struct {
	Identifier    *IdentifierOrCall
	Expression    *Expression
	CastedTerm    *CastedTerm
	Parenthesized *Expression // This will be nil if the term is not parenthesized
	String  string
	Bool    bool
	Integer int
	Float   float64
}

type IdentifierOrCall struct {
	Identifier string
	// Arguments  []*Argument // This will be nil if no arguments are present
}

type CastedTerm struct {
	Term       *Term
	Annotation *TypeAnnotation // This will be nil if no type annotation is present
}

type TypeAnnotation struct {
	Type lexer.TokenTypes // This will store one of: INT, FLOAT, BOOL, STRING
}
