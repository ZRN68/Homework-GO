package main

import "fmt"

// isLetter сообщает, является ли руна буквой латиницы или кириллицы.
func isLetter(r rune) bool {
	switch {
	case r >= 'a' && r <= 'z':
		return true
	case r >= 'A' && r <= 'Z':
		return true
	case r >= 'а' && r <= 'я':
		return true
	case r >= 'А' && r <= 'Я':
		return true
	case r == 'ё' || r == 'Ё':
		return true
	}
	return false
}

// toLower приводит букву к нижнему регистру (латиница и кириллица).
func toLower(r rune) rune {
	switch {
	case r >= 'A' && r <= 'Z':
		return r + ('a' - 'A')
	case r >= 'А' && r <= 'Я':
		return r + ('а' - 'А')
	case r == 'Ё':
		return 'ё'
	}
	return r
}

// IsPalindrome проверяет, читается ли строка одинаково в обе стороны.
// Регистр и все небуквенные символы игнорируются.
func IsPalindrome(s string) bool {
	runes := make([]rune, 0, len(s))
	for _, r := range s {
		if isLetter(r) {
			runes = append(runes, toLower(r))
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
