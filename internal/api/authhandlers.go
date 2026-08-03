package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dakaii/vibegopher/internal/domain"
)

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type MeResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
	IsBot    bool   `json:"is_bot"`
}

// Signup is demoted password registration (kept for tests / legacy).
func (h *Handlers) Signup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.Password == "" {
		h.writeError(w, "Password is required", http.StatusBadRequest)
		return
	}

	user := domain.User{
		Username: req.Username,
		Password: req.Password,
	}

	authToken, err := h.controllers.UserController.Signup(user)
	if err != nil {
		h.writeErr(w, err)
		return
	}

	h.writeJSON(w, authToken, http.StatusCreated)
}

// Login is demoted password login (kept for tests / legacy).
func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		h.writeError(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	user := domain.User{
		Username: req.Username,
		Password: req.Password,
	}

	authToken, err := h.controllers.UserController.Login(user)
	if err != nil {
		h.writeErr(w, err)
		return
	}

	h.writeJSON(w, authToken, http.StatusOK)
}

func (h *Handlers) Me(w http.ResponseWriter, r *http.Request, user domain.User) {
	full, err := h.controllers.UserController.GetByID(user.ID)
	if err != nil {
		// Fall back to JWT claims if the user row is temporarily unavailable.
		h.writeJSON(w, MeResponse{ID: user.ID.String(), Username: user.Username}, http.StatusOK)
		return
	}
	h.writeJSON(w, MeResponse{
		ID:       full.ID.String(),
		Username: full.Username,
		Email:    full.Email,
		IsBot:    full.IsBot,
	}, http.StatusOK)
}
