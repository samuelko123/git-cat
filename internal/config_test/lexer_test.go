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
		"[ core ]": {
			config.NewToken(config.Position{Line: 1, Column: 3}, config.SECTION, "core"),
			config.NewToken(config.Position{Line: 1, Column: 9}, config.EOF, ""),
		},
		"\n[core]\n": {
			config.NewToken(config.Position{Line: 2, Column: 2}, config.SECTION, "core"),
			config.NewToken(config.Position{Line: 3, Column: 1}, config.EOF, ""),
		},
		" \t [core] \t ": {
			config.NewToken(config.Position{Line: 1, Column: 5}, config.SECTION, "core"),
			config.NewToken(config.Position{Line: 1, Column: 13}, config.EOF, ""),
		},
		"[remote \"origin\"]": {
			config.NewToken(config.Position{Line: 1, Column: 2}, config.SECTION, "remote"),
			config.NewToken(config.Position{Line: 1, Column: 10}, config.SUBSECTION, "origin"),
			config.NewToken(config.Position{Line: 1, Column: 18}, config.EOF, ""),
		},
		"[ \t remote \"origin\" \t ]": {
			config.NewToken(config.Position{Line: 1, Column: 5}, config.SECTION, "remote"),
			config.NewToken(config.Position{Line: 1, Column: 13}, config.SUBSECTION, "origin"),
			config.NewToken(config.Position{Line: 1, Column: 24}, config.EOF, ""),
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
		"[core \t ":             "missing ] character (1:9)",
		"[\ncore]":              "unexpected newline character (1:2)",
		"[co\nre]":              "unexpected newline character (1:4)",
		"[core\n]":              "unexpected newline character (1:6)",
		"[remote origin]":       "missing \" character (1:9)",
		"[remote \"ori\ngin\"]": "unexpected newline character (1:13)",
		"[remote \"ori":         "unexpected EOF character (1:13)",
		"[remote \"origin\"":    "missing ] character (1:17)",
		"[\nremote \"origin\"]": "unexpected newline character (1:2)",
		"[remote\n\"origin\"]":  "unexpected newline character (1:8)",
		"[remote \"origin\"\n]": "unexpected newline character (1:17)",
	}

	for input, expected := range testcases {
		fn := func() {
			lexer := config.NewLexer(input)
			lexer.Lex()
		}

		assert.PanicsWithError(t, expected, fn, "Input:\n"+input)
	}
}
