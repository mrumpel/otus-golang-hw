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

	if unicode.IsDigit([]rune(str)[0]) {
		return "", ErrInvalidString
	}

	var prev rune
	var count int
	var err error
	unpacked := true

	res := strings.Builder{}

	for _, curr := range str {
		if unicode.IsDigit(curr) {

			if unicode.IsDigit(prev) {
				return "", ErrInvalidString
			}

			count, err = strconv.Atoi(string(curr))
			if err != nil {
				return "", err
			}

			unpacked = true

			if count != 0 {
				res.WriteString(strings.Repeat(string(prev), count))
			}
		} else {
			if !unpacked {
				res.WriteRune(prev)
			}
			unpacked = false
		}

		prev = curr
		count = 1
	}

	// for last unpair item
	if !unpacked && count == 1 {
		res.WriteRune(prev)
	}

	return res.String(), nil
}
