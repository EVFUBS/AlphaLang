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
	errors    []ParserError
}

type ParserError string

type InfixFunction func(node ast.AstExpression) ast.AstExpression
type PrefixFunction func() ast.AstExpression

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken = p.l.NextToken()
	return p
}

func (p *Parser) Errors() []ParserError {
	return p.errors
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, ParserError(msg))
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

	p.AdvanceToken()

	switch p.curToken.Type {
	case token.VAR:
		statement = p.ParseVarStatement()
	case token.RETURN:
		statement = p.ParseReturnStatement()
	case token.IF:
		statement = p.ParseIfStatement()
	case token.FOR:
		statement = p.ParseForStatement()
	case token.WHILE:
		statement = p.parseWhileStatement()
	default:
		statement = p.ParseExpressionStatement()
	}

	return &statement
}

// Expression Parsing
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
	case token.LBRACKET:
		node = p.ParseArrayLiteral()
	case token.LBRACE:
		node = p.parseHashLiteral()
	}

	var prefixOperations = map[token.TokenType]PrefixFunction{
		token.BANG:  p.parsePrefixExpression,
		token.MINUS: p.parsePrefixExpression,
	}

	if prefixOperations[p.curToken.Type] != nil {
		node = prefixOperations[p.curToken.Type]()
	}

	// precedence parsing
	var infixOperations = map[token.TokenType]InfixFunction{
		token.PLUS:      p.parseInfixExpression,
		token.MINUS:     p.parseInfixExpression,
		token.ASTERISK:  p.parseInfixExpression,
		token.SLASH:     p.parseInfixExpression,
		token.MODULUS:   p.parseInfixExpression,
		token.EQUAL:     p.parseInfixExpression,
		token.NOTEQUAL:  p.parseInfixExpression,
		token.LTHAN:     p.parseInfixExpression,
		token.GTHAN:     p.parseInfixExpression,
		token.GEQUAL:    p.parseInfixExpression,
		token.LEQUAL:    p.parseInfixExpression,
		token.LPAREN:    p.parseCallExpression,
		token.INCREMENT: p.parseInfixExpression,
		token.DECREMENT: p.parseInfixExpression,
		token.LBRACKET:  p.parseIndexExpression,
	}
	if _, ok := infixOperations[p.nextToken.Type]; ok {
		node = infixOperations[p.nextToken.Type](node)
		if infix, ok := node.(*ast.InfixExpression); ok {
			node = p.sortPrecedence(infix)
		}
		return node
	}

	return node
}

