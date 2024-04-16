package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/maffkipp/golox/errors"
	"github.com/maffkipp/golox/lexer"
	"github.com/maffkipp/golox/parser"
)

func main() {

	if len(os.Args) > 2 {
		fmt.Println("Usage: golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		err := RunFile(os.Args[1])
		if err != nil {
			fmt.Println("Unable to run file")
		}
	} else {
		err := RunPrompt()
		if err != nil {
			fmt.Println("Unable to read from prompt")
		}
	}
}

func RunFile(path string) error {
	if bytes, err := os.ReadFile(path); err != nil {
		return err
	} else if err := run(string(bytes)); err != nil {
		os.Exit(65)
	}
	return nil
}

func RunPrompt() error {

	reader := bufio.NewReader(os.Stdin)

	lineNumber := 0
	for {
		fmt.Print("> ")
		lineNumber++

		if line, err := reader.ReadString('\n'); err != nil {
			return err
		} else {
			// user can type "exit" or submit an empty line to close repl
			if len(line) == 1 || line == "exit\n" {
				break
			} else if err := run(line); err != nil {
				errors.Error(lineNumber, err.Error())
			}
			// newline after each output
			fmt.Println("")
		}
	}

	return nil
}

func run(source string) error {

	s := lexer.NewScanner(source)
	tokens, hadErrors := s.ScanTokens()

	if hadErrors {
		return fmt.Errorf("encountered errors while scanning")
	}

	p := parser.NewParser(tokens)
	statements, hadErrors := p.Parse()

	if hadErrors {
		return fmt.Errorf("encountered errors while parsing")
	}

	i := parser.NewInterpreter()

	if hadErrors := i.Interpret(statements); hadErrors {
		return fmt.Errorf("encountered runtime errors")
	}

	return nil
}
