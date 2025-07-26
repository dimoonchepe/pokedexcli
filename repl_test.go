package main

import (
    "testing"
)

func TestCleanInput(t *testing.T) {
    cases := []struct {
        input       string
        expected    []string
    }{
        {
            input: "  hello world  ",
            expected: []string {"hello", "world"},
        },
        {
            input: "  Battlestar GALACTICA  ",
            expected: []string {"battlestar", "galactica"},
        },
        {
            input: "",
            expected: []string {},
        },
    }

    for _, c := range cases {
        actual := cleanInput(c.input)
        if len(actual) != len(c.expected) {
            t.Errorf("Length of output is %d, want %d", len(actual), len(c.expected))
        }

        for i := range actual {
            word := actual[i]
            expectedWord := c.expected[i]
            if word != expectedWord {
                t.Errorf("Got word '%s', want '%s'", word, expectedWord)
            }
        }
    }
}
