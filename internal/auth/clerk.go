package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkjwt "github.com/clerk/clerk-sdk-go/v2/jwt"
	clerkuser "github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/dakaii/vibegopher/internal/envvar"
)

// ClerkTokenValidator validates Clerk session JWTs. Tests can inject a fake.
type ClerkTokenValidator interface {
	// Validate verifies the session JWT and returns at least Sub.
	Validate(ctx context.Context, rawToken string) (*domain.ClerkProfile, error)
	// Enrich loads email/name/username from the Clerk Backend API.
	Enrich(ctx context.Context, clerkUserID string) (*domain.ClerkProfile, error)
}

type clerkValidator struct{}

// NewClerkTokenValidator returns the production Clerk session validator.
func NewClerkTokenValidator() ClerkTokenValidator {
	return clerkValidator{}
}

// ConfigureClerk sets the Clerk secret key used by the SDK (JWKS + Backend API).
func ConfigureClerk() {
	if key := envvar.ClerkSecretKey(); key != "" {
		clerk.SetKey(key)
	}
}

func (clerkValidator) Validate(ctx context.Context, rawToken string) (*domain.ClerkProfile, error) {
	if envvar.ClerkSecretKey() == "" {
		return nil, fmt.Errorf("CLERK_SECRET_KEY is not configured")
	}
	clerk.SetKey(envvar.ClerkSecretKey())

	params := &clerkjwt.VerifyParams{Token: rawToken}
	if parties := envvar.ClerkAuthorizedParties(); len(parties) > 0 {
		allowed := make(map[string]struct{}, len(parties))
		for _, p := range parties {
			allowed[p] = struct{}{}
		}
		params.AuthorizedPartyHandler = func(azp string) bool {
			if azp == "" {
				return true
			}
			_, ok := allowed[azp]
			return ok
		}
	}

	claims, err := clerkjwt.Verify(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("invalid clerk session token: %w", err)
	}
	if claims.Subject == "" {
		return nil, fmt.Errorf("clerk token missing subject")
	}
	return &domain.ClerkProfile{Sub: claims.Subject}, nil
}

func (clerkValidator) Enrich(ctx context.Context, clerkUserID string) (*domain.ClerkProfile, error) {
	if envvar.ClerkSecretKey() == "" {
		return nil, fmt.Errorf("CLERK_SECRET_KEY is not configured")
	}
	clerk.SetKey(envvar.ClerkSecretKey())

	usr, err := clerkuser.Get(ctx, clerkUserID)
	if err != nil {
		return nil, fmt.Errorf("fetch clerk user: %w", err)
	}
	return profileFromClerkUser(usr), nil
}

func profileFromClerkUser(usr *clerk.User) *domain.ClerkProfile {
	p := &domain.ClerkProfile{Sub: usr.ID}
	if usr.Username != nil {
		p.Username = strings.TrimSpace(*usr.Username)
	}
	if usr.FirstName != nil {
		p.GivenName = strings.TrimSpace(*usr.FirstName)
	}
	if usr.PrimaryEmailAddressID != nil {
		for _, ea := range usr.EmailAddresses {
			if ea != nil && ea.ID == *usr.PrimaryEmailAddressID {
				p.Email = strings.TrimSpace(ea.EmailAddress)
				break
			}
		}
	}
	if p.Email == "" && len(usr.EmailAddresses) > 0 && usr.EmailAddresses[0] != nil {
		p.Email = strings.TrimSpace(usr.EmailAddresses[0].EmailAddress)
	}
	return p
}
