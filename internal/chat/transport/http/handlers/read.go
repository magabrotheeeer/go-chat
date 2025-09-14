package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/magabrotheeeer/go-chat/internal/chat/domain"
)

type Repo interface {
	FindByRoom(ctx context.Context, roomID string) ([]*domain.Message, error)
}

type Handler struct {
	repo Repo
}

func NewHandler(repo Repo) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Read(c *gin.Context) {
	roomID := c.Param("roomID")
	msgs, err := h.repo.FindByRoom(context.Background(), roomID)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, msgs)
}
