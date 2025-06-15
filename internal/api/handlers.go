package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/dakaii/vibegopher/internal/auth"
	"github.com/dakaii/vibegopher/internal/controller"
	"github.com/dakaii/vibegopher/internal/domain"
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

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// withAuth is a middleware that extracts and validates the user from JWT token
func (h *Handlers) withAuth(fn func(w http.ResponseWriter, r *http.Request, user domain.User)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := h.getCurrentUser(r)
		if err != nil {
			h.writeError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fn(w, r, user)
	}
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
