package ast

import (
	"fmt"
	"strconv"
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
func (bs *BlockStatement) String() string { return "bs" }

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

type IfStatement struct {
	Condition   AstExpression
	Consequence BlockStatement
	Alternative BlockStatement
}

func (is *IfStatement) StatementNode() {}
func (is *IfStatement) String() string { return "is" }

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

//expressions
type InfixExpression struct {
	Left     AstExpression
	Operator string
	Right    AstExpression
}

func (ie *InfixExpression) ExpressionNode() {}
func (ie *InfixExpression) String() string  { return "ie" }

type PrefixExpression struct {
	Prefix     string
	Expression AstExpression
}

func (pe *PrefixExpression) ExpressionNode() {}
func (pe *PrefixExpression) String() string  { return "pe" }

type IdentiferLiteral struct {
	Ident string
}

func (idl *IdentiferLiteral) ExpressionNode() {}
func (idl *IdentiferLiteral) String() string {
	var output string

	output += idl.Ident

	return output
}

type IntegerLiteral struct {
	Value int64
}

func (il *IntegerLiteral) ExpressionNode() {}
func (il *IntegerLiteral) String() string {
	var output string

	output += strconv.Itoa(int(il.Value))

	return output
}

type FloatLiteral struct {
	Value float64
}

func (fl *FloatLiteral) ExpressionNode() {}
func (fl *FloatLiteral) String() string {
	var output string

	output += fmt.Sprintf("%v", fl.Value)

	return output
}

type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) ExpressionNode() {}
func (sl *StringLiteral) String() string {
	var output string

	output += sl.Value

	return output
}

type BooleanLiteral struct {
	Value bool
}

func (bl *BooleanLiteral) ExpressionNode() {}
func (sl *BooleanLiteral) String() string {
	var output string

	output += strconv.FormatBool(sl.Value)

	return output
}
