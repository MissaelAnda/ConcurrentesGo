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

func process(close chan bool, resume chan Data, toClose chan Data) {
	active := [5]bool{true, true, true, true, true}
	processes := [5]uint64{0, 0, 0, 0, 0}

	for {
		select {
		case <-close:
			done := false
			for i := 0; i <= 4; i++ {
				if active[i] {
					toClose <- Data{
						Id:    i,
						Count: processes[i],
					}
					active[i] = false
					done = true
					break
				}
			}
			if !done {
				toClose <- Data{
					Id:    -1,
					Count: 0,
				}
			}
			break

		case resumeData := <-resume:
			if resumeData.Id < 0 || resumeData.Id > 4 {
				break
			}
			active[resumeData.Id] = true
			processes[resumeData.Id] = resumeData.Count
			break

		default:
		}

		fmt.Println("---------------")
		for i := 0; i <= 4; i++ {
			if active[i] {
				fmt.Println(i, ":", processes[i])
				processes[i]++
			}
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func manageClient(conn net.Conn, toClose chan Data, resume chan Data) {
	item := <-toClose
	err := gob.NewEncoder(conn).Encode(item)

	if err != nil {
		fmt.Println("Hubo un error al enviar informacion al cliente: ", err)
		resume <- item
		conn.Close()
		return
	}

	if item.Id < 0 {
		return
	}

	var rd Data
	err = gob.NewDecoder(conn).Decode(&rd)
	if err != nil {
		fmt.Println("Error al recibir el proceso", item.Id, "del cliente:", err)
		resume <- item
	} else {
		resume <- rd
	}

	fmt.Println("Cliente desconectado")
}

func main() {
	toClose := make(chan Data)
	closeChannel := make(chan bool)
	resumeChannel := make(chan Data)
	ln, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println(err)
		return
	}

	go process(closeChannel, resumeChannel, toClose)
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		closeChannel <- true
		go manageClient(conn, toClose, resumeChannel)
	}
}
