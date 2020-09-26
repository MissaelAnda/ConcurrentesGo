package main

import (
	"fmt"
)

func main() {
	var mes int
	var dia int
	var signo string

	fmt.Println("Signo zodiacal")
	fmt.Print("Ingrese mes: ")
	fmt.Scan(&mes)

	if mes < 1 || mes > 12 {
		println("Mes incorrecto")
		return
	}

	fmt.Print("Ingrese dia: ")
	fmt.Scan(&dia)

	if dia < 1 {
		println("Dia invalido")
		return
	}

	switch mes {
	case 1:
		if dia > 31 {
			println("Dia invalido.")
			return
		}
		if dia <= 20 {
			signo = "Capricornio"
		} else {
			signo = "Acuario"
		}
		break

	case 2:
		if dia > 29 {
			println("Dia invalido.")
			return
		}
		if dia <= 19 {
			signo = "Acuario"
		} else {
			signo = "Piscis"
		}
		break

	case 3:
		if dia > 31 {
			println("Dia invalido.")
			return
		}
		if dia <= 20 {
			signo = "Piscis"
		} else {
			signo = "Aries"
		}
		break

	case 4:
		if dia > 30 {
			println("Dia invalido.")
			return
		}
		if dia <= 20 {
			signo = "Aries"
		} else {
			signo = "Tauro"
		}
		break

	case 5:
		if dia > 31 {
			println("Dia invalido.")
			return
		}
		if dia <= 21 {
			signo = "Tauro"
		} else {
			signo = "Geminis"
		}
		break

	case 6:
		if dia > 30 {
			println("Dia invalido.")
			return
		}
		if dia <= 21 {
			signo = "Geminis"
		} else {
			signo = "Cancer"
		}
		break

	case 7:
		if dia > 31 {
			println("Dia invalido.")
			return
		}
		if dia <= 23 {
			signo = "Cancer"
		} else {
			signo = "Leo"
		}
		break

	case 8:
		if dia > 31 {
			println("Dia invalido.")
			return
		}
		if dia <= 23 {
			signo = "Leo"
		} else {
			signo = "Virgo"
		}
		break

	case 9:
		if dia > 30 {
			println("Dia invalido.")
			return
		}
		if dia <= 23 {
			signo = "Virgo"
		} else {
			signo = "Libra"
		}
		break

	case 10:
		if dia > 31 {
			println("Dia invalido.")
			return
		}
		if dia <= 23 {
			signo = "Libra"
		} else {
			signo = "Escorpio"
		}
		break

	case 11:
		if dia > 30 {
			println("Dia invalido.")
			return
		}
		if dia <= 22 {
			signo = "Escorpio"
		} else {
			signo = "Sagitario"
		}
		break

	case 12:
		if dia > 31 {
			println("Dia invalido.")
			return
		}
		if dia <= 21 {
			signo = "Sagitario"
		} else {
			signo = "Capricornio"
		}
		break
	}

	fmt.Println("Tu signo es: ", signo)
}
