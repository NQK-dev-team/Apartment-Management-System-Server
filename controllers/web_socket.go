package controllers

import (
	"api/structs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketController struct{}

func NewWebSocketController() *WebSocketController {
	return &WebSocketController{}
}

var clients = make(map[*websocket.Conn]bool)                  // A map to track clients
var NotificationBroadcast = make(chan structs.NotificationWS) // A channel to broadcast notifications

// Configure the WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *WebSocketController) HandleNotificationConnection(ctx *gin.Context) {
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

func HandleNotificationBroadcast() {
	for {
		// Wait for a new notification to come in
		notification := <-NotificationBroadcast

		// Loop through all connected clients and send the notification
		for client := range clients {
			err := client.WriteJSON(notification)
			if err != nil {
				// log.Printf("Error sending notification to client: %v", err)
				client.Close()
				delete(clients, client) // Clean up on error
			}
		}
	}
}
