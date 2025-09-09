package domain

import (
	"context"

	entities "github.com/magabrotheeeer/go-chat/internal/chat/domain/entities"
)

type MessageRepository interface {
	Save(ctx context.Context, msg *entities.Message) error
	FindByRoom(ctx context.Context, roomID string) ([]*entities.Message, error)
}

type RoomRepository interface {
	Create(ctx context.Context, room *string) error
	FindByID(ctx context.Context, roomID string) (*string, error)
}
