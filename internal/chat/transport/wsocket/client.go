package wsocket

import (
	"context"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/magabrotheeeer/go-chat/internal/chat/domain"
)

type Client struct {
	RoomID string
	Send   chan *domain.Message
}

type MessageRepository interface {
	Save(ctx context.Context, msg *domain.Message) error
	FindByRoom(ctx context.Context, roomID string) ([]*domain.Message, error)
}

type RoomRepository interface {
	Create(ctx context.Context, room *string) error
	FindByID(ctx context.Context, roomID string) (*string, error)
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

// ReadPump перекачивает сообщения из WebSocket соединения в hub
func (c *Client) ReadPump(conn *websocket.Conn, msgRepo MessageRepository, hub *Hub) {
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
		var msg domain.Message
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
