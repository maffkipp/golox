package parser

import "github.com/maffkipp/golox/lexer"

type Environment struct {
	values map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{values: make(map[string]any)}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(name lexer.Token) (any, error) {
	if val, ok := e.values[name.Lexeme]; ok {
		return val, nil
	}
	return nil, NewRuntimeError(name, "undefined variable '"+name.Lexeme+"'.")
}

func (e *Environment) Assign(name lexer.Token, value any) error {
	if _, ok := e.values[name.Lexeme]; ok {
		e.Define(name.Lexeme, value)
		return nil
	} else {
		return NewRuntimeError(name, "undefined variable '"+name.Lexeme+"'.")
	}
}
