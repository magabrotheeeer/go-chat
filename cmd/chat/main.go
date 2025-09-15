package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/magabrotheeeer/go-chat/internal/chat/storage/postgres"
	"github.com/magabrotheeeer/go-chat/internal/chat/storage/postgres/migrations"
	http "github.com/magabrotheeeer/go-chat/internal/chat/transport/http/handlers"
	"github.com/magabrotheeeer/go-chat/internal/chat/transport/wsocket"
	"github.com/magabrotheeeer/go-chat/internal/config"
	"github.com/magabrotheeeer/go-chat/internal/lib/sl"
)



func main() {
	ctx := context.Background()
	cfg := config.MustLoad()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger.Info("starting go-chat")

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	db := postgres.ConnectDB(cfg.Database.Connection, ctx)
	err := migrations.RunMigration(ctx, db, "./migrations")
	if err != nil {
		logger.Error("failed to run migrations", sl.Err(err))
		return
	}
	msgRepo := postgres.NewPostgresMessageRepository(db)
	_ = postgres.NewPostgresRoomRepository(db)

	hub := wsocket.NewHub()
	go hub.Run()

	router.GET("/rooms/:roomID/messages", http.NewHandler(msgRepo, logger).Read)
	router.GET("/ws/:roomID", wsocket.NewHandler(hub, msgRepo, logger).HandleWebSocket)
	router.Run(cfg.Server.Port)
}
