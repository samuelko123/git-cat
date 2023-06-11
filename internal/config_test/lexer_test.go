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
		" [core] \n ignorecase = false ": {
			config.NewToken(config.Position{Line: 1, Column: 2}, config.LEFT_SQUARE_BRACKET, "["),
			config.NewToken(config.Position{Line: 1, Column: 3}, config.EXPRESSION, "core"),
			config.NewToken(config.Position{Line: 1, Column: 7}, config.RIGHT_SQUARE_BRACKET, "]"),
			config.NewToken(config.Position{Line: 1, Column: 9}, config.LINE_BREAK, "\n"),
			config.NewToken(config.Position{Line: 2, Column: 2}, config.EXPRESSION, "ignorecase "),
			config.NewToken(config.Position{Line: 2, Column: 13}, config.EQUAL_SIGN, "="),
			config.NewToken(config.Position{Line: 2, Column: 15}, config.EXPRESSION, "false "),
			config.NewToken(config.Position{Line: 2, Column: 21}, config.EOF, ""),
		},
		" [remote \"origin\"] # comment 1 \n url = https://pkg.go.dev ; comment 2": {
			config.NewToken(config.Position{Line: 1, Column: 2}, config.LEFT_SQUARE_BRACKET, "["),
			config.NewToken(config.Position{Line: 1, Column: 3}, config.EXPRESSION, "remote "),
			config.NewToken(config.Position{Line: 1, Column: 10}, config.DOUBLE_QUOTE, "\""),
			config.NewToken(config.Position{Line: 1, Column: 11}, config.EXPRESSION, "origin"),
			config.NewToken(config.Position{Line: 1, Column: 17}, config.DOUBLE_QUOTE, "\""),
			config.NewToken(config.Position{Line: 1, Column: 18}, config.RIGHT_SQUARE_BRACKET, "]"),
			config.NewToken(config.Position{Line: 1, Column: 20}, config.HASH_SIGN, "#"),
			config.NewToken(config.Position{Line: 1, Column: 22}, config.EXPRESSION, "comment 1 "),
			config.NewToken(config.Position{Line: 1, Column: 32}, config.LINE_BREAK, "\n"),
			config.NewToken(config.Position{Line: 2, Column: 2}, config.EXPRESSION, "url "),
			config.NewToken(config.Position{Line: 2, Column: 6}, config.EQUAL_SIGN, "="),
			config.NewToken(config.Position{Line: 2, Column: 8}, config.EXPRESSION, "https://pkg.go.dev "),
			config.NewToken(config.Position{Line: 2, Column: 27}, config.SEMI_COLON, ";"),
			config.NewToken(config.Position{Line: 2, Column: 29}, config.EXPRESSION, "comment 2"),
			config.NewToken(config.Position{Line: 2, Column: 38}, config.EOF, ""),
		},
	}

	for input, expected := range testcases {
		lexer := config.NewLexer(input)

		token := lexer.Lex()

		assert.Equal(t, expected, token, "Failed input: "+input)
	}
}
