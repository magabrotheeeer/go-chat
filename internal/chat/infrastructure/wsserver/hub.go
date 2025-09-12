package wsserver

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	entities "github.com/magabrotheeeer/go-chat/internal/chat/domain/entities"
	repositories "github.com/magabrotheeeer/go-chat/internal/chat/domain/repositories"
)

type Client struct {
	RoomID string
	Send   chan *entities.Message
}

const (
	// Время ожидания записи сообщения в peer
	writeWait = 10 * time.Second

	// Время ожидания pong от peer
	pongWait = 60 * time.Second

	// Отправка ping в peer. Должно быть меньше pongWait
	pingPeriod = (pongWait * 9) / 10

	// Максимальный размер сообщения
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ReadPump перекачивает сообщения из WebSocket соединения в hub
func (c *Client) ReadPump(conn *websocket.Conn, msgRepo repositories.MessageRepository, hub *Hub) {
	defer func() {
		hub.UnregisterClient(c)
		conn.Close()
	}()

	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var msg entities.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Сохраняем сообщение в базу данных
		msg.RoomID = c.RoomID
		err = msgRepo.Save(context.Background(), &msg)
		if err != nil {
			log.Printf("Error saving message: %v", err)
			continue
		}

		// Отправляем сообщение всем клиентам в комнате
		hub.BroadcastMessage(&msg)
	}
}

// WritePump перекачивает сообщения из hub в WebSocket соединение
func (c *Client) WritePump(conn *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := conn.WriteJSON(message); err != nil {
				return
			}

		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

type Hub struct {
	mu         sync.RWMutex
	rooms      map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *entities.Message
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *entities.Message),
	}
}

func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

func (h *Hub) BroadcastMessage(msg *entities.Message) {
	h.broadcast <- msg
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.rooms[client.RoomID] == nil {
				h.rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.rooms[client.RoomID][client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.rooms[client.RoomID]; ok {
				delete(clients, client)
				close(client.Send)
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.RLock()
			if clients, ok := h.rooms[msg.RoomID]; ok {
				for client := range clients {
					select {
					case client.Send <- msg:
					default:
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}
