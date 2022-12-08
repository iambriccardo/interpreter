package lexer

import (
	"fmt"
	"monkey/token"
    "strings"
)

type Lexer struct {
	input        []string
	inputLine    int  // current line
	position     int  // current position in input
	peekPosition int  // current peeking position in input
	ch           byte // current char under examination
}

// TODO: support unicode with "rune" type
// Gives us the next character and advances our position.
func (l *Lexer) readChar() {
	if l.peekPosition >= len(l.input[l.inputLine]) {
        // If we are at the end of the last line, we will return the NUL character, otherwise
        // we should silently go to the next line.
        if l.inputLine == len(l.input) - 1 {
            // ASCII code for "NUL" character.
            l.ch = 0
        } else {
            // ASCII code for "SOH" character.
            l.ch = 1
        }
	} else {
        l.ch = l.input[l.inputLine][l.peekPosition]
	}

	l.position = l.peekPosition
	l.peekPosition += 1
}

func (l *Lexer) readCharUntil(block func(byte) bool) string {
	position := l.position
	for block(l.ch) {
		l.readChar()
	}
    return l.input[l.inputLine][position:l.position]
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
	if l.peekPosition >= len(l.input) {
		return 0
	} else {
        return l.input[l.inputLine][l.peekPosition]
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

func (l *Lexer) nextLine() {
    if l.inputLine < len(l.input) - 1 {
        l.inputLine++
        // We reset the position.
        l.peekPosition = 0
        l.position = 0
    }
}

func getLinesFromInput(input string) []string {
    return strings.Split(input, "\n")
}

func New(input string) *Lexer {
	l := &Lexer{input: getLinesFromInput(input)}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token = token.Token{}

	// We continue reading until we find l.ch that is not a whitespace.
	l.skipWhitespace()

    // We add column and line information.
    tok.Line = l.inputLine
    tok.Column = l.position

	switch l.ch {
	case '=':
		if str, err := l.peekAheadAndReturnIf("="); err == nil {
			tok.AddStr(token.EQ, str)
		} else {
			tok.AddChar(token.ASSIGN, l.ch)
		}
	case '+':
        tok.AddChar(token.PLUS, l.ch)
	case '-':
        tok.AddChar(token.MINUS, l.ch)
	case '!':
		if str, err := l.peekAheadAndReturnIf("="); err == nil {
            tok.AddStr(token.NOT_EQ, str)
		} else {
            tok.AddChar(token.BANG, l.ch)
		}
	case '/':
        tok.AddChar(token.SLASH, l.ch)
	case '*':
        tok.AddChar(token.ASTERISK, l.ch)
	case '<':
        tok.AddChar(token.LT, l.ch)
	case '>':
        tok.AddChar(token.GT, l.ch)
	case ';':
        tok.AddChar(token.SEMICOLON, l.ch)
	case ',':
        tok.AddChar(token.COMMA, l.ch)
	case '(':
        tok.AddChar(token.LPAREN, l.ch)
	case ')':
        tok.AddChar(token.RPAREN, l.ch)
	case '{':
        tok.AddChar(token.LBRACE, l.ch)
	case '}':
        tok.AddChar(token.RBRACE, l.ch)
	case 1:
        // If we ended in a new line, we will go to the next line, read a new character and then call the next token on it.
        l.nextLine()
        l.readChar()
        return l.NextToken()
	case 0:
        tok.AddStr(token.EOF, "")
	default:
		if isLetter(l.ch) {
			literal := l.readCharUntil(isLetter)
			tokeyType := token.LookupIdent(literal)
            tok.AddStr(tokeyType, literal)
			return tok
		} else if isDigit(l.ch) {
			literal := l.readCharUntil(isDigit)
            tok.AddStr(token.INT, literal)
			return tok
		} else {
            tok.AddChar(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()

	return tok
}
