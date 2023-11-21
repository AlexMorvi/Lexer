package main

//autor: De La Cruz Luis
import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)

type Token int

const (
	EOF = iota
	ILLEGAL
	IDENT
	INT
	SEMI   // ;
	STRING // "string"

	// Infix ops
	ADD // +
	SUB // -
	MUL // *
	DIV // /

	ASSIGN // =

	LPAREN // (
	RPAREN // )

)

var tokens = []string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",
	IDENT:   "IDENT",
	INT:     "INT",
	SEMI:    ";",
	STRING:  "STRING",

	// Infix ops
	ADD:    "+",
	SUB:    "-",
	MUL:    "*",
	DIV:    "/",
	ASSIGN: "=",
	LPAREN: "(",
	RPAREN: ")",
}

// String returns the string representation of the token
func (t Token) String() string {
	return tokens[t]
}

// Position represents the position of a token in the input
type Position struct {
	line   int
	column int
}

// Lexer is the lexer for the language
type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

// backup backs up the lexer by one rune
func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{line: 1, column: 0},
		reader: bufio.NewReader(reader),
	}
}

// Lex returns the next token from the input
func (l *Lexer) Lex() (Position, Token, string) {
	// keep looping until we return a token
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, ""
			}

			// at this point there isn't much we can do, and the compiler
			// should just return the raw error to the user
			panic(err)
		}

		// update the column to the position of the newly read in rune
		l.pos.column++

		switch r {
		case '\n':
			l.resetPosition()
		case ';':
			return l.pos, SEMI, ";"
		case '+':
			return l.pos, ADD, "+"
		case '-':
			return l.pos, SUB, "-"
		case '*':
			return l.pos, MUL, "*"
		case '/':
			return l.pos, DIV, "/"
		case '=':
			return l.pos, ASSIGN, "="
		case '(':
			return l.pos, LPAREN, "("
		case ')':
			return l.pos, RPAREN, ")"

		default:
			if unicode.IsSpace(r) {
				continue // nothing to do here, just move on
			} else if unicode.IsDigit(r) {
				// backup and let lexInt rescan the beginning of the int
				startPos := l.pos
				l.backup()
				lit := l.lexInt()
				return startPos, INT, lit
			} else if unicode.IsLetter(r) {
				// backup and let lexIdent rescan the beginning of the ident
				startPos := l.pos
				l.backup()
				lit := l.lexIdent()
				return startPos, IDENT, lit
			} else if r == '"' {
				// backup and let lexString rescan the beginning of the string
				startPos := l.pos
				l.backup()
				lit := l.lexString()
				return startPos, STRING, lit
			} else {
				return l.pos, ILLEGAL, string(r)
			}
		}
	}
}

func (l *Lexer) resetPosition() {
	l.pos.line++
	l.pos.column = 0
}

func (l *Lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}

	l.pos.column--
}

// lexInt scans the input until the end of an integer and then returns the
// literal.
func (l *Lexer) lexInt() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return lit
			}

			panic(err)
		}

		l.pos.column++

		if unicode.IsDigit(r) {
			lit += string(r)
		} else {
			l.backup()
			return lit
		}
	}
}

// lexIdent scans the input until the end of an identifier and then returns the
// literal.
func (l *Lexer) lexIdent() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return lit
			}

			panic(err)
		}

		l.pos.column++

		if unicode.IsLetter(r) {
			lit += string(r)
		} else {
			l.backup()
			return lit
		}
	}
}

// lexString scans the input until the end of a string and then returns the
// literal.
func (l *Lexer) lexString() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return lit
			}

			panic(err)
		}

		l.pos.column++

		if r == '"' && len(lit) >= 1 {
			return lit
		} else {
			if r != '"' {
				lit += string(r)
			}
		}
	}
}

func main() {
	file, err := os.Open("input.test")
	if err != nil {
		panic(err)
	}

	lexer := NewLexer(file)
	for {
		pos, tok, lit := lexer.Lex()
		if tok == EOF {
			break
		}

		fmt.Printf("%d:%d\t%s\t%s\n", pos.line, pos.column, tok, lit)
	}
}
