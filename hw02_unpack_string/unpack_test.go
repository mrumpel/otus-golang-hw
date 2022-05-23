package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},

		// asterisk test
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		{input: `qwe\\`, expected: `qwe\`},

		// Zero check
		{input: "aabb0", expected: "aab"},
		{input: "a0ab", expected: "ab"},
		{input: "a0", expected: ""},

		// Corner check
		{input: "aacb2", expected: "aacbb"},
		{input: "a2cbb", expected: "aacbb"},
		{input: "a2cb2", expected: "aacbb"},

		// Unicode + variable rune len check
		{input: "去有趣", expected: "去有趣"},
		{input: "去有趣", expected: "去有趣"},
		{input: "去2有趣3", expected: "去去有趣趣趣"},
		{input: "去ё有й趣", expected: "去ё有й趣"},
		{input: "去ё2有3й趣", expected: "去ёё有有有й趣"},
		{input: "去f有й趣", expected: "去f有й趣"},
		{input: "去f2有3й2趣", expected: "去ff有有有йй趣"},

		// Spec. symbols check
		{input: "\tabc", expected: "\tabc"},
		{input: "\t3abc", expected: "\t\t\tabc"},
		{input: "d\nabc", expected: "d\nabc"},
		{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},
		{input: "a-3b", expected: "a---b"},
		{input: "ab 5", expected: "ab     "},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{
		"3abc",
		"45",
		"aaa10b",
		`aa\a`,
		`aaa\`,
	}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
