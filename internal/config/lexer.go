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
	SUBSECTION
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
	l.readNextNonSpaceRune()
	pos := l.pos
	l.unreadRune()

	for {
		r := l.readNextRune()

		if r == '\n' {
			panic(errors.New(fmt.Sprintf("unexpected newline character (%d:%d)", l.pos.Line, l.pos.Column)))
		} else if unicode.IsDigit(r) || unicode.IsLetter(r) || r == '-' || r == '.' {
			literal += string(r)
		} else if unicode.IsSpace(r) {
			l.tokens = append(l.tokens, NewToken(pos, SECTION, literal))

			r := l.readNextNonSpaceRune()
			if r == ']' {
				return
			} else if r == '\n' || r == rune(0) {
				panic(errors.New(fmt.Sprintf("missing ] character (%d:%d)", l.pos.Line, l.pos.Column)))
			} else if r != '"' {
				panic(errors.New(fmt.Sprintf("missing \" character (%d:%d)", l.pos.Line, l.pos.Column)))
			}

			l.lexSubSection()

			return
		} else if r == ']' {
			l.tokens = append(l.tokens, NewToken(pos, SECTION, literal))
			return
		} else {
			panic(errors.New(fmt.Sprintf("missing ] character (%d:%d)", l.pos.Line, l.pos.Column)))
		}
	}
}

func (l *Lexer) lexSubSection() {
	literal := ""
	pos := l.getNextPos()

	for {
		r := l.readNextRune()
		switch r {
		case '\n':
			panic(errors.New(fmt.Sprintf("unexpected newline character (%d:%d)", l.pos.Line, l.pos.Column)))
		case rune(0):
			panic(errors.New(fmt.Sprintf("unexpected EOF character (%d:%d)", l.pos.Line, l.pos.Column)))
		case '\\':
			r = l.readNextRune()
			if r == '"' {
				literal += "\""
			} else if r == '\\' {
				literal += "\\"
			} else {
				literal += string(r)
			}
		case '"':
			l.tokens = append(l.tokens, NewToken(pos, SUBSECTION, literal))

			r = l.readNextNonSpaceRune()
			if r == '\n' {
				panic(errors.New(fmt.Sprintf("unexpected newline character (%d:%d)", l.pos.Line, l.pos.Column)))
			} else if r != ']' {
				panic(errors.New(fmt.Sprintf("missing ] character (%d:%d)", l.pos.Line, l.pos.Column)))
			}

			return
		default:
			literal += string(r)
		}
	}
}

func (l *Lexer) readNextNonSpaceRune() rune {
	var r rune

	for {
		r = l.readNextRune()
		if !unicode.IsSpace(r) || r == '\n' {
			break
		}
	}

	return r
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

func (l *Lexer) getNextPos() Position {
	pos := l.pos
	pos.Column++
	return pos
}
