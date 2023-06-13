package config_test

import (
	"testing"

	"github.com/samuelko123/git-cat/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLex(t *testing.T) {
	testcases := map[string][]config.Token{
		"":     {config.NewToken(config.Position{Line: 1, Column: 1}, config.EOF, "")},
		" \t ": {config.NewToken(config.Position{Line: 1, Column: 4}, config.EOF, "")},
		"[core]": {
			config.NewToken(config.Position{Line: 1, Column: 2}, config.SECTION, "core"),
			config.NewToken(config.Position{Line: 1, Column: 7}, config.EOF, ""),
		},
		"[core-123]": {
			config.NewToken(config.Position{Line: 1, Column: 2}, config.SECTION, "core-123"),
			config.NewToken(config.Position{Line: 1, Column: 11}, config.EOF, ""),
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
		"[core] ; comment": {
			config.NewToken(config.Position{Line: 1, Column: 2}, config.SECTION, "core"),
			config.NewToken(config.Position{Line: 1, Column: 9}, config.COMMENT, " comment"),
			config.NewToken(config.Position{Line: 1, Column: 17}, config.EOF, ""),
		},
		" # comment \n [core] ": {
			config.NewToken(config.Position{Line: 1, Column: 3}, config.COMMENT, " comment "),
			config.NewToken(config.Position{Line: 2, Column: 3}, config.SECTION, "core"),
			config.NewToken(config.Position{Line: 2, Column: 9}, config.EOF, ""),
		},
		"[remote \"origin\"] \n \t ### comment ": {
			config.NewToken(config.Position{Line: 1, Column: 2}, config.SECTION, "remote"),
			config.NewToken(config.Position{Line: 1, Column: 10}, config.SUBSECTION, "origin"),
			config.NewToken(config.Position{Line: 2, Column: 5}, config.COMMENT, "## comment "),
			config.NewToken(config.Position{Line: 2, Column: 16}, config.EOF, ""),
		},
		"user": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.KEY, "user"),
			config.NewToken(config.Position{Line: 1, Column: 5}, config.EOF, ""),
		},
		"user-123": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.KEY, "user-123"),
			config.NewToken(config.Position{Line: 1, Column: 9}, config.EOF, ""),
		},
		" user ": {
			config.NewToken(config.Position{Line: 1, Column: 2}, config.KEY, "user"),
			config.NewToken(config.Position{Line: 1, Column: 7}, config.EOF, ""),
		},
		"\nuser\n": {
			config.NewToken(config.Position{Line: 2, Column: 1}, config.KEY, "user"),
			config.NewToken(config.Position{Line: 3, Column: 1}, config.EOF, ""),
		},
		" \t user \t ": {
			config.NewToken(config.Position{Line: 1, Column: 4}, config.KEY, "user"),
			config.NewToken(config.Position{Line: 1, Column: 11}, config.EOF, ""),
		},
		"user # comment": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.KEY, "user"),
			config.NewToken(config.Position{Line: 1, Column: 7}, config.COMMENT, " comment"),
			config.NewToken(config.Position{Line: 1, Column: 15}, config.EOF, ""),
		},
		"user=\n": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.KEY, "user"),
			config.NewToken(config.Position{Line: 1, Column: 6}, config.VALUE, ""),
			config.NewToken(config.Position{Line: 2, Column: 1}, config.EOF, ""),
		},
		"user=john": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.KEY, "user"),
			config.NewToken(config.Position{Line: 1, Column: 6}, config.VALUE, "john"),
			config.NewToken(config.Position{Line: 1, Column: 10}, config.EOF, ""),
		},
		"user=john\n": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.KEY, "user"),
			config.NewToken(config.Position{Line: 1, Column: 6}, config.VALUE, "john"),
			config.NewToken(config.Position{Line: 2, Column: 1}, config.EOF, ""),
		},
		"user=\"john\"": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.KEY, "user"),
			config.NewToken(config.Position{Line: 1, Column: 7}, config.VALUE, "john"),
			config.NewToken(config.Position{Line: 1, Column: 12}, config.EOF, ""),
		},
		"user=\"jo#hn\"": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.KEY, "user"),
			config.NewToken(config.Position{Line: 1, Column: 7}, config.VALUE, "jo#hn"),
			config.NewToken(config.Position{Line: 1, Column: 13}, config.EOF, ""),
		},
		"user=\"jo;hn\"": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.KEY, "user"),
			config.NewToken(config.Position{Line: 1, Column: 7}, config.VALUE, "jo;hn"),
			config.NewToken(config.Position{Line: 1, Column: 13}, config.EOF, ""),
		},
		"user=john;comment": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.KEY, "user"),
			config.NewToken(config.Position{Line: 1, Column: 6}, config.VALUE, "john"),
			config.NewToken(config.Position{Line: 1, Column: 11}, config.COMMENT, "comment"),
			config.NewToken(config.Position{Line: 1, Column: 18}, config.EOF, ""),
		},
		"user=john#comment": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.KEY, "user"),
			config.NewToken(config.Position{Line: 1, Column: 6}, config.VALUE, "john"),
			config.NewToken(config.Position{Line: 1, Column: 11}, config.COMMENT, "comment"),
			config.NewToken(config.Position{Line: 1, Column: 18}, config.EOF, ""),
		},
		"user=john\\\\": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.KEY, "user"),
			config.NewToken(config.Position{Line: 1, Column: 6}, config.VALUE, "john\\"),
			config.NewToken(config.Position{Line: 1, Column: 12}, config.EOF, ""),
		},
		"user=john\\\"": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.KEY, "user"),
			config.NewToken(config.Position{Line: 1, Column: 6}, config.VALUE, "john\""),
			config.NewToken(config.Position{Line: 1, Column: 12}, config.EOF, ""),
		},
		"user=john\\t": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.KEY, "user"),
			config.NewToken(config.Position{Line: 1, Column: 6}, config.VALUE, "john\t"),
			config.NewToken(config.Position{Line: 1, Column: 12}, config.EOF, ""),
		},
		"user=john\\b": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.KEY, "user"),
			config.NewToken(config.Position{Line: 1, Column: 6}, config.VALUE, "john\b"),
			config.NewToken(config.Position{Line: 1, Column: 12}, config.EOF, ""),
		},
		"user=john\\n": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.KEY, "user"),
			config.NewToken(config.Position{Line: 1, Column: 6}, config.VALUE, "john\n"),
			config.NewToken(config.Position{Line: 1, Column: 12}, config.EOF, ""),
		},
		"user=john \\\n doe": {
			config.NewToken(config.Position{Line: 1, Column: 1}, config.KEY, "user"),
			config.NewToken(config.Position{Line: 1, Column: 6}, config.VALUE, "john  doe"),
			config.NewToken(config.Position{Line: 2, Column: 5}, config.EOF, ""),
		},
	}

	for input, expected := range testcases {
		fn := func() {
			lexer := config.NewLexer(input)
			lexer.Lex()
		}
		require.NotPanics(t, fn, "Input failed:\n"+input)

		lexer := config.NewLexer(input)
		token := lexer.Lex()
		assert.Equal(t, expected, token, "Input:\n"+input)
	}
}

