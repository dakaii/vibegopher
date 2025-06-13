package database

import (
	"fmt"
	"log"

	"github.com/dakaii/vibegopher/internal/envvar"
	"github.com/dakaii/vibegopher/internal/repository/userrepo"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GetDatabase returns a database instance.
// If isTest is true, it will auto-migrate the schema for testing.
func GetDatabase(isTest ...bool) *gorm.DB {
	user := envvar.DBUser()
	password := envvar.DBPassword()
	dbname := envvar.DBName()
	dbhost := envvar.DBHost()
	dbport := envvar.DBPort()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo",
		dbhost, user, password, dbname, dbport)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Auto-migrate schema if this is for testing
	if len(isTest) > 0 && isTest[0] {
		err = db.AutoMigrate(
			&userrepo.UserEntity{},
			// Add other entities here as you create them
		)
		if err != nil {
			log.Fatal("Failed to auto-migrate database schema:", err)
		}
	}

	return db
}
