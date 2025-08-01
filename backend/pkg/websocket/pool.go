package websocket

import "fmt"

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("size of connection pool:", len(pool.Clients))
			for client, _ := range pool.Clients {
				fmt.Println(client)
				client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..."})
			}
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("size of connection pool:", len(pool.Clients))
			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
			}
			break
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in the pool")
			for client, _ := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}

// package websocket

// import "fmt"

// type Pool struct {
// 	Register   chan *Client
// 	Unregister chan *Client
// 	Clients    map[*Client]bool
// 	Broadcast  chan Message
// }

// func NewPool() *Pool {
// 	return &Pool{
// 		Register:   make(chan *Client),
// 		Unregister: make(chan *Client),
// 		Clients:    make(map[*Client]bool),
// 		Broadcast:  make(chan Message),
// 	}
// }

// func (pool *Pool) Start() {
// 	for {
// 		select {

// 		case client := <-pool.Register:
// 			pool.Clients[client] = true
// 			fmt.Println("✅ New client registered. Pool size:", len(pool.Clients))
// 			// Notify others
// 			for c := range pool.Clients {
// 				if c != client {
// 					c.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..."})
// 				}
// 			}

// 		case client := <-pool.Unregister:
// 			if _, ok := pool.Clients[client]; ok {
// 				delete(pool.Clients, client)
// 				fmt.Println("⚠️ Client unregistered. Pool size:", len(pool.Clients))
// 				// Notify others
// 				for c := range pool.Clients {
// 					if c != client {
// 						c.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
// 					}
// 				}
// 			}

// 		case message := <-pool.Broadcast:
// 			fmt.Println("📢 Broadcasting message to all clients")
// 			for client := range pool.Clients {
// 				if err := client.Conn.WriteJSON(message); err != nil {
// 					fmt.Println("❌ Write error:", err)
// 					pool.Unregister <- client
// 				}
// 			}
// 		}
// 	}
// }
