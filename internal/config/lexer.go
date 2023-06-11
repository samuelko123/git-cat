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
		c, err := l.readNextRune()

		if err != nil {
			if err == io.EOF {
				l.pushToken(NewToken(l.pos, EOF, ""))
				return l.tokens
			}

			panic(err)
		}

		switch c {
		case '[':
			l.pushToken(NewToken(l.pos, LEFT_SQUARE_BRACKET, "["))
			continue
		case ']':
			l.pushToken(NewToken(l.pos, RIGHT_SQUARE_BRACKET, "]"))
			continue
		case '"':
			l.pushToken(NewToken(l.pos, DOUBLE_QUOTE, "\""))
			continue
		case '\n':
			l.pos.Line += 1
			l.pos.Column = 0
			l.pushToken(NewToken(l.pos, LINE_BREAK, "\n"))
			continue
		case '=':
			l.pushToken(NewToken(l.pos, EQUAL_SIGN, "="))
			continue
		case ';':
			l.pushToken(NewToken(l.pos, SEMI_COLON, ";"))
			continue
		case '#':
			l.pushToken(NewToken(l.pos, HASH_SIGN, "#"))
			continue
		case '\\':
			l.pushToken(NewToken(l.pos, BACKSLASH, "\\"))
			continue
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

func (l *Lexer) pushToken(t Token) {
	l.tokens = append(l.tokens, t)
}
