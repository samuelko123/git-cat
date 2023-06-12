package config

import (
	"bufio"
	"strings"
)

type Position struct {
	Line   int
	Column int
}

type TokenType int

const (
	EOF TokenType = iota
)

type Token struct {
	Pos    Position
	TType  TokenType
	TValue string
}

func NewToken(pos Position, tType TokenType, tValue string) Token {
	return Token{
		Pos:    pos,
		TType:  tType,
		TValue: tValue,
	}
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
	tokens []Token
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		pos:    Position{Line: 1, Column: 0},
		reader: bufio.NewReader(strings.NewReader(input)),
	}
}

func (l *Lexer) Lex() []Token {
	for {
		r := l.readNextRune()

		if r == rune(0) {
			l.tokens = append(l.tokens, NewToken(l.pos, EOF, ""))
			return l.tokens
		}
	}
}

func (l *Lexer) readNextRune() rune {
	r, _, _ := l.reader.ReadRune()
	l.pos.Column += 1

	return r
}
