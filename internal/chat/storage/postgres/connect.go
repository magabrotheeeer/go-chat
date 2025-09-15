package postgres

import (
	"context"
	"database/sql"
	"log"
)

func ConnectDB(connection string, ctx context.Context) *sql.DB {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatalf("Unable to connect to DB: %v", err)
	}

	// Проверяем подключение
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Unable to ping DB: %v", err)
	}

	return db
}