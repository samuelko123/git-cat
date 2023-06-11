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
		c, err := l.readNextRune()

		if err != nil {
			if err == io.EOF {
				t := NewToken(l.pos, EOF, "")
				tokens = append(tokens, t)
				return tokens
			}

			panic(err)
		}

		switch c {
		case '[':
			t := NewToken(l.pos, LEFT_SQUARE_BRACKET, "[")
			tokens = append(tokens, t)
		case ']':
			t := NewToken(l.pos, RIGHT_SQUARE_BRACKET, "]")
			tokens = append(tokens, t)
		case '"':
			t := NewToken(l.pos, DOUBLE_QUOTE, "\"")
			tokens = append(tokens, t)
		case '\n':
			pos := l.pos
			l.pos.Line += 1
			l.pos.Column = 0
			t := NewToken(pos, LINE_BREAK, "\n")
			tokens = append(tokens, t)
		case '=':
			t := NewToken(l.pos, EQUAL_SIGN, "=")
			tokens = append(tokens, t)
		case ';':
			t := NewToken(l.pos, SEMI_COLON, ";")
			tokens = append(tokens, t)
		case '#':
			t := NewToken(l.pos, HASH_SIGN, "#")
			tokens = append(tokens, t)
		case '\\':
			t := NewToken(l.pos, BACKSLASH, "\\")
			tokens = append(tokens, t)
		default:
			if unicode.IsSpace(c) {
				continue
			} else {
				t := l.getNextExprToken()
				tokens = append(tokens, t)
			}
		}
	}
}

func (l *Lexer) readNextRune() (rune, error) {
	c, _, err := l.reader.ReadRune()
	l.pos.Column += 1

	return c, err
}

func (l *Lexer) unreadRune() {
	l.reader.UnreadRune()
	l.pos.Column -= 1
}

func (l *Lexer) getNextExprToken() Token {
	startPos := l.pos
	literal := ""

	l.unreadRune()

	for {
		c, err := l.readNextRune()

		if err != nil {
			if err == io.EOF {
				l.unreadRune()
				return NewToken(startPos, EXPRESSION, literal)
			}

			panic(err)
		}

		switch c {
		case '"', '\n':
			l.unreadRune()
			return NewToken(startPos, EXPRESSION, literal)
		default:
			literal += string(c)
		}
	}
}
