package main

import (
	"fmt"
	"unicode"
)

func IsPalindrome(s string) bool {
	runes := make([]rune, 0, len(s))
	for _, r := range s {
		if unicode.IsLetter(r) {
			runes = append(runes, unicode.ToLower(r))
		}
	}

	i, j := 0, len(runes)-1
	for i < j {
		if runes[i] != runes[j] {
			return false
		}
		i++
		j--
	}
	return true
}

func main() {
	tests := []struct {
		s    string
		want bool
	}{
		{"шалаш", true},
		{"А роза упала на лапу Азора", true},
		{"Was it a car or a cat I saw?", true},
		{"hello", false},
	}

	for _, t := range tests {
		got := IsPalindrome(t.s)
		fmt.Printf("IsPalindrome(%q) = %v (ожидалось %v)\n", t.s, got, t.want)
	}
}
