package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

type clientCommandId uint

const (
	MESSAGE clientCommandId = iota
	FILE
	QUIT
)

type ClientMessage struct {
	Type  clientCommandId
	Text  string
	Bytes []byte
}

type Message struct {
	Id       uint
	Type     clientCommandId
	Nickname string
	Text     string
	Bytes    []byte
}

func listener(conn net.Conn, messages chan Message) {
	decoder := gob.NewDecoder(conn)

	for {
		msg := Message{}
		err := decoder.Decode(&msg)

		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "closed network connection") {
				conn.Close()
				return
			}
			fmt.Println("Error al recibir mensaje:", err)
			continue
		}

		if msg.Type == FILE {
			f, err := os.Create(msg.Text)

			if err != nil {
				fmt.Println("Error al crear el archivo", msg.Text, "enviado:", err)
			} else {
				_, err := f.Write(msg.Bytes)

				if err != nil {
					fmt.Println("Error al guardar los datos del archivo:", err)
				}

				f.Close()
			}
		}

		messages <- msg
	}
}

type MessageManager struct {
	Messages []Message
}

func (mm *MessageManager) Run(messages chan Message) {
	for msg := range messages {
		mm.Messages = append(mm.Messages, msg)
		mm.PrintMessages()
	}
}

func (mm *MessageManager) PrintMessages() {
	fmt.Println()
	for _, m := range mm.Messages {
		if m.Type == MESSAGE {
			var nick string
			if m.Id != 0 {
				nick = m.Nickname
			} else {
				nick = "Tú"
			}
			fmt.Println(nick+":", m.Text)
		} else {
			var nick string
			if m.Id == 0 {
				nick = "Has"
			} else {
				nick = m.Nickname + " ha"
			}
			fmt.Println(nick, "enviado un archivo:", m.Text)
		}
	}
	fmt.Println()
}

func main() {
	messages := make(chan Message)
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
	go messageManager.Run(messages)
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
		fmt.Println("2. Enviar archivo")
		fmt.Println("3. Ver mensajes")
		fmt.Println("0. Salir")
		fmt.Scan(&op)

		switch op {
		case 1:
			fmt.Print("> ")
			text, _ := consoleReader.ReadString('\n')
			text = text[:len(text)-1]

			err = encoder.Encode(ClientMessage{
				Type:  MESSAGE,
				Text:  text,
				Bytes: nil,
			})

			if err != nil {
				fmt.Println("Error al enviar el mensaje:", err)
			} else {
				messages <- Message{
					Id:       0,
					Type:     MESSAGE,
					Nickname: nickname,
					Text:     text,
					Bytes:    nil,
				}
			}
			break

		case 2:
			fmt.Print("Nombre del archivo (con la direccion): ")
			input, _ := consoleReader.ReadString('\n')
			input = input[:len(input)-1]

			data, err := ioutil.ReadFile(input)

			if err != nil {
				fmt.Println("Error al leer el archivo:", err)
				break
			}

			last := strings.LastIndex("/", input)
			if last < 0 {
				last = 0
			}
			input = input[last:]

			err = encoder.Encode(ClientMessage{
				Type:  FILE,
				Text:  input,
				Bytes: data,
			})

			if err != nil {
				fmt.Println("Error al enviar el archivo:", err)
			} else {
				messages <- Message{
					Id:       0,
					Type:     FILE,
					Nickname: nickname,
					Text:     input,
				}
			}
			break

		case 3:
			messageManager.PrintMessages()
			break

		case 0:
			err = encoder.Encode(ClientMessage{
				Type:  QUIT,
				Text:  "",
				Bytes: nil,
			})

			if err != nil {
				fmt.Println("Error al intentar salir:", err)
			}
			return

		default:
		}
	}
}
