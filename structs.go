package main

import "fmt"

type Imagen struct {
	titulo  string
	formato string
	canales uint
}

type Audio struct {
	titulo   string
	formato  string
	duracion uint
}

type Video struct {
	titulo  string
	formato string
	frames  uint
}

type Multimedia interface {
	mostrar()
}

func (i Imagen) mostrar() {
	fmt.Println(i.titulo)
	fmt.Println(i.formato)
	fmt.Println(i.canales)
}

func (a Audio) mostrar() {
	fmt.Println(a.titulo)
	fmt.Println(a.formato)
	fmt.Println(a.duracion)
}

func (v Video) mostrar() {
	fmt.Println(v.titulo)
	fmt.Println(v.formato)
	fmt.Println(v.frames)
}

type ContenidoWeb struct {
	multimedia []Multimedia
}

func (cw *ContenidoWeb) push(m Multimedia) {
	cw.multimedia = append(cw.multimedia, m)
}

func (cw *ContenidoWeb) mostrar() {
	for _, m := range cw.multimedia {
		m.mostrar()
		fmt.Println("---------------------")
	}
}

func main() {
	var menu uint
	var string1 string
	var string2 string
	var number uint

	var multi ContenidoWeb
	end := false

	for !end {
		fmt.Println("Menu:")
		fmt.Println("1. Crear Imagen")
		fmt.Println("2. Crear Audio")
		fmt.Println("3. Crear Video")
		fmt.Println("4. Mostrar")
		fmt.Println("0. Salir")
		fmt.Scan(&menu)

		switch menu {
		case 1:
			fmt.Print("Titulo: ")
			fmt.Scan(&string1)
			fmt.Print("Formato: ")
			fmt.Scan(&string2)
			fmt.Print("Canales: ")
			fmt.Scan(&number)

			multi.push(Imagen{string1, string2, number})
			break

		case 2:
			fmt.Print("Titulo: ")
			fmt.Scan(&string1)
			fmt.Print("Formato: ")
			fmt.Scan(&string2)
			fmt.Print("Duracion: ")
			fmt.Scan(&number)

			multi.push(Audio{string1, string2, number})
			break

		case 3:
			fmt.Print("Titulo: ")
			fmt.Scan(&string1)
			fmt.Print("Formato: ")
			fmt.Scan(&string2)
			fmt.Print("Frames: ")
			fmt.Scan(&number)

			multi.push(Video{string1, string2, number})
			break

		case 4:
			fmt.Println()
			multi.mostrar()
			fmt.Println()
			break

		case 0:
		default:
			end = true
		}
	}
}
