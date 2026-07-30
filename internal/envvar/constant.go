package envvar

import (
	"fmt"
	"os"
	"strconv"
)

func Port() string {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8081"
	}
	return port
}

// GetSecret returns the jwt secret.
func AuthSecret() string {
	secret, exists := os.LookupEnv("AUTH_SECRET")
	if !exists {
		secret = "secret_key"
	}
	return secret
}

// DatabaseURL returns a full Postgres URL when set (Neon / Cloud Run Secret Manager).
// When empty, callers should fall back to discrete POSTGRES_* variables via PostgresDSN.
func DatabaseURL() string {
	return os.Getenv("DATABASE_URL")
}

// PostgresDSN returns DATABASE_URL if set, otherwise a lib/pq keyword DSN from POSTGRES_*.
func PostgresDSN() string {
	if u := DatabaseURL(); u != "" {
		return u
	}
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo",
		DBHost(), DBUser(), DBPassword(), DBName(), DBPort(),
	)
}

// GoogleOAuthClientID is the OAuth client ID used as the audience for Google ID tokens.
func GoogleOAuthClientID() string {
	return os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
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

func HashCost() int {
	costString, _ := os.LookupEnv("HASH_COST")
	res, err := strconv.Atoi(costString)
	if err != nil {
		return 8
	}

	return res
}
