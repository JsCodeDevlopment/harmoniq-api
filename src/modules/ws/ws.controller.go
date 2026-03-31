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
		return true // Adjust for production
	},
}

func InitModule(router *gin.RouterGroup) {
	hub := InitHub()
	go hub.Run()

	router.GET("/ws", func(c *gin.Context) {
		ServeWS(hub, c)
	})
}

func ServeWS(hub *Hub, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade to websocket: %v", err)
		return
	}

	// Optional: Extract user ID from context if authentication guard was used
	userID, _ := c.Get("user_id")
	idStr, _ := userID.(string)
	if idStr == "" {
		idStr = "anonymous"
	}

	client := &Client{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 256),
		ID:   idStr,
	}

	client.Hub.register <- client

	go client.WritePump()
	go client.ReadPump()
}
