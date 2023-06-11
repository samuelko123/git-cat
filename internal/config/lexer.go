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
	tokens := make([]Token, 0)

	for {
		c, err := l.readNextRune()

		if err != nil {
			if err == io.EOF {
				if expr != "" {
					tokens = append(tokens, NewToken(exprPos, EXPRESSION, expr))
				}
				tokens = append(tokens, NewToken(l.pos, EOF, ""))
				return tokens
			}

			panic(err)
		}

		if slices.Contains(specialRunes, c) {
			if expr != "" {
				tokens = append(tokens, NewToken(exprPos, EXPRESSION, expr))
			}
			expr = ""

			if c == '\n' {
				pos := l.pos
				l.pos.Line += 1
				l.pos.Column = 0
				tokens = append(tokens, NewToken(pos, runeMap[c], string(c)))
			} else {
				tokens = append(tokens, NewToken(l.pos, runeMap[c], string(c)))
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

func (l *Lexer) readNextRune() (rune, error) {
	c, _, err := l.reader.ReadRune()
	l.pos.Column += 1

	return c, err
}

func (l *Lexer) unreadRune() {
	l.reader.UnreadRune()
	l.pos.Column -= 1
}
