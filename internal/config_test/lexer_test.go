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
		"[remote \"origin\"]": {
			config.NewToken(config.Position{Line: 1, Column: 2}, config.SECTION, "remote"),
			config.NewToken(config.Position{Line: 1, Column: 10}, config.SUBSECTION, "origin"),
			config.NewToken(config.Position{Line: 1, Column: 18}, config.EOF, ""),
		},
		"[remote \"ori\\gin\"]": {
			config.NewToken(config.Position{Line: 1, Column: 2}, config.SECTION, "remote"),
			config.NewToken(config.Position{Line: 1, Column: 10}, config.SUBSECTION, "origin"),
			config.NewToken(config.Position{Line: 1, Column: 19}, config.EOF, ""),
		},
		"[remote \"ori\\\"gin\"]": {
			config.NewToken(config.Position{Line: 1, Column: 2}, config.SECTION, "remote"),
			config.NewToken(config.Position{Line: 1, Column: 10}, config.SUBSECTION, "ori\"gin"),
			config.NewToken(config.Position{Line: 1, Column: 20}, config.EOF, ""),
		},
		"[remote \"ori\\\\gin\"]": {
			config.NewToken(config.Position{Line: 1, Column: 2}, config.SECTION, "remote"),
			config.NewToken(config.Position{Line: 1, Column: 10}, config.SUBSECTION, "ori\\gin"),
			config.NewToken(config.Position{Line: 1, Column: 20}, config.EOF, ""),
		},
		"[remote \"ori]gin\"]": {
			config.NewToken(config.Position{Line: 1, Column: 2}, config.SECTION, "remote"),
			config.NewToken(config.Position{Line: 1, Column: 10}, config.SUBSECTION, "ori]gin"),
			config.NewToken(config.Position{Line: 1, Column: 19}, config.EOF, ""),
		},
	}

	for input, expected := range testcases {
		lexer := config.NewLexer(input)

		token := lexer.Lex()

		assert.Equal(t, expected, token, "Input:\n"+input)
	}
}

func TestLex_Panics(t *testing.T) {
	testcases := map[string]string{
		"[core":                 "missing ] character (1:6)",
		"[core\n]":              "missing ] character (1:6)",
		"[remote \t ":           "missing \" character (1:11)",
		"[remote origin]":       "missing \" character (1:9)",
		"[remote \"ori\ngin\"]": "unexpected newline character (1:13)",
		"[remote \"ori":         "unexpected EOF character (1:13)",
		"[remote \"origin\"":    "missing ] character (1:17)",
		"[remote \"origin\" ]":  "missing ] character (1:17)",
	}

	for input, expected := range testcases {
		fn := func() {
			lexer := config.NewLexer(input)
			lexer.Lex()
		}

		assert.PanicsWithError(t, expected, fn, "Input:\n"+input)
	}
}
