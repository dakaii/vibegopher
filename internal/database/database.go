package database

import (
	"context"
	"log"
	"time"

	"github.com/dakaii/vibegopher/db"
	"github.com/dakaii/vibegopher/internal/envvar"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GetDatabase returns a database instance.
// If isTest is true, goose migrations are applied first (same SQL as CI/prod).
func GetDatabase(isTest ...bool) *gorm.DB {
	dsn := envvar.PostgresDSN()

	if len(isTest) > 0 && isTest[0] {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		if err := db.Up(ctx, dsn); err != nil {
			log.Fatal("Failed to apply database migrations:", err)
		}
	}

	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return gdb
}
