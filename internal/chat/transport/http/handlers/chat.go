package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/magabrotheeeer/go-chat/internal/chat/domain"
	"github.com/magabrotheeeer/go-chat/internal/lib/sl"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *domain.Chat) error
	FindByUserID(ctx context.Context, userID string) ([]*domain.Chat, error)
	FindByID(ctx context.Context, chatID string) (*domain.Chat, error)
	FindByUsers(ctx context.Context, user1ID, user2ID string) (*domain.Chat, error)
}

type ChatHandler struct {
	chatRepo ChatRepository
	logger   *slog.Logger
}

func NewChatHandler(chatRepo ChatRepository, logger *slog.Logger) *ChatHandler {
	return &ChatHandler{
		chatRepo: chatRepo,
		logger:   logger,
	}
}

// CreateChat создает новый чат между двумя пользователями
func (h *ChatHandler) CreateChat(c *gin.Context) {
	var req domain.CreateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind request", sl.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Проверяем, что пользователи не одинаковые
	if req.User1ID == req.User2ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot create chat with yourself"})
		return
	}

	// Проверяем, существует ли уже чат между этими пользователями
	existingChat, err := h.chatRepo.FindByUsers(c.Request.Context(), req.User1ID, req.User2ID)
	if err == nil && existingChat != nil {
		c.JSON(http.StatusOK, gin.H{"chat": existingChat})
		return
	}

	// Создаем новый чат
	chat := &domain.Chat{
		ID:        uuid.New().String(),
		User1ID:   req.User1ID,
		User2ID:   req.User2ID,
		CreatedAt: time.Now(),
	}

	if err := h.chatRepo.Create(c.Request.Context(), chat); err != nil {
		h.logger.Error("failed to create chat", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"chat": chat})
}

// GetUserChats получает все чаты пользователя
func (h *ChatHandler) GetUserChats(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	chats, err := h.chatRepo.FindByUserID(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("failed to get user chats", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chats": chats})
}

// GetChat получает конкретный чат по ID
func (h *ChatHandler) GetChat(c *gin.Context) {
	chatID := c.Param("chatID")
	if chatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chat ID is required"})
		return
	}

	chat, err := h.chatRepo.FindByID(c.Request.Context(), chatID)
	if err != nil {
		h.logger.Error("failed to get chat", sl.Err(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chat": chat})
}
