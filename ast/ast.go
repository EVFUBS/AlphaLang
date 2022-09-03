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

func (p *Program) String() string {
	var out string
	for _, s := range p.Statements {
		out += s.String()
	}
	return out
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

type ForStatement struct {
	Initializer AstStatement
	Conditional AstExpression
	Increment   AstStatement
	Body        BlockStatement
}

func (fs *ForStatement) StatementNode() {}
func (fs *ForStatement) String() string {
	var output string

	output += "ForStatement("
	output += fs.Initializer.String()
	output += fs.Conditional.String()
	output += fs.Increment.String()
	output += ")"
	output += fs.Body.String()
	output += ")"
	return output
}

//expressions
type InfixExpression struct {
	Left     AstExpression
	Operator *token.Token
	Right    AstExpression
}

func (ie *InfixExpression) ExpressionNode() {}
func (ie *InfixExpression) String() string {
	var output string

	output += "InfixExpr("
	output += ie.Left.String()
	output += ie.Operator.Literal
	output += ie.Right.String()
	output += ")"

	return output
}

type PrefixExpression struct {
	Prefix     *token.Token
	Expression AstExpression
}

func (pe *PrefixExpression) ExpressionNode() {}
func (pe *PrefixExpression) String() string {
	var output string

	output += "Prefix("
	output += pe.Prefix.Literal
	output += pe.Expression.String()
	output += ")"
	return output
}

type CallExpression struct {
	Function  AstExpression
	Arguments []AstExpression
}

func (ce *CallExpression) ExpressionNode() {}
func (ce *CallExpression) String() string {
	var output string

	output += "Call("
	output += ce.Function.String()
	output += "Arguments("
	for _, arg := range ce.Arguments {
		output += arg.String()
	}
	output += ")"
	return output
}

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
	Name       AstExpression
	Parameters []IdentiferLiteral
	Body       BlockStatement
}

func (fl *FunctionLiteral) ExpressionNode() {}
func (fl *FunctionLiteral) String() string {
	var output string

	output += "Func( \n"
	output += fl.Name.String()
	output += "Params("
	for _, ident := range fl.Parameters {
		output += ident.String()
	}
	output += ") \n"
	output += fl.Body.String()
	output += ")"

	return output
}

type ArrayLiteral struct {
	Elements []AstExpression
}

func (al *ArrayLiteral) ExpressionNode() {}
func (al *ArrayLiteral) String() string {
	var output string

	output += "Array("
	for _, elem := range al.Elements {
		output += elem.String()
	}
	output += ")"

	return output
}

type IndexExpression struct {
	Left  AstExpression
	Index AstExpression
}

func (ie *IndexExpression) ExpressionNode() {}
func (ie *IndexExpression) String() string {
	var output string

	output += "Index("
	output += ie.Left.String()
	output += ie.Index.String()
	output += ")"

	return output
}

type HashLiteral struct {
	Pairs map[AstExpression]AstExpression
}

func (hl *HashLiteral) ExpressionNode() {}
func (hl *HashLiteral) String() string {
	var output string

	output += "Hash("
	for key, value := range hl.Pairs {
		output += key.String()
		output += value.String()
	}
	output += ")"

	return output
}
