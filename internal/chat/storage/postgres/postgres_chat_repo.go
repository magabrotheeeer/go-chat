package postgres

import (
	"context"
	"database/sql"

	"github.com/magabrotheeeer/go-chat/internal/chat/domain"
)

type PostgresChatRepository struct {
	db *sql.DB
}

func NewPostgresChatRepository(db *sql.DB) *PostgresChatRepository {
	return &PostgresChatRepository{db: db}
}

func (r *PostgresChatRepository) Create(ctx context.Context, chat *domain.Chat) error {
	query := `
		INSERT INTO chats (id, user1_id, user2_id, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user1_id, user2_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, chat.ID, chat.User1ID, chat.User2ID, chat.CreatedAt)
	return err
}

func (r *PostgresChatRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.Chat, error) {
	query := `
		SELECT id, user1_id, user2_id, created_at
		FROM chats 
		WHERE user1_id = $1 OR user2_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []*domain.Chat
	for rows.Next() {
		var chat domain.Chat
		var id string
		if err := rows.Scan(&id, &chat.User1ID, &chat.User2ID, &chat.CreatedAt); err != nil {
			return nil, err
		}
		chat.ID = id
		chats = append(chats, &chat)
	}
	return chats, nil
}

func (r *PostgresChatRepository) FindByID(ctx context.Context, chatID string) (*domain.Chat, error) {
	query := `SELECT id, user1_id, user2_id, created_at FROM chats WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, chatID)
	var chat domain.Chat
	var id string
	err := row.Scan(&id, &chat.User1ID, &chat.User2ID, &chat.CreatedAt)
	if err != nil {
		return nil, err
	}
	chat.ID = id
	return &chat, nil
}

func (r *PostgresChatRepository) FindByUsers(ctx context.Context, user1ID, user2ID string) (*domain.Chat, error) {
	query := `
		SELECT id, user1_id, user2_id, created_at 
		FROM chats 
		WHERE (user1_id = $1 AND user2_id = $2) OR (user1_id = $2 AND user2_id = $1)
	`
	row := r.db.QueryRowContext(ctx, query, user1ID, user2ID)
	var chat domain.Chat
	var id string
	err := row.Scan(&id, &chat.User1ID, &chat.User2ID, &chat.CreatedAt)
	if err != nil {
		return nil, err
	}
	chat.ID = id
	return &chat, nil
}
