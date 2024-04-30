package parser

import (
	"fmt"
	lex "tkom/lexer"
)

type Parser struct {
	lexer        *lex.Lexer
	ErrorHandler func(error)
	token        lex.Token
}

func NewParser(lexer *lex.Lexer, errHandler func(error)) *Parser {
	p := &Parser{
		lexer:        lexer,
		ErrorHandler: errHandler,
	}
	p.consumeToken()
	return p
}

func (p *Parser) consumeToken() {
	p.token = *p.lexer.GetNextToken()
	// log.Printf("token type: %v, token value: %v", p.token.Type, p.token.Value)
}

func (p *Parser) requierAndConsume(tokenType lex.TokenType, syntaxErrMessage string) lex.Token {
	token := p.token
	if token.Type != tokenType {
		panic(fmt.Errorf(syntaxErrMessage, token.Position.Column, token.Position.Line))
	}
	p.consumeToken()
	return token
}

func (p *Parser) parseIdentifier() (string, lex.Position) {
	name := p.token.Value.(string)
	possition := p.token.Position
	p.consumeToken()
	return name, possition
}

func (p *Parser) recoverFromPanic() {
	if err := recover(); err != nil {
		p.ErrorHandler(err.(error))
	}
}

// program = { function_definition } ;
func (p *Parser) ParseProgram() *Program {
	defer p.recoverFromPanic()

	functions := map[string]FunDef{}

	for funDef := p.parseFunDef(); funDef != nil; {
		if f, ok := functions[funDef.Name]; ok {
			tokenCol := p.token.Position.Column
			tokenLine := p.token.Position.Line
			panic(fmt.Errorf(SYNTAX_ERROR_FUNCTION_REDEFINITION, tokenCol, tokenLine, f.Position.Column, f.Position.Line))
		}
	}

	if p.token.Type != lex.ETX {
		panic(NO_ETX_TOKEN)
	}
	return NewProgram(functions)
}

// function_definition = identifier , "(", [ parameters ], ")", [ type_annotation ] , block ;
func (p *Parser) parseFunDef() *FunDef {
	if p.token.Type != lex.IDENTIFIER {
		return nil
	}
	name := p.token.Value.(string)
	possition := p.token.Position

	p.consumeToken()
	p.requierAndConsume(lex.LEFT_PARENTHESIS, SYNTAX_ERROR_FUNC_DEF_NO_PARENTHASIS)
	params := p.parseParameters()
	p.requierAndConsume(lex.RIGHT_PARENTHESIS, SYNTAX_ERROR_FUNC_DEF_NO_PARENTHASIS)
	funcType := p.parseTypeAnnotation()
	block := p.parseBlock()

	return NewFunctionDefinition(name, params, funcType, block, possition)
}

// parameters = parameter_group , { "," , parameter_group } ;
func (p *Parser) parseParameters() []*Variable {
	paramGroup := p.parseParameterGroup()
	if paramGroup == nil {
		return nil
	}
	parameters := []*Variable{}
	parameters = append(parameters, paramGroup...)

	for p.token.Type == lex.COMMA {
		p.consumeToken()
		paramGroup := p.parseParameterGroup()
		parameters = append(parameters, paramGroup...)
		if paramGroup == nil {
			panic(fmt.Errorf("syntax error, no parameters after comma"))
		}
	}
	return parameters
}

// parameter_group = identifier , { ",", identifier }, type_annotation ;
func (p *Parser) parseParameterGroup() []*Variable {
	if p.token.Type != lex.IDENTIFIER {
		return nil
	}

	type Tuple struct {
		Name     string
		Position lex.Position
	}

	name := p.token.Value.(string)
	possition := p.token.Position

	namesAndPositions := []Tuple{}
	namesAndPositions = append(namesAndPositions, Tuple{Name: name, Position: possition})

	p.consumeToken()
	for p.token.Type == lex.COMMA {
		p.consumeToken()
		if p.token.Type != lex.IDENTIFIER {
			panic(fmt.Errorf(SYNTAX_ERROR_NO_IDENTIFIER, p.token.Position.Column, p.token.Position.Line))
		}
		name := p.token.Value.(string)
		possition := p.token.Position
		namesAndPositions = append(namesAndPositions, Tuple{Name: name, Position: possition})
		p.consumeToken()
	}
	paramsType := p.parseTypeAnnotation()

	if paramsType == nil {
		panic(fmt.Errorf(SYNTAX_ERROR_NO_TYPE, p.token.Position.Column, p.token.Position.Line))
	}
	params := []*Variable{}

	for _, t := range namesAndPositions {
		params = append(params, NewVariable(*paramsType, NewIdentifier(t.Name, t.Position), nil))
	}
	return params
}

