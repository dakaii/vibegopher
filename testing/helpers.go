package testing

import (
	"log"

	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/dakaii/vibegopher/internal/envvar"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TruncateAllTables cleans up application tables after tests.
// The goose version table is left intact; the AI bot user is re-seeded.
func TruncateAllTables() {
	gormDB, err := gorm.Open(postgres.Open(envvar.PostgresDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		db, _ := gormDB.DB()
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	err = gormDB.Exec("SET session_replication_role = replica;").Error
	if err != nil {
		log.Println("Failed to disable foreign key constraints:", err)
		return
	}

	err = gormDB.Exec(`
		DO $$
		DECLARE
			r RECORD;
		BEGIN
			FOR r IN (
				SELECT tablename
				FROM pg_tables
				WHERE schemaname = current_schema()
				  AND tablename <> 'goose_db_version'
			) LOOP
				EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' RESTART IDENTITY CASCADE';
			END LOOP;
		END $$;
	`).Error
	if err != nil {
		log.Println("Failed to truncate tables:", err)
	}

	err = gormDB.Exec("SET session_replication_role = DEFAULT;").Error
	if err != nil {
		log.Println("Failed to re-enable foreign key constraints:", err)
	}

	// Re-seed bot account required for FK when the worker posts replies.
	err = gormDB.Exec(`
		INSERT INTO users (id, created_at, updated_at, username, password, google_sub, email, is_bot)
		VALUES (?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, ?, NULL, NULL, NULL, TRUE)
		ON CONFLICT (id) DO NOTHING
	`, domain.BotUserID, domain.BotUsername).Error
	if err != nil {
		log.Println("Failed to re-seed bot user:", err)
	}
}
