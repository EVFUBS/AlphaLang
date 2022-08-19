package token

const (
	// Special tokens
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifier
	IDENT = "IDENT"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	// Brackets
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNC"
	VAR      = "VAR"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELIF     = "ELIF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"

	// Types
	STRING = "STRING"
	NUMBER = "NUMBER"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	
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
}

func keywordLookUp(word string) TokenType {
	if tok, ok := keywords[word]; ok {
		return tok
	}

	return IDENT
}
