package main

import "fmt"

type CustomerStatus int

const (
	Regular CustomerStatus = iota
	Silver
	Gold
)

const (
	deliveryPrice    = 299.0  // цена доставки
	freeDeliveryFrom = 3000.0 // порог для бесплатной доставки
)

// statusName возвращает имя статуса для вывода
func statusName(s CustomerStatus) string {
	switch s {
	case Regular:
		return "Regular"
	case Silver:
		return "Silver"
	case Gold:
		return "Gold"
	default:
		return "Unknown"
	}
}

// process считает и печатает скидку, доставку и итог
func process(sum float64, status CustomerStatus) {
	// скидка по статусу
	var discountPercent float64
	switch status {
	case Regular:
		discountPercent = 0
	case Silver:
		discountPercent = 5
	case Gold:
		discountPercent = 10
	default:
		fmt.Println("неизвестный статус")
		discountPercent = 0
	}

	discount := sum * discountPercent / 100
	afterDiscount := sum - discount

	// доставка: бесплатно, если сумма после скидки от 3000
	var delivery float64
	if afterDiscount >= freeDeliveryFrom {
		delivery = 0
	} else {
		delivery = deliveryPrice
	}

	total := afterDiscount + delivery

	// вывод
	fmt.Printf("Сумма: %.2f, статус: %s\n", sum, statusName(status))
	fmt.Printf("Скидка: %.2f\n", discount)
	fmt.Printf("Доставка: %.2f\n", delivery)
	fmt.Printf("Итого: %.2f\n\n", total)
}

func main() {
	// Проверяем все примеры из условия
	process(4500.00, Gold)
	process(2900.00, Silver)
	process(1000.00, Regular)
}
