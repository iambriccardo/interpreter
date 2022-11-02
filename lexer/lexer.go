package lexer

import (
	"fmt"
	"monkey/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input
	readPosition int  // current reading position in input
	ch           byte // current char under examination
}

// TODO: support unicode with "rune" type
// Gives us the next character and advances our position.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		// ASCII code for "NUL" character
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func newTokenWithStr(tokenType token.TokenType, str string) token.Token {
	return token.Token{Type: tokenType, Literal: str}
}

func (l *Lexer) readCharUntil(block func(byte) bool) string {
	position := l.position
	for block(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) peekAhead() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) peekAheadAndReturnIf(expectedChar string) (string, error) {
	if l.peekAhead() == []byte(expectedChar)[0] {
		previousChar := l.ch
		l.readChar()
		currentChar := l.ch
		return string([]byte{previousChar, currentChar}), nil
	} else {
		return "", fmt.Errorf("the next char doesn't match the expected char")
	}
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	// We continue reading until we find l.ch that is not a whitespace.
	l.skipWhitespace()

	switch l.ch {
	case '=':
		if str, err := l.peekAheadAndReturnIf("="); err == nil {
			tok = newTokenWithStr(token.EQ, str)
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if str, err := l.peekAheadAndReturnIf("="); err == nil {
			tok = newTokenWithStr(token.NOT_EQ, str)
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readCharUntil(isLetter)
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readCharUntil(isDigit)
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()

	return tok
}
