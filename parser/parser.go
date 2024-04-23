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

func (p *Parser) mustBe(tokenType lex.TokenTypes, syntaxErrMsg string) any {
	token := p.token
	if token.Type != tokenType {
		panic(fmt.Errorf(syntaxErrMsg, token.Position.Column, token.Position.Line))
	}
	value := p.token.Value

	p.consumeToken()
	return value
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

	// wiemy że to string, ale musze tak zrobić żeby compilator był szczęśliwy
	name := p.token.Value.(string)
	pos := p.token.Position

	p.consumeToken()

	_ = p.mustBe(lex.LEFT_PARENTHESIS, SYNTAX_ERROR_FUNC_DEF_NO_PARENTHASIS)
	params := p.parseParameters()
	_ = p.mustBe(lex.RIGHT_PARENTHESIS, SYNTAX_ERROR_FUNC_DEF_NO_PARENTHASIS)
	funcType := p.parseType()
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
	params = append(params, newParameter(p.token.Value.(string), p.token.Position))

	p.consumeToken()
	// param1, param2, param3 string
	for p.token.Type == lex.COMMA {
		p.consumeToken()
		if p.token.Type != lex.IDENTIFIER {
			panic(fmt.Errorf(SYNTAX_ERROR_NO_IDENTIFIER, p.token.Position.Column, p.token.Position.Line))
		}
		params = append(params, newParameter(p.token.Value.(string), p.token.Position))
		p.consumeToken()
	}

	paramsType := p.parseType()

	if paramsType == nil {
		panic(fmt.Errorf(SYNTAX_ERROR_NO_TYPE, p.token.Position.Column, p.token.Position.Line))
	}

	// warning ma racje, ale musze tak zrobic bo jak inaczej dodam typ dla każdego identifiera z grupy huh?
	for i := range params {
		params[i].Type = *paramsType
	}

	return params
}

// type_annotation = "int" | "float" | "bool" | "str" ;
func (p *Parser) parseType() *lex.TokenTypes {
	switch token := p.token.Type; token {
	case lex.INT, lex.FLOAT, lex.BOOL, lex.STRING:
		typ := p.token.Type
		p.consumeToken()
		return &typ
	}

	// if p.token.Type == lex.INT || p.token.Type == lex.FLOAT || p.token.Type == lex.BOOL || p.token.Type == lex.STRING {
	// 	typ := p.token.Type
	// 	p.consumeToken()
	// 	return &typ
	// }
    
	return nil
}

func (p *Parser) parseBlock() *Block {
	return nil
}
