package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthToken is returned to clients after password login (tests/local).
type AuthToken struct {
	TokenType string `json:"tokenType"`
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expiresIn"`
}

// AuthTokenClaim is the app session JWT used by password auth. Never embed Password or full User.
type AuthTokenClaim struct {
	jwt.RegisteredClaims
	UserID   string `json:"uid"`
	Username string `json:"username"`
}

// ClerkProfile is the verified Clerk identity used to upsert a user.
type ClerkProfile struct {
	Sub       string
	Email     string
	Username  string
	GivenName string
}

func (c AuthTokenClaim) User() (User, error) {
	id, err := uuid.Parse(c.UserID)
	if err != nil {
		return User{}, err
	}
	return User{ID: id, Username: c.Username}, nil
}
