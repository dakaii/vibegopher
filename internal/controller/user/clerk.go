package user

import (
	"context"
	"errors"

	"github.com/dakaii/vibegopher/internal/domain"
)

// AuthenticateWithClerk verifies a Clerk session JWT and returns the local user,
// creating one on first sight of a Clerk user id.
func (c *Controller) AuthenticateWithClerk(ctx context.Context, sessionToken string) (domain.User, error) {
	if sessionToken == "" {
		return domain.User{}, errors.New("session token is required")
	}

	profile, err := c.clerkValidator.Validate(ctx, sessionToken)
	if err != nil {
		return domain.User{}, err
	}

	if existing, err := c.service.GetByClerkUserID(profile.Sub); err == nil {
		return *existing, nil
	}

	if profile.Email == "" && profile.GivenName == "" && profile.Username == "" {
		enriched, err := c.clerkValidator.Enrich(ctx, profile.Sub)
		if err != nil {
			return domain.User{}, err
		}
		profile = enriched
	}

	user, err := c.service.UpsertClerkUser(*profile)
	if err != nil {
		return domain.User{}, err
	}
	return *user, nil
}
