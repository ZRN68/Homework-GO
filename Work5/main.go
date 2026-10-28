package main

import "fmt"

func CountBits(n uint8) int {
	count := 0
	for i := 0; i < 8; i++ {
		if n&1 == 1 {
			count++
		}
		n >>= 1
	}
	return count
}

func IsPowerOfTwo(n uint) bool {
	return n != 0 && (n&(n-1)) == 0
}

func main() {
	// 1. Проверка CountBits
	fmt.Println("CountBits(0) =", CountBits(0))     // 0
	fmt.Println("CountBits(5) =", CountBits(5))     // 2
	fmt.Println("CountBits(12) =", CountBits(12))   // 2
	fmt.Println("CountBits(255) =", CountBits(255)) // 8

	// 2. Проверка IsPowerOfTwo
	fmt.Println("IsPowerOfTwo(1) =", IsPowerOfTwo(1))   // true
	fmt.Println("IsPowerOfTwo(6) =", IsPowerOfTwo(6))   // false
	fmt.Println("IsPowerOfTwo(64) =", IsPowerOfTwo(64)) // true
	fmt.Println("IsPowerOfTwo(0) =", IsPowerOfTwo(0))   // false

	// 3. Демонстрация подсказки
	fmt.Println("\nПодсказка (n и n-1):")
	for _, n := range []uint{8, 6} {
		fmt.Printf("n = %d: %08b, n-1: %08b, результат n & (n-1): %08b\n", 
			n, n, n-1, n&(n-1))
	}
}