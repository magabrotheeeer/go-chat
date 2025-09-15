package wsocket

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/magabrotheeeer/go-chat/internal/chat/domain"
	"github.com/magabrotheeeer/go-chat/internal/lib/sl"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Handler struct {
	logger  *slog.Logger
	hub     *Hub
	msgRepo MessageRepository
}

func NewHandler(hub *Hub, msgRepo MessageRepository, logger *slog.Logger) *Handler {
	return &Handler{
		logger:  logger,
		hub:     hub,
		msgRepo: msgRepo,
	}
}

func (h *Handler) HandleWebSocket(c *gin.Context) {
	chatID := c.Param("chatID")
	if chatID == "" {
		h.logger.Error("failed to find param chatID")
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("failed to upgrade connection to web socket", sl.Err(err))
		return
	}

	client := &Client{
		ChatID: chatID,
		Send:   make(chan *domain.Message),
	}
	h.hub.RegisterClient(client)

	go client.WritePump(conn)
	go client.ReadPump(conn, h.msgRepo, h.hub)
}
