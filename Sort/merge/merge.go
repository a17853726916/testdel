package main

import "fmt"

func Mergesort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}
	mid := len(arr) / 2
	left := Mergesort(arr[:mid])
	right := Mergesort(arr[mid:])
	return Merge(left, right)
}

func Merge(left, right []int) []int {
	p1, p2 := 0, 0
	var res []int
	for p1 < len(left) && p2 < len(right) {
		if left[p1] < right[p2] {
			res = append(res, left[p1])
			p1++
		} else {
			res = append(res, right[p2])
			p2++
		}
	}
	res = append(res, left[p1:]...)
	res = append(res, right[p2:]...)
	return res
}

func main() {
	s := []int{1, 2, 3, 4, 5, 6, 7, 0}
	fmt.Println(Mergesort(s))
}
