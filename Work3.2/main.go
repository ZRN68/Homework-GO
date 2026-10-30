package main

import (
	"fmt"
	"strings"
)

// normalize приводит к нижнему регистру и убирает все пробелы.
func normalize(s string) []rune {
	s = strings.ToLower(s)                  // регистр игнорируем
	s = strings.Join(strings.Fields(s), "") // Fields режет по любым пробелам, Join склеивает
	return []rune(s)                        // работаем с рунами, а не байтами
}

func AreAnagrams(s1, s2 string) bool {
	r1 := normalize(s1)
	r2 := normalize(s2)

	if len(r1) != len(r2) {
		return false
	}

	counts := make(map[rune]int, len(r1))
	for _, r := range r1 {
		counts[r]++
	}
	for _, r := range r2 {
		counts[r]--
		if counts[r] < 0 {
			return false
		}
	}
	return true
}

func main() {
	tests := []struct {
		a, b string
		want bool
	}{
		{"listen", "silent", true},
		{"Eleven plus two", "Twelve plus one", true},
		{"Апельсин", "спаниель", true},
		{"Hello", "World", false},
	}

	for _, t := range tests {
		got := AreAnagrams(t.a, t.b)
		fmt.Printf("AreAnagrams(%q, %q) = %v (ожидалось %v)\n", t.a, t.b, got, t.want)
	}
}
