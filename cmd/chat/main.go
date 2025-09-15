package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/magabrotheeeer/go-chat/internal/chat/storage/postgres"
	"github.com/magabrotheeeer/go-chat/internal/chat/storage/postgres/migrations"
	"github.com/magabrotheeeer/go-chat/internal/chat/transport/http/handlers"
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
	chatRepo := postgres.NewPostgresChatRepository(db)

	chatHandler := handlers.NewChatHandler(chatRepo, logger)

	hub := wsocket.NewHub()
	go hub.Run()

	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	api := router.Group("/api")
	{
		api.POST("/chats", chatHandler.CreateChat)
		api.GET("/users/:userID/chats", chatHandler.GetUserChats)
		api.GET("/chats/:chatID", chatHandler.GetChat)
	}

	router.GET("/ws/chat/:chatID", wsocket.NewHandler(hub, msgRepo, logger).HandleWebSocket)

	logger.Info("server starting on port " + cfg.Server.Port)
	router.Run(cfg.Server.Port)
}