// type_annotation = "int" | "float" | "bool" | "str" ;
func (p *Parser) parseTypeAnnotation() *lex.TokenType {
	switch token := p.token.Type; token {
	case lex.INT, lex.FLOAT, lex.BOOL, lex.STRING:
		typ := p.token.Type
		p.consumeToken()
		return &typ
	}
	return nil
}

// block = "{" , { statement } , "}" ;
func (p *Parser) parseBlock() []Statement {
	p.requierAndConsume(lex.LEFT_BRACE, SYNTAX_ERROR_NO_BLOCK)

	statements := []Statement{}

	for p.token.Type != lex.RIGHT_BRACE {
		statement := p.parseStatement()

		// czy my chcemy żeby statement był nil??
		if statement != nil {
			statements = append(statements, statement)
		}
	}
	p.requierAndConsume(lex.RIGHT_BRACE, SYNTAX_ERROR_EXPECTED_RIGHT_BRACE)

	return statements
}

// statement = variable_declaration | assigment | conditional_statement | loop_statement | switch_statement | return_statement ;
func (p *Parser) parseStatement() Statement {
	switch p.token.Type {
	case lex.INT, lex.FLOAT, lex.BOOL, lex.STRING:
		return p.parseVariableDeclaration()
	case lex.IDENTIFIER:
		return p.parseAssignment()
	case lex.IF:
		return p.parseConditionalStatement()
	case lex.WHILE:
		return p.parseWhileStatement()
	// case lex.SWITCH:
	// 	return p.parseSwitchStatement()
	// case lex.RETURN:
	// 	return p.parseReturnStatement()
	default:
		// czy na pewno chcemy nil?? moze jednak panic?
		return nil // panic(fmt.Errorf(SYNTAX_ERROR_UNKNOWN_STATEMENT, p.token.Position.Column, p.token.Position.Line))
	}
}

// variable_declaration  = type_annotation, identifier, ":=", expression ;
func (p *Parser) parseVariableDeclaration() Statement {
	typeAnnotation := p.parseTypeAnnotation()
	if typeAnnotation == nil {
		panic(fmt.Sprintf(SYNTAX_ERROR_NO_TYPE_IN_DECLARATION, p.token.Position.Line, p.token.Position.Column))
	}
	identifierToken := p.requierAndConsume(lex.IDENTIFIER, SYNTAX_ERROR_NO_VARIABLE_IDETIFIER)
	p.requierAndConsume(lex.DECLARE, SYNTAX_ERROR_MISSING_COLON_ASSIGN)

	expression := p.parseExpression()

	if expression == nil {
		panic(fmt.Sprintf(SYNTAX_ERROR_NO_EXPRESSION_IN_VARIABLE_DECLARATION, p.token.Position.Column, p.token.Position.Line))
	}

	name := identifierToken.Value.(string)
	position := identifierToken.Position
	identifier := NewIdentifier(name, position)
	variable := NewVariable(*typeAnnotation, identifier, expression)
	return variable
}

// assigment = identifier_or_call,  [ "=", expression ] ;
func (p *Parser) parseAssignment() Statement {
	if p.token.Type != lex.IDENTIFIER {
		return nil
	}

	identifier := NewIdentifier(p.token.Value.(string), p.token.Position)
	p.consumeToken()

	if functionCall := p.parseFunctionCall(identifier); functionCall != nil {
		return functionCall
	}

	if p.token.Type != lex.ASSIGN {
		return identifier
	}

	p.consumeToken()

	expression := p.parseExpression()

	return NewAssignment(identifier, expression)
}

// identifier_or_call = identifier, [ "(", [ argumets ], ")" ] ;
func (p *Parser) parseIdentifierOrCall() Statement {
	if p.token.Type != lex.IDENTIFIER {
		return nil
	}

	identifier := NewIdentifier(p.token.Value.(string), p.token.Position)
	p.consumeToken()

	if functionCall := p.parseFunctionCall(identifier); functionCall != nil {
		return functionCall
	}

	return identifier
}

func (p *Parser) parseFunctionCall(identifier Identifier) Statement {
	if p.token.Type != lex.LEFT_PARENTHESIS {
		return nil
	}

	arguments := p.parseArguments()

	if p.token.Type != lex.RIGHT_PARENTHESIS {
		panic(SYNTAX_ERROR_FUNC_CALL_NOT_CLOSED)
	}

	return newFunctionCall(identifier, arguments, identifier.Position)
}

// arguments = expression , { "," , expression } ;
func (p *Parser) parseArguments() []Expression { // return []Variable czy []Expression ??
	expressions := []Expression{}

	expression := p.parseExpression()
	// na pewno?
	if expression == nil {
		return expressions
	}

	expressions = append(expressions, expression)

	for p.token.Type == lex.COMMA {
		expression := p.parseExpression()
		expressions = append(expressions, expression)
		// p.consumeToken()
	}
	return expressions
}

