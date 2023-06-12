package config

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
	"unicode"
)

type Position struct {
	Line   int
	Column int
}

type TokenType int

const (
	EOF TokenType = iota
	SECTION
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

		if unicode.IsSpace(r) {
			continue
		}

		if r == '[' {
			l.lexSection()
		}
	}
}

func (l *Lexer) lexSection() {
	literal := ""
	l.readNextRune()
	pos := l.pos
	l.unreadRune()

	for {
		r := l.readNextRune()

		if unicode.IsDigit(r) || unicode.IsLetter(r) || r == '-' || r == '.' {
			literal += string(r)
		} else if r == ']' {
			l.tokens = append(l.tokens, NewToken(pos, SECTION, literal))
			return
		} else {
			panic(errors.New(fmt.Sprintf("missing ] character (%d:%d)", l.pos.Line, l.pos.Column)))
		}
	}
}

func (l *Lexer) readNextRune() rune {
	r, _, _ := l.reader.ReadRune()
	l.pos.Column += 1

	return r
}

func (l *Lexer) unreadRune() {
	l.reader.UnreadRune()
	l.pos.Column -= 1
}
