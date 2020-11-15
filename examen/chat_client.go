package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

type clientCommandId uint

const (
	MESSAGE clientCommandId = iota
	QUIT
)

type ClientMessage struct {
	Type clientCommandId
	Text string
}

type Message struct {
	Id       uint
	Nickname string
	Text     string
}

func listener(conn net.Conn, inputs chan Message) {
	decoder := gob.NewDecoder(conn)
	var msg Message
	for {
		err := decoder.Decode(&msg)

		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "closed network connection") {
				conn.Close()
				return
			}
			fmt.Println("Error al recibir mensaje:", err)
			continue
		}

		inputs <- msg
	}
}

type MessageManager struct {
	Messages []Message
}

func (mm *MessageManager) Run(messages chan Message, signal chan string) {
	for {
		select {
		case msg := <-messages:
			mm.Messages = append(mm.Messages, msg)
			fmt.Println()
			for _, m := range mm.Messages {
				fmt.Println(m.Nickname+":", m.Text)
			}
			fmt.Println()
			break

		case file := <-signal:
			file += ".json"
			f, err := os.Create(file)

			if err != nil {
				fmt.Println("Hubo un error al intentar guardar el chat:", err)
				break
			}

			clientsLenght := len(mm.Messages)

			f.WriteString("[\n")
			for i, m := range mm.Messages {
				str := "  {\n" +
					"    \"Nickname\": \"" + m.Nickname + "\",\n" +
					"    \"Message\": \"" + m.Text + "\"\n" +
					"  }"

				if i+1 != clientsLenght {
					str += ","
				}
				str += "\n"

				f.WriteString(str)
			}
			f.WriteString("]\n")

			f.Close()

			fmt.Println("Se ha guardado el chat en '" + file + "'.")
			break

		default:
		}
	}
}

func main() {
	messages := make(chan Message)
	signal := make(chan string)
	messageManager := MessageManager{
		Messages: make([]Message, 0),
	}
	conn, err := net.Dial("tcp", ":8080")

	if err != nil {
		fmt.Println("Error al intentar conectarse al servidor:", err)
		return
	}

	var nickname string
	fmt.Print("Nickname: ")
	fmt.Scan(&nickname)

	go listener(conn, messages)
	go messageManager.Run(messages, signal)
	defer conn.Close()

	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(nickname)

	if err != nil {
		fmt.Println("Error al intentar iniciar sesion:", err)
		return
	}

	var op int
	consoleReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println()
		fmt.Println("1. Enviar mensaje")
		fmt.Println("2. Guardar chat")
		fmt.Println("0. Salir")
		fmt.Scan(&op)

		switch op {
		case 1:
			fmt.Print("> ")
			input, _ := consoleReader.ReadString('\n')
			input = input[:len(input)-1]

			err = encoder.Encode(ClientMessage{
				Type: MESSAGE,
				Text: input,
			})

			if err != nil {
				fmt.Println("Error al enviar el mensaje:", err)
			} else {
				messages <- Message{
					Id:       0,
					Nickname: nickname,
					Text:     input,
				}
			}
			break

		case 2:
			fmt.Print("Nombre del archivo: ")
			input, _ := consoleReader.ReadString('\n')
			input = input[:len(input)-1]

			signal <- input
			break

		case 0:
			err = encoder.Encode(ClientMessage{
				Type: QUIT,
				Text: "",
			})

			if err != nil {
				fmt.Println("Error al intentar salir:", err)
			}
			return

		default:
		}
	}
}
