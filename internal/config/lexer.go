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
	RIGHT_SQUARE_BRACKET
	DOUBLE_QUOTE
	LINE_BREAK
	EQUAL_SIGN
	SEMI_COLON
	HASH_SIGN
	BACKSLASH
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
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		pos:    Position{Line: 1, Column: 0},
		reader: bufio.NewReader(strings.NewReader(input)),
	}
}

func (l *Lexer) Lex() []Token {
	tokens := make([]Token, 0)

	for {
		t := l.getNextToken()
		tokens = append(tokens, t)
		if t.TType == EOF {
			return tokens
		}
	}
}

func (l *Lexer) getNextToken() Token {
	for {
		c, err := l.readNextRune()

		if err != nil {
			if err == io.EOF {
				return NewToken(l.pos, EOF, "")
			}

			panic(err)
		}

		switch c {
		case '[':
			return NewToken(l.pos, LEFT_SQUARE_BRACKET, "[")
		case ']':
			return NewToken(l.pos, RIGHT_SQUARE_BRACKET, "]")
		case '"':
			return NewToken(l.pos, DOUBLE_QUOTE, "\"")
		case '\n':
			pos := l.pos
			l.pos.Line += 1
			l.pos.Column = 0
			return NewToken(pos, LINE_BREAK, "\n")
		case '=':
			return NewToken(l.pos, EQUAL_SIGN, "=")
		case ';':
			return NewToken(l.pos, SEMI_COLON, ";")
		case '#':
			return NewToken(l.pos, HASH_SIGN, "#")
		case '\\':
			return NewToken(l.pos, BACKSLASH, "\\")
		default:
			if unicode.IsSpace(c) {
				continue
			}
		}
	}
}

func (l *Lexer) readNextRune() (rune, error) {
	c, _, err := l.reader.ReadRune()
	l.pos.Column += 1

	return c, err
}