// expression = conjunction_term, { "or", conjunction_term } ;
func (p *Parser) parseExpression() Expression {
	leftExpression := p.parseAndCondition()
	if leftExpression == nil {
		return nil
	}

	for p.token.Type == lex.OR {
		p.consumeToken()
		rightExpression := p.parseAndCondition()
		if rightExpression == nil {
			panic(fmt.Sprintf(ERROR_MISSING_EXPRESSION, p.token.Position.Column, p.token.Position.Line, "OR"))
		}
		// leftExpression = NewExpression(leftExpression, OR, rightExpression)
		leftExpression = NewOrExpression(leftExpression, rightExpression)
	}
	return leftExpression
}

// conjunction_term = relation_term, { "and", relation_term } ;
func (p *Parser) parseAndCondition() Expression {
	leftExpression := p.parseRelationCondition()
	if leftExpression == nil {
		return nil
	}

	for p.token.Type == lex.AND {
		p.consumeToken()
		rightExpression := p.parseRelationCondition()
		if rightExpression == nil {
			panic(fmt.Sprintf(ERROR_MISSING_EXPRESSION, p.token.Position.Column, p.token.Position.Line, "AND"))
		}
		// leftExpression = NewExpression(leftExpression, AND, rightExpression)
		leftExpression = NewAndExpression(leftExpression, rightExpression)
	}

	return leftExpression
}

// relation_term = additive_term, [ relation_operator, additive_term ] ;
func (p *Parser) parseRelationCondition() Expression {
	leftExpression := p.parsePlusOrMinus()
	if leftExpression == nil {
		return nil
	}

	validOperators := map[lex.TokenType]bool{
		lex.EQUALS:           true,
		lex.NOT_EQUALS:       true,
		lex.GREATER_THAN:     true,
		lex.GREATER_OR_EQUAL: true,
		lex.LESS_OR_EQUAL:    true,
		lex.LESS_THAN:        true,
	}

	if _, ok := validOperators[p.token.Type]; !ok {
		return leftExpression
	}

	operation := p.token.Type
	p.consumeToken()

	rightExpression := p.parsePlusOrMinus()
	if rightExpression == nil {
		panic(fmt.Sprintf(ERROR_MISSING_EXPRESSION, p.token.Position.Column, p.token.Position.Line, operation))
	}

	switch operation {
	case lex.EQUALS:
		return NewEqualsExpression(leftExpression, rightExpression)
	case lex.NOT_EQUALS:
		return NewNotEqualsExpression(leftExpression, rightExpression)
	case lex.GREATER_THAN:
		return NewGreaterThanExpression(leftExpression, rightExpression)
	case lex.GREATER_OR_EQUAL:
		return NewGreaterOrEqualExpression(leftExpression, rightExpression)
	case lex.LESS_OR_EQUAL:
		return NewLessOrEqualExpression(leftExpression, rightExpression)
	case lex.LESS_THAN:
		return NewLessThanExpression(leftExpression, rightExpression)
	default:
		// TODO: should panic
		return nil
	}
}

// additive_term = multiplicative_term, { ("+" | "-"), multiplicative_term } ;
func (p *Parser) parsePlusOrMinus() Expression {
	leftExpression := p.parseMultiplyCondition()
	if leftExpression == nil {
		return nil
	}

	for p.token.Type == lex.PLUS || p.token.Type == lex.MINUS {
		operation := p.token.Type
		p.consumeToken()
		rightExpression := p.parseMultiplyCondition()
		if rightExpression == nil {
			panic(fmt.Sprintf(ERROR_MISSING_EXPRESSION, p.token.Position.Column, p.token.Position.Line, "+ or -"))
		}

		switch operation {
		case lex.PLUS:
			leftExpression = NewSumExpression(leftExpression, rightExpression)
		case lex.MINUS:
			leftExpression = NewSubstractExpression(leftExpression, rightExpression)
		}
	}

	return leftExpression
}

// multiplicative_term = casted_term, { ("*" | "/"), casted_term } ;
func (p *Parser) parseMultiplyCondition() Expression {
	leftExpression := p.parseCastCondition()
	if leftExpression == nil {
		return nil
	}

	for p.token.Type == lex.MULTIPLY || p.token.Type == lex.DIVIDE {
		operation := p.token.Type
		p.consumeToken()
		rightExpression := p.parseCastCondition()

		if rightExpression == nil {
			panic(fmt.Sprintf(ERROR_MISSING_EXPRESSION, p.token.Position.Column, p.token.Position.Line, "* or /"))
		}

		switch operation {
		case lex.MULTIPLY:
			leftExpression = NewMultiplyExpression(leftExpression, rightExpression)
		case lex.DIVIDE:
			leftExpression = NewDivideExpression(leftExpression, rightExpression)
		}
	}

	return leftExpression
}

