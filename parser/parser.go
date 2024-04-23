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
}

func (p *Parser) requierAndConsume(tokenType lex.TokenTypes, syntaxErrMessage string) lex.Token {
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
	name, pos := p.parseIdentifier()
	p.requierAndConsume(lex.LEFT_PARENTHESIS, SYNTAX_ERROR_FUNC_DEF_NO_PARENTHASIS)
	params := p.parseParameters()
	p.requierAndConsume(lex.RIGHT_PARENTHESIS, SYNTAX_ERROR_FUNC_DEF_NO_PARENTHASIS)
	funcType := p.parseTypeAnnotation()
	block := p.parseBlock()
	if block == nil {
		panic(fmt.Errorf(SYNTAX_ERROR_NO_BLOCK, p.token.Position.Column, p.token.Position.Line))
	}

	return NewFunctionDefinition(name, params, funcType, block, pos)
}

// parameters = parameter_group , { "," , parameter_group } ;
func (p *Parser) parseParameters() []Parameter {
	paramGroup := p.parseParameterGroup()
	if paramGroup == nil {
		return nil
	}
	parameters := []Parameter{}
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
func (p *Parser) parseParameterGroup() []Parameter {
	if p.token.Type != lex.IDENTIFIER {
		return nil
	}
	params := []Parameter{}
	params = append(params, NewParameter(p.parseIdentifier()))

	for p.token.Type == lex.COMMA {
		p.consumeToken()
		if p.token.Type != lex.IDENTIFIER {
			panic(fmt.Errorf(SYNTAX_ERROR_NO_IDENTIFIER, p.token.Position.Column, p.token.Position.Line))
		}
		params = append(params, NewParameter(p.parseIdentifier()))
	}

	paramsType := p.parseTypeAnnotation()

	if paramsType == nil {
		panic(fmt.Errorf(SYNTAX_ERROR_NO_TYPE, p.token.Position.Column, p.token.Position.Line))
	}

	for i := range params {
		params[i].Type = *paramsType
	}
	return params
}

// type_annotation = "int" | "float" | "bool" | "str" ;
func (p *Parser) parseTypeAnnotation() *lex.TokenTypes {
	switch token := p.token.Type; token {
	case lex.INT, lex.FLOAT, lex.BOOL, lex.STRING:
		typ := p.token.Type
		p.consumeToken()
		return &typ
	}

	return nil
}

// block = "{" , { statement } , "}" ;
func (p *Parser) parseBlock() *Block {
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

	return &Block{Statements: statements}
}

// statement = variable_declaration | assigment | conditional_statement | loop_statement | switch_statement | return_statement ;
func (p *Parser) parseStatement() Statement {
	switch p.token.Type {
	case lex.INT, lex.FLOAT, lex.BOOL, lex.STRING:
		return p.parseVariableDeclaration()
	// case lex.IDENTIFIER:
	// 	return p.parseAssignment()
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
	identifier := p.requierAndConsume(lex.IDENTIFIER, SYNTAX_ERROR_NO_VARIABLE_IDETIFIER)
	p.requierAndConsume(lex.DECLARE, SYNTAX_ERROR_MISSING_COLON_ASSIGN)
	expression := p.parseExpression()

	variable := newVariable(*typeAnnotation, identifier.Value.(string), expression)
	return &variable
}

// identifier_or_call = identifier, [ "(", [ argumets ], ")" ] ;
func (p *Parser) parseIdentifierOrCall() {}

// assigment = identifier_or_call,  [ "=", expression ] ;
func (p *Parser) parseAssignment() {}

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
func (p *Parser) parseExpression() int { return 0 }
