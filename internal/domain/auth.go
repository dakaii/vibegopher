package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthToken is returned to clients after login / Google sign-in.
type AuthToken struct {
	TokenType string `json:"tokenType"`
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expiresIn"`
}

// AuthTokenClaim is the app session JWT. Never embed Password or full User.
type AuthTokenClaim struct {
	jwt.RegisteredClaims
	UserID   string `json:"uid"`
	Username string `json:"username"`
}

// GoogleAuthRequest is the payload for POST /api/auth/google.
type GoogleAuthRequest struct {
	IDToken string `json:"id_token"`
}

// GoogleProfile is the verified Google identity used to upsert a user.
type GoogleProfile struct {
	Sub           string
	Email         string
	EmailVerified bool
	Name          string
	GivenName     string
	Picture       string
}

func (c AuthTokenClaim) User() (User, error) {
	id, err := uuid.Parse(c.UserID)
	if err != nil {
		return User{}, err
	}
	return User{ID: id, Username: c.Username}, nil
}
