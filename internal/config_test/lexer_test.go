package config_test

import (
	"testing"

	"github.com/samuelko123/git-cat/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLex(t *testing.T) {
	testcases := map[string][]config.Token{
		"": {config.NewToken(config.Position{Line: 1, Column: 1}, config.EOF, "")},
		"[": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.LEFT_SQUARE_BRACKET, "["),
			config.NewToken(config.Position{Line: 1, Column: 2}, config.EOF, ""),
		},
		"]": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.RIGHT_SQUARE_BRACKET, "]"),
			config.NewToken(config.Position{Line: 1, Column: 2}, config.EOF, ""),
		},
		"\"": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.DOUBLE_QUOTE, "\""),
			config.NewToken(config.Position{Line: 1, Column: 2}, config.EOF, ""),
		},
		"=": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.EQUAL_SIGN, "="),
			config.NewToken(config.Position{Line: 1, Column: 2}, config.EOF, ""),
		},
		";": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.SEMI_COLON, ";"),
			config.NewToken(config.Position{Line: 1, Column: 2}, config.EOF, ""),
		},
		"#": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.HASH_SIGN, "#"),
			config.NewToken(config.Position{Line: 1, Column: 2}, config.EOF, ""),
		},
		"\\": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.BACKSLASH, "\\"),
			config.NewToken(config.Position{Line: 1, Column: 2}, config.EOF, ""),
		},
		"A": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.EXPRESSION, "A"),
			config.NewToken(config.Position{Line: 1, Column: 2}, config.EOF, ""),
		},
		"A;": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.EXPRESSION, "A"),
			config.NewToken(config.Position{Line: 1, Column: 2}, config.SEMI_COLON, ";"),
			config.NewToken(config.Position{Line: 1, Column: 3}, config.EOF, ""),
		},
		"A \t ": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.EXPRESSION, "A"),
			config.NewToken(config.Position{Line: 1, Column: 2}, config.WHITESPACE, " \t "),
			config.NewToken(config.Position{Line: 1, Column: 5}, config.EOF, ""),
		},
		" \t abc \t \n123": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.WHITESPACE, " \t "),
			config.NewToken(config.Position{Line: 1, Column: 4}, config.EXPRESSION, "abc"),
			config.NewToken(config.Position{Line: 1, Column: 7}, config.WHITESPACE, " \t "),
			config.NewToken(config.Position{Line: 1, Column: 10}, config.LINE_BREAK, "\n"),
			config.NewToken(config.Position{Line: 2, Column: 1}, config.EXPRESSION, "123"),
			config.NewToken(config.Position{Line: 2, Column: 4}, config.EOF, ""),
		},
	}

	for input, expected := range testcases {
		lexer := config.NewLexer(input)

		token := lexer.Lex()

		assert.Equal(t, expected, token, "Failed input: "+input)
	}
}
