package main

import (
	"fmt"
	"strings"
)

func WordCount(s string) map[string]int {
	// strings.Fields разбивает строку по пробелам,
	// игнорируя множественные пробелы и пробелы по краям
	words := strings.Fields(s)

	// Создаем map для хранения результатов
	counts := make(map[string]int)

	for _, word := range words {
		// Увеличиваем счетчик для каждого слова
		counts[word]++
	}

	return counts
}

func main() {
	fmt.Println(WordCount("go is fun go"))
	fmt.Println(WordCount("Hello hello World"))
	fmt.Println(WordCount(" мир мир труд "))
}
