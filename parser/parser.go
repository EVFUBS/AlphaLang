package parser

import (
	"strconv"

	"github.com/EVFUBS/AlphaLang/ast"
	"github.com/EVFUBS/AlphaLang/lexer"
	"github.com/EVFUBS/AlphaLang/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  *token.Token
	nextToken *token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken = p.l.NextToken()
	return p
}

func (p *Parser) AdvanceToken() {
	p.curToken = p.nextToken
	p.nextToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	var program ast.Program

	for p.nextToken.Type != token.EOF {
		statement := p.Parse()
		program.Statements = append(program.Statements, *statement)
	}

	return &program
}

func (p *Parser) Parse() *ast.AstStatement {
	var statement ast.AstStatement

	switch p.nextToken.Type {
	case token.VAR:
		statement = p.ParseVarStatement()
	case token.RETURN:
		statement = p.ParseReturnStatement()
	}

	p.AdvanceToken()

	return &statement
}

func (p *Parser) ParseExpression() ast.AstExpression {
	var node ast.AstExpression

	switch p.nextToken.Type {
	case token.IDENT:
		node = &ast.IdentiferLiteral{
			Ident: p.nextToken.Literal,
		}
	case token.INTEGER:
		node = p.ParseIntegerLiteral()
	case token.TRUE:
		node = p.ParseBoolLiteral()
	case token.FALSE:
		node = p.ParseBoolLiteral()
	}

	return node
}

func (p *Parser) ParseVarStatement() *ast.VarStatement {
	p.AdvanceToken()

	statement := &ast.VarStatement{
		Identifer: p.ParseExpression(),
	}

	p.AdvanceToken()
	p.AdvanceToken()

	statement.Value = p.ParseExpression()

	return statement
}

func (p *Parser) ParseReturnStatement() *ast.ReturnStatement {
	p.AdvanceToken()

	statement := &ast.ReturnStatement{
		ReturnValue: p.ParseExpression(),
	}

	return statement
}

//Literal Parsing
func (p *Parser) ParseIntegerLiteral() *ast.IntegerLiteral {
	intVal, err := strconv.ParseInt(p.nextToken.Literal, 0, 64)
	if err != nil {
		//return some kind of error
	}
	return &ast.IntegerLiteral{
		Value: intVal,
	}
}

func (p *Parser) ParseBoolLiteral() *ast.BooleanLiteral {
	if p.nextToken.Type == token.TRUE {
		return &ast.BooleanLiteral{
			Value: true,
		}
	}
	return &ast.BooleanLiteral{
		Value: false,
	}

}
