package testing

import (
	"log"

	"github.com/dakaii/vibegopher/internal/envvar"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TruncateAllTables cleans up application tables after tests.
// The goose version table is left intact so migrations are not re-applied from scratch every test.
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

	// Disable foreign key constraints temporarily
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
}
