package domain

import (
	"github.com/golang-jwt/jwt/v5"
)

// AuthToken struct
type AuthToken struct {
	TokenType string `json:"tokenType"`
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expiresIn"`
}

// AuthTokenClaim struct
type AuthTokenClaim struct {
	jwt.RegisteredClaims
	User
}
