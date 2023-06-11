package config_test

import (
	"testing"

	"github.com/samuelko123/git-cat/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLex(t *testing.T) {
	testcases := map[string][]config.Token{
		"":             {config.NewToken(config.Position{Line: 1, Column: 1}, config.EOF, "")},
		"♥":            {config.NewToken(config.Position{Line: 1, Column: 2}, config.EOF, "")},
		"  \n \t \n  ": {config.NewToken(config.Position{Line: 3, Column: 3}, config.EOF, "")},
	}

	for input, expected := range testcases {
		lexer := config.NewLexer(input)

		token := lexer.Lex()

		assert.Equal(t, expected, token, "Failed: "+input)
	}
}
