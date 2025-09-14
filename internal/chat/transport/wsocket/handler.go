package wsocket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/magabrotheeeer/go-chat/internal/chat/domain"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Handler struct {
	hub     *Hub
	msgRepo MessageRepository
}

func NewHandler(hub *Hub, msgRepo MessageRepository) *Handler {
	return &Handler{
		hub:     hub,
		msgRepo: msgRepo,
	}
}

func (h *Handler) HandleWebSocket(c *gin.Context) {
	roomID := c.Param("roomID")
	if roomID == "" {
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &Client{
		RoomID: roomID,
		Send:   make(chan *domain.Message),
	}
	h.hub.RegisterClient(client)

	go client.WritePump(conn)
	go client.ReadPump(conn, h.msgRepo, h.hub)
}
