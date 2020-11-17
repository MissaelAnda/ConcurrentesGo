package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
)

type MessageChannel chan Message

type Message struct {
	Id       uint
	Type     clientCommandId
	Nickname string
	Text     string
	Bytes    []byte
}

type commandID uint

const (
	CONECTION_FAILED commandID = iota
	CONECTION_STABLISHED
	CLIENT_FILE
	CLIENT_MESSAGE
	CLIENT_QUIT
)

type ServerMessage struct {
	Type     commandID
	ClientID uint
	Nickname string
	Text     string
	Bytes    []byte
}

type ServerChannel chan ServerMessage

type ClientsManager struct {
	Clients  []*Client
	Output   ServerChannel
	Messages []Message
}

func (cm *ClientsManager) RemoveClient(id uint) {
	found := false
	var i int
	for idx, v := range cm.Clients {
		if v.Id == id {
			found = true
			i = idx
			break
		}
	}

	if !found {
		return
	}

	cm.Clients[len(cm.Clients)-1], cm.Clients[i] = cm.Clients[i], cm.Clients[len(cm.Clients)-1]
	cm.Clients = cm.Clients[:len(cm.Clients)-1]
}

func (c *ClientsManager) Append(client *Client) {
	c.Clients = append(c.Clients, client)
}

func (cm *ClientsManager) SendMessage(msg Message) {
	cm.Messages = append(cm.Messages, msg)

	for _, c := range cm.Clients {
		if c.SignedUp && c.Id != msg.Id {
			c.Input <- msg
		}
	}
}

func server(clientsManager *ClientsManager) {
	ln, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println("Error al intentar iniciar el servidor:", err)
		return
	}

	defer ln.Close()

	var idCounter uint = 1

	go messageDispatcher(clientsManager)

	fmt.Println("Servidor iniciado")
	fmt.Println("1. Mostrar registro de mensajes/archivos")
	fmt.Println("2. Guardar registro de mensajes/archivos")
	fmt.Println("3. Apagar servidor")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error al intentar conectar un cliente:", err)
			continue
		}

		// Registra nuevo cliente
		client := Client{
			Id:       idCounter,
			Nickname: "",
			SignedUp: false,
			Conn:     conn,
			Input:    make(MessageChannel),
			Output:   clientsManager.Output,
		}
		// Aumenta el id
		idCounter++

		clientsManager.Append(&client)

		go client.Run()
	}
}

func messageDispatcher(clientManager *ClientsManager) {
	for msg := range clientManager.Output {
		switch msg.Type {
		case CONECTION_FAILED:
			clientManager.RemoveClient(msg.ClientID)
			break

		case CONECTION_STABLISHED:
			fmt.Println("Se ha conectado el usuario", msg.Nickname)
			clientManager.SendMessage(Message{
				Id:       msg.ClientID,
				Type:     MESSAGE,
				Nickname: "Servidor",
				Text:     "Se ha conectado '" + msg.Nickname + "'.",
				Bytes:    nil,
			})
			break

		case CLIENT_QUIT:
			fmt.Println("Se ha desconectado el usuario", msg.Nickname)
			clientManager.SendMessage(Message{
				Id:       msg.ClientID,
				Type:     MESSAGE,
				Nickname: "Servidor",
				Text:     "Se ha desconectado '" + msg.Nickname + "'.",
				Bytes:    nil,
			})
			clientManager.RemoveClient(msg.ClientID)
			break

		case CLIENT_MESSAGE:
			fmt.Println(msg.ClientID, "-", msg.Nickname+":", msg.Text)

			clientManager.SendMessage(Message{
				Id:       msg.ClientID,
				Type:     MESSAGE,
				Nickname: msg.Nickname,
				Text:     msg.Text,
				Bytes:    nil,
			})
			break

		case CLIENT_FILE:
			fmt.Println(msg.ClientID, "-", msg.Nickname, "ha enviado un archivo:", msg.Text)

			clientManager.SendMessage(Message{
				Id:       msg.ClientID,
				Type:     FILE,
				Nickname: msg.Nickname,
				Text:     msg.Text,
				Bytes:    msg.Bytes,
			})
		}
	}
}

