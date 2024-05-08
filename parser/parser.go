package parser

import (
	"fmt"
	. "tkom/ast"
	lex "tkom/lexer"
)

func (p *Parser) recoverFromPanic() {
	if err := recover(); err != nil {
		p.ErrorHandler(err.(error))
	}
}

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
}

func (p *Parser) requierAndConsume(tokenType lex.TokenType, syntaxErrMessage string) lex.Token {
	token := p.token
	if token.Type != tokenType {
		panic(NewParserError(fmt.Sprintf(syntaxErrMessage, token.Position.Line, token.Position.Column)))
	}
	p.consumeToken()
	return token
}

// program = { function_definition } ;
func (p *Parser) ParseProgram() *Program {
	defer p.recoverFromPanic()

	functions := map[string]*FunDef{}

	for funDef := p.parseFunDef(); funDef != nil; funDef = p.parseFunDef() {
		if f, ok := functions[funDef.Name]; ok {
			tokenCol := p.token.Position.Column
			tokenLine := p.token.Position.Line
			panic(NewParserError(fmt.Sprintf(SYNTAX_ERROR_FUNCTION_REDEFINITION, tokenCol, tokenLine, f.Position.Line, f.Position.Column)))
		} else {
			functions[funDef.Name] = funDef
		}
	}

	if p.token.Type != lex.ETX {
		panic(NewParserError(fmt.Sprintf(ERROR_NO_ETX_TOKEN, p.token.Position.Line, p.token.Position.Line)))
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

	var funcType *TypeAnnotation
	if t, ok := ValidTypeAnnotation[p.token.Type]; !ok { // p.parseTypeAnnotation()
		funcType = nil
	} else {
		funcType = &t
	}
	block := p.parseBlock()
	if block == nil {
		panic(NewParserError(fmt.Sprintf(SYNTA_ERROR_NO_BLOCK_DEFINED, p.token.Position.Line, p.token.Position.Column)))
	}

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
			panic(NewParserError(fmt.Sprintf(SYNTAX_ERROR_NO_PARAMETERS_AFTER_COMMA, p.token.Position.Line, p.token.Position.Column)))
		}
	}
	return parameters
}

// (a, b int)

// parameter_group = identifier , { ",", identifier }, type_annotation ;
func (p *Parser) parseParameterGroup() []*Variable {
	if p.token.Type != lex.IDENTIFIER {
		return nil
	}

	type Parameter struct {
		Name     string
		Position lex.Position
	}

	name := p.token.Value.(string)
	possition := p.token.Position

	namesAndPositions := []Parameter{}
	namesAndPositions = append(namesAndPositions, Parameter{Name: name, Position: possition})

	p.consumeToken()
	for p.token.Type == lex.COMMA {
		p.consumeToken()
		if p.token.Type != lex.IDENTIFIER {
			panic(NewParserError(fmt.Sprintf(SYNTAX_ERROR_NO_IDENTIFIER, p.token.Position.Line, p.token.Position.Column)))
		}
		name := p.token.Value.(string)
		possition := p.token.Position
		namesAndPositions = append(namesAndPositions, Parameter{Name: name, Position: possition})
		p.consumeToken()
	}
	paramsType := p.parseTypeAnnotation()

	if paramsType == nil {
		panic(NewParserError(fmt.Sprintf(SYNTAX_ERROR_NO_TYPE, p.token.Position.Line, p.token.Position.Column)))
	}
	params := []*Variable{}

	for _, t := range namesAndPositions {
		params = append(params, NewVariable(*paramsType, t.Name, nil, t.Position))
	}
	return params
}

// type_annotation = "int" | "float" | "bool" | "str" ;
func (p *Parser) parseTypeAnnotation() *TypeAnnotation {
	var typeAnnotation *TypeAnnotation
	if t, ok := ValidTypeAnnotation[p.token.Type]; !ok {
		return nil
	} else {
		p.consumeToken()
		typeAnnotation = &t
	}
	return typeAnnotation
}

