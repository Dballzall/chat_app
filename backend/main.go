package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Client represents a connected user
type Client struct {
	conn     *websocket.Conn
	username string
}

// Message represents the structure of a chat message
type Message struct {
	Type     string `json:"type"`
	Content  string `json:"content"`
	Username string `json:"username"`
}

var (
	// Configure the upgrader
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for development
		},
	}

	// Maintain active clients
	clients    = make(map[*Client]bool)
	clientsMux sync.Mutex

	// Broadcast channel
	broadcast = make(chan Message)
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}
	defer ws.Close()

	// Create new client
	client := &Client{
		conn:     ws,
		username: r.URL.Query().Get("username"),
	}

	// Register new client
	clientsMux.Lock()
	clients[client] = true
	clientsMux.Unlock()

	// Announce new user
	broadcast <- Message{
		Type:     "join",
		Content:  fmt.Sprintf("%s joined the chat", client.username),
		Username: "System",
	}

	for {
		var msg Message
		// Read new message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			clientsMux.Lock()
			delete(clients, client)
			clientsMux.Unlock()
			broadcast <- Message{
				Type:     "leave",
				Content:  fmt.Sprintf("%s left the chat", client.username),
				Username: "System",
			}
			break
		}
		// Add username to message
		msg.Username = client.username
		// Send message to broadcast channel
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Grab next message from broadcast channel
		msg := <-broadcast

		// Send to every client
		clientsMux.Lock()
		for client := range clients {
			err := client.conn.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.conn.Close()
				delete(clients, client)
			}
		}
		clientsMux.Unlock()
	}
}

func main() {
	// Create a simple file server
	fs := http.FileServer(http.Dir("../frontend/dist"))
	http.Handle("/", fs)

	// Configure websocket route
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		handleConnections(w, r)
	})

	// Start listening for incoming chat messages
	go handleMessages()

	// Start the server
	log.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
