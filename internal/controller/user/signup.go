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
	// Validate username
	if err := validateUsername(user.Username); err != nil {
		return domain.AuthToken{}, err
	}
	
	// Validate password
	if err := validatePassword(user.Password); err != nil {
		return domain.AuthToken{}, err
	}
	
	// Check if username already exists
	existingUser, _ := c.service.GetExistingUser(user.Username)
	if existingUser != nil {
		return domain.AuthToken{}, errors.New("this username is already in use")
	}

	createdUser, err := c.service.CreateUser(user)
	if err != nil {
		return domain.AuthToken{}, err
	}

	token := auth.GenerateJWT(*createdUser)
	return token, nil
}

func validateUsername(username string) error {
	// Trim whitespace
	username = strings.TrimSpace(username)
	
	// Check length (minimum 3, maximum 50 characters)
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	if len(username) > 50 {
		return errors.New("username cannot be longer than 50 characters")
	}
	
	// Check for valid characters (alphanumeric, underscore, hyphen)
	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validUsername.MatchString(username) {
		return errors.New("username can only contain letters, numbers, underscores, and hyphens")
	}
	
	// Username cannot start or end with underscore or hyphen
	if strings.HasPrefix(username, "_") || strings.HasPrefix(username, "-") ||
		strings.HasSuffix(username, "_") || strings.HasSuffix(username, "-") {
		return errors.New("username cannot start or end with underscore or hyphen")
	}
	
	return nil
}

func validatePassword(password string) error {
	// Check minimum length
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	
	// Check maximum length to prevent DoS
	if len(password) > 128 {
		return errors.New("password cannot be longer than 128 characters")
	}
	
	// Check for at least one uppercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	
	// Check for at least one lowercase letter
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	
	// Check for at least one digit
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasDigit {
		return errors.New("password must contain at least one digit")
	}
	
	return nil
}

// Legacy function for backward compatibility with existing tests
func isValidUsername(username string) bool {
	return validateUsername(username) == nil
}
