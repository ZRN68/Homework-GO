package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	k := flag.Int("k", 5, "максимальное число одновременных запросов")
	q := flag.String("q", "Go", "искомая строка (с учётом регистра)")
	flag.Parse()

	if *k <= 0 {
		fmt.Fprintln(os.Stderr, "k должно быть > 0")
		os.Exit(1)
	}

	client := &http.Client{}
	sem := make(chan struct{}, *k) // семафор на k слотов
	var wg sync.WaitGroup

	// Мьютекс и общая сумма — защищаем от гонок.
	var mu sync.Mutex
	total := 0

	scanner := bufio.NewScanner(os.Stdin)
	// Увеличим буфер, чтобы не падать на длинных строках.
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url == "" {
			continue
		}

		// Захватываем слот ДО запуска горутины — это ограничивает
		// одновременно работающих обработчиков ровно до k.
		// Пока ввод читается быстрее, чем обрабатывается, мы блокируемся здесь,
		// но при этом горутин не больше k — простоя нет.
		sem <- struct{}{}

		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			defer func() { <-sem }() // освобождаем слот

			count, err := countInURL(client, u, *q)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error for %s: %v\n", u, err)
				return
			}

			fmt.Printf("Count for %s: %d\n", u, count)

			mu.Lock()
			total += count
			mu.Unlock()
		}(url)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "stdin error: %v\n", err)
	}

	wg.Wait()
	fmt.Printf("Total: %d\n", total)
}

// countInURL делает GET-запрос и считает число неперекрывающихся вхождений q в тело.
func countInURL(client *http.Client, url, q string) (int, error) {
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("bad status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	return strings.Count(string(body), q), nil
}