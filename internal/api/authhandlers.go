package api

import (
	"encoding/json"
	"net/http"

	"github.com/dakaii/vibegopher/internal/domain"
)

// SignupRequest represents the signup request payload
type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// MeResponse represents the me endpoint response
type MeResponse struct {
	Username string `json:"username"`
}

// Signup handles user registration
func (h *Handlers) Signup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user := domain.User{
		Username: req.Username,
		Password: req.Password,
	}

	authToken, err := h.controllers.UserController.Signup(user)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.writeJSON(w, authToken, http.StatusCreated)
}

// Login handles user authentication
func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user := domain.User{
		Username: req.Username,
		Password: req.Password,
	}

	authToken, err := h.controllers.UserController.Login(user)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	h.writeJSON(w, authToken, http.StatusOK)
}

// Me handles getting current user info
func (h *Handlers) Me(w http.ResponseWriter, r *http.Request, user domain.User) {
	response := MeResponse{
		Username: user.Username,
	}

	h.writeJSON(w, response, http.StatusOK)
}
