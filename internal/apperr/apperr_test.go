package apperr

import (
	"errors"
	"net/http"
	"testing"
)

func TestHTTPStatusAuthBeforeHeuristics(t *testing.T) {
	cases := []struct {
		err    error
		status int
		msg    string
	}{
		{errors.New("invalid credentials"), http.StatusUnauthorized, "Invalid credentials"},
		{errors.New("no user found with the inputted username"), http.StatusUnauthorized, "Invalid credentials"},
		{errors.New("google id token invalid: aud mismatch"), http.StatusUnauthorized, "Invalid Google credentials"},
		{errors.New("no post found with ID: x"), http.StatusNotFound, "Not found"},
		{errors.New("content cannot be empty and must be 280 characters or less"), http.StatusBadRequest, "content cannot be empty and must be 280 characters or less"},
		{ErrRateLimited, http.StatusTooManyRequests, "Too many requests"},
		{errors.New("database exploded"), http.StatusInternalServerError, "Something went wrong"},
	}
	for _, tc := range cases {
		status, msg := HTTPStatus(tc.err)
		if status != tc.status || msg != tc.msg {
			t.Fatalf("%v: got %d %q, want %d %q", tc.err, status, msg, tc.status, tc.msg)
		}
	}
}
