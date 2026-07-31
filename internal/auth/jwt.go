package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/dakaii/vibegopher/internal/envvar"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrUnauthorized = errors.New("unauthorized")

const tokenTTL = 24 * time.Hour

func GenerateJWT(user domain.User) (domain.AuthToken, error) {
	secret := envvar.AuthSecret()
	now := time.Now()
	expiresAt := now.Add(tokenTTL)

	claims := &domain.AuthTokenClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		UserID:   user.ID.String(),
		Username: user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return domain.AuthToken{}, fmt.Errorf("sign jwt: %w", err)
	}
	return domain.AuthToken{
		Token:     tokenString,
		TokenType: "Bearer",
		ExpiresIn: int64(tokenTTL.Seconds()),
	}, nil
}

func VerifyJWT(tknStr string) (domain.User, error) {
	claims := &domain.AuthTokenClaim{}
	token, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(envvar.AuthSecret()), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return domain.User{}, errors.New("signature invalid")
		}
		return domain.User{}, errors.New("could not parse the auth token")
	}
	if !token.Valid {
		return domain.User{}, errors.New("invalid token")
	}

	id, err := uuid.Parse(claims.UserID)
	if err != nil {
		return domain.User{}, errors.New("invalid user ID in token")
	}
	if claims.Username == "" {
		return domain.User{}, errors.New("invalid username in token")
	}

	return domain.User{
		ID:       id,
		Username: claims.Username,
	}, nil
}