// block = "{" , { statement } , "}" ;
func (p *Parser) parseBlock() []Statement {
	if p.token.Type != lex.LEFT_BRACE {
		return nil
	}
	p.consumeToken()

	statements := []Statement{}

	for statement := p.parseStatement(); statement != nil; statement = p.parseStatement() {
		statements = append(statements, statement)
	}
	// if statment is nil just return empty Statements slice
	p.requierAndConsume(lex.RIGHT_BRACE, SYNTAX_ERROR_EXPECTED_RIGHT_BRACE)

	return statements
}

// statement = variable_declaration | assigment | conditional_statement | loop_statement | switch_statement | return_statement ;
func (p *Parser) parseStatement() Statement {
	if statement := p.parseVariableDeclaration(); statement != nil {
		return statement
	}
	if statement := p.parseAssignment(); statement != nil {
		return statement
	}
	if statement := p.parseConditionalStatement(); statement != nil {
		return statement
	}
	if statement := p.parseWhileStatement(); statement != nil {
		return statement
	}
	if statement := p.parseSwitchStatement(); statement != nil {
		return statement
	}
	if statement := p.parseReturnStatement(); statement != nil {
		return statement
	}
	return nil
}

// variable_declaration  = type_annotation, identifier, ":=", expression ;
func (p *Parser) parseVariableDeclaration() Statement {
	typeAnnotation := p.parseTypeAnnotation()
	if typeAnnotation == nil {
		return nil
	}
	identifierToken := p.requierAndConsume(lex.IDENTIFIER, SYNTAX_ERROR_NO_VARIABLE_IDETIFIER)
	p.requierAndConsume(lex.DECLARE, SYNTAX_ERROR_MISSING_COLON_ASSIGN)

	expression := p.parseExpression()

	if expression == nil {
		panic(NewParserError(fmt.Sprintf(SYNTAX_ERROR_NO_EXPRESSION_IN_VARIABLE_DECLARATION, p.token.Position.Line, p.token.Position.Column)))
	}

	name := identifierToken.Value.(string)
	position := identifierToken.Position
	variable := NewVariable(*typeAnnotation, name, expression, position)
	return variable
}

// assigment = identifier_or_call,  [ "=", expression ] ;
func (p *Parser) parseAssignment() Statement {
	if p.token.Type != lex.IDENTIFIER {
		return nil
	}

	name := p.token.Value.(string)
	position := p.token.Position
	p.consumeToken()

	if functionCall := p.parseFunctionCall(name, position); functionCall != nil {
		return functionCall
	}

	if p.token.Type != lex.ASSIGN {
		return NewIdentifier(name, position)
	}

	p.consumeToken()

	expression := p.parseExpression()

	return NewAssignment(NewIdentifier(name, position), expression)
}

// identifier_or_call = identifier, [ "(", [ argumets ], ")" ] ;
func (p *Parser) parseIdentifierOrCall() Statement {
	if p.token.Type != lex.IDENTIFIER {
		return nil
	}

	name := p.token.Value.(string)
	position := p.token.Position
	p.consumeToken()

	if functionCall := p.parseFunctionCall(name, position); functionCall != nil {
		return functionCall
	}

	return NewIdentifier(name, position)
}

func (p *Parser) parseFunctionCall(name string, position lex.Position) Statement {
	if p.token.Type != lex.LEFT_PARENTHESIS {
		return nil
	}
	p.consumeToken()

	arguments := p.parseArguments()

	p.requierAndConsume(lex.RIGHT_PARENTHESIS, SYNTAX_ERROR_FUNC_CALL_NOT_CLOSED)

	return NewFunctionCall(name, position, arguments)
}

