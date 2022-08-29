package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/EVFUBS/AlphaLang/lexer"
	"github.com/EVFUBS/AlphaLang/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)

		/* tok := l.NextToken()
		for tok.Type != token.EOF {
			println(tok.String())
			tok = l.NextToken()
		} */

		p := parser.New(l)
		ast := p.ParseProgram()
		for _, statement := range ast.Statements {
			fmt.Println(statement.String())
		}
	}
}
