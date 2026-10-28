package main

import (
	"errors"
	"fmt"
)

// Operation — тип функции-операции: принимает два числа, возвращает результат и ошибку.
type Operation func(a, b float64) (float64, error)

// operations — map функций. Ключ — символ операции, значение — сама операция.
var operations = map[string]Operation{
	"+": func(a, b float64) (float64, error) {
		return a + b, nil
	},
	"-": func(a, b float64) (float64, error) {
		return a - b, nil
	},
	"*": func(a, b float64) (float64, error) {
		return a * b, nil
	},
	"/": func(a, b float64) (float64, error) {
		if b == 0 {
			return 0, errors.New("деление на ноль")
		}
		return a / b, nil
	},
}

// Calculate находит операцию в map и вызывает её.
// Если операции нет — возвращает ошибку, не вызывая функцию.
func Calculate(a float64, op string, b float64) (float64, error) {
	fn, ok := operations[op]
	if !ok {
		return 0, fmt.Errorf("неизвестная операция %q", op)
	}
	return fn(a, b)
}

// printResult — вспомогательная функция для красивого вывода.
func printResult(a float64, op string, b float64) {
	res, err := Calculate(a, op, b)
	if err != nil {
		// Например: 1 / 0: деление на ноль
		fmt.Printf("%v %s %v: %v\n", a, op, b, err)
		return
	}
	// Например: 10 + 5 = 15.00
	fmt.Printf("%v %s %v = %.2f\n", a, op, b, res)
}

func main() {
	printResult(10, "+", 5) // 10 + 5 = 15.00
	printResult(10, "/", 4) // 10 / 4 = 2.50
	printResult(7, "*", 6)  // 7 * 6 = 42.00
	printResult(1, "/", 0)  // 1 / 0: деление на ноль
	printResult(2, "%", 3)  // 2 % 3: неизвестная операция "%"
}