// arguments = expression , { "," , expression } ;
func (p *Parser) parseArguments() []Expression { // return []Variable czy []Expression ??
	expressions := []Expression{}

	expression := p.parseExpression()
	if expression == nil {
		return expressions
	}

	expressions = append(expressions, expression)

	for p.token.Type == lex.COMMA {
		expression := p.parseExpression()
		if expression == nil {
			panic(NewParserError(fmt.Sprintf(ERROR_MISSING_EXPRESSION, p.token.Position.Line, p.token.Position.Column, p.token.Type.TypeName())))
		}
		expressions = append(expressions, expression)
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
		position := p.token.Position
		p.consumeToken()
		rightExpression := p.parseAndCondition()
		if rightExpression == nil {
			panic(NewParserError(fmt.Sprintf(ERROR_MISSING_EXPRESSION, p.token.Position.Line, p.token.Position.Column, "OR")))
		}
		// leftExpression = NewExpression(leftExpression, OR, rightExpression)
		leftExpression = NewOrExpression(leftExpression, rightExpression, position)
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
		position := p.token.Position
		p.consumeToken()
		rightExpression := p.parseRelationCondition()
		if rightExpression == nil {
			panic(NewParserError(fmt.Sprintf(ERROR_MISSING_EXPRESSION, p.token.Position.Line, p.token.Position.Column, "AND")))
		}
		// leftExpression = NewExpression(leftExpression, AND, rightExpression)
		leftExpression = NewAndExpression(leftExpression, rightExpression, position)
	}

	return leftExpression
}

// relation_term = additive_term, [ relation_operator, additive_term ] ;
func (p *Parser) parseRelationCondition() Expression {
	leftExpression := p.parsePlusOrMinus()
	if leftExpression == nil {
		return nil
	}

	operatorToFactory := map[lex.TokenType]func(Expression, Expression, lex.Position) Expression{
		lex.EQUALS:           NewEqualsExpression,
		lex.NOT_EQUALS:       NewNotEqualsExpression,
		lex.GREATER_THAN:     NewGreaterThanExpression,
		lex.GREATER_OR_EQUAL: NewGreaterOrEqualExpression,
		lex.LESS_OR_EQUAL:    NewLessOrEqualExpression,
		lex.LESS_THAN:        NewLessThanExpression,
	}

	if factory, ok := operatorToFactory[p.token.Type]; ok {
		operationType := p.token.Type.TypeName()
		position := p.token.Position
		p.consumeToken()

		rightExpression := p.parsePlusOrMinus()
		if rightExpression == nil {
			panic(NewParserError(fmt.Sprintf(ERROR_MISSING_EXPRESSION, p.token.Position.Line, p.token.Position.Column, operationType)))
		}

		leftExpression = factory(leftExpression, rightExpression, position)
	}

	return leftExpression
}

// additive_term = multiplicative_term, { ("+" | "-"), multiplicative_term } ;
func (p *Parser) parsePlusOrMinus() Expression {
	leftExpression := p.parseMultiplyCondition()
	if leftExpression == nil {
		return nil
	}

	additiveToFactory := map[lex.TokenType]func(Expression, Expression, lex.Position) Expression{
		lex.PLUS:  NewSumExpression,
		lex.MINUS: NewSubstractExpression,
	}

	for {
		if factory, ok := additiveToFactory[p.token.Type]; ok {
			position := p.token.Position
			p.consumeToken()
			rightExpression := p.parseMultiplyCondition()
			if rightExpression == nil {
				panic(NewParserError(fmt.Sprintf(ERROR_MISSING_EXPRESSION, p.token.Position.Line, p.token.Position.Column, "additive operator")))
			}
			leftExpression = factory(leftExpression, rightExpression, position)
		} else {
			return leftExpression
		}
	}
}

// multiplicative_term = casted_term, { ("*" | "/"), casted_term } ;
func (p *Parser) parseMultiplyCondition() Expression {
	leftExpression := p.parseCastCondition()
	if leftExpression == nil {
		return nil
	}

	multiplicativeToFactory := map[lex.TokenType]func(Expression, Expression, lex.Position) Expression{
		lex.MULTIPLY: NewMultiplyExpression,
		lex.DIVIDE:   NewDivideExpression,
	}

	for {
		if factory, ok := multiplicativeToFactory[p.token.Type]; ok {
			position := p.token.Position
			p.consumeToken()
			rightExpression := p.parseCastCondition()
			if rightExpression == nil {
				panic(NewParserError(fmt.Sprintf(ERROR_MISSING_EXPRESSION, p.token.Position.Line, p.token.Position.Column, "* or /")))
			}
			leftExpression = factory(leftExpression, rightExpression, position)
		} else {
			return leftExpression
		}
	}
}

