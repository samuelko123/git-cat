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
	WHITESPACE
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
	for {
		c := l.readNextRune()

		if c == rune(0) {
			l.tokens = append(l.tokens, NewToken(l.pos, EOF, ""))
			return l.tokens
		}

		if slices.Contains(specialRunes, c) {
			if c == '\n' {
				pos := l.pos
				l.pos.Line += 1
				l.pos.Column = 0
				l.tokens = append(l.tokens, NewToken(pos, runeMap[c], string(c)))
			} else {
				l.tokens = append(l.tokens, NewToken(l.pos, runeMap[c], string(c)))
			}
		} else if unicode.IsSpace(c) {
			l.lexSpace()
		} else {
			l.lexExpression()
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

func (l *Lexer) lexExpression() {
	expr := ""
	pos := l.pos
	l.unreadRune()

	for {
		c := l.readNextRune()
		if c == rune(0) || slices.Contains(specialRunes, c) || unicode.IsSpace(c) {
			l.unreadRune()
			l.tokens = append(l.tokens, NewToken(pos, EXPRESSION, expr))
			return
		}
		expr += string(c)
	}
}

func (l *Lexer) lexSpace() {
	expr := ""
	pos := l.pos
	l.unreadRune()

	for {
		c := l.readNextRune()
		if !unicode.IsSpace(c) || c == '\n' {
			l.unreadRune()
			l.tokens = append(l.tokens, NewToken(pos, WHITESPACE, expr))
			return
		}
		expr += string(c)
	}
}
