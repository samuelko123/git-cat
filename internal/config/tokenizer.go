package config

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/samuelko123/git-cat/internal/utils"
)

// https://git-scm.com/docs/git-config#_syntax

var (
	SECTION_NAME_REGEX           = regexp.MustCompile("^[a-z-.]+$")
	SUBSECTION_NAME_REGEX        = regexp.MustCompile("^[^\n\u0000]*$")
	SUBSECTION_NAME_ESCAPE_REGEX = regexp.MustCompile("([\\\\])([^\\\\]|[\\\\])")
	VARIABLE_NAME_REGEX          = regexp.MustCompile("^[a-z-]+$")
)

type TokenType string

const (
	UNDEFINED  TokenType = "undefined"
	SECTION    TokenType = "section"
	SUBSECTION TokenType = "subsection"
	NAME       TokenType = "name"
	VALUE      TokenType = "value"
	COMMENT    TokenType = "comment"
)

type Token struct {
	Type  TokenType
	Value string
}

type Tokenizer struct {
	tokens    []Token
	prevToken Token
	currToken Token
	lineNum   int
	colNum    int
}

func (t *Tokenizer) Tokenize(input string) (_ []Token, err error) {
	// TODO
	// variable value - escape sequence, backslash next line
	defer utils.ReturnError(&err)

	t.tokens = make([]Token, 0)
	t.setCurrToken(UNDEFINED)
	t.lineNum = 1
	t.colNum = 1

	if strings.TrimSpace(input) == "" {
		return t.tokens, nil
	}

	firstChar := strings.TrimSpace(input)[0]
	if firstChar != '[' {
		panic(errors.New("config should begin with the [ character, got " + string(firstChar)))
	}

	for _, c := range input {
		t.colNum++
		if c == '\n' {
			t.handleLineBreak(c)
		} else if unicode.IsSpace(c) {
			t.handleWhiteSpace(c)
		} else if c == '[' {
			t.handleOpeningSquareBracket(c)
		} else if c == ']' {
			t.handleClosingSquareBracket(c)
		} else if c == '=' {
			t.handleEqualSign(c)
		} else if c == '#' || c == ';' {
			t.handleCommentChar(c)
		} else if c == '"' {
			t.handleDoubleQuote(c)
		} else {
			t.handleOtherChar(c)
		}
	}

	t.flushCurrToken()

	return t.tokens, nil
}

func (t *Tokenizer) handleWhiteSpace(c rune) {
	if t.currToken.Type == UNDEFINED {
		if t.prevToken.Type == SUBSECTION {
			t.doPanic()
		} else {
			return
		}
	} else if t.currToken.Type == SECTION {
		if t.currToken.Value == "" {
			t.doPanic()
		} else {
			t.flushCurrToken()
			t.setCurrToken(UNDEFINED)
		}
	}

	t.appendRune(c)
}

func (t *Tokenizer) handleOpeningSquareBracket(c rune) {
	if t.currToken.Type == UNDEFINED {
		t.setCurrToken(SECTION)
		return
	}

	t.appendRune(c)
}

func (t *Tokenizer) handleClosingSquareBracket(c rune) {
	if t.currToken.Type == UNDEFINED {
		if t.prevToken.Type == SECTION {
			t.doPanic()
		}
	} else if t.currToken.Type == SECTION {
		t.flushCurrToken()
		t.setCurrToken(UNDEFINED)
	} else if t.currToken.Type == NAME || t.currToken.Type == SUBSECTION {
		t.doPanic()
	} else {
		t.appendRune(c)
	}
}

func (t *Tokenizer) handleLineBreak(c rune) {
	if t.currToken.Type == VALUE || t.currToken.Type == COMMENT {
		t.flushCurrToken()
		t.setCurrToken(UNDEFINED)
	} else if t.currToken.Type != UNDEFINED {
		t.doPanic()
	}

	t.lineNum++
	t.colNum = 1
}

func (t *Tokenizer) handleEqualSign(c rune) {
	if t.currToken.Type == NAME {
		t.flushCurrToken()
		t.setCurrToken(VALUE)
		return
	}

	t.appendRune(c)
}

func (t *Tokenizer) handleCommentChar(c rune) {
	if t.currToken.Type == UNDEFINED {
		t.setCurrToken(COMMENT)
		return
	}

	t.appendRune(c)
}

func (t *Tokenizer) handleDoubleQuote(c rune) {
	if t.currToken.Type == SUBSECTION {
		val := t.currToken.Value
		if val == "" || val[len(val)-1:] != "\\" {
			t.flushCurrToken()
			t.setCurrToken(UNDEFINED)
			return
		}
	} else if t.prevToken.Type == SECTION {
		t.flushCurrToken()
		t.setCurrToken(SUBSECTION)
		return
	}

	t.appendRune(c)
}

func (t *Tokenizer) handleOtherChar(c rune) {
	if t.currToken.Type == UNDEFINED {
		t.setCurrToken(NAME)
	}

	t.appendRune(c)
}

func (t *Tokenizer) appendRune(c rune) {
	t.currToken.Value += string(c)
}

func (t *Tokenizer) flushCurrToken() {
	if t.currToken.Type == SECTION {
		name := strings.ToLower(strings.TrimSpace(t.currToken.Value))
		if !SECTION_NAME_REGEX.MatchString(name) {
			panic(errors.New(fmt.Sprintf("invalid section name %s on line %d", name, t.lineNum)))
		}
		t.currToken.Value = name
	} else if t.currToken.Type == SUBSECTION {
		name := SUBSECTION_NAME_ESCAPE_REGEX.ReplaceAllString(t.currToken.Value, "$2")
		if !SUBSECTION_NAME_REGEX.MatchString(name) {
			panic(errors.New(fmt.Sprintf("invalid subsection name %s on line %d", name, t.lineNum)))
		}
		t.currToken.Value = name
	} else if t.currToken.Type == NAME {
		name := strings.ToLower(strings.TrimSpace(t.currToken.Value))
		if !VARIABLE_NAME_REGEX.MatchString(name) {
			panic(errors.New(fmt.Sprintf("invalid variable name %s on line %d", name, t.lineNum)))
		}
		t.currToken.Value = name
	} else if t.currToken.Type == VALUE {
		val := t.currToken.Value
		val = strings.TrimSpace(val)
		if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
			val = val[1 : len(val)-1]
		}

		if regexp.MustCompile("[^\\\\][\"]|^[\"]").FindString(val) != "" {
			panic(errors.New(fmt.Sprintf("invalid variable value %s on line %d", val, t.lineNum)))
		}

		val = strings.ReplaceAll(val, "\\\\", "\\")
		val = strings.ReplaceAll(val, "\\\"", "\"")
		val = strings.ReplaceAll(val, "\\n", "\n")
		val = strings.ReplaceAll(val, "\\t", "\t")
		val = strings.ReplaceAll(val, "\\b", "\b")

		t.currToken.Value = val
	} else if t.currToken.Type == UNDEFINED {
		return
	}

	t.prevToken = t.currToken
	t.tokens = append(t.tokens, t.currToken)
}

func (t *Tokenizer) setCurrToken(tokenType TokenType) {
	t.currToken = Token{
		Type:  tokenType,
		Value: "",
	}
}

func (t *Tokenizer) doPanic() {
	panic(errors.New(fmt.Sprintf("unexpected character on line %d column %d", t.lineNum, t.colNum)))
}
