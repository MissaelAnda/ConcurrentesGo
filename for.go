package main

import (
	"fmt"
)

func main() {
	var presicion int
	var e float32 = 0.0

	fmt.Println("Calcular constante de euler")
	fmt.Print("Ingrese presicion: ")
	fmt.Scan(&presicion)

	for i := 0; i <= presicion; i++ {
		e += 1.0 / float32(factorial(uint(i)))
	}

	fmt.Println("e = ", e)
}

func factorial(n uint) uint {
	if n == 0 {
		return 1
	}
	var f uint = 1
	var i uint
	for i = 1; i <= n; i++ {
		f *= i
	}
	return f
}
