package envvar

import (
	"strings"
	"testing"
)

func TestDBSSLMode(t *testing.T) {
	t.Setenv("POSTGRES_SSLMODE", "")
	if got := DBSSLMode(); got != "disable" {
		t.Fatalf("empty POSTGRES_SSLMODE: got %q, want disable", got)
	}

	t.Setenv("POSTGRES_SSLMODE", "   ")
	if got := DBSSLMode(); got != "disable" {
		t.Fatalf("whitespace POSTGRES_SSLMODE: got %q, want disable", got)
	}

	t.Setenv("POSTGRES_SSLMODE", "require")
	if got := DBSSLMode(); got != "require" {
		t.Fatalf("POSTGRES_SSLMODE=require: got %q, want require", got)
	}
}

func TestPostgresDSNIncludesSSLMode(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("POSTGRES_HOST", "db.example")
	t.Setenv("POSTGRES_USER", "u")
	t.Setenv("POSTGRES_PASSWORD", "p")
	t.Setenv("POSTGRES_DB", "d")
	t.Setenv("POSTGRES_PORT", "5432")
	t.Setenv("POSTGRES_SSLMODE", "require")

	dsn := PostgresDSN()
	if !strings.Contains(dsn, "sslmode=require") {
		t.Fatalf("PostgresDSN() = %q, want sslmode=require", dsn)
	}

	t.Setenv("POSTGRES_SSLMODE", "")
	dsn = PostgresDSN()
	if !strings.Contains(dsn, "sslmode=disable") {
		t.Fatalf("PostgresDSN() with empty SSLMODE = %q, want sslmode=disable", dsn)
	}
}
