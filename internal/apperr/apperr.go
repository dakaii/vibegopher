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
	ErrRateLimited  = errors.New("rate limited")
)

// HTTPStatus maps an error to a safe client status + message.
// Auth failures are checked before broad "invalid"/"not found" heuristics.
func HTTPStatus(err error) (int, string) {
	if err == nil {
		return http.StatusInternalServerError, "internal error"
	}
	msg := err.Error()
	lower := strings.ToLower(msg)

	switch {
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized, "Unauthorized"
	case errors.Is(err, ErrRateLimited) || strings.Contains(lower, "rate limit"):
		return http.StatusTooManyRequests, "Too many requests"
	case strings.Contains(lower, "invalid credentials") ||
		strings.Contains(lower, "no user found with the inputted"):
		return http.StatusUnauthorized, "Invalid credentials"
	case strings.Contains(lower, "google id token") ||
		strings.Contains(lower, "google_oauth") ||
		strings.Contains(lower, "id_token"):
		return http.StatusUnauthorized, "Invalid Google credentials"
	case errors.Is(err, ErrNotFound) ||
		strings.HasPrefix(lower, "no post found") ||
		strings.HasPrefix(lower, "no comment found") ||
		strings.HasPrefix(lower, "no user found"):
		return http.StatusNotFound, "Not found"
	case errors.Is(err, ErrConflict) ||
		strings.Contains(lower, "already taken") ||
		strings.Contains(lower, "already in use"):
		return http.StatusConflict, msg
	case errors.Is(err, ErrValidation) ||
		strings.Contains(lower, "required") ||
		strings.Contains(lower, "cannot be empty") ||
		strings.Contains(lower, "must be") ||
		strings.Contains(lower, "invalid") ||
		strings.Contains(lower, "too short") ||
		strings.Contains(lower, "too long"):
		return http.StatusBadRequest, msg
	default:
		return http.StatusInternalServerError, "Something went wrong"
	}
}
