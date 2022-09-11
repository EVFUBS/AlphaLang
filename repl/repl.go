package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/EVFUBS/AlphaLang/evaluator"
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

		if len(p.Errors()) > 0 {
			for _, err := range p.Errors() {
				fmt.Println(err)
			}
		} else {
			for _, statement := range ast.Statements {
				fmt.Println(statement.String())
			}
		}

		/* e := evaluator.New()
		evaluated := e.Eval(ast, e.Env)
		if evaluated != nil {
			fmt.Println(evaluated.Inspect())
		} */
	}
}

func RunFile(file io.Reader) {
	scanner := bufio.NewScanner(file)
	var code string
	for scanner.Scan() {
		code += scanner.Text()
	}
	l := lexer.New(code)
	p := parser.New(l)
	ast := p.ParseProgram()
	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			fmt.Println(err)
		}
	} else {
		for _, statement := range ast.Statements {
			fmt.Println(statement.String())
		}
	}
	e := evaluator.New()
	evaluated := e.Eval(ast, e.Env)
	if evaluated != nil {
		fmt.Println(evaluated.Inspect())
	}
}
