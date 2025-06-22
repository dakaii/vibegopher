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

func GenerateJWT(user domain.User) domain.AuthToken {
	secret := envvar.AuthSecret()
	expiresAt := time.Now().Add(time.Minute * 15).Unix()

	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims = &domain.AuthTokenClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expiresAt, 0)),
		},
		User: user,
	}

	tokenString, error := token.SignedString([]byte(secret))
	if error != nil {
		fmt.Println(error)
	}
	return domain.AuthToken{
		Token:     tokenString,
		TokenType: "Bearer",
		ExpiresIn: expiresAt,
	}
}

func VerifyJWT(tknStr string) (domain.User, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method to prevent "none algorithm" attack
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

	decoded := make(map[string]interface{})
	for key, val := range claims {
		decoded[key] = val
	}

	var username string
	if keyExists(decoded, "username") {
		username = decoded["username"].(string)
	}

	var userID string
	if keyExists(decoded, "id") {
		userID = decoded["id"].(string)
	}

	// Parse UUID from string
	id, err := uuid.Parse(userID)
	if err != nil {
		return domain.User{}, errors.New("invalid user ID in token")
	}

	return domain.User{
		ID:       id,
		Username: username,
	}, nil
}

func keyExists(decoded map[string]interface{}, key string) bool {
	val, ok := decoded[key]
	return ok && val != nil
}

// func parseTimeString(param interface{}) (time.Time, error) {
// 	timeStringTypes := []string{time.RFC3339, time.RFC3339Nano, "2006-01-02T15:04:05.999999999Z07:00", "2024-04-27T02:56:51.45722294Z"}
// 	var t time.Time
// 	var err error

// 	switch v := param.(type) {
// 	case string:
// 		for _, timeStringType := range timeStringTypes {
// 			t, err = time.Parse(timeStringType, v)
// 		}
// 		if t.IsZero() {
// 			return time.Time{}, fmt.Errorf("could not parse string to time: %w", err)
// 		}
// 	case time.Time:
// 		t = v
// 	default:
// 		return time.Time{}, errors.New("unsupported type for time parsing")
// 	}

// 	return t, nil
// }
