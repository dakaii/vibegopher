package apperr

import (
	"errors"
	"net/http"
	"strings"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrNotFound     = errors.New("not found")
	ErrValidation   = errors.New("validation failed")
	ErrConflict     = errors.New("conflict")
)

// HTTPStatus maps an error to a safe client status + message.
func HTTPStatus(err error) (int, string) {
	if err == nil {
		return http.StatusInternalServerError, "internal error"
	}
	msg := err.Error()
	switch {
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized, "Unauthorized"
	case errors.Is(err, ErrNotFound) || strings.HasPrefix(msg, "no post found") || strings.HasPrefix(msg, "no comment found") || strings.HasPrefix(msg, "no user found"):
		return http.StatusNotFound, "Not found"
	case errors.Is(err, ErrConflict) || strings.Contains(msg, "already taken") || strings.Contains(msg, "already in use"):
		return http.StatusConflict, msg
	case errors.Is(err, ErrValidation) ||
		strings.Contains(msg, "required") ||
		strings.Contains(msg, "cannot be empty") ||
		strings.Contains(msg, "must be") ||
		strings.Contains(msg, "invalid") ||
		strings.Contains(msg, "too short") ||
		strings.Contains(msg, "too long"):
		return http.StatusBadRequest, msg
	case strings.Contains(msg, "invalid credentials") || strings.Contains(msg, "no user found with the inputted"):
		return http.StatusUnauthorized, "Invalid credentials"
	case strings.Contains(msg, "google id token") || strings.Contains(msg, "GOOGLE_OAUTH"):
		return http.StatusUnauthorized, "Invalid Google credentials"
	default:
		return http.StatusInternalServerError, "Something went wrong"
	}
}
