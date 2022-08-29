package ast

import (
	"fmt"
	"strconv"

	"github.com/EVFUBS/AlphaLang/token"
)

//interfaces
type AstNode interface {
	String() string
}

type AstExpression interface {
	AstNode
	ExpressionNode()
}

type AstStatement interface {
	AstNode
	StatementNode()
}

//program
type Program struct {
	Statements []AstStatement
}

//statements
type BlockStatement struct {
	Statements []AstStatement
}

func (bs *BlockStatement) StatementNode() {}
func (bs *BlockStatement) String() string {
	var output string

	output += "BlockStatement(\n"
	println(len(bs.Statements))
	for _, statement := range bs.Statements {
		output += statement.String()
		output += "\n"
	}
	output += ")\n"

	return output
}

type VarStatement struct {
	Identifer AstExpression
	Value     AstExpression
}

func (vs *VarStatement) StatementNode() {}
func (vs *VarStatement) String() string {
	var output string

	output += "VarStatement("
	output += "var "
	output += vs.Identifer.String()
	output += " = "
	output += vs.Value.String()
	output += ")"

	return output
}

type Conditional struct {
	Condition   AstExpression
	Consequence BlockStatement
}

func (con *Conditional) StatementNode() {}
func (con *Conditional) String() string {
	var output string

	output += "Condition(\n"
	output += con.Condition.String()
	output += ")"
	output += "Consequence(\n"
	output += con.Consequence.String()
	output += ")"

	return output
}

type IfStatement struct {
	If   Conditional
	Elif []Conditional
	Else BlockStatement
}

func (is *IfStatement) StatementNode() {}
func (is *IfStatement) String() string {
	var output string

	output += "If(\n"
	output += is.If.String()
	output += ")"

	if is.Elif != nil {
		for _, elif := range is.Elif {
			output += "ELif(\n"
			output += elif.String()
			output += ")"
		}
	}

	if is.Else.Statements != nil {
		output += "Else(\n"
		output += is.Else.String()
		output += ")"
	}

	return output
}

type ReturnStatement struct {
	ReturnValue AstExpression
}

func (rs *ReturnStatement) StatementNode() {}
func (rs *ReturnStatement) String() string {
	var output string

	output += "ReturnStatement("
	output += rs.ReturnValue.String()
	output += ")"

	return output
}

type ExpressionStatement struct {
	Token      token.Token
	Expression AstExpression
}

func (es *ExpressionStatement) StatementNode() {}
func (es *ExpressionStatement) String() string {
	var output string

	output += "ExprStatement(\n"
	output += es.Expression.String()
	output += ")\n"

	return output
}

//expressions
type InfixExpression struct {
	Left     AstExpression
	Operator string
	Right    AstExpression
}

func (ie *InfixExpression) ExpressionNode() {}
func (ie *InfixExpression) String() string {
	var output string

	output += "InfixExpr("
	output += ie.Left.String()
	output += ie.Operator
	output += ie.Right.String()
	output += ")"

	return output
}

type PrefixExpression struct {
	Prefix     string
	Expression AstExpression
}

func (pe *PrefixExpression) ExpressionNode() {}
func (pe *PrefixExpression) String() string  { return "pe" }

type IdentiferLiteral struct {
	Token token.Token
	Ident string
}

func (idl *IdentiferLiteral) ExpressionNode() {}
func (idl *IdentiferLiteral) String() string {
	var output string

	output += "Ident("
	output += idl.Ident
	output += ")"

	return output
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) ExpressionNode() {}
func (il *IntegerLiteral) String() string {
	var output string

	output += "Integer("
	output += strconv.Itoa(int(il.Value))
	output += ")"

	return output
}

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) ExpressionNode() {}
func (fl *FloatLiteral) String() string {
	var output string

	output += "Float("
	output += fmt.Sprintf("%v", fl.Value)
	output += ")"

	return output
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) ExpressionNode() {}
func (sl *StringLiteral) String() string {
	var output string

	output += "String("
	output += sl.Value
	output += ")"

	return output
}

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) ExpressionNode() {}
func (bl *BooleanLiteral) String() string {
	var output string

	output += "Bool("
	output += strconv.FormatBool(bl.Value)
	output += ")"

	return output
}

type FunctionLiteral struct {
	Name       string
	Parameters []IdentiferLiteral
	Body       BlockStatement
}

func (fl *FunctionLiteral) ExpressionNode() {}
func (fl *FunctionLiteral) String() string {
	var output string

	output += "Func( \n"
	output += "Params("
	for _, ident := range fl.Parameters {
		output += ident.String()
	}
	output += ") \n"
	output += fl.Body.String()
	output += ")"

	return output
}
