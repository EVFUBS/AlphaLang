package lexer

type Lexer struct {
	input         string
	curPosition   int
	nextPostition int
	ch            byte
}

func (l *Lexer) readChar() {
	if l.curPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.curPosition]
	}

	l.curPosition = l.nextPostition
	l.nextPostition += 1
}
