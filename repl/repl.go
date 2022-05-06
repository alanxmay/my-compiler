package repl

import (
	"io"
	"my-compiler/evaluator"
	"my-compiler/lexer"
	"my-compiler/object"
	"my-compiler/parser"

	"github.com/chzyer/readline"
)

const PROMPT = ">>> "

func Start(in io.Reader, out io.Writer) {
	l, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[34m»\033[0m ",
		HistoryFile:     "/tmp/my_compiler_readline.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})

	if err != nil {
		panic(err)
	}
	defer l.Close()

	env := object.NewEnvironment()

	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				return
			} else {
				continue
			}
		} else if err == io.EOF {
			return
		}

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
