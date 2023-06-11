package config

import (
	"bufio"
	"io"
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
	EXPRESSION
	LEFT_SQUARE_BRACKET
	RIGHT_SQUARE_BRAKCET
	DOUBLE_QUOTE
	LINE_BREAK
	EQUAL_SIGN
	SEMI_COLON
	HASH_SIGN
)

var TOKEN_TYPES = []rune{
	LEFT_SQUARE_BRACKET:  '[',
	RIGHT_SQUARE_BRAKCET: ']',
	DOUBLE_QUOTE:         '"',
	LINE_BREAK:           '\n',
	EQUAL_SIGN:           '=',
	SEMI_COLON:           ';',
	HASH_SIGN:            '#',
}

func (t TokenType) Rune() rune {
	return TOKEN_TYPES[t]
}

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
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		pos:    Position{Line: 1, Column: 0},
		reader: bufio.NewReader(strings.NewReader(input)),
	}
}

func (l *Lexer) Lex() Token {
	for {
		c, err := l.readNextRune()

		if c == LINE_BREAK.Rune() {
			l.pos.Line += 1
			l.pos.Column = 0
		}

		if unicode.IsSpace(c) {
			continue
		}

		if err != nil {
			if err == io.EOF {
				return NewToken(l.pos, EOF, "")
			}

			panic(err)
		}
	}
}

func (l *Lexer) readNextRune() (rune, error) {
	c, _, err := l.reader.ReadRune()
	l.pos.Column += 1

	return c, err
}
