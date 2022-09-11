package main

import (
	"os"
	"strings"

	"github.com/EVFUBS/AlphaLang/repl"
)

func main() {
	// take file with al extension as input and evaluate it
	if len(os.Args) > 1 {
		if strings.HasSuffix(os.Args[1], ".al") {
			file, err := os.Open(os.Args[1])
			if err != nil {
				panic(err)
			}
			defer file.Close()
			repl.RunFile(file)
		} else {
			panic("File must have .al extension")
		}
	} else {
		repl.Start(os.Stdin, os.Stdout)
	}
}
