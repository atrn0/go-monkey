package repl

import (
	"bufio"
	"fmt"
	"github.com/atrn0/go-monkey/evaluator"
	"github.com/atrn0/go-monkey/lexer"
	"github.com/atrn0/go-monkey/object"
	"github.com/atrn0/go-monkey/parser"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t")
		io.WriteString(out, msg)
		io.WriteString(out, "\n")
	}
}
