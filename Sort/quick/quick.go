package main

import "fmt"

func quick(arr []int, start, end int) {
	if start > end {
		return
	}
	pivoIndex := partion(arr, start, end)
	quick(arr, start, pivoIndex-1)
	quick(arr, pivoIndex+1, end)
}
func partion(arr []int, start, end int) int {
	//选第一个元素为基准元素
	pivot := arr[start]
	left := start
	right := end
	index := start
	for right > left {
		for right > left {
			if arr[right] < pivot {
				arr[left] = arr[right]
				index = right
				left++
				break
			}
			right--
		}
		for right > left {
			if arr[left] > pivot {
				arr[right] = arr[left]
				index = left
				right--
				break
			}
			left++
		}

	}
	arr[index] = pivot
	return index
}
func main() {
	s := []int{1, 2, 3, 4, 5, 6, 7, 0}
	quick(s, 0, len(s)-1)
	fmt.Println(s)
}
