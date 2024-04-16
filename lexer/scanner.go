package lexer

import (
	"fmt"
	"strconv"

	"github.com/maffkipp/golox/errors"
)

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{source, []Token{}, 0, 0, 1}
}

func (s *Scanner) ScanTokens() (tokens []Token, hadErrors bool) {

	hadErrors = false

	for !s.isAtEnd() {
		// We are at the beginning of the next lexeme.
		s.start = s.current

		if err := s.scanToken(); err != nil {
			errors.Error(s.line, err.Error())
			hadErrors = true
		}
	}

	s.tokens = append(s.tokens, *NewToken(EOF, "", nil, s.line))
	return s.tokens, false
}

func (s *Scanner) scanToken() error {
	char := s.advance()
	switch char {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '!':
		s.addTokenOnCondition(s.match('='), BANG_EQUAL, BANG)
	case '=':
		s.addTokenOnCondition(s.match('='), EQUAL_EQUAL, EQUAL)
	case '<':
		s.addTokenOnCondition(s.match('='), LESS_EQUAL, LESS)
	case '>':
		s.addTokenOnCondition(s.match('='), GREATER_EQUAL, GREATER)
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
	case '\n':
		s.line++
	// Do nothing for whitespace
	case ' ':
	case '\r':
	case '\t':
	case '"':
		if err := s.string(); err != nil {
			return err
		}
	default:
		if isDigit(char) {
			if err := s.number(); err != nil {
				return err
			}
		} else if isAlpha(char) {
			if err := s.identifier(); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unexpected character %s", string(char))
		}
	}
	return nil
}

func (s *Scanner) advance() byte {
	char := s.source[s.current]
	s.current++
	return char
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func (s *Scanner) addToken(t TokenType) {
	s.addTokenWithLiteral(t, nil)
}

func (s *Scanner) addTokenWithLiteral(t TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, *NewToken(t, text, literal, s.line))
}

func (s *Scanner) addTokenOnCondition(condition bool, ifTrue TokenType, ifFalse TokenType) {
	if condition {
		s.addToken(ifTrue)
	} else {
		s.addToken(ifFalse)
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() || s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) string() error {
	for s.peek() != '"' && !s.isAtEnd() {

		// Lox supports multiline strings
		if s.peek() == '\n' {
			s.line++
		}

		s.advance()
	}

	if s.isAtEnd() {
		return fmt.Errorf("unterminated string")
	}

	// Account for closing quote
	s.advance()

	// Trim quotes
	text := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(STRING, text)
	return nil
}

func (s *Scanner) number() error {
	for isDigit(s.peek()) {
		s.advance()
	}

	// Look for a decimal
	if s.peek() == '.' && isDigit(s.peekNext()) {
		// discard decimal
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	num, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		return fmt.Errorf("unable to parse number")
	}

	s.addTokenWithLiteral(NUMBER, num)
	return nil
}

func (s *Scanner) identifier() error {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]

	if tokenType, ok := keywords[text]; ok {
		s.addToken(tokenType)
	} else {
		s.addToken(IDENTIFIER)
	}

	return nil
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func isAlpha(char byte) bool {
	isLowercaseAlpha := char >= 'a' && char <= 'z'
	isUppercaseAlpha := char >= 'A' && char <= 'Z'
	isUnderscore := char == '_'

	return isLowercaseAlpha || isUppercaseAlpha || isUnderscore
}

func isAlphaNumeric(char byte) bool {
	return isAlpha(char) || isDigit(char)
}
