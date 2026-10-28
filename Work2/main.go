package main

import (
	"fmt"
)

// kopecksInRuble — сколько копеек в одном рубле.
// Баланс хранится в копейках, поэтому для вывода в рублях делим на эту константу.
const kopecksInRuble = 100

// NewAccount создаёт счёт с балансом в копейках и возвращает три замыкания,
// которые работают с одним и тем же балансом.
func NewAccount(initial int64) (deposit func(int64), withdraw func(int64) error, balance func() int64) {
	bal := initial // общая переменная для всех трёх замыканий

	deposit = func(amount int64) {
		bal += amount
	}

	withdraw = func(amount int64) error {
		if amount > bal {
			// Переиспользуем formatMoney — не дублируем перевод копеек в рубли.
			return fmt.Errorf("недостаточно средств (баланс %s)", formatMoney(bal))
		}
		bal -= amount
		return nil
	}

	balance = func() int64 {
		return bal
	}

	return deposit, withdraw, balance
}

// formatMoney переводит копейки в строку рублей с двумя знаками.
func formatMoney(kopecks int64) string {
	return fmt.Sprintf("%.2f", float64(kopecks)/kopecksInRuble)
}

func main() {
	depositA, withdrawA, balanceA := NewAccount(100000) // 1000.00 ₽
	depositB, _, balanceB := NewAccount(0)              // 0.00 ₽

	fmt.Printf("Счёт A: %s\n", formatMoney(balanceA()))

	depositA(50000) // +500.00
	fmt.Printf("Пополнение 500.00 -> %s\n", formatMoney(balanceA()))

	// Снятие 2000.00 — не хватает: печатаем ошибку.
	err := withdrawA(200000)
	if err != nil {
		fmt.Printf("Снятие 2000.00: ошибка: %v\n", err)
	} else {
		fmt.Printf("Снятие 2000.00 -> %s\n", formatMoney(balanceA()))
	}

	// Снятие 300.50 — успех: печатаем новый баланс.
	err = withdrawA(30050)
	if err != nil {
		fmt.Printf("Снятие 300.50: ошибка: %v\n", err)
	} else {
		fmt.Printf("Снятие 300.50 -> %s\n", formatMoney(balanceA()))
	}

	// Счёт B независим от A
	fmt.Printf("Счёт B: %s\n", formatMoney(balanceB()))

	// Для наглядности: пополним B и проверим, что A не изменился.
	depositB(10000) // +100.00
	fmt.Printf("Пополнение B на 100.00 -> %s\n", formatMoney(balanceB()))
	fmt.Printf("Счёт A после этого: %s (не изменился)\n", formatMoney(balanceA()))
}
