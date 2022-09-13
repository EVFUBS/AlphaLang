package token

const (
	// Special tokens
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifier
	IDENT = "IDENT"

	GTHAN = ">"
	LTHAN = "<"

	GEQUAL = "GEQUAL"
	LEQUAL = "LEQUAL"

	INCREMENT = "INCREMENT"
	DECREMENT = "DECREMENT"

	// Operators
	ASSIGN   = "="
	EQUAL    = "EQUAL"
	NOTEQUAL = "NOTEQUAL"
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	MODULUS  = "%"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	// Brackets
	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNC"
	VAR      = "VAR"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELIF     = "ELIF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	FOR      = "FOR"
	WHILE    = "WHILE"

	// Types
	STRING  = "STRING"
	INTEGER = "INTEGER"
	FLOAT   = "FLOAT"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func (t *Token) String() string {
	token := "Type: " + string(t.Type) + " Literal: " + t.Literal
	return token
}

var keywords = map[string]TokenType{
	"if":     IF,
	"elif":   ELIF,
	"else":   ELSE,
	"func":   FUNCTION,
	"return": RETURN,
	"var":    VAR,
	"true":   TRUE,
	"false":  FALSE,
	"for":    FOR,
	"while":  WHILE,
}

func KeywordLookUp(word string) TokenType {
	if tok, ok := keywords[word]; ok {
		return tok
	}

	return IDENT
}
