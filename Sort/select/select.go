package main

import "fmt"

func selectsort(arr []int) {
	if len(arr) == 0 {
		return
	}
	n := len(arr)
	for i := 0; i < n-1; i++ {
		minIndex := i
		for j := i + 1; j < n; j++ {
			if arr[j] < arr[minIndex] {
				minIndex = j
			}
		}
		if i != minIndex {
			arr[i], arr[minIndex] = arr[minIndex], arr[i]
		}
	}
}
func main() {
	s := []int{1, 2, 3, 4, 5, 6, 7, 0}
	selectsort(s)
	fmt.Println(s)
}