// casted_term = unary_operator, [ "as", type_annotation ] ;
func (p *Parser) parseCastCondition() Expression {
	unaryTerm := p.parseUnaryOperator()

	if unaryTerm == nil {
		return nil
	}

	if p.token.Type != lex.AS {
		return unaryTerm
	}
	position := p.token.Position

	typeAnnotation := p.parseTypeAnnotation()
	if typeAnnotation == nil {
		panic(NewParserError(fmt.Sprintf(SYNTAX_ERROR_NO_TYPE_IN_CAST, p.token.Position.Line, p.token.Position.Column)))
	} else {
		return NewCastExpression(unaryTerm, typeAnnotation, position)
	}
}

// unary_operator = [ ("-" | "!") ], term ;
func (p *Parser) parseUnaryOperator() Expression {
	if p.token.Type != lex.MINUS && p.token.Type != lex.NEGATE {
		return p.parseTerm()
	}

	position := p.token.Position
	p.consumeToken()
	term := p.parseTerm()

	return NewNegateExpression(term, position)
}

// term = integer | float | bool | string | identifier_or_call | "(" , expression , ")" ;
func (p *Parser) parseTerm() Expression {
	switch p.token.Type {
	case lex.IDENTIFIER:
		value := p.parseIdentifierOrCall()
		if value == nil {
			panic(NewParserError(fmt.Sprintf(ERROR_MISSING_EXPRESSION, p.token.Position.Line, p.token.Position.Column, p.token.Type.TypeName())))
		}
		expr, ok := value.(Expression)
		if !ok {
			// Handle the case where the value couldn't be converted to Expression
			panic(NewParserError(fmt.Sprintf("Unexpected type: %T", value)))
		}
		return expr
		// return value
	case lex.CONST_INT:
		value := p.token.Value.(int)
		position := p.token.Position
		p.consumeToken()
		return NewIntExpression(value, position)
	case lex.CONST_FLOAT:
		value := p.token.Value.(float64)
		position := p.token.Position
		p.consumeToken()
		return NewFloatExpression(value, position)
	case lex.CONST_TRUE:
		position := p.token.Position
		p.consumeToken()
		return NewBoolExpression(true, position)
	case lex.CONST_FALSE:
		position := p.token.Position
		p.consumeToken()
		return NewBoolExpression(false, position)
	case lex.CONST_STRING:
		value := p.token.Value.(string)
		position := p.token.Position
		p.consumeToken()
		return NewStringExpression(value, position)
	case lex.LEFT_PARENTHESIS:
		p.consumeToken()
		expression := p.parseExpression()
		if expression == nil {
			panic(NewParserError(fmt.Sprintf(ERROR_MISSING_EXPRESSION, p.token.Position.Line, p.token.Position.Column, p.token.Type.TypeName())))
		}
		p.requierAndConsume(lex.RIGHT_PARENTHESIS, SYNTAX_ERROR_NO_RIGHT_PARENTHESIS_IN_NESTED_EXPRESSION)
		return expression
	default:
		return nil
	}
}

// conditional_statement = "if" , expression , block , [ "else" , block ] ;
func (p *Parser) parseConditionalStatement() *IfStatement {
	if p.token.Type != lex.IF {
		return nil
	}
	p.consumeToken()

	condition := p.parseExpression()
	if condition == nil {
		panic(NewParserError(fmt.Sprintf(ERROR_MISSING_EXPRESSION, p.token.Position.Line, p.token.Position.Column, "if")))
	}

	instructions := p.parseBlock()
	if instructions == nil {
		panic(NewParserError(fmt.Sprintf(SYNTAX_ERROR_EMPTY_BLOCK_IN_IF_STATEMENT, p.token.Position.Line, p.token.Position.Column)))
	}

	if p.token.Type != lex.ELSE {
		return NewIfStatement(condition, instructions, nil)
	}
	p.consumeToken()

	elseInstructions := p.parseBlock()
	if elseInstructions == nil {
		panic(NewParserError(fmt.Sprintf(SYNTAX_ERROR_EMPTY_BLOCK_IN_IF_STATEMENT, p.token.Position.Line, p.token.Position.Column)))
	}

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
		panic(NewParserError(fmt.Sprintf(ERROR_MISSING_EXPRESSION, p.token.Position.Line, p.token.Position.Column, lex.WHILE.TypeName())))
	}

	instructions := p.parseBlock()
	if instructions == nil {
		panic(NewParserError(fmt.Sprintf(SYNTAX_ERROR_EMPTY_BLOCK_IN_WHILE_STATEMENT, p.token.Position.Line, p.token.Position.Column)))
	}

	return NewWhileStatement(condition, instructions)
}

