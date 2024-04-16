package parser

import (
	"fmt"

	"github.com/maffkipp/golox/errors"
	"github.com/maffkipp/golox/lexer"
)

type ParseError struct {
	Token   lexer.Token
	Message string
}

type RuntimeError struct {
	Token   lexer.Token
	Message string
}

func NewParseError(token lexer.Token, message string) *ParseError {
	return &ParseError{Token: token, Message: message}
}

func NewRuntimeError(token lexer.Token, message string) *ParseError {
	return &ParseError{Token: token, Message: message}
}

func (p *ParseError) Report() {
	if p.Token.TokenType == lexer.EOF {
		errors.Report(p.Token.Line, " at end", p.Message)
	} else {
		errors.Report(p.Token.Line, "at '"+p.Token.Lexeme+"'", p.Message)
	}
}

func (p ParseError) Error() string {
	return errors.LoxErrorFmt(p.Token.Line, "", p.Message)
}

func (r RuntimeError) Report() {
	fmt.Printf("%s/n[line %d]", r.Error(), r.Token.Line)
}

func (r RuntimeError) Error() string {
	return r.Message
}
