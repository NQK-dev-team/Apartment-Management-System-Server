package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// type WebSocketController struct {
// 	clients   map[*websocket.Conn]bool
// 	broadcast chan interface{}
// 	upgrader  websocket.Upgrader
// }

// func NewWebSocketController() *WebSocketController {
// 	return &WebSocketController{
// 		clients:   make(map[*websocket.Conn]bool),
// 		broadcast: make(chan interface{}),
// 		upgrader: websocket.Upgrader{
// 			CheckOrigin: func(r *http.Request) bool {
// 				return true // Allow all origins for simplicity; adjust as needed
// 			},
// 		},
// 	}
// }

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan interface{})
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleNotificationConnection(ctx *gin.Context) {
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		// log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer ws.Close()

	clients[ws] = true
	// log.Println("New client connected")

	// This loop keeps the WebSocket connection open.
	// You might read messages here if clients can send data back.
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			// log.Printf("Client disconnected: %v", err)
			delete(clients, ws)
			break
		}
	}
}

func HandleBroadcast() {
	for {
		// Wait for a new message to come in
		message := <-broadcast

		// Loop through all connected clients and send the message
		for client := range clients {
			err := client.WriteJSON(message)
			if err != nil {
				// log.Printf("Error sending notification to client: %v", err)
				client.Close()
				delete(clients, client) // Clean up on error
			}
		}
	}
}

func AddBroadcast(signal interface{}) {
	broadcast <- signal
}
