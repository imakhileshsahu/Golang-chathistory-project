package websocket

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
	mu   sync.Mutex
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()
	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		message := Message{Type: messageType, Body: string(p)}
		c.Pool.Broadcast <- message
		fmt.Printf("Message received:%+v\n", message)
	}
}

// package websocket

// import (
// 	"fmt"
// 	"log"
// 	"sync"

// 	"github.com/gorilla/websocket"
// )

// type Client struct {
// 	ID   string
// 	Conn *websocket.Conn
// 	Pool *Pool
// 	mu   sync.Mutex
// }

// type Message struct {
// 	Type   int    `json:"type"`
// 	Body   string `json:"body"`
// 	Sender string `json:"sender,omitempty"`
// }

// func (c *Client) Read() {
// 	defer func() {
// 		c.Pool.Unregister <- c
// 		c.Conn.Close()
// 		log.Printf("Client %s disconnected", c.ID)
// 	}()

// 	for {
// 		messageType, p, err := c.Conn.ReadMessage()
// 		if err != nil {
// 			log.Printf("Read error from client %s: %v", c.ID, err)
// 			return
// 		}

// 		message := Message{
// 			Type:   messageType,
// 			Body:   string(p),
// 			Sender: c.ID,
// 		}

// 		fmt.Printf("📩 Message received from %s: %+v\n", c.ID, message)
// 		c.Pool.Broadcast <- message
// 	}
// }
