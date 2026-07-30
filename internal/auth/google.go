package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/dakaii/vibegopher/internal/envvar"
	"google.golang.org/api/idtoken"
)

// GoogleTokenValidator validates Google ID tokens. Tests can inject a fake.
type GoogleTokenValidator interface {
	Validate(ctx context.Context, rawToken, audience string) (*domain.GoogleProfile, error)
}

type googleValidator struct{}

// NewGoogleTokenValidator returns the production Google ID token validator.
func NewGoogleTokenValidator() GoogleTokenValidator {
	return googleValidator{}
}

func (googleValidator) Validate(ctx context.Context, rawToken, audience string) (*domain.GoogleProfile, error) {
	if audience == "" {
		return nil, fmt.Errorf("GOOGLE_OAUTH_CLIENT_ID is not configured")
	}
	payload, err := idtoken.Validate(ctx, rawToken, audience)
	if err != nil {
		return nil, fmt.Errorf("invalid google id token: %w", err)
	}

	email, _ := payload.Claims["email"].(string)
	emailVerified, _ := payload.Claims["email_verified"].(bool)
	name, _ := payload.Claims["name"].(string)
	givenName, _ := payload.Claims["given_name"].(string)
	picture, _ := payload.Claims["picture"].(string)

	if payload.Subject == "" {
		return nil, fmt.Errorf("google token missing subject")
	}

	return &domain.GoogleProfile{
		Sub:           payload.Subject,
		Email:         strings.TrimSpace(email),
		EmailVerified: emailVerified,
		Name:          name,
		GivenName:     givenName,
		Picture:       picture,
	}, nil
}

// ValidateGoogleIDToken validates using the configured OAuth client ID.
func ValidateGoogleIDToken(ctx context.Context, rawToken string, validator GoogleTokenValidator) (*domain.GoogleProfile, error) {
	if validator == nil {
		validator = NewGoogleTokenValidator()
	}
	return validator.Validate(ctx, rawToken, envvar.GoogleOAuthClientID())
}
