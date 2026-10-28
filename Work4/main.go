package main

import (
	"fmt"
	"strings"
)

const (
	SecondsInMinute = 60
	SecondsInHour   = 60 * SecondsInMinute
	SecondsInDay    = 24 * SecondsInHour
)

func uptime(total int) string {
	days := total / SecondsInDay
	total = total % SecondsInDay

	hours := total / SecondsInHour
	total = total % SecondsInHour

	minutes := total / SecondsInMinute
	seconds := total % SecondsInMinute

	// Используем слайс для накопления частей строки
	var parts []string

	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dд", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dч", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dм", minutes))
	}
	if seconds > 0 {
		parts = append(parts, fmt.Sprintf("%dс", seconds))
	}
	if len(parts) == 0 {
		return "0с"
	}
	return strings.Join(parts, " ")
}

func main() {
	fmt.Printf("%d -> %s\n", 93784, uptime(93784))
	fmt.Printf("%d -> %s\n", 86405, uptime(86405))
	fmt.Printf("%d -> %s\n", 3600, uptime(3600))
	fmt.Printf("%d -> %s\n", 59, uptime(59))
	fmt.Printf("%d -> %s\n", 0, uptime(0))
}
