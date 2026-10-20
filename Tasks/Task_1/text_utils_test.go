package task1

import (
	"reflect"
	"testing"
)

func TestWordFrequency(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  map[string]int
	}{
		{"empty string", "", map[string]int{}},
		{"punctuation only", "!!! ... ,,, ???", map[string]int{}},
		{"simple sentence", "the cat and the hat", map[string]int{"the": 2, "cat": 1, "and": 1, "hat": 1}},
		{"mixed case", "Go go GO gO", map[string]int{"go": 4}},
		{"punctuation ignored", "Hello, world! Hello.", map[string]int{"hello": 2, "world": 1}},
		{"numbers count as words", "route 66 and route 66", map[string]int{"route": 2, "66": 2, "and": 1}},
		{"extra whitespace", "  a   b\ta\nb  ", map[string]int{"a": 2, "b": 2}},
		{"amharic words", "ሰላም ሰላም ዓለም።", map[string]int{"ሰላም": 2, "ዓለም": 1}},
		{"unicode mixed case", "Éa éA ÉA", map[string]int{"éa": 3}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := WordFrequency(tc.input)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("WordFrequency(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty string", "", true},
		{"punctuation only", "?!., ", true},
		{"single character", "a", true},
		{"simple palindrome", "racecar", true},
		{"not a palindrome", "hello", false},
		{"mixed case", "RaceCar", true},
		{"spaces and punctuation", "A man, a plan, a canal: Panama!", true},
		{"almost palindrome", "A man, a plan, a canal: Panamo", false},
		{"digits", "12321", true},
		{"digits not palindrome", "12345", false},
		{"amharic palindrome", "ሰበሰ", true},
		{"amharic not palindrome", "ሰላም", false},
		{"unicode accents", "Été", true},
		{"even length", "abba", true},
		{"two different chars", "ab", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsPalindrome(tc.input); got != tc.want {
				t.Errorf("IsPalindrome(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}
