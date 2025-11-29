package errors

import (
	"fmt"
	"os"
)


func ReportError(line int, message string) {
	Report(line, "Report", message)
}

func ReportParser(message string, code int) {
	fmt.Fprintf(os.Stderr, "\t|\n")
	fmt.Fprintf(os.Stderr, "\t| Parser::Error -> %s\n", message)
	fmt.Fprintf(os.Stderr, "\t|\n")

	os.Exit(code)
}

func Report(line int, where, message string) {
	fmt.Fprintf(os.Stderr, "     |\n")
	fmt.Fprintf(os.Stderr, "%4d | %s::Error -> %s\n", line, where, message)
	fmt.Fprintf(os.Stderr, "     |\n")
}

func Exit(code int) {
	fmt.Printf("\n[Process exited with code: %d]\n", code)
	os.Exit(code)
}
