package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketController struct {
	clients   map[*websocket.Conn]bool
	broadcast chan interface{}
	upgrader  websocket.Upgrader
}

func NewWebSocketController() *WebSocketController {
	return &WebSocketController{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan interface{}),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for simplicity; adjust as needed
			},
		},
	}
}

func (c *WebSocketController) HandleNotificationConnection(ctx *gin.Context) {
	ws, err := c.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		// log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer ws.Close()

	c.clients[ws] = true
	// log.Println("New client connected")

	// This loop keeps the WebSocket connection open.
	// You might read messages here if clients can send data back.
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			// log.Printf("Client disconnected: %v", err)
			delete(c.clients, ws)
			break
		}
	}
}

func (c *WebSocketController) HandleBroadcast() {
	for {
		// Wait for a new message to come in
		message := <-c.broadcast

		// Loop through all connected clients and send the message
		for client := range c.clients {
			err := client.WriteJSON(message)
			if err != nil {
				// log.Printf("Error sending notification to client: %v", err)
				client.Close()
				delete(c.clients, client) // Clean up on error
			}
		}
	}
}
