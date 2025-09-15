package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magabrotheeeer/go-chat/internal/chat/database"
	http "github.com/magabrotheeeer/go-chat/internal/chat/transport/http/handlers"
	"github.com/magabrotheeeer/go-chat/internal/chat/transport/wsocket"
	"github.com/magabrotheeeer/go-chat/internal/config"
)

func connectDB(connection string, ctx context.Context) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, connection)
	if err != nil {
		log.Fatalf("Unable to connect to DB: %v", err)
	}
	return pool
}

func main() {
	ctx := context.Background()
	cfg := config.MustLoad()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger.Info("starting go-chat")

	router := gin.Default()

	dbpool := connectDB(cfg.Database.Connection, ctx)
	// err := migrations.RunMigration(ctx, dbpool, "./migrations")
	// if err != nil {
	// 	logger.Error("failed to run migrations", sl.Err(err))
	// 	return		
	// }
	msgRepo := database.NewPostgresMessageRepository(dbpool)
	_ = database.NewPostgresRoomRepository(dbpool)

	hub := wsocket.NewHub()
	go hub.Run()

	router.GET("/rooms/:roomID/messages", http.NewHandler(msgRepo, logger).Read)
	router.GET("/ws/:roomID", wsocket.NewHandler(hub, msgRepo, logger).HandleWebSocket)
	router.Run(cfg.Server.Port)
}
