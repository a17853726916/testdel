package main

import "fmt"

func bubble(arr []int) {
	if len(arr) == 0 {
		return
	}
	for i := 0; i < len(arr); i++ {
		for j := i; j < len(arr)-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}

func main() {
	arr := []int{1, 5, 4, 2, 6, 3}
	bubble(arr)
	fmt.Println(arr)
}
