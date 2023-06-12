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
	COMMENT
	KEY
)

const (
	ERR_MISSING_CLOSING_BRACKET string = "missing ] character (%d:%d)"
	ERR_MISSING_QUOTE           string = "missing \" character (%d:%d)"
	ERR_INVALID_CHARACTER       string = "invalid character %s (%d:%d)"
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
	currPos  Position
	prevPos  Position
	reader   *bufio.Reader
	tokens   []Token
	currRune rune
	prevRune rune
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		currPos: Position{Line: 1, Column: 0},
		reader:  bufio.NewReader(strings.NewReader(input)),
	}
}

func (l *Lexer) Lex() []Token {
	for {
		r := l.readNextRune()

		if r == rune(0) {
			l.tokens = append(l.tokens, NewToken(l.currPos, EOF, ""))
			return l.tokens
		}

		if unicode.IsSpace(r) {
			continue
		}

		if r == '[' {
			l.lexSection()
		} else if r == ';' || r == '#' {
			l.lexComment()
		} else if r == '=' {
			l.lexValue()
		} else {
			l.unreadRune()
			l.lexKey()
		}
	}
}

func (l *Lexer) lexSection() {
	literal := ""
	l.readNextNonSpaceRune()
	pos := l.currPos
	l.unreadRune()

	for {
		r := l.readNextRune()

		if r == '\n' {
			panic(errors.New(fmt.Sprintf(ERR_MISSING_CLOSING_BRACKET, l.currPos.Line, l.currPos.Column)))
		} else if unicode.IsDigit(r) || unicode.IsLetter(r) || r == '-' || r == '.' {
			literal += string(r)
		} else if unicode.IsSpace(r) {
			l.tokens = append(l.tokens, NewToken(pos, SECTION, literal))

			r := l.readNextNonSpaceRune()
			if r == ']' {
				return
			} else if r == '\n' || r == rune(0) {
				panic(errors.New(fmt.Sprintf(ERR_MISSING_CLOSING_BRACKET, l.currPos.Line, l.currPos.Column)))
			} else if r != '"' {
				panic(errors.New(fmt.Sprintf(ERR_MISSING_QUOTE, l.currPos.Line, l.currPos.Column)))
			}

			l.lexSubSection()

			return
		} else if r == ']' {
			l.tokens = append(l.tokens, NewToken(pos, SECTION, literal))
			return
		} else {
			panic(errors.New(fmt.Sprintf(ERR_MISSING_CLOSING_BRACKET, l.currPos.Line, l.currPos.Column)))
		}
	}
}

func (l *Lexer) lexSubSection() {
	literal := ""
	pos := l.getNextPos()

	for {
		r := l.readNextRune()
		switch r {
		case '\n', rune(0):
			panic(errors.New(fmt.Sprintf(ERR_MISSING_QUOTE, l.currPos.Line, l.currPos.Column)))
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
				panic(errors.New(fmt.Sprintf(ERR_MISSING_CLOSING_BRACKET, l.currPos.Line, l.currPos.Column)))
			} else if r != ']' {
				panic(errors.New(fmt.Sprintf(ERR_MISSING_CLOSING_BRACKET, l.currPos.Line, l.currPos.Column)))
			}

			return
		default:
			literal += string(r)
		}
	}
}

func (l *Lexer) lexComment() {
	literal := ""
	pos := l.getNextPos()

	for {
		r := l.readNextRune()
		switch r {
		case '\n', rune(0):
			l.unreadRune()
			l.tokens = append(l.tokens, NewToken(pos, COMMENT, literal))
			return
		default:
			literal += string(r)
		}
	}
}

func (l *Lexer) lexKey() {
	literal := ""
	pos := l.getNextPos()

	for {
		r := l.readNextRune()

		if r == '\n' || r == rune(0) || r == '=' || unicode.IsSpace(r) {
			l.unreadRune()
			l.tokens = append(l.tokens, NewToken(pos, KEY, literal))
			return
		} else if unicode.IsDigit(r) || unicode.IsLetter(r) || r == '-' {
			literal += string(r)
		} else {
			panic(errors.New(fmt.Sprintf(ERR_INVALID_CHARACTER, string(r), l.currPos.Line, l.currPos.Column)))
		}
	}
}

func (l *Lexer) lexValue() {
	// TODO
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
	l.prevPos = l.currPos
	l.prevRune = l.currRune
	l.currRune = r
	if l.prevRune == '\n' {
		l.currPos.Line += 1
		l.currPos.Column = 1
	} else {
		l.currPos.Column += 1
	}

	return r
}

func (l *Lexer) unreadRune() {
	l.reader.UnreadRune()
	l.currPos = l.prevPos
	l.currRune = l.prevRune
}

func (l *Lexer) getNextPos() Position {
	l.readNextRune()
	pos := l.currPos
	l.unreadRune()
	return pos
}
