package envvar

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Port() string {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8081"
	}
	return port
}

// AppEnv returns APP_ENV (development|test|production). Empty means development.
func AppEnv() string {
	v := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	if v == "" {
		return "development"
	}
	return v
}

func IsProduction() bool {
	e := AppEnv()
	return e == "production" || e == "prod"
}

// AuthSecret returns the JWT signing secret used by password auth (tests/local).
// Prefer ValidateRuntimeConfig() at process start — do not rely on the weak default when password auth is enabled.
func AuthSecret() string {
	secret, exists := os.LookupEnv("AUTH_SECRET")
	if !exists {
		secret = "secret_key"
	}
	return secret
}

// ClerkSecretKey returns the Clerk Backend API secret used to verify session JWTs.
func ClerkSecretKey() string {
	return strings.TrimSpace(os.Getenv("CLERK_SECRET_KEY"))
}

// ClerkAuthorizedParties returns origins allowed in the Clerk session JWT azp claim.
// Prefers CLERK_AUTHORIZED_PARTIES (comma-separated); otherwise uses CORS_ORIGIN when it is not "*".
func ClerkAuthorizedParties() []string {
	if raw := strings.TrimSpace(os.Getenv("CLERK_AUTHORIZED_PARTIES")); raw != "" {
		return splitCSV(raw)
	}
	cors := CORSOrigin()
	if cors == "" || cors == "*" {
		return nil
	}
	return splitCSV(cors)
}

func splitCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func DatabaseURL() string {
	return os.Getenv("DATABASE_URL")
}

func PostgresDSN() string {
	if u := DatabaseURL(); u != "" {
		return u
	}
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Tokyo",
		DBHost(), DBUser(), DBPassword(), DBName(), DBPort(), DBSSLMode(),
	)
}


func GeminiAPIKey() string {
	return os.Getenv("GEMINI_API_KEY")
}

func envTruthy(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

// CriticWorkerEnabled starts the in-process critic poller.
// Prefers CRITIC_WORKER_ENABLED when set to a non-empty value; falls back to BOT_WORKER_ENABLED.
func CriticWorkerEnabled() bool {
	if v, ok := os.LookupEnv("CRITIC_WORKER_ENABLED"); ok && strings.TrimSpace(v) != "" {
		return envTruthy(v)
	}
	return envTruthy(os.Getenv("BOT_WORKER_ENABLED"))
}

// PasswordAuthEnabled exposes legacy /api/signup and /api/login (tests/local only).
func PasswordAuthEnabled() bool {
	return envTruthy(os.Getenv("ENABLE_PASSWORD_AUTH"))
}

func CORSOrigin() string {
	return strings.TrimSpace(os.Getenv("CORS_ORIGIN"))
}

func DBHost() string {
	host, exists := os.LookupEnv("POSTGRES_HOST")
	if !exists {
		host = "localhost"
	}
	return host
}

func DBUser() string {
	user, exists := os.LookupEnv("POSTGRES_USER")
	if !exists {
		user = "postgres"
	}
	return user
}

func DBPassword() string {
	password, exists := os.LookupEnv("POSTGRES_PASSWORD")
	if !exists {
		password = "postgres"
	}
	return password
}

func DBName() string {
	dbname, exists := os.LookupEnv("POSTGRES_DB")
	if !exists {
		dbname = "postgres"
	}
	return dbname
}

func DBPort() string {
	port, exists := os.LookupEnv("POSTGRES_PORT")
	if !exists {
		port = "5432"
	}
	return port
}

// DBSSLMode returns the Postgres sslmode. Defaults to disable for local Docker.
// Empty POSTGRES_SSLMODE is treated like unset.
func DBSSLMode() string {
	sslmode := strings.TrimSpace(os.Getenv("POSTGRES_SSLMODE"))
	if sslmode == "" {
		return "disable"
	}
	return sslmode
}

func HashCost() int {
	costString, _ := os.LookupEnv("HASH_COST")
	res, err := strconv.Atoi(costString)
	if err != nil {
		return 8
	}
	return res
}

// ValidateRuntimeConfig fails fast on unsafe API-server production settings.
func ValidateRuntimeConfig() error {
	if err := validateSharedSecrets(); err != nil {
		return err
	}
	if IsProduction() {
		if CORSOrigin() == "" || CORSOrigin() == "*" {
			return fmt.Errorf("CORS_ORIGIN must be set to explicit frontend origin(s) in production")
		}
		if ClerkSecretKey() == "" {
			return fmt.Errorf("CLERK_SECRET_KEY is required in production")
		}
		if PasswordAuthEnabled() {
			return fmt.Errorf("ENABLE_PASSWORD_AUTH must not be enabled in production")
		}
	}
	if CriticWorkerEnabled() {
		if err := validateGeminiKey(); err != nil {
			return err
		}
	}
	return nil
}

// ValidateWorkerConfig checks settings for the standalone critic worker (no HTTP/CORS/OAuth).
func ValidateWorkerConfig() error {
	return validateGeminiKey()
}

func validateSharedSecrets() error {
	// AUTH_SECRET is only required when password auth issues HS256 JWTs.
	if !PasswordAuthEnabled() {
		return nil
	}
	secret := AuthSecret()
	if secret == "" || secret == "secret_key" || len(secret) < 16 {
		return fmt.Errorf("AUTH_SECRET must be set to a strong value when ENABLE_PASSWORD_AUTH is on (min 16 chars)")
	}
	return nil
}

func validateGeminiKey() error {
	if GeminiAPIKey() != "" {
		return nil
	}
	if IsProduction() {
		return fmt.Errorf("GEMINI_API_KEY is required for the critic worker in production")
	}
	fmt.Println("warning: critic worker enabled but GEMINI_API_KEY is empty; critic jobs will fail")
	return nil
}
