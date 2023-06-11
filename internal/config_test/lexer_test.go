package config_test

import (
	"testing"

	"github.com/samuelko123/git-cat/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLex(t *testing.T) {
	testcases := map[string][]config.Token{
		"":  {config.NewToken(config.Position{Line: 1, Column: 1}, config.EOF, "")},
		"♥": {config.NewToken(config.Position{Line: 1, Column: 2}, config.EOF, "")},
		"  \n \t \n  ": {
			config.NewToken(config.Position{Line: 2, Column: 0}, config.LINE_BREAK, "\n"),
			config.NewToken(config.Position{Line: 3, Column: 0}, config.LINE_BREAK, "\n"),
			config.NewToken(config.Position{Line: 3, Column: 3}, config.EOF, ""),
		},
		"[]\"\n = ;#\\": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.LEFT_SQUARE_BRACKET, "["),
			config.NewToken(config.Position{Line: 1, Column: 2}, config.RIGHT_SQUARE_BRACKET, "]"),
			config.NewToken(config.Position{Line: 1, Column: 3}, config.DOUBLE_QUOTE, "\""),
			config.NewToken(config.Position{Line: 2, Column: 0}, config.LINE_BREAK, "\n"),
			config.NewToken(config.Position{Line: 2, Column: 2}, config.EQUAL_SIGN, "="),
			config.NewToken(config.Position{Line: 2, Column: 4}, config.SEMI_COLON, ";"),
			config.NewToken(config.Position{Line: 2, Column: 5}, config.HASH_SIGN, "#"),
			config.NewToken(config.Position{Line: 2, Column: 6}, config.BACKSLASH, "\\"),
			config.NewToken(config.Position{Line: 2, Column: 7}, config.EOF, ""),
		},
	}

	for input, expected := range testcases {
		lexer := config.NewLexer(input)

		token := lexer.Lex()

		assert.Equal(t, expected, token, "Failed input: "+input)
	}
}
