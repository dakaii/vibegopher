package testing

import (
	"fmt"
	"log"

	"github.com/dakaii/vibegopher/internal/envvar"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TruncateAllTables cleans up all tables after tests
func TruncateAllTables() {
	password := envvar.DBPassword()
	dbname := envvar.DBName()
	dbhost := envvar.DBHost()
	dbport := envvar.DBPort()
	user := envvar.DBUser()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo",
		dbhost, user, password, dbname, dbport)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		db, _ := gormDB.DB()
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// Disable foreign key constraints temporarily
	err = gormDB.Exec("SET session_replication_role = replica;").Error
	if err != nil {
		log.Println("Failed to disable foreign key constraints:", err)
		return
	}

	// Delete all rows from all tables
	err = gormDB.Exec(`
		DO $$
		DECLARE
			r RECORD;
		BEGIN
			FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = current_schema()) LOOP
				EXECUTE 'DELETE FROM ' || quote_ident(r.tablename);
			END LOOP;
		END $$;
	`).Error
	if err != nil {
		log.Println("Failed to delete all rows from all tables:", err)
	}

	// Re-enable foreign key constraints
	err = gormDB.Exec("SET session_replication_role = DEFAULT;").Error
	if err != nil {
		log.Println("Failed to re-enable foreign key constraints:", err)
	}
}
