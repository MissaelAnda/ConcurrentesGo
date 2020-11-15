package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

type Data struct {
	Id    int
	Count uint64
}

func exit(exit chan bool) {
	var input string
	fmt.Scanln(&input)

	exit <- true
}

func main() {
	exitChannel := make(chan bool)
	conn, err := net.Dial("tcp", ":8080")

	if err != nil {
		fmt.Println(err)
		return
	}

	go exit(exitChannel)

	var data Data
	err = gob.NewDecoder(conn).Decode(&data)

	if err != nil {
		fmt.Println(err)
		return
	}

	if data.Id < 0 {
		fmt.Println("No hay mas procesos disponibles")
		return
	}

	for {
		select {
		case <-exitChannel:
			err := gob.NewEncoder(conn).Encode(data)
			if err != nil {
				fmt.Println("Error al enviar el proceso al servidor: ", err)
			}
			fmt.Println("Exit")
			return

		default:
			fmt.Println(data.Id, ":", data.Count)
			data.Count++
			time.Sleep(time.Millisecond * 500)
		}
	}
}
