package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	domain "github.com/magabrotheeeer/go-chat/internal/chat/domain/entities"
)

type PostgresMessageRepository struct {
	db *pgxpool.Pool
}

func NewPostgresMessageRepository(db *pgxpool.Pool) *PostgresMessageRepository {
	return &PostgresMessageRepository{db: db}
}

func (r *PostgresMessageRepository) Save(ctx context.Context, msg *domain.Message) error {
	query := `
        INSERT INTO messages (id, room_id, author_id, content, created_at)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err := r.db.Exec(ctx, query, msg.ID, msg.RoomID, msg.AuthorID, msg.Content, msg.CreatedAt)
	return err
}

func (r *PostgresMessageRepository) FindByRoom(ctx context.Context, roomID string) ([]*domain.Message, error) {
	query := `
        SELECT id, room_id, author_id, content, created_at
        FROM messages WHERE room_id = $1 ORDER BY created_at DESC LIMIT 50
    `
	rows, err := r.db.Query(ctx, query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*domain.Message
	for rows.Next() {
		var m domain.Message
		var id string
		if err := rows.Scan(&id, &m.RoomID, &m.AuthorID, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		m.ID = (id)
		messages = append(messages, &m)
	}
	return messages, nil
}
