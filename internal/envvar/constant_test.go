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

func TestCriticWorkerEnabled(t *testing.T) {
	t.Run("prefers CRITIC_WORKER_ENABLED", func(t *testing.T) {
		t.Setenv("CRITIC_WORKER_ENABLED", "true")
		t.Setenv("BOT_WORKER_ENABLED", "false")
		if !CriticWorkerEnabled() {
			t.Fatal("expected CRITIC_WORKER_ENABLED to win")
		}
	})
	t.Run("falls back to BOT_WORKER_ENABLED", func(t *testing.T) {
		t.Setenv("CRITIC_WORKER_ENABLED", "")
		t.Setenv("BOT_WORKER_ENABLED", "true")
		if !CriticWorkerEnabled() {
			t.Fatal("expected BOT_WORKER_ENABLED fallback")
		}
	})
}

func TestValidateRuntimeConfigRejectsPasswordAuthInProd(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("AUTH_SECRET", "a-sufficiently-long-secret")
	t.Setenv("CORS_ORIGIN", "https://example.com")
	t.Setenv("GOOGLE_OAUTH_CLIENT_ID", "client.apps.googleusercontent.com")
	t.Setenv("ENABLE_PASSWORD_AUTH", "true")
	t.Setenv("CRITIC_WORKER_ENABLED", "false")
	t.Setenv("BOT_WORKER_ENABLED", "")

	if err := ValidateRuntimeConfig(); err == nil {
		t.Fatal("expected password auth to be rejected in production")
	}
}

func TestValidateRuntimeConfigOKInProd(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("AUTH_SECRET", "a-sufficiently-long-secret")
	t.Setenv("CORS_ORIGIN", "https://example.com")
	t.Setenv("GOOGLE_OAUTH_CLIENT_ID", "client.apps.googleusercontent.com")
	t.Setenv("ENABLE_PASSWORD_AUTH", "false")
	t.Setenv("CRITIC_WORKER_ENABLED", "false")
	t.Setenv("BOT_WORKER_ENABLED", "")

	if err := ValidateRuntimeConfig(); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}
