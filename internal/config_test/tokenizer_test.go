package config

import (
	"testing"

	"github.com/samuelko123/git-cat/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	testcases := map[string][]config.Token{
		"[remote]": {
			{Type: config.SECTION, Value: "remote"},
		},
		"[ReMote]": {
			{Type: config.SECTION, Value: "remote"},
		},
		" \t [remote] \t ": {
			{Type: config.SECTION, Value: "remote"},
		},
		"[ remote ]":   nil,
		"[\tremote\t]": nil,
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
		"[ remote \"origin\" ]":   nil,
		"[\tremote \"origin\"\t]": nil,
		"ignorecase=true": {
			{Type: config.NAME, Value: "ignorecase"},
			{Type: config.VALUE, Value: "true"},
		},
		" \t ignorecase \t = \t true \t ": {
			{Type: config.NAME, Value: "ignorecase"},
			{Type: config.VALUE, Value: "true"},
		},
	}

	for input, expected := range testcases {
		tokenizer := config.Tokenizer{}

		tokens, err := tokenizer.Tokenize(input)
		if expected == nil {
			assert.NotNil(t, err, "This input should fail: "+input)
		} else {
			assert.Equal(t, expected, tokens, "This input should not fail: "+input)
		}
	}
}
