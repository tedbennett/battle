package relay

import (
	"golang.org/x/net/websocket"
)

type Relay struct {
	// Map of clients and whether or not they're connected
	clients map[*Client]bool

	Broadcast chan []byte
	// Channels - sending Client pointers
	Register   chan *Client
	Unregister chan *Client
}

func NewRelay() *Relay {
	return &Relay{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

type Client struct {
	conn     *websocket.Conn
	Messages chan []byte
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{conn, make(chan []byte)}
}

func (r *Relay) Run() {
	for {
		select {
		case msg := <-r.Broadcast:
			for client, active := range r.clients {
				if !active {
					continue
				}
				client.Messages <- msg
			}
		case client := <-r.Unregister:
			r.clients[client] = false

		case client := <-r.Register:
			r.clients[client] = true
		}
	}
}
