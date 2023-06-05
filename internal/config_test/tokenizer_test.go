package config

import (
	"testing"

	"github.com/samuelko123/git-cat/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_ValidCases(t *testing.T) {
	testcases := map[string][]config.Token{
		"": {},
		"[remote]": {
			{Type: config.SECTION, Value: "remote"},
		},
		"[ReMote]": {
			{Type: config.SECTION, Value: "remote"},
		},
		" \t [remote] \t ": {
			{Type: config.SECTION, Value: "remote"},
		},
		"[remote \"origin\"]": {
			{Type: config.SECTION, Value: "remote"},
			{Type: config.SUBSECTION, Value: "origin"},
		},
		"[remote \t \"origin\"]": {
			{Type: config.SECTION, Value: "remote"},
			{Type: config.SUBSECTION, Value: "origin"},
		},
		"[remote \"origin\"];some comments": {
			{Type: config.SECTION, Value: "remote"},
			{Type: config.SUBSECTION, Value: "origin"},
			{Type: config.COMMENT, Value: "some comments"},
		},
		"[remote \"origin\"];some comments[core]": {
			{Type: config.SECTION, Value: "remote"},
			{Type: config.SUBSECTION, Value: "origin"},
			{Type: config.COMMENT, Value: "some comments[core]"},
		},
		"[remote \" OriGin\t\"]": {
			{Type: config.SECTION, Value: "remote"},
			{Type: config.SUBSECTION, Value: " OriGin\t"},
		},
		"[core]\nignorecase=true": {
			{Type: config.SECTION, Value: "core"},
			{Type: config.NAME, Value: "ignorecase"},
			{Type: config.VALUE, Value: "true"},
		},
		"[core]\nignore-case=false": {
			{Type: config.SECTION, Value: "core"},
			{Type: config.NAME, Value: "ignore-case"},
			{Type: config.VALUE, Value: "false"},
		},
		"[core]\nignoreCase=false": {
			{Type: config.SECTION, Value: "core"},
			{Type: config.NAME, Value: "ignorecase"},
			{Type: config.VALUE, Value: "false"},
		},
		"[core]\n \t ignorecase \t = \t true \t ": {
			{Type: config.SECTION, Value: "core"},
			{Type: config.NAME, Value: "ignorecase"},
			{Type: config.VALUE, Value: "true"},
		},
	}

	for input, expected := range testcases {
		tokenizer := config.Tokenizer{}

		tokens, err := tokenizer.Tokenize(input)
		require.Nil(t, err, input)
		assert.Equal(t, expected, tokens, "This input should not fail: "+input)
	}
}

func TestParse_InvalidCases(t *testing.T) {
	testcases := map[string]string{
		"abc":                   "config should begin with the [ character, got a",
		"[ remote]":             "unexpected character on line 1 column 3",
		"[\tremote]":            "unexpected character on line 1 column 3",
		"[remote ]":             "unexpected character on line 1 column 10",
		"[remote\t]":            "unexpected character on line 1 column 10",
		"[remote \"origin\" ]":  "unexpected character on line 1 column 18",
		"[remote \"origin\"\t]": "unexpected character on line 1 column 18",
		"[core]\nvar1=true":     "invalid variable name var1 on line 2",
	}

	for input, expected := range testcases {
		tokenizer := config.Tokenizer{}

		_, err := tokenizer.Tokenize(input)
		require.NotNil(t, err, input)
		assert.Equal(t, expected, err.Error(), "This input gives incorrect error message: "+input)
	}
}
