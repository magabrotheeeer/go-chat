package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magabrotheeeer/go-chat/internal/config"
	"github.com/magabrotheeeer/go-chat/internal/chat/database"
	http "github.com/magabrotheeeer/go-chat/internal/chat/transport/http/handlers"
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

	router.GET("/rooms/:roomID/messages", http.NewHandler(msgRepo).Read)
	router.GET("/ws/:roomID", wsocket.NewHandler(hub, msgRepo).HandleWebSocket)
	router.Run(cfg.Server.Port)
}
