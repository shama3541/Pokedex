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
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "   ",
			expected: []string{},
		},
		{
			input:    "Golang\tIs\nGreat",
			expected: []string{"golang", "is", "great"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		// Check length
		if len(actual) != len(c.expected) {
			t.Errorf("cleanInput(%q) returned %d words, expected %d\nGot: %v\nWant: %v",
				c.input, len(actual), len(c.expected), actual, c.expected)
			continue
		}

		// Check each word
		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("cleanInput(%q) word mismatch at index %d: got %q, expected %q",
					c.input, i, actual[i], c.expected[i])
			}
		}
	}
}