// Infix Parsing
// need to consider precedence
func (p *Parser) parseInfixExpression(node ast.AstExpression) ast.AstExpression {
	left := node
	p.AdvanceToken()
	operator := p.curToken
	p.AdvanceToken()
	right := p.ParseExpression()
	newNode := &ast.InfixExpression{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
	return newNode
}

func (p *Parser) parsePrefixExpression() ast.AstExpression {
	operator := p.curToken
	p.AdvanceToken()
	expression := p.ParseExpression()
	return &ast.PrefixExpression{
		Prefix:     operator,
		Expression: expression,
	}
}

var precedence = map[token.TokenType]int{
	token.PLUS:     1,
	token.MINUS:    2,
	token.ASTERISK: 3,
	token.SLASH:    4,
}

func (p *Parser) sortPrecedence(node *ast.InfixExpression) *ast.InfixExpression {
	// if right node is infix expression
	if right, ok := node.Right.(*ast.InfixExpression); ok {
		// if right node has higher precedence
		if precedence[node.Operator.Type] > precedence[right.Operator.Type] {
			// swap
			node.Right = right.Left
			right.Left = node
			return p.sortPrecedence(right)
		}
	}
	return node
}

func (p *Parser) parseCallExpression(node ast.AstExpression) ast.AstExpression {
	p.AdvanceToken()
	p.AdvanceToken()
	list := p.parseExpressionList(token.RPAREN)
	return &ast.CallExpression{
		Function:  node,
		Arguments: list,
	}
}

func (p *Parser) parseIndexExpression(node ast.AstExpression) ast.AstExpression {
	p.AdvanceToken()
	p.AdvanceToken()
	index := p.ParseExpression()
	p.AdvanceToken()
	return &ast.IndexExpression{
		Left:  node,
		Index: index,
	}
}

func (p *Parser) parseHashLiteral() ast.AstExpression {
	hash := make(map[ast.AstExpression]ast.AstExpression)
	p.AdvanceToken()
	for p.nextToken.Type != token.RBRACE {
		key := p.ParseExpression()
		p.AdvanceToken()
		p.AdvanceToken()
		value := p.ParseExpression()
		hash[key] = value
		if p.nextToken.Type == token.COMMA {
			p.AdvanceToken()
			p.AdvanceToken()
		}
	}
	p.AdvanceToken()
	return &ast.HashLiteral{
		Pairs: hash,
	}
}

// Statement Parsing
func (p *Parser) ParseVarStatement() *ast.VarStatement {
	p.AdvanceToken()
	var statement *ast.VarStatement

	if p.curToken.Type != token.IDENT {
		p.errors = append(p.errors, ParserError("Expected identifier"))
		return nil
	}

	statement = &ast.VarStatement{
		Identifer: *p.ParseIdentiferLiteral(),
	}

	p.CheckTokenAdvance(token.ASSIGN)
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

	if p.curToken.Type == token.IDENT && p.nextToken.Type == token.ASSIGN {
		exprStatement.Expression = p.parseReassignment()
	}

	switch p.nextToken.Type {
	case token.LPAREN:
		exprStatement.Expression = p.ParseExpression()
	case token.INCREMENT:
		exprStatement.Expression = p.ParseExpression()
	case token.DECREMENT:
		exprStatement.Expression = p.ParseExpression()
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

func (p *Parser) ParseForStatement() *ast.ForStatement {
	var ForStatement ast.ForStatement

	p.AdvanceToken()
	ForStatement.Initializer = p.ParseVarStatement()
	p.CheckTokenAdvance(token.SEMICOLON)
	p.CheckTokenAdvance(token.IDENT)
	ForStatement.Conditional = p.ParseExpression()
	p.CheckTokenAdvance(token.SEMICOLON)
	ForStatement.Increment = *p.Parse()
	p.CheckTokenAdvance(token.LBRACE)
	ForStatement.Body = p.ParseBlockStatement()
	p.CheckTokenAdvance(token.RBRACE)
	return &ForStatement
}

//for var x = 10; x < 11; x += 1 { println(x) }

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	var WhileStatement ast.WhileStatement
	p.AdvanceToken()
	WhileStatement.Condition = p.ParseExpression()
	p.CheckTokenAdvance(token.LBRACE)
	WhileStatement.Body = p.ParseBlockStatement()
	p.CheckTokenAdvance(token.RBRACE)
	return &WhileStatement
}

// while x < 10 {println(x)}

// Literal Parsing
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
		Value: p.curToken.Literal,
	}
}

func (p *Parser) ParseFunctionLiteral() *ast.FunctionLiteral {
	var function ast.FunctionLiteral

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

	return &function
}

func (p *Parser) ParseArrayLiteral() *ast.ArrayLiteral {
	var array ast.ArrayLiteral

	p.AdvanceToken()
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return &array
}

func (p *Parser) parseReassignment() *ast.Reassignment {
	Ident := p.ParseExpression()
	identifer := Ident.(*ast.IdentiferLiteral)
	p.CheckTokenAdvance(token.ASSIGN)
	p.AdvanceToken()
	Value := p.ParseExpression()
	return &ast.Reassignment{
		Ident: *identifer,
		Value: Value,
	}
}

// Helper Functions
func (p *Parser) CheckTokenAdvance(wanted token.TokenType) {
	if p.nextToken.Type != wanted {
		//some kind of error handle
		var newError string = "Expected " + string(wanted) + " but got " + p.nextToken.String()
		p.addError(newError)
	}
	p.AdvanceToken()
}

func (p *Parser) nextTokenIs(wanted token.TokenType) bool {
	return p.nextToken.Type == wanted
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.AstExpression {
	var expressions []ast.AstExpression

	for p.curToken.Type != end {
		if p.curToken.Type == token.COMMA {
			p.AdvanceToken()
		}
		expressions = append(expressions, p.ParseExpression())
		p.AdvanceToken()
	}

	return expressions
}

// func test(a,b){return a}
