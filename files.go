package main

import (
	"fmt"
	"os"
	"sort"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var str string
	var opt string
	var strings []string
	end := false

	fmt.Println("Introduce las cadenas")
	for !end {
		fmt.Scan(&str)
		strings = append(strings, str)
		fmt.Print("Introducir otra cadena? (s/n): ")
		fmt.Scan(&opt)

		if opt == "n" {
			end = true
		}
	}

	f, err := os.Create("asecendente.txt")
	check(err)

	sort.Strings(strings)
	for _, s := range strings {
		_, err := f.WriteString(s)
		check(err)
		_, err2 := f.WriteString("\n")
		check(err2)
	}

	f.Close()

	f2, err := os.Create("descendente.txt")
	check(err)

	sort.Slice(strings, func(i, j int) bool {
		return strings[i][0] > strings[j][0]
	})

	for _, s := range strings {
		_, err := f2.WriteString(s)
		check(err)
		_, err2 := f2.WriteString("\n")
		check(err2)
	}

	f2.Close()
}
