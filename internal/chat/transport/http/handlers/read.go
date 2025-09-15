package http

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/magabrotheeeer/go-chat/internal/chat/domain"
	"github.com/magabrotheeeer/go-chat/internal/lib/sl"
)

type Repo interface {
	FindByRoom(ctx context.Context, roomID string) ([]*domain.Message, error)
}

type Handler struct {
	logger *slog.Logger
	repo Repo
}

func NewHandler(repo Repo, logger *slog.Logger) *Handler {
	return &Handler{
		logger: logger,
		repo: repo,
	}
}

func (h *Handler) Read(c *gin.Context) {
	roomID := c.Param("roomID")
	msgs, err := h.repo.FindByRoom(context.Background(), roomID)
	if err != nil {
		h.logger.Error("failed to find messages by room", sl.Err(err))
		return
	}
	c.JSON(http.StatusOK, msgs)
}