func TestLex_Panics(t *testing.T) {
	testcases := map[string]string{
		"[core":                    "missing ] character (1:6)",
		"[core \t ":                "missing ] character (1:9)",
		"[\ncore]":                 "missing ] character (1:2)",
		"[co\nre]":                 "missing ] character (1:4)",
		"[core\n]":                 "missing ] character (1:6)",
		"[remote origin]":          "missing \" character (1:9)",
		"[remote \"ori\ngin\"]":    "missing \" character (1:13)",
		"[remote \"ori":            "missing \" character (1:13)",
		"[remote \"origin\"":       "missing ] character (1:17)",
		"[\nremote \"origin\"]":    "missing ] character (1:2)",
		"[remote\n\"origin\"]":     "missing ] character (1:8)",
		"[remote \"origin\"\n]":    "missing ] character (1:17)",
		"user]123":                 "invalid character ] (1:5)",
		"user=\"":                  "missing \" character (1:7)",
		"user=john\"":              "unexpected \" character (1:10)",
		"user=\"john":              "missing \" character (1:11)",
		"user=\"john\n":            "missing \" character (1:11)",
		"user=\"john\\a\"":         "invalid escape sequence \\a (1:11)",
		"user=\"john\\a ; comment": "invalid escape sequence \\a (1:11)",
	}

	for input, expected := range testcases {
		fn := func() {
			lexer := config.NewLexer(input)
			lexer.Lex()
		}

		assert.PanicsWithError(t, expected, fn, "Input:\n"+input)
	}
}