// casted_term = unary_operator, [ "as", type_annotation ] ;
func (p *Parser) parseCastCondition() Expression {
	unaryTerm := p.parseUnaryOperator()

	// TODO: nie moze byc nil, to oznacza ze nic nie sparsowalismy nie?
	if unaryTerm == nil {
		return nil
	}

	if p.token.Type != lex.AS {
		return unaryTerm
	}

	typeToken := p.parseTypeAnnotation()
	if typeAnnotation, ok := validTypes[*typeToken]; !ok {
		panic(fmt.Sprintf(SYNTAX_ERROR_NO_TYPE_IN_CAST, p.token.Position.Column, p.token.Position.Line))
	} else {
		// TODO: czy na pewno casted terma moge zwracac jako expression???
		return NewCastExpression(unaryTerm, typeAnnotation)
	}
}

// unary_operator = [ ("-" | "!") ], term ;
func (p *Parser) parseUnaryOperator() Expression {
	if p.token.Type != lex.MINUS && p.token.Type != lex.NEGATE {
		return p.parseTerm()
	}

	p.consumeToken()
	term := p.parseTerm()

	// TODO: should i check whether it is int float or different, if int and token is "!" == syntax error?
	return NewNegateExpression(term)
}

// term = integer | float | bool | string | identifier_or_call | "(" , expression , ")" ;
func (p *Parser) parseTerm() Expression {
	var value any
	switch p.token.Type {
	case lex.IDENTIFIER:
		name, position := p.parseIdentifier()
		identifier := NewIdentifier(name, position)
		if functionCall := p.parseFunctionCall(identifier); functionCall != nil {
			return functionCall
		}
		return identifier
	case lex.CONST_INT:
		value = p.token.Value.(int)
		p.consumeToken()
		return value
	case lex.CONST_FLOAT:
		value = p.token.Value.(float64)
		p.consumeToken()
		return value
	case lex.CONST_TRUE:
		p.consumeToken()
		return true
	case lex.CONST_FALSE:
		p.consumeToken()
		return false
	case lex.CONST_STRING:
		value = p.token.Value.(string)
		p.consumeToken()
		return value
	case lex.LEFT_PARENTHESIS:
		p.consumeToken()
		value = p.parseExpression()
		if value == nil {
			panic(fmt.Sprintf(ERROR_MISSING_EXPRESSION, p.token.Position.Column, p.token.Position.Line, p.token.Type.TypeName()))
		}
		p.requierAndConsume(lex.RIGHT_PARENTHESIS, SYNTAX_ERROR_NO_RIGHT_PARENTHESIS_IN_NESTED_EXPRESSION)
		return value
	default:
		panic(fmt.Sprintf(SYNTAX_ERROR_NO_TERM, p.token.Position.Column, p.token.Position.Line))
	}
}

// conditional_statement = "if" , expression , block , [ "else" , block ] ;
func (p *Parser) parseConditionalStatement() *IfStatement {
	if p.token.Type == lex.IF {
		p.consumeToken()
	}
	condition := p.parseExpression()

	if condition == nil {
		panic(fmt.Sprintf(ERROR_MISSING_EXPRESSION, p.token.Position.Column, p.token.Position.Line, "if"))
	}

	instructions := p.parseBlock()

	if p.token.Type != lex.ELSE {
		return NewIfStatement(condition, instructions, nil)
	}
	p.consumeToken()

	elseInstructions := p.parseBlock()

	return NewIfStatement(condition, instructions, elseInstructions)
}

// loop_statement = "while" , expression, block ;
func (p *Parser) parseWhileStatement() *WhileStatement {
	if p.token.Type != lex.WHILE {
		return nil
	}
	p.consumeToken()

	condition := p.parseExpression()
	if condition == nil {
		panic(fmt.Sprintf(ERROR_MISSING_EXPRESSION, p.token.Position.Column, p.token.Position.Line, lex.WHILE.TypeName()))
	}
	instructions := p.parseBlock()

	return NewWhileStatement(condition, instructions)
}

// switch_statement = "switch", ( variable_declaration, { ",", variable_declaraion } ) | expression, "{", switch_case, { ",", switch_case "}" ;
func (p *Parser) parseSwitchStatement() {}

// switch_case = ( ( [relation_operator], expression ) | "default" ), "=>", ( expression | block ), } ;
func (p *Parser) parseSwitchCase() {}

// return_statement = "return" , [ expression ] ;
func (p *Parser) parseReturnStatement() {}
