package main

import "fmt"

func insert(arr []int) {

	for i := 1; i < len(arr); i++ {
		insertValue := arr[i] //暂存元素
		j := i - 1
		for ; j >= 0 && insertValue < arr[j]; j-- {
			arr[j+1] = arr[j] //向后备份元素
		}
		arr[j+1] = insertValue
	}
}

func main() {
	s := []int{1, 2, 3, 4, 5, 6, 7, 0}
	insert(s)
	fmt.Println(s)
}
