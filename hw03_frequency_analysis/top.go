package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
)

var re *regexp.Regexp

func init() {
	re = regexp.MustCompilePOSIX(`[!"#$%&'()*+,./:;<=>?@\[\\\]^_‘{|}~]|[[:punct:][:space:]]-|-[[:punct:][:space:]$]|^-|-$`)
}

func Top10(source string) []string {
	// Clear text for asterisk task
	source = re.ReplaceAllString(source, " ")
	source = strings.ToLower(source)

	// Getting stats
	text := strings.FieldsFunc(source, unicode.IsSpace)

	if len(text) == 0 {
		return nil
	}

	ratemap := make(map[string]int)

	for _, word := range text {
		ratemap[word]++
	}

	// Sortable entity
	type wordcount struct {
		word  string
		count int
	}

	rate := make([]wordcount, 0, len(ratemap))

	for word, count := range ratemap {
		rate = append(rate, wordcount{word: word, count: count})
	}

	sort.Slice(rate, func(i, j int) bool {
		if rate[i].count == rate[j].count {
			return rate[i].word < rate[j].word
		}

		return rate[i].count > rate[j].count
	})

	// Result prep.
	reslen := 10

	if len(rate) < 10 {
		reslen = len(rate)
	}

	res := make([]string, 0, reslen)

	for i := 0; i < reslen; i++ {
		res = append(res, rate[i].word)
	}

	return res
}
