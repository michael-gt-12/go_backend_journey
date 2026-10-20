// Package task1 contains text-processing helpers: word frequency counting
// and palindrome checking. Both work on runes, so non-ASCII text
// (for example Amharic) is handled correctly.
package task1

import (
	"strings"
	"unicode"
)

// isWordRune reports whether r is part of a word (a letter or a digit).
func isWordRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r)
}

// WordFrequency counts how many times each word appears in s.
// Counting is case-insensitive and punctuation is ignored,
// so "Go, go!" yields {"go": 2}.
func WordFrequency(s string) map[string]int {
	freq := make(map[string]int)

	// FieldsFunc splits on every rune for which the function returns true.
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !isWordRune(r)
	})

	for _, w := range words {
		freq[strings.ToLower(w)]++
	}
	return freq
}

// IsPalindrome reports whether s reads the same forwards and backwards,
// ignoring spaces, punctuation and letter case.
// An empty string (or one with no letters/digits) counts as a palindrome.
func IsPalindrome(s string) bool {
	cleaned := make([]rune, 0, len(s))
	for _, r := range s {
		if isWordRune(r) {
			cleaned = append(cleaned, unicode.ToLower(r))
		}
	}

	// Compare from both ends moving toward the middle.
	for i, j := 0, len(cleaned)-1; i < j; i, j = i+1, j-1 {
		if cleaned[i] != cleaned[j] {
			return false
		}
	}
	return true
}
