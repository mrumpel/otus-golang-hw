package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if len(str) == 0 {
		return "", nil
	}

	var prev rune
	var count int
	var err error
	pass, escape := true, false
	res := strings.Builder{}

	for _, curr := range str {
		// Asterisk escape logic

		if escape {
			if !(unicode.IsDigit(curr) || curr == '\\') {
				return "", ErrInvalidString
			}

			prev = curr
			pass, escape = false, false

			continue
		}

		if curr == '\\' {
			escape = true
		}

		// Main work logic
		if pass && !escape {
			if unicode.IsDigit(curr) {
				return "", ErrInvalidString
			}

			prev = curr
			pass = false

			continue
		}

		if unicode.IsDigit(curr) {
			count, err = strconv.Atoi(string(curr))
			if err != nil {
				return "", err
			}

			pass = true
		} else {
			count = 1
			pass = false
		}

		res.WriteString(strings.Repeat(string(prev), count))

		prev = curr
	}

	// Last item: check escape error & add
	if escape {
		return "", ErrInvalidString
	}

	if !pass {
		res.WriteRune(prev)
	}

	return res.String(), nil
}
