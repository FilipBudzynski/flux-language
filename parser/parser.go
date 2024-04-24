package parser

import (
	"fmt"
	"reflect"
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
func (p *Parser) parseParameters() []Variable {
	paramGroup := p.parseParameterGroup()
	if paramGroup == nil {
		return nil
	}
	parameters := []Variable{}
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
func (p *Parser) parseParameterGroup() []Variable {
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
	params := []Variable{}

	for _, t := range namesAndPositions {
		params = append(params, newVariable(*paramsType, newIdentifier(t.Name, t.Position), nil))
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
func (p *Parser) parseBlock() Block {
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

	return Block{Statements: statements}
}

// statement = variable_declaration | assigment | conditional_statement | loop_statement | switch_statement | return_statement ;
func (p *Parser) parseStatement() Statement {
	switch p.token.Type {
	case lex.INT, lex.FLOAT, lex.BOOL, lex.STRING:
		return p.parseVariableDeclaration()
	case lex.IDENTIFIER:
		return p.parseIdentifierOrCall()
	// case lex.IF:
	// 	return p.parseConditionalStatement()
	// case lex.WHILE:
	// 	return p.parseLoopStatement()
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
	identifierToken := p.requierAndConsume(lex.IDENTIFIER, SYNTAX_ERROR_NO_VARIABLE_IDETIFIER)
	p.requierAndConsume(lex.DECLARE, SYNTAX_ERROR_MISSING_COLON_ASSIGN)
	expression := p.parseExpression()

	// na pewno???
	// if expression == nil {
	//     panic(fmt.Sprintf(SYNTAX_ERROR_NO_EXPRESSION_IN_DECLARATION, p.token.Position.Column, p.token.Position.Line))
	// }

	name := identifierToken.Value.(string)
	position := identifierToken.Position
	identifier := newIdentifier(name, position)
	variable := newVariable(*typeAnnotation, identifier, expression)
	return variable
}

// assigment = identifier_or_call,  [ "=", expression ] ;
func (p *Parser) parseAssignment() Statement {
	identifierOrCall := p.parseIdentifierOrCall()

	if p.token.Type != lex.ASSIGN {
		return identifierOrCall
	}

	// illigal assigment to function call
	if reflect.TypeOf(identifierOrCall) != reflect.TypeOf(Identifier{}) {
		panic(fmt.Sprintf(ERROR_ASIGNMENT_TO_FUNCTION_CALL, p.token.Position.Column, p.token.Position.Line))
	}

	return identifierOrCall
}

// identifier_or_call = identifier, [ "(", [ argumets ], ")" ] ;
func (p *Parser) parseIdentifierOrCall() Statement {
	if p.token.Type != lex.IDENTIFIER {
		return nil
	}
	identifier := newIdentifier(p.token.Value.(string), p.token.Position)

	if statement := p.parseFunctionCall(identifier); statement == nil {
		return identifier
	}
	return newFunctionCall(identifier, nil, identifier.Position)
}

func (p *Parser) parseFunctionCall(identifier Identifier) Statement {
	if p.token.Type != lex.RIGHT_PARENTHESIS {
		return nil
	}
	p.consumeToken()
	arguments := p.parseArguments()

	if p.token.Type != lex.RIGHT_PARENTHESIS {
		panic(fmt.Sprintf(SYNTAX_ERROR_FUNC_CALL_NOT_CLOSED, p.token.Position.Column, p.token.Position.Line))
	}

	return newFunctionCall(identifier, arguments, identifier.Position)
}

// arguments = expression , { "," , expression } ;
func (p *Parser) parseArguments() []Variable { // return []Variable czy []Expression ??
	validExpressionTypes := map[lex.TokenType]bool{
		lex.CONST_INT:    true,
		lex.CONST_FLOAT:  true,
		lex.CONST_FALSE:  true,
		lex.CONST_TRUE:   true,
		lex.CONST_STRING: true,
		lex.IDENTIFIER:   true,
	}
	if !validExpressionTypes[p.token.Type] {
		return nil
	}

	expressions := []Variable{}
	expression := p.parseExpression()
	expressions = append(expressions, expression)

	for p.token.Type == lex.COMMA {
		p.consumeToken()
	}
	return expressions
}

// conditional_statement = "if" , expression , block , [ "else" , block ] ;
func (p *Parser) parseConditionalStatement() {}

// loop_statement = "while" , expression, block ;
func (p *Parser) parseLoopStatement() {}

// switch_statement = "switch", ( variable_declaration, { ",", variable_declaraion } ) | expression, "{", switch_case, { ",", switch_case "}" ;
func (p *Parser) parseSwitchStatement() {}

// switch_case = ( ( [relation_operator], expression ) | "default" ), "=>", ( expression | block ), } ;
func (p *Parser) parseSwitchCase() {}

// return_statement = "return" , [ expression ] ;
func (p *Parser) parseReturnStatement() {}

// expression = conjunction_term, { "or", conjunction_term } ;
func (p *Parser) parseExpression() Variable { return Variable{} }
