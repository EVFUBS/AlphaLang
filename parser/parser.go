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

	println("PROGRAM COMPLETE")
	return &program
}

func (p *Parser) Parse() *ast.AstStatement {
	var statement ast.AstStatement

	p.AdvanceToken()

	switch p.curToken.Type {
	case token.VAR:
		statement = p.ParseVarStatement()
	case token.RETURN:
		statement = p.ParseReturnStatement()
	case token.IF:
		statement = p.ParseIfStatement()
	default:
		statement = p.ParseExpressionStatement()
	}

	return &statement
}

//Expression Parsing
func (p *Parser) ParseExpression() ast.AstExpression {
	var node ast.AstExpression

	switch p.curToken.Type {
	case token.IDENT:
		node = p.ParseIdentiferLiteral()
	case token.INTEGER:
		node = p.ParseIntegerLiteral()
	case token.TRUE:
		node = p.ParseBoolLiteral()
	case token.FALSE:
		node = p.ParseBoolLiteral()
	case token.STRING:
		node = p.ParseStringLiteral()
	}

	if _, ok := InfixExpressions[p.nextToken.Type]; ok {
		node = p.parseInfixExpression(node)
	}

	return node
}

//change to cleaner method later
//add ( infix to check for function calls
var InfixExpressions = map[token.TokenType]string{
	token.PLUS:     "+",
	token.MINUS:    "-",
	token.ASTERISK: "*",
	token.SLASH:    "/",
	token.GTHAN:    "<",
	token.LTHAN:    ">",
}

// need to consider precedence
func (p *Parser) parseInfixExpression(node ast.AstExpression) *ast.InfixExpression {
	left := node
	p.AdvanceToken()
	operator := p.curToken.Literal
	p.AdvanceToken()
	right := p.ParseExpression()
	newNode := &ast.InfixExpression{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
	return newNode
}

//Statement Parsing
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

func (p *Parser) ParseExpressionStatement() *ast.ExpressionStatement {
	var exprStatement ast.ExpressionStatement

	switch p.curToken.Type {
	case token.FUNCTION:
		exprStatement.Expression = p.ParseFunctionLiteral()
	}

	return &exprStatement
}

func (p *Parser) ParseBlockStatement() *ast.BlockStatement {
	var statements ast.BlockStatement
	for p.nextToken.Type != token.RBRACE {
		statement := p.Parse()
		statements.Statements = append(statements.Statements, *statement)
	}

	return &statements
}

// make sure expressions rseult in a condition ie < not +
func (p *Parser) ParseIfStatement() *ast.IfStatement {
	var IfStatement ast.IfStatement

	p.AdvanceToken()
	IfStatement.If.Condition = p.ParseExpression()
	p.CheckTokenAdvance(token.LBRACE)
	IfStatement.If.Consequence = *p.ParseBlockStatement()
	p.CheckTokenAdvance(token.RBRACE)

	if p.nextTokenIs(token.ELIF) {
		for {
			var conditional ast.Conditional

			p.CheckTokenAdvance(token.ELIF)
			p.AdvanceToken()
			conditional.Condition = p.ParseExpression()
			p.CheckTokenAdvance(token.LBRACE)
			conditional.Consequence = *p.ParseBlockStatement()
			p.CheckTokenAdvance(token.RBRACE)

			IfStatement.Elif = append(IfStatement.Elif, conditional)

			if !p.nextTokenIs(token.ELIF) {
				break
			}
		}
	}

	if p.nextTokenIs(token.ELSE) {
		p.AdvanceToken()
		p.CheckTokenAdvance(token.LBRACE)
		IfStatement.Else = *p.ParseBlockStatement()
		p.CheckTokenAdvance(token.RBRACE)
	}

	return &IfStatement

}

//Literal Parsing
func (p *Parser) ParseIdentiferLiteral() *ast.IdentiferLiteral {
	return &ast.IdentiferLiteral{
		Token: *p.curToken,
		Ident: p.curToken.Literal,
	}
}

func (p *Parser) ParseIntegerLiteral() *ast.IntegerLiteral {
	intVal, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		//return some kind of error
	}
	return &ast.IntegerLiteral{
		Token: *p.curToken,
		Value: intVal,
	}
}

func (p *Parser) ParseBoolLiteral() *ast.BooleanLiteral {
	if p.curToken.Type == token.TRUE {
		return &ast.BooleanLiteral{
			Token: *p.curToken,
			Value: true,
		}
	}
	return &ast.BooleanLiteral{
		Token: *p.curToken,
		Value: false,
	}
}

func (p *Parser) ParseStringLiteral() *ast.StringLiteral {
	return &ast.StringLiteral{
		Token: *p.curToken,
		Value: string(p.curToken.Literal),
	}
}

func (p *Parser) ParseFunctionLiteral() *ast.FunctionLiteral {
	var function ast.FunctionLiteral

	//p.AdvanceToken()
	p.CheckTokenAdvance(token.IDENT)
	function.Name = p.curToken.Literal
	p.AdvanceToken()
	p.AdvanceToken()

	for p.curToken.Type != token.RPAREN {
		if p.curToken.Type == token.COMMA {
			p.AdvanceToken()
		}
		if p.curToken.Type == token.IDENT {
			function.Parameters = append(function.Parameters, *p.ParseIdentiferLiteral())
			p.AdvanceToken()
		}
	}
	p.AdvanceToken()
	function.Body = *p.ParseBlockStatement()
	p.AdvanceToken()
	//println(function.Body.String())

	return &function
}

//Helper Functions
func (p *Parser) CheckTokenAdvance(wanted token.TokenType) {
	if p.nextToken.Type != wanted {
		//some kind of error handle
		println("we got a problem here", "wanted: ", wanted, "got: ", p.nextToken.Type)
	}
	p.AdvanceToken()
}

func (p *Parser) nextTokenIs(wanted token.TokenType) bool {
	if p.nextToken.Type == wanted {
		return true
	}
	return false
}

// func test(a,b){return a}
