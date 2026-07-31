package user

import (
	"context"
	"errors"

	"github.com/dakaii/vibegopher/internal/auth"
	"github.com/dakaii/vibegopher/internal/domain"
)

// LoginWithGoogle verifies a Google ID token and returns an app session JWT.
func (c *Controller) LoginWithGoogle(ctx context.Context, idToken string) (domain.AuthToken, error) {
	if idToken == "" {
		return domain.AuthToken{}, errors.New("id_token is required")
	}

	profile, err := auth.ValidateGoogleIDToken(ctx, idToken, c.googleValidator)
	if err != nil {
		return domain.AuthToken{}, err
	}

	user, err := c.service.UpsertGoogleUser(*profile)
	if err != nil {
		return domain.AuthToken{}, err
	}

	return auth.GenerateJWT(*user)
}
