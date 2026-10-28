package main

import (
	"fmt"
	"time"
)

// Функция move — это коллбэк (вызывающий код), который решает, что делать с каждым ходом.
func Hanoi(n int, from, to, via string, move func(from, to string)) {
	if n == 0 {
		return
	}
	// 1. Переносим n-1 дисков с from на via (используя to как вспомогательный)
	Hanoi(n-1, from, via, to, move)

	// 2. Переносим самый большой диск с from на to
	move(from, to)

	// 3. Переносим n-1 дисков с via на to (используя from как вспомогательный)
	Hanoi(n-1, via, to, from, move)
}

func main() {
	// --- ЧАСТЬ 1: Для n = 3 ---
	fmt.Println("1. Ходы для n = 3:")
	step := 0 // Локальная переменная, которая будет захвачена замыканием
	Hanoi(3, "A", "C", "B", func(from, to string) {
		step++
		fmt.Printf("%d: %s -> %s\n", step, from, to)
	})

	// --- ЧАСТЬ 2: Для n = 20 (только подсчет и время) ---
	fmt.Println("\n2. Для n = 20:")

	var movesCount int

	// Замер времени через defer
	start := time.Now()
	defer func() {
		// Этот блок выполнится в самом конце функции main
		elapsed := time.Since(start)
		fmt.Printf("n=20: %d ходов\n", movesCount)
		fmt.Printf("время: %v\n", elapsed)
	}()

	// Запускаем алгоритм для 20 дисков, просто считая ходы
	Hanoi(20, "A", "C", "B", func(from, to string) {
		movesCount++
	})
}
