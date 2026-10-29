package main

import "fmt"

// RemoveElement возвращает новый слайс без элемента с индексом index.
// Порядок остальных элементов сохраняется. Исходный слайс не меняется.
// При недопустимом индексе возвращает nil и ошибку.
func RemoveElement(nums []int, index int) ([]int, error) {
	if index < 0 || index >= len(nums) {
		return nil, fmt.Errorf("индекс %d вне диапазона [0, %d)", index, len(nums))
	}

	result := make([]int, 0, len(nums)-1)
	result = append(result, nums[:index]...)
	result = append(result, nums[index+1:]...)

	return result, nil
}

func main() {
	nums := []int{1, 2, 2, 3, 1}

	res, err := RemoveElement(nums, 3)
	fmt.Println(res, err) // [1 2 2 1] <nil>
	fmt.Println(nums)     // [1 2 2 3 1] — исходный не изменился

	res, err = RemoveElement([]int{1, 2, 3}, 5)
	fmt.Println(res, err) // [] индекс 5 вне диапазона [0, 3)

	res, err = RemoveElement([]int{1, 2, 3}, -1)
	fmt.Println(res, err) // [] индекс -1 вне диапазона [0, 3)
}


//append не создаёт копию среза, а дописывает элементы в свободную ёмкость того же массива
