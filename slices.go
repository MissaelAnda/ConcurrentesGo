package main

import (
	"fmt"
)

func main() {
	var size uint
	var helper int

	fmt.Scan(&size)

	var arr []int

	var sum int = 0

	for i := 0; i < int(size); i++ {
		fmt.Scan(&helper)
		arr = append(arr, helper)
		sum += helper
	}

	fmt.Println(sum)
}
