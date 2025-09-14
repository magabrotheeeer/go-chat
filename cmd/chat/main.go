package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magabrotheeeer/go-chat/internal/chat/config"
	"github.com/magabrotheeeer/go-chat/internal/chat/database"
	"github.com/magabrotheeeer/go-chat/internal/chat/domain"
	"github.com/magabrotheeeer/go-chat/internal/chat/transport/wsocket"
)

func connectDB(connection string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), connection)
	if err != nil {
		log.Fatalf("Unable to connect to DB: %v", err)
	}
	return pool
}

func main() {
	cfg := config.MustLoad()
	router := gin.Default()

	dbpool := connectDB(cfg.Database.Connection)
	msgRepo := database.NewPostgresMessageRepository(dbpool)
	_ = database.NewPostgresRoomRepository(dbpool)

	hub := wsocket.NewHub()
	go hub.Run()

	router.GET("/rooms/:roomID/messages", func(c *gin.Context) {
		roomID := c.Param("roomID")
		msgs, err := msgRepo.FindByRoom(context.Background(), roomID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, msgs)
	})

	router.GET("/ws/:roomID", func(c *gin.Context) {
		roomID := c.Param("roomID")
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Failed to set WebSocket upgrade:", err)
			return
		}
		client := &wsocket.Client{RoomID: roomID, Send: make(chan *domain.Message)}
		hub.RegisterClient(client)

		// Запускаем горутины для чтения и записи
		go client.WritePump(conn)
		go client.ReadPump(conn, msgRepo, hub)
	})

	router.Run(":8080")
}
