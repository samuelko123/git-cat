package config_test

import (
	"testing"

	"github.com/samuelko123/git-cat/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLex(t *testing.T) {
	testcases := map[string][]config.Token{
		"": {config.NewToken(config.Position{Line: 1, Column: 1}, config.EOF, "")},
		"  \n \t \n  ": {
			config.NewToken(config.Position{Line: 1, Column: 3}, config.LINE_BREAK, "\n"),
			config.NewToken(config.Position{Line: 2, Column: 4}, config.LINE_BREAK, "\n"),
			config.NewToken(config.Position{Line: 3, Column: 3}, config.EOF, ""),
		},
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
		"\"\"": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.DOUBLE_QUOTE, "\""),
			config.NewToken(config.Position{Line: 1, Column: 2}, config.DOUBLE_QUOTE, "\""),
			config.NewToken(config.Position{Line: 1, Column: 3}, config.EOF, ""),
		},
		"\n\n": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.LINE_BREAK, "\n"),
			config.NewToken(config.Position{Line: 2, Column: 1}, config.LINE_BREAK, "\n"),
			config.NewToken(config.Position{Line: 3, Column: 1}, config.EOF, ""),
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
		"♥": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.EXPRESSION, "♥"),
			config.NewToken(config.Position{Line: 1, Column: 2}, config.EOF, ""),
		},
		"\"abc 123\"": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.DOUBLE_QUOTE, "\""),
			config.NewToken(config.Position{Line: 1, Column: 2}, config.EXPRESSION, "abc 123"),
			config.NewToken(config.Position{Line: 1, Column: 9}, config.DOUBLE_QUOTE, "\""),
			config.NewToken(config.Position{Line: 1, Column: 10}, config.EOF, ""),
		},
		"\nabc 123\n": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.LINE_BREAK, "\n"),
			config.NewToken(config.Position{Line: 2, Column: 1}, config.EXPRESSION, "abc 123"),
			config.NewToken(config.Position{Line: 2, Column: 8}, config.LINE_BREAK, "\n"),
			config.NewToken(config.Position{Line: 3, Column: 1}, config.EOF, ""),
		},
	}

	for input, expected := range testcases {
		lexer := config.NewLexer(input)

		token := lexer.Lex()

		assert.Equal(t, expected, token, "Failed input: "+input)
	}
}
