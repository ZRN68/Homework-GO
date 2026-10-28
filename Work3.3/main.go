package main

import "fmt"

// FirstUnique возвращает первый символ, встречающийся в строке ровно один раз.
// Если такого нет — возвращает 0 (нулевую руну).
// Поиск регистрозависимый.
func FirstUnique(s string) rune {
	// Первый проход: считаем частоты символов.
	counts := make(map[rune]int)
	for _, r := range s {
		counts[r]++
	}

	// Второй проход: ищем первый символ с частотой 1 в исходном порядке.
	for _, r := range s {
		if counts[r] == 1 {
			return r
		}
	}

	// Уникальных символов нет.
	return 0
}

func main() {
	tests := []struct {
		s    string
		want rune
	}{
		{"abacabad", 'c'},
		{"abcabcd", 'd'},
		{"абракадабра", 'к'},
		{"aabbcc", 0},
	}

	for _, t := range tests {
		got := FirstUnique(t.s)
		fmt.Printf("FirstUnique(%q) = %q (ожидалось %q)\n", t.s, got, t.want)
	}
}
