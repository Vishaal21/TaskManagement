package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Production mein specific origins set karo
	},
}

// ServeWs Gin handler hai jo WebSocket connections handle karta hai
func ServeWs(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		// HTTP request ko WebSocket mein upgrade karo
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("WebSocket upgrade failed:", err)
			c.JSON(500, gin.H{"error": "WebSocket upgrade failed"})
			return
		}

		// Naya client banao
		client := &Client{
			conn: conn,
			send: make(chan []byte, 256),
		}

		// Client ko hub mein register karo
		hub.register <- client

		// Goroutines start karo
		go client.readPump(hub)
		go client.writePump()
	}
}

// readPump aur writePump same rahenge
func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		// Message ko broadcast karo (optional, agar chahiye)
		hub.Broadcast <- message
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for message := range c.send {
		c.conn.WriteMessage(websocket.TextMessage, message)
	}
}
