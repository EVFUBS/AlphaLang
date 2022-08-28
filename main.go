package main

import (
	"os"

	"github.com/EVFUBS/AlphaLang/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
