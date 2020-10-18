package main

import "fmt"

func main() {
	fmt.Println(fibonacci(10))

	arr := []int{1, 4, 2, 5, 12, 3, 64, 2}
	fmt.Println(max(arr...))

	impares := generadorImpares()
	fmt.Println(impares())
	fmt.Println(impares())
	fmt.Println(impares())
	fmt.Println(impares())

	var a int = 15
	var b int = 51
	fmt.Println(a, b)
	intercambia(&a, &b)
	fmt.Println(a, b)
}

func fibonacci(cont int) int {
	if cont == 0 || cont == 1 {
		return 1
	}
	return fibonacci(cont-1) + fibonacci(cont-2)
}

func max(args ...int) int {
	var max int = args[0]
	for i := 1; i < len(args); i++ {
		if args[i] > max {
			max = args[i]
		}
	}
	return max
}

func generadorImpares() func() uint {
	i := uint(1) // i permanecerá en el clousure de la función anónima a retornar
	return func() uint {
		var par = i
		i += 2
		return par
	}
}

func intercambia(a *int, b *int) {
	temp := *b
	*b = *a
	*a = temp
}
