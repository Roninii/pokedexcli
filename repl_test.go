package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "HeLlO WoRlD",
			expected: []string{"hello", "world"},
		},
		{
			input:    "world",
			expected: []string{"world"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("Expected length of %d but got %d", len(c.expected), len(actual))
			continue
		}
		for i := range actual {
			actualWord := actual[i]
			expectedWord := c.expected[i]

			if actualWord != expectedWord {
				t.Errorf("Expected %s but got %s", expectedWord, actualWord)
			}
		}
	}
}
