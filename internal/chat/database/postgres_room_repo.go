package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	domain "github.com/magabrotheeeer/go-chat/internal/chat/domain/entities"
)

type PostgresRoomRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRoomRepository(db *pgxpool.Pool) *PostgresRoomRepository {
	return &PostgresRoomRepository{db: db}
}

func (r *PostgresRoomRepository) Create(ctx context.Context, room *domain.Room) error {
	query := `INSERT INTO rooms (id, name) VALUES ($1, $2)`
	_, err := r.db.Exec(ctx, query, room.ID, room.Name)
	return err
}

func (r *PostgresRoomRepository) FindByID(ctx context.Context, roomID string) (*domain.Room, error) {
	query := `SELECT id, name FROM rooms WHERE id = $1`
	row := r.db.QueryRow(ctx, query, roomID)
	var room domain.Room
	var id string
	err := row.Scan(&id, &room.Name)
	if err != nil {
		return nil, err
	}
	room.ID = id
	return &room, nil
}
