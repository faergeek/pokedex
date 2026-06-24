package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "goodbye mars",
			expected: []string{"goodbye", "mars"},
		},
		{
			input:    "justasingleword",
			expected: []string{"justasingleword"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "    awordsurroundedbyspaces      ",
			expected: []string{"awordsurroundedbyspaces"},
		},
		{
			input:    "    \ttabbing  \t tabbing\nand newlining,\there and\rthere, here\nand there",
			expected: []string{"tabbing", "tabbing", "and", "newlining,", "here", "and", "there,", "here", "and", "there"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("result length doesn't match. Expected: %d, got: %d", len(c.expected), len(actual))
		}

		for i := range actual {
			word := actual[i]
			expectedword := c.expected[i]

			if word != expectedword {
				t.Errorf("word %d does not match. Expected: %s, got: %s", i, expectedword, word)
			}
		}
	}
}