func saveToFile(file string, msgs []Message) {
	file += ".json"
	f, err := os.Create(file)

	if err != nil {
		fmt.Println("Hubo un error al intentar guardar el chat:", err)
		return
	}

	clientsLenght := len(msgs)

	f.WriteString("[\n")
	for i, m := range msgs {
		var message string
		if m.Type == FILE {
			message = "Archivo enviado: '" + m.Text + "'"
		} else {
			message = m.Text
		}

		str := "  {\n" +
			"    \"Nickname\": \"" + m.Nickname + "\",\n" +
			"    \"Message\": \"" + message + "\"\n" +
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
}

func main() {
	clientsManager := ClientsManager{
		Clients:  make([]*Client, 0),
		Output:   make(ServerChannel),
		Messages: make([]Message, 0),
	}

	go server(&clientsManager)

	var op int
	consoleReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Scan(&op)

		switch op {
		case 1:
			fmt.Println()
			for _, m := range clientsManager.Messages {
				var msg string
				if m.Type == FILE {
					msg = "Envio el archivo: '" + m.Text + "'."
				} else {
					msg = m.Text
				}
				fmt.Println(m.Nickname, "-", msg)
			}
			fmt.Println()
			break

		case 2:
			fmt.Print("File name: ")
			input, _ := consoleReader.ReadString('\n')
			input = input[:len(input)-1]

			// Save messages
			saveToFile(input, clientsManager.Messages)
			break

		case 3:
			return

		default:
		}
	}

}

// Client Data
type Client struct {
	Id       uint
	Conn     net.Conn
	Nickname string
	SignedUp bool
	Input    MessageChannel
	Output   ServerChannel
}

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

func (c *Client) Run() {
	decoder := gob.NewDecoder(c.Conn)

	var nickname string
	err := decoder.Decode(&nickname)

	if err != nil {
		fmt.Println("Error al iniciar sesion:", c.Conn.RemoteAddr)
		c.Conn.Close()
		c.Output <- ServerMessage{
			Type:     CONECTION_FAILED,
			ClientID: c.Id,
			Nickname: "",
			Text:     "",
		}
		return
	}

	c.Nickname = nickname
	c.SignedUp = true

	c.Output <- ServerMessage{
		Type:     CONECTION_STABLISHED,
		ClientID: c.Id,
		Nickname: nickname,
		Text:     "",
	}

	go c.Sender()

	for {
		msg := ClientMessage{}
		err := decoder.Decode(&msg)

		if err != nil {
			fmt.Println("Error al recibir mensaje del cliente ", c.Id, "-", c.Nickname, "\nError:", err)
			if err == io.EOF {
				c.Close()
				c.Output <- ServerMessage{
					Type:     CONECTION_FAILED,
					ClientID: c.Id,
					Nickname: "",
					Text:     "",
				}
				return
			}
			continue
		}

		switch msg.Type {
		case MESSAGE:
			c.Output <- ServerMessage{
				Type:     CLIENT_MESSAGE,
				ClientID: c.Id,
				Nickname: c.Nickname,
				Text:     msg.Text,
			}
			break

		case QUIT:
			c.Output <- ServerMessage{
				Type:     CLIENT_QUIT,
				ClientID: c.Id,
				Nickname: c.Nickname,
				Text:     "",
			}
			c.Close()
			return

		case FILE:
			c.Output <- ServerMessage{
				Type:     CLIENT_FILE,
				ClientID: c.Id,
				Nickname: c.Nickname,
				Text:     msg.Text,
				Bytes:    msg.Bytes,
			}
			break

		default:
		}
	}
}

func (c *Client) Sender() {
	encoder := gob.NewEncoder(c.Conn)

	for msg := range c.Input {
		err := encoder.Encode(msg)
		if err != nil {
			fmt.Println("Error al enviar el mensaje a:", c.Nickname, "error:", err)
		}
	}
}

func (c *Client) Close() {
	c.Conn.Close()
}
