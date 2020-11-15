package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"
)

type MessageChannel chan Message

type Message struct {
	Id       uint
	Nickname string
	Text     string
}

type commandID uint

const (
	CONECTION_FAILED commandID = iota
	CONECTION_STABLISHED
	CLIENT_MESSAGE
	CLIENT_QUIT
)

type ServerMessage struct {
	Type     commandID
	ClientID uint
	Nickname string
	Text     string
}

type ServerChannel chan ServerMessage

type ClientsManager struct {
	Clients []*Client
	Output  ServerChannel
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
	for _, c := range cm.Clients {
		if c.SignedUp && c.Id != msg.Id {
			c.Input <- msg
		}
	}
}

func server() {
	ln, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println("Error al intentar iniciar el servidor:", err)
		return
	}

	defer ln.Close()

	var idCounter uint = 0
	clientsManager := ClientsManager{
		Clients: make([]*Client, 0),
		Output:  make(ServerChannel),
	}

	go messageDispatcher(&clientsManager)

	fmt.Println("Servidor iniciado...")

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
				Nickname: "Servidor",
				Text:     "Se ha conectado '" + msg.Nickname + "'.",
			})
			break

		case CLIENT_QUIT:
			fmt.Println("Se ha desconectado el usuario", msg.Nickname)
			clientManager.SendMessage(Message{
				Id:       msg.ClientID,
				Nickname: "Servidor",
				Text:     "Se ha desconectado '" + msg.Nickname + "'.",
			})
			clientManager.RemoveClient(msg.ClientID)
			break

		case CLIENT_MESSAGE:
			fmt.Println(msg.ClientID, "-", msg.Nickname+":", msg.Text)
			clientManager.SendMessage(Message{
				Id:       msg.ClientID,
				Nickname: msg.Nickname,
				Text:     msg.Text,
			})
		}
	}
}

func main() {
	go server()

	var input string
	fmt.Scan(&input)
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
	QUIT
)

type ClientMessage struct {
	Type clientCommandId
	Text string
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

	var msg ClientMessage
	for {
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
