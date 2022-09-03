package objects

import (
	"hash/fnv"
	"strconv"
	"strings"

	"github.com/EVFUBS/AlphaLang/ast"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type ObjectType string

const (
	INT          = "INTEGER"
	STRING       = "STRING"
	FLOAT        = "FLOAT"
	BOOLEAN      = "BOOLEAN"
	NULL         = "NULL"
	RETURN_VALUE = "RETURN_VALUE"
	ERROR        = "ERROR"
	FUNCTION     = "FUNCTION"
	BUILTIN      = "BUILTIN"
	ARRAY        = "ARRAY"
	HASH         = "HASH"
)

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INT }
func (i *Integer) Inspect() string  { return strconv.Itoa(int(i.Value)) }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type Float struct {
	Value float64
}

func (f *Float) Type() ObjectType { return FLOAT }
func (f *Float) Inspect() string  { return strconv.FormatFloat(f.Value, 'f', -1, 64) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN }
func (b *Boolean) Inspect() string  { return strconv.FormatBool(b.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type Function struct {
	Parameters []*ast.IdentiferLiteral
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION }
func (f *Function) Inspect() string  { return f.String() }

func (f *Function) String() string {
	var out string

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out += "fn"
	out += "("
	out += strings.Join(params, ", ")
	out += ")"
	out += f.Body.String()

	return out
}

type BuiltinFunction func(args ...Object) Object

func (bf *BuiltinFunction) Type() ObjectType { return BUILTIN }
func (bf *BuiltinFunction) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY }
func (ao *Array) Inspect() string  { return ao.String() }

func (ao *Array) String() string {
	var out string

	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out += "["
	out += strings.Join(elements, ", ")
	out += "]"

	return out
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Hashable interface {
	HashKey() HashKey
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH }
func (h *Hash) Inspect() string  { return h.String() }

func (h *Hash) String() string {
	var out string

	pairs := []string{}
	for _, p := range h.Pairs {
		pairs = append(pairs, p.Key.Inspect()+":"+p.Value.Inspect())
	}

	out += "{"
	out += strings.Join(pairs, ", ")
	out += "}"

	return out
}
