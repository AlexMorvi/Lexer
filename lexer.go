package main

import (
	"bufio"
	"io"
)

type Token int

const (
	EOF = iota
	ILLEGAL
	IDENT
	INT
	SEMI // ;

	// Infix ops
	ADD // +
	SUB // -
	MUL // *
	DIV // /

	ASSIGN // =
)

var tokens = []string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",
	IDENT:   "IDENT",
	INT:     "INT",
	SEMI:    ";",

	// Infix ops
	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",

	ASSIGN: "=",
}

func (t Token) String() string {
	return tokens[t]
}

type Position struct {
	line   int
	column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{line: 1, column: 0},
		reader: bufio.NewReader(reader),
	}
}
