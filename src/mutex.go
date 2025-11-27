package src

import (
	"bufio"
	"fmt"
	"os"
)

type Mutex struct {
	hadError bool
}

var mutex = Mutex{}

func Main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: m <filepath>")
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
		mutex.run(line)
		mutex.hadError = false
	}
}

func (m *Mutex) run(sourceCode string) {
	scanner := NewScanner(sourceCode)
	tokens := scanner.ScanTokens()

	for _, token := range tokens {
		fmt.Println(token)
	}
}

func (m *Mutex) reportError(line int, message string) {
	mutex.report(line, "", message)
}

func (m *Mutex) report(line int, where, message string) {
	fmt.Printf("[ LINE %d ] Error -> %s: %s\n", line, where, message)
	mutex.hadError = true
}
