package lexer

import (
	"strings"

	"github.com/EVFUBS/AlphaLang/token"
)

type Lexer struct {
	input         string
	curPosition   int
	nextPostition int
	ch            byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.nextPostition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.nextPostition]
	}

	l.curPosition = l.nextPostition
	l.nextPostition += 1
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) isChar(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) isNum(ch byte) bool {
	if ch >= '0' && ch <= '9' {
		return true
	}
	return false
}

func (l *Lexer) peekChar() byte {
	if l.nextPostition < len(l.input) {
		return l.input[l.nextPostition]
	}
	return 0
}

func (l *Lexer) readIdent() string {
	position := l.curPosition
	for l.isChar(l.peekChar()) {
		l.readChar()
	}

	return l.input[position:l.nextPostition]
}

func (l *Lexer) readNum() string {
	position := l.curPosition
	for l.isNum(l.peekChar()) || l.peekChar() == '.' {
		l.readChar()
	}
	return l.input[position:l.nextPostition]
}

func (l *Lexer) readString() string {
	position := l.curPosition
	for l.peekChar() != '"' {
		l.readChar()
	}
	l.readChar()
	//return l.input[position:l.nextPostition]
	return l.input[position+1 : l.nextPostition-1]
}

func (l *Lexer) NextToken() *token.Token {
	l.skipWhitespace()
	var newToken *token.Token

	switch l.ch {
	case '{':
		newToken = &token.Token{Type: token.LBRACE, Literal: "{"}
	case '}':
		newToken = &token.Token{Type: token.RBRACE, Literal: "}"}
	case '(':
		newToken = &token.Token{Type: token.LPAREN, Literal: "("}
	case ')':
		newToken = &token.Token{Type: token.RPAREN, Literal: ")"}
	case '[':
		newToken = &token.Token{Type: token.LBRACKET, Literal: "["}
	case ']':
		newToken = &token.Token{Type: token.RBRACKET, Literal: "]"}
	case ';':
		newToken = &token.Token{Type: token.SEMICOLON, Literal: ";"}
	case ':':
		newToken = &token.Token{Type: token.COLON, Literal: ":"}
	case ',':
		newToken = &token.Token{Type: token.COMMA, Literal: ","}
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			newToken = &token.Token{Type: token.EQUAL, Literal: "=="}
		} else {
			newToken = &token.Token{Type: token.ASSIGN, Literal: "="}
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			newToken = &token.Token{Type: token.LEQUAL, Literal: "<="}
		} else {
			newToken = &token.Token{Type: token.LTHAN, Literal: "<"}
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			newToken = &token.Token{Type: token.GEQUAL, Literal: ">="}
		} else {
			newToken = &token.Token{Type: token.GTHAN, Literal: ">"}
		}
	case '+':
		if l.peekChar() == '=' {
			l.readChar()
			newToken = &token.Token{Type: token.INCREMENT, Literal: "+="}
		} else {
			newToken = &token.Token{Type: token.PLUS, Literal: "+"}
		}
	case '-':
		if l.peekChar() == '=' {
			l.readChar()
			newToken = &token.Token{Type: token.DECREMENT, Literal: "-="}
		} else {
			newToken = &token.Token{Type: token.MINUS, Literal: "-"}
		}
	case '*':
		newToken = &token.Token{Type: token.ASTERISK, Literal: "*"}
	case '/':
		newToken = &token.Token{Type: token.SLASH, Literal: "/"}
	case '%':
		newToken = &token.Token{Type: token.MODULUS, Literal: "%"}
	case 0:
		newToken = &token.Token{Type: token.EOF, Literal: "EOF"}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			newToken = &token.Token{Type: token.NOTEQUAL, Literal: "!="}
		} else {
			newToken = &token.Token{Type: token.BANG, Literal: "!"}
		}
	case '"':
		newString := l.readString()
		newToken = &token.Token{Type: token.STRING, Literal: newString}
	default:
		if l.isChar(l.ch) {
			ident := l.readIdent()
			tokType := token.KeywordLookUp(ident)
			newToken = &token.Token{Type: tokType, Literal: ident}
		} else if l.isNum(l.ch) {
			num := l.readNum()
			if strings.Contains(num, ".") {
				newToken = &token.Token{Type: token.FLOAT, Literal: num}
			} else {
				newToken = &token.Token{Type: token.INTEGER, Literal: num}
			}
		} else {
			newToken = &token.Token{Type: token.ILLEGAL, Literal: "ILLEGAL"}
		}
	}
	l.readChar()
	return newToken
}
