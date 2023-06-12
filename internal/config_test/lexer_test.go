package config_test

import (
	"testing"

	"github.com/samuelko123/git-cat/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLex(t *testing.T) {
	testcases := map[string][]config.Token{
		"":     {config.NewToken(config.Position{Line: 1, Column: 1}, config.EOF, "")},
		" \t ": {config.NewToken(config.Position{Line: 1, Column: 4}, config.EOF, "")},
		"[core]": {
			config.NewToken(config.Position{Line: 1, Column: 2}, config.SECTION, "core"),
			config.NewToken(config.Position{Line: 1, Column: 7}, config.EOF, ""),
		},
	}

	for input, expected := range testcases {
		lexer := config.NewLexer(input)

		token := lexer.Lex()

		assert.Equal(t, expected, token, "Failed input:\n"+input)
	}
}

func TestLex_Panics(t *testing.T) {
	testcases := map[string]string{
		"[core": "missing ] character (1:6)",
	}

	for input, expected := range testcases {
		fn := func() {
			lexer := config.NewLexer(input)
			lexer.Lex()
		}

		assert.PanicsWithError(t, expected, fn, "Input does not panic:\n"+input)
	}
}
