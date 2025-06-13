package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/dakaii/graphyy/internal/auth"
	"github.com/dakaii/graphyy/internal/controller"
	"github.com/dakaii/graphyy/internal/domain"
)

// Handlers contains all HTTP handlers
type Handlers struct {
	controllers *controller.Controllers
}

// NewHandlers creates a new Handlers instance
func NewHandlers(controllers *controller.Controllers) *Handlers {
	return &Handlers{
		controllers: controllers,
	}
}

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

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
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
func (h *Handlers) Me(w http.ResponseWriter, r *http.Request) {
	user, err := h.getCurrentUser(r)
	if err != nil {
		h.writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response := MeResponse{
		Username: user.Username,
	}

	h.writeJSON(w, response, http.StatusOK)
}

// getCurrentUser extracts and validates the current user from JWT token
func (h *Handlers) getCurrentUser(r *http.Request) (domain.User, error) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return domain.User{}, auth.ErrUnauthorized
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authorization, bearerPrefix) {
		log.Println("Invalid authorization format")
		return domain.User{}, auth.ErrUnauthorized
	}

	token := strings.TrimPrefix(authorization, bearerPrefix)
	user, err := auth.VerifyJWT(token)
	if err != nil {
		log.Printf("Failed to verify JWT: %v", err)
		return domain.User{}, err
	}

	return user, nil
}

// writeJSON writes a JSON response
func (h *Handlers) writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// writeError writes an error response
func (h *Handlers) writeError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
