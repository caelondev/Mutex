package src

import (
	"bufio"
	"fmt"
	"os"

	"github.com/caelondev/mutex/src/errors"
	"github.com/caelondev/mutex/src/frontend/parser"
	"github.com/caelondev/mutex/src/runtime"
	// "github.com/sanity-io/litter"
)

type Mutex struct {
	hadError bool
}

var mutex = Mutex{}
var env = runtime.NewEnvironment(nil)

func Main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: mutex <filepath>")
		os.Exit(64)
	}

	if len(os.Args) == 2 {
		mutex.runFile(os.Args[1])
	} else {
		mutex.runRepl()
	}
}

func (m *Mutex) runFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	mutex.run(string(bytes))

	if m.hadError {
		os.Exit(65)
	}

	return nil
}

func (m *Mutex) runRepl() {
	input := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf(">> ")
		if ok := input.Scan(); !ok {
			break
		}

		line := input.Text()

		if line == "@exit" {
			m.Exit(0)
		}

		mutex.run(line)
		mutex.hadError = false
	}
}

func (m *Mutex) run(sourceCode string) {
	scanner := NewScanner(sourceCode)
	tokens := scanner.ScanTokens()
	ast := parser.ProduceAST(tokens)

	var result runtime.RuntimeValue
	for _, stmt := range ast.Body {
		result = runtime.EvaluateStatement(stmt, env)
	}

	fmt.Printf("%v\n\n", result)
	// litter.Dump(ast)
}

func (m *Mutex) ReportError(line int, message string) {
	mutex.Report(line, "Report", message)
}

func (m *Mutex) Report(line int, where, message string) {
	m.hadError = true

	errors.Report(line, where, message)
}

func (m *Mutex) Exit(code int) {
	fmt.Printf("\n[Process exited with code: %d]\n", code)
	os.Exit(code)
}

func GetMutex() *Mutex {
	return &mutex
}
