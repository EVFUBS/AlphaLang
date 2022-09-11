package evaluator

import (
	"fmt"

	"github.com/EVFUBS/AlphaLang/ast"
	"github.com/EVFUBS/AlphaLang/builtins"
	"github.com/EVFUBS/AlphaLang/objects"
	"github.com/EVFUBS/AlphaLang/token"
)

type Evaluator struct {
	Env *objects.Environment
}

func New() *Evaluator {
	return &Evaluator{Env: objects.NewEnvironment()}
}

func (e *Evaluator) Eval(node ast.AstNode, env *objects.Environment) objects.Object {
	//println(node.String())
	switch node := node.(type) {
	case *ast.Program:
		return e.evalProgram(node)
	case *ast.VarStatement:
		return e.evalVarStatement(node, env)
	case *ast.IdentiferLiteral:
		return e.evalIdentiferLiteral(node, env)
	case *ast.IntegerLiteral:
		return e.evalIntegerLiteral(node)
	case *ast.FloatLiteral:
		return e.evalFloatLiteral(node)
	case *ast.StringLiteral:
		return e.evalStringLiteral(node)
	case *ast.BooleanLiteral:
		return e.evalBooleanLiteral(node)
	case *ast.ArrayLiteral:
		return e.evalArrayLiteral(node, env)
	case *ast.HashLiteral:
		return e.evalHashLiteral(node, env)
	case *ast.IndexExpression:
		return e.evalIndexExpression(node, env)
	case *ast.FunctionLiteral:
		return e.evalFunctionLiteral(node, env)
	case *ast.CallExpression:
		return e.evalCallExpression(node, env)
	case *ast.InfixExpression:
		return e.evalInfixExpression(node, env)
	case *ast.PrefixExpression:
		return e.evalPrefixExpression(node, env)
	case *ast.ExpressionStatement:
		return e.evalExprStatement(node, env)
	case *ast.ReturnStatement:
		return e.evalReturnStatement(node, env)
	case *ast.BlockStatement:
		return e.evalBlockStatement(node, env)
	case *ast.IfStatement:
		return e.evalIfStatement(node, env)
	case *ast.ForStatement:
		return e.evalForStatement(node, env)
	}
	return nil
}

func (e *Evaluator) evalProgram(node *ast.Program) objects.Object {
	var result objects.Object
	for _, statement := range node.Statements {
		result = e.Eval(statement, e.Env)
		if result != nil {
			return result
		}
	}
	return result
}

func (e *Evaluator) evalBlockStatement(node *ast.BlockStatement, env *objects.Environment) objects.Object {
	for _, statement := range node.Statements {
		result := e.Eval(statement, env)
		if result != nil {
			return result
		}
	}
	return nil
}

func (e *Evaluator) evalVarStatement(node *ast.VarStatement, env *objects.Environment) objects.Object {
	value := e.Eval(node.Value, env)
	ident := node.Identifer.Ident

	env.Set(ident, value)
	return nil
}

func (e *Evaluator) evalExprStatement(node *ast.ExpressionStatement, env *objects.Environment) objects.Object {
	switch node := node.Expression.(type) {
	case *ast.Reassignment:
		return e.evalReassignement(node, env)
	case *ast.FunctionLiteral:
		return e.evalFunctionLiteral(node, env)
	case *ast.CallExpression:
		return e.evalCallExpression(node, env)
	case *ast.InfixExpression:
		return e.evalInfixExpression(node, env)
	}
	return nil
}

func (e *Evaluator) evalReturnStatement(node *ast.ReturnStatement, env *objects.Environment) objects.Object {
	value := e.Eval(node.ReturnValue, env)
	return &objects.ReturnValue{Value: value}
}

func (e *Evaluator) evalIfStatement(node *ast.IfStatement, env *objects.Environment) objects.Object {
	if e.Eval(node.If.Condition, env).(*objects.Boolean).Value {
		return e.Eval(&node.If.Consequence, env)
	} else if node.Elif != nil {
		for _, elif := range node.Elif {
			if e.Eval(elif.Condition, env).(*objects.Boolean).Value {
				return e.Eval(&elif.Consequence, env)
			}
		}
	} else {
		return e.Eval(&node.Else, env)
	}
	return nil
}

func (e *Evaluator) evalForStatement(node *ast.ForStatement, env *objects.Environment) objects.Object {
	e.Eval(node.Initializer, env)
	for {
		condition := e.Eval(node.Conditional, env)
		if condition.(*objects.Boolean).Value {
			e.Eval(node.Body, env)
		} else {
			break
		}
		e.Eval(node.Increment, env)
	}
	return nil
}

