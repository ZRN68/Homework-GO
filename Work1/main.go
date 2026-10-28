package main

import "fmt"

const N = 8

func draw(n int) {
	for row := 0; row < n; row++ {
		for col := 0; col < n; col++ {
			if (row+col)%2 == 1 {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func main() {
	// Основное значение из условия
	fmt.Printf("N = %d:\n", N)
	draw(N)

	// Проверка остальных значений из условия
	fmt.Println("N=5")
	draw(5)

	fmt.Println("N=1")
	draw(1)
}
