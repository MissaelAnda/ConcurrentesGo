package main

import (
	"fmt"
)

func main() {
	var farenheit float64

	fmt.Println("Transformar de fahrenheit a celsius")
	fmt.Print("Ingrese farenheits: ")
	fmt.Scanf("%f", &farenheit)

	salida := (farenheit - 32.0) * (5.0 / 9.0)

	fmt.Println("Celsius = ", salida)
}