// eval expression statements
func (e *Evaluator) evalReassignement(node *ast.Reassignment, env *objects.Environment) objects.Object {
	if _, ok := env.Get(node.Ident.Ident); ok {
		env.Set(node.Ident.Ident, e.Eval(node.Value, env))
	}
	return nil
}

func (e *Evaluator) evalFunctionLiteral(node *ast.FunctionLiteral, env *objects.Environment) objects.Object {
	fn := &objects.Function{
		Name:       node.Name,
		Parameters: node.Parameters,
		Body:       node.Body,
		Env:        *env,
	}
	env.Set(node.Name, fn)
	return nil
}

func (e *Evaluator) evalCallExpression(node *ast.CallExpression, env *objects.Environment) objects.Object {
	function := e.Eval(node.Function, env)

	var args []objects.Object

	for _, arg := range node.Arguments {
		args = append(args, e.Eval(arg, env))
	}

	return e.applyFunction(function, args, env)
}

func (e *Evaluator) applyFunction(fn objects.Object, args []objects.Object, env *objects.Environment) objects.Object {
	switch fn := fn.(type) {

	case *objects.Function:
		extendedEnv := objects.NewEnclosedEnvironment(env)
		for i, param := range fn.Parameters {
			extendedEnv.Set(param.Ident, args[i])
		}
		evaluated := e.Eval(&fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *objects.Builtin:
		return fn.Fn(args...)
	}

	return nil
}

func unwrapReturnValue(obj objects.Object) objects.Object {
	if returnValue, ok := obj.(*objects.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

// eval expressions
func (e *Evaluator) evalInfixExpression(node *ast.InfixExpression, env *objects.Environment) objects.Object {

	leftVal := e.Eval(node.Left, env)
	rightVal := e.Eval(node.Right, env)

	//could take in operator instead of whole node
	if leftVal.Type() == objects.INT && rightVal.Type() == objects.INT || leftVal.Type() == objects.FLOAT && rightVal.Type() == objects.FLOAT {
		return e.evalNumericInfixExpression(node, leftVal, rightVal, env)
	} else if leftVal.Type() == objects.STRING && rightVal.Type() == objects.STRING {
		return e.evalStringInfixExpression(node, leftVal, rightVal)
	} else {
		return NewError("type mismatch: %s %s %s", leftVal.Type(), node.Operator, rightVal.Type())
	}
}

func (e *Evaluator) evalNumericInfixExpression(node *ast.InfixExpression, left, right objects.Object, env *objects.Environment) objects.Object {

	var leftVal float64
	var rightVal float64

	if left.Type() == objects.INT {
		leftVal = float64(left.(*objects.Integer).Value)
		rightVal = float64(right.(*objects.Integer).Value)
	} else {
		leftVal = left.(*objects.Float).Value
		rightVal = right.(*objects.Float).Value
	}

	switch node.Operator.Literal {
	case "+":
		if left.Type() == objects.INT {
			return &objects.Integer{Value: int64(leftVal + rightVal)}
		}
		return &objects.Float{Value: leftVal + rightVal}
	case "-":
		if left.Type() == objects.INT {
			return &objects.Integer{Value: int64(leftVal - rightVal)}
		}
		return &objects.Float{Value: leftVal - rightVal}
	case "*":
		if left.Type() == objects.INT {
			return &objects.Integer{Value: int64(leftVal * rightVal)}
		}
		return &objects.Float{Value: leftVal * rightVal}
	case "/":
		if left.Type() == objects.INT {
			return &objects.Integer{Value: int64(leftVal / rightVal)}
		}
		return &objects.Float{Value: leftVal / rightVal}
	case "%":
		if left.Type() == objects.INT {
			return &objects.Integer{Value: int64(leftVal) % int64(rightVal)}
		}
		return &objects.Float{Value: leftVal / rightVal}
	case "<":
		return &objects.Boolean{Value: leftVal < rightVal}
	case ">":
		return &objects.Boolean{Value: leftVal > rightVal}
	case "==":
		return &objects.Boolean{Value: leftVal == rightVal}
	case "!=":
		return &objects.Boolean{Value: leftVal != rightVal}
	case "<=":
		return &objects.Boolean{Value: leftVal <= rightVal}
	case ">=":
		return &objects.Boolean{Value: leftVal >= rightVal}
	case "+=":
		if left.Type() == objects.INT {
			env.Set(node.Left.(*ast.IdentiferLiteral).Ident, &objects.Integer{Value: int64(leftVal) + int64(rightVal)})
		} else {
			return NewError("Increment has to be an integer")
		}
	case "-=":
		if left.Type() == objects.INT {
			env.Set(node.Left.(*ast.IdentiferLiteral).Ident, &objects.Integer{Value: int64(leftVal) - int64(rightVal)})
		} else {
			return NewError("decrement has to be an integer")
		}
	}
	//error here
	return nil
}

func (e *Evaluator) evalStringInfixExpression(node *ast.InfixExpression, left, right objects.Object) objects.Object {
	leftVal := left.(*objects.String).Value
	rightVal := right.(*objects.String).Value

	switch node.Operator.Literal {
	case "+":
		return &objects.String{Value: leftVal + rightVal}
	case "==":
		return &objects.Boolean{Value: leftVal == rightVal}
	case "!=":
		return &objects.Boolean{Value: leftVal != rightVal}
	}
	//return error
	return nil
}

func (e *Evaluator) evalPrefixExpression(node *ast.PrefixExpression, env *objects.Environment) objects.Object {
	expr := e.Eval(node.Expression, env)
	//minus prefix and bang prefix
	if node.Prefix.Type == token.MINUS {

		if expr, ok := expr.(*objects.Integer); ok {
			//convert postive to negative
			return &objects.Integer{Value: -expr.Value}
		}

		if expr, ok := expr.(*objects.Float); ok {
			return &objects.Float{Value: -expr.Value}
		}

		return NewError("Expected a float or an integer")

	} else if node.Prefix.Type == token.BANG {

		if expr, ok := expr.(*objects.Boolean); ok {
			return &objects.Boolean{Value: !expr.Value}
		}
	}
	//error
	return nil
}

func (e *Evaluator) evalIntegerLiteral(node *ast.IntegerLiteral) *objects.Integer {
	return &objects.Integer{
		Value: node.Value,
	}
}

func (e *Evaluator) evalFloatLiteral(node *ast.FloatLiteral) *objects.Float {
	return &objects.Float{
		Value: node.Value,
	}
}

func (e *Evaluator) evalStringLiteral(node *ast.StringLiteral) *objects.String {
	return &objects.String{
		Value: node.Value,
	}
}

func (e *Evaluator) evalBooleanLiteral(node *ast.BooleanLiteral) *objects.Boolean {
	return &objects.Boolean{
		Value: node.Value,
	}
}

func (e *Evaluator) evalIdentiferLiteral(node *ast.IdentiferLiteral, env *objects.Environment) objects.Object {
	if val, ok := env.Get(node.Ident); ok {
		return val
	} else if val, ok := builtins.BuiltIns[node.Ident]; ok {
		return &val
	}
	//error
	return nil
}

func (e *Evaluator) evalArrayLiteral(node *ast.ArrayLiteral, env *objects.Environment) objects.Object {
	var elements []objects.Object
	for _, elem := range node.Elements {
		elements = append(elements, e.Eval(elem, env))
	}
	return &objects.Array{Elements: elements}
}

func (e *Evaluator) evalHashLiteral(node *ast.HashLiteral, env *objects.Environment) objects.Object {
	pairs := make(map[objects.HashKey]objects.HashPair)
	for keyNode, valueNode := range node.Pairs {
		key := e.Eval(keyNode, env)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(objects.Hashable)
		if !ok {
			return &objects.Error{Message: fmt.Sprintf("unusable as hash key: %s", key.Type())}
		}
		value := e.Eval(valueNode, env)
		if isError(value) {
			return value
		}
		hashed := hashKey.HashKey()
		pairs[hashed] = objects.HashPair{Key: key, Value: value}
	}
	return &objects.Hash{Pairs: pairs}
}

func (e *Evaluator) evalIndexExpression(node *ast.IndexExpression, env *objects.Environment) objects.Object {
	left := e.Eval(node.Left, env)
	index := e.Eval(node.Index, env)

	if isError(left) {
		return left
	}
	if isError(index) {
		return index
	}

	switch {
	case left.Type() == objects.ARRAY && index.Type() == objects.INT:
		return e.evalArrayIndexExpression(left, index)
	case left.Type() == objects.HASH:
		return e.evalHashIndexExpression(left, index)
	}
	//error
	return nil
}

func (e *Evaluator) evalArrayIndexExpression(array, index objects.Object) objects.Object {
	arr := array.(*objects.Array)
	idx := index.(*objects.Integer).Value
	max := int64(len(arr.Elements) - 1)

	if idx < 0 || idx > max {
		return &objects.Null{}
	}
	return arr.Elements[idx]
}

func (e *Evaluator) evalHashIndexExpression(hash, index objects.Object) objects.Object {
	hashObj := hash.(*objects.Hash)
	key, ok := index.(objects.Hashable)
	if !ok {
		//error
		return nil
	}
	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return &objects.Null{}
	}
	return pair.Value
}

func NewError(format string, a ...interface{}) *objects.Error {
	return &objects.Error{
		Message: fmt.Sprintf(format, a...),
	}
}

func isError(obj objects.Object) bool {
	if obj != nil {
		return obj.Type() == objects.ERROR
	}
	return false
}
