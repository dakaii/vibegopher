package envvar

import "testing"

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

func TestValidateWorkerConfigSkipsAPIOnlySettings(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("AUTH_SECRET", "")
	t.Setenv("CORS_ORIGIN", "")
	t.Setenv("GOOGLE_OAUTH_CLIENT_ID", "")
	t.Setenv("GEMINI_API_KEY", "test-key")

	if err := ValidateWorkerConfig(); err != nil {
		t.Fatalf("worker should only require Gemini: %v", err)
	}
}
