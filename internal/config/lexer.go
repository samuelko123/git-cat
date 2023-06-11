package config

import (
	"bufio"
	"io"
	"strings"
	"unicode"

	"golang.org/x/exp/slices"
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

var runeMap = map[rune]TokenType{
	'[':  LEFT_SQUARE_BRACKET,
	']':  RIGHT_SQUARE_BRACKET,
	'"':  DOUBLE_QUOTE,
	'\n': LINE_BREAK,
	'=':  EQUAL_SIGN,
	';':  SEMI_COLON,
	'#':  HASH_SIGN,
	'\\': BACKSLASH,
}

var specialRunes = []rune{
	'[', ']', '"', '\n', '=', ';', '#', '\\',
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
	tokens []Token
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		pos:    Position{Line: 1, Column: 0},
		reader: bufio.NewReader(strings.NewReader(input)),
	}
}

func (l *Lexer) Lex() []Token {
	expr := ""
	var exprPos Position

	for {
		c := l.readNextRune()

		if c == rune(0) {
			if expr != "" {
				l.tokens = append(l.tokens, NewToken(exprPos, EXPRESSION, expr))
			}
			l.tokens = append(l.tokens, NewToken(l.pos, EOF, ""))
			return l.tokens
		}

		if slices.Contains(specialRunes, c) {
			if expr != "" {
				l.tokens = append(l.tokens, NewToken(exprPos, EXPRESSION, expr))
			}
			expr = ""

			if c == '\n' {
				pos := l.pos
				l.pos.Line += 1
				l.pos.Column = 0
				l.tokens = append(l.tokens, NewToken(pos, runeMap[c], string(c)))
			} else {
				l.tokens = append(l.tokens, NewToken(l.pos, runeMap[c], string(c)))
			}
		} else if unicode.IsSpace(c) && expr == "" {
			continue
		} else {
			if expr == "" {
				exprPos = l.pos
			}
			expr += string(c)
		}
	}
}

func (l *Lexer) readNextRune() rune {
	c, _, err := l.reader.ReadRune()
	l.pos.Column += 1

	if err != nil && err != io.EOF {
		panic(err)
	}

	return c
}

func (l *Lexer) unreadRune() {
	l.reader.UnreadRune()
	l.pos.Column -= 1
}
