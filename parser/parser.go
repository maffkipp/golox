package parser

import (
	"fmt"

	"github.com/maffkipp/golox/lexer"
)

type Parser struct {
	tokens  []lexer.Token
	current int
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{tokens, 0}
}

func (p *Parser) Parse() (statements []Stmt, hadErrors bool) {

	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	return statements, false
}

func (p *Parser) declaration() Stmt {
	// Error boundary should be at each statement
	defer func() {
		if err := recover(); err != nil {
			if parseErr, ok := err.(ParseError); ok {
				parseErr.Report()
				p.synchronize()
				return
			} else {
				panic(err)
			}
		}
	}()

	if p.match(lexer.VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) statement() Stmt {
	if p.match(lexer.PRINT) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) varDeclaration() Stmt {
	var initializer Expr
	name := p.consume(lexer.IDENTIFIER, "Expect variable name.")

	if p.match(lexer.EQUAL) {
		initializer = p.expression()
	}

	p.consume(lexer.SEMICOLON, "Expect ';' after variable declaration.")
	return NewVarStmt(name, initializer)
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(lexer.SEMICOLON, "Expect ';' after value.")
	return NewPrintStmt(value)
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(lexer.SEMICOLON, "Expect ';' after expression.")
	return NewExpressionStmt(expr)
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.equality()

	if p.match(lexer.EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if _, ok := expr.(VariableExpr); ok {
			name := expr.(VariableExpr).Name
			return NewAssignExpr(name, value)
		}
		// Don't need to panic here
		err := NewParseError(equals, "Invalid assignment target.")
		err.Report()
	}

	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(lexer.BANG_EQUAL, lexer.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(lexer.GREATER, lexer.GREATER_EQUAL, lexer.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(lexer.MINUS, lexer.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr
}

func (p *Parser) factor() Expr {

	expr := p.unary()

	for p.match(lexer.SLASH, lexer.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr
}

func (p *Parser) unary() Expr {

	if p.match(lexer.BANG, lexer.MINUS) {
		operator := p.previous()
		right := p.unary()
		return NewUnaryExpr(operator, right)
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(lexer.FALSE) {
		return NewLiteralExpr(false)
	} else if p.match(lexer.TRUE) {
		return NewLiteralExpr(true)
	} else if p.match(lexer.NIL) {
		return NewLiteralExpr(nil)
	}

	if p.match(lexer.NUMBER, lexer.STRING) {
		return NewLiteralExpr(p.previous().Literal)
	}

	if p.match(lexer.IDENTIFIER) {
		return NewVariableExpr(p.previous())
	}

	if p.match(lexer.LEFT_PAREN) {
		expr := p.expression()
		p.consume(lexer.RIGHT_PAREN, "Expect ')' after expression.")
		return NewGroupingExpr(expr)
	}

	err := NewParseError(p.peek(), "Expect expression.")

	if err != nil {
		fmt.Printf("%+v", p.peek())
	}
	panic(err)
}

func (p *Parser) consume(tokenType lexer.TokenType, message string) lexer.Token {
	if p.check(tokenType) {
		return p.advance()
	}

	err := NewParseError(p.peek(), message)
	panic(err)
}

func (p *Parser) match(tokenTypes ...lexer.TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType lexer.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == tokenType
}

func (p *Parser) advance() lexer.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == lexer.EOF
}

func (p *Parser) peek() lexer.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() lexer.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == lexer.SEMICOLON {
			return
		}

		switch p.peek().TokenType {
		case lexer.CLASS:
		case lexer.FUN:
		case lexer.VAR:
		case lexer.FOR:
		case lexer.IF:
		case lexer.WHILE:
		case lexer.PRINT:
		case lexer.RETURN:
			return
		}
		p.advance()
	}
}
