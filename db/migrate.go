// Package db owns versioned Postgres schema migrations (goose).
// GORM is used only as an ORM for queries — not for schema changes outside tests'
// historical AutoMigrate path (now also replaced by these migrations).
package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/lock"

	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func openDB(databaseURL string) (*sql.DB, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("database URL is empty (set DATABASE_URL or POSTGRES_* env vars)")
	}
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return db, nil
}

func newProvider(databaseURL string) (*goose.Provider, *sql.DB, error) {
	sqlDB, err := openDB(databaseURL)
	if err != nil {
		return nil, nil, err
	}

	fsys, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		_ = sqlDB.Close()
		return nil, nil, fmt.Errorf("migrations fs: %w", err)
	}

	sessionLocker, err := lock.NewPostgresSessionLocker()
	if err != nil {
		_ = sqlDB.Close()
		return nil, nil, fmt.Errorf("create session locker: %w", err)
	}

	provider, err := goose.NewProvider(
		goose.DialectPostgres,
		sqlDB,
		fsys,
		goose.WithSessionLocker(sessionLocker),
	)
	if err != nil {
		_ = sqlDB.Close()
		return nil, nil, fmt.Errorf("create goose provider: %w", err)
	}
	return provider, sqlDB, nil
}

// Up applies all pending migrations. Safe to run concurrently (Postgres advisory lock).
func Up(ctx context.Context, databaseURL string) error {
	provider, sqlDB, err := newProvider(databaseURL)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	if _, err := provider.Up(ctx); err != nil {
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil
}

// Down rolls back one migration. Prefer fix-forward in production.
func Down(ctx context.Context, databaseURL string) error {
	provider, sqlDB, err := newProvider(databaseURL)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	if _, err := provider.Down(ctx); err != nil {
		return fmt.Errorf("migrate down: %w", err)
	}
	return nil
}

// Status prints migration status to stdout.
func Status(ctx context.Context, databaseURL string) error {
	provider, sqlDB, err := newProvider(databaseURL)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	results, err := provider.Status(ctx)
	if err != nil {
		return fmt.Errorf("migrate status: %w", err)
	}
	for _, r := range results {
		fmt.Printf("%-12s %s\n", r.State, r.Source.Path)
	}
	return nil
}
