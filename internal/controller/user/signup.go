package user

import (
	"errors"
	"regexp"
	"strings"

	"github.com/dakaii/vibegopher/internal/auth"
	"github.com/dakaii/vibegopher/internal/domain"
)

// Signup lets users sign up for this application and returns a jwt.
func (c *Controller) Signup(user domain.User) (domain.AuthToken, error) {
	if err := validateUsername(user.Username); err != nil {
		return domain.AuthToken{}, err
	}
	if err := validatePassword(user.Password); err != nil {
		return domain.AuthToken{}, err
	}

	existingUser, _ := c.service.GetExistingUser(user.Username)
	if existingUser != nil {
		return domain.AuthToken{}, errors.New("this username is already in use")
	}

	createdUser, err := c.service.CreateUser(user)
	if err != nil {
		return domain.AuthToken{}, err
	}

	return auth.GenerateJWT(*createdUser)
}

func validateUsername(username string) error {
	username = strings.TrimSpace(username)
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	if len(username) > 50 {
		return errors.New("username cannot be longer than 50 characters")
	}

	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validUsername.MatchString(username) {
		return errors.New("username can only contain letters, numbers, underscores, and hyphens")
	}

	if strings.HasPrefix(username, "_") || strings.HasPrefix(username, "-") ||
		strings.HasSuffix(username, "_") || strings.HasSuffix(username, "-") {
		return errors.New("username cannot start or end with underscore or hyphen")
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if len(password) > 128 {
		return errors.New("password cannot be longer than 128 characters")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return errors.New("password must contain at least one digit")
	}
	return nil
}

// Legacy function for backward compatibility with existing tests
func isValidUsername(username string) bool {
	return validateUsername(username) == nil
}