func (p *Parser) parseSwitchVariables() (variables []*Variable) {
	variable := p.parseVariableDeclaration()
	if variable == nil {
		return nil
	}
	if variable, ok := variable.(*Variable); ok {
		variables = append(variables, variable)
	} else {
		panic(NewParserError(fmt.Sprintf(SYNTAX_ERROR_BAD_VARIABLE_DECLARATION, p.token.Position.Line, p.token.Position.Column)))
	}

	for p.token.Type == lex.COMMA {
		p.consumeToken()
		variableDeclaration := p.parseVariableDeclaration()
		if variableDeclaration == nil {
			panic(NewParserError(fmt.Sprintf(SYNTAX_ERROR_NO_VARIABLE_AFTER_COMMA, p.token.Position.Line, p.token.Position.Column)))
		}
		if variable, ok := variableDeclaration.(*Variable); ok {
			variables = append(variables, variable)
		} else {
			panic(NewParserError(fmt.Sprintf(SYNTAX_ERROR_BAD_VARIABLE_DECLARATION, p.token.Position.Line, p.token.Position.Column)))
		}
	}
	return variables
}

// switch_statement = "switch", (( variable_declaration, { ",", variable_declaraion } ) | expression ), "{", switch_case, { ",", switch_case "}" ;
func (p *Parser) parseSwitchStatement() *SwitchStatement {
	if p.token.Type != lex.SWITCH {
		return nil
	}
	p.consumeToken()

	var expression Expression

	variables := p.parseSwitchVariables()
	if variables == nil {
		expression = p.parseExpression()
	}
	// if expression nil it means that the switch case is empty and we allow it

	p.requierAndConsume(lex.LEFT_BRACE, SYNTAX_ERROR_NO_LEFT_CURLY_BRACKET_IN_SWITCH)

	cases := []Statement{}

	caseStatement := p.parseSwitchCase()
	if caseStatement == nil {
		panic(NewParserError(fmt.Sprintf(ERROR_MISSING_SWITCH_CASE, p.token.Position.Line, p.token.Position.Column)))
	}
	cases = append(cases, caseStatement)

	for p.token.Type == lex.COMMA {
		p.consumeToken()
		caseStatement := p.parseSwitchCase()
		if caseStatement == nil {
			panic(NewParserError(fmt.Sprintf(ERROR_MISSING_SWITCH_CASE, p.token.Position.Line, p.token.Position.Column)))
		}
		cases = append(cases, caseStatement)
	}

	p.requierAndConsume(lex.RIGHT_BRACE, SYNTAX_ERROR_NOT_CLOSED_SWITCH)

	return NewSwitchStatement(variables, expression, cases)
}

// switch_case = ( ( [relation_operator], expression ) | "default" ), "=>", ( expression | block ), } ;
func (p *Parser) parseSwitchCase() Statement {
	if p.token.Type == lex.DEFAULT {
		p.consumeToken()
		p.requierAndConsume(lex.CASE_ARROW, SYNTAX_ERROR_NO_ARROW)
		expression := p.parseExpression()

		return NewDefaultCase(expression)
	}

	expression := p.parseExpression()

	if expression == nil {
		panic(NewParserError(fmt.Sprintf(ERROR_MISSING_SWITCH_CASE, p.token.Position.Line, p.token.Position.Column)))
	}

	p.requierAndConsume(lex.CASE_ARROW, SYNTAX_ERROR_NO_ARROW)

	instructions := p.parseExpression()
	return NewSwitchCase(expression, instructions)
}

// return_statement = "return" , [ expression ] ;
func (p *Parser) parseReturnStatement() *ReturnStatement {
	if p.token.Type != lex.RETURN {
		return nil
	}
	p.consumeToken()

	expression := p.parseExpression()
	return NewReturnStatement(expression)
}
