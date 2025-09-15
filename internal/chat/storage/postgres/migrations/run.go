package migrations

import (
	"context"
	"database/sql"
	"fmt"
	pgxv5 "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigration выполняет миграции из указанной директории
func RunMigration(ctx context.Context, db *sql.DB, migrationsPath string) error {
	driver, err := pgxv5.WithInstance(db, &pgxv5.Config{})
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"pgx_v5",
		driver,
	)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("%w", err)
	}

	return nil
}
