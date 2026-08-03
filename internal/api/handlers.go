package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/dakaii/vibegopher/internal/apperr"
	"github.com/dakaii/vibegopher/internal/auth"
	"github.com/dakaii/vibegopher/internal/controller"
	"github.com/dakaii/vibegopher/internal/domain"
)

type Handlers struct {
	controllers *controller.Controllers
}

func NewHandlers(controllers *controller.Controllers) *Handlers {
	return &Handlers{controllers: controllers}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

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

func (h *Handlers) getCurrentUser(r *http.Request) (domain.User, error) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return domain.User{}, auth.ErrUnauthorized
	}
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authorization, bearerPrefix) {
		return domain.User{}, auth.ErrUnauthorized
	}
	token := strings.TrimPrefix(authorization, bearerPrefix)
	user, err := auth.VerifyJWT(token)
	if err != nil {
		log.Printf("Failed to verify JWT: %v", err)
		return domain.User{}, auth.ErrUnauthorized
	}
	return user, nil
}

func (h *Handlers) writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("encode json: %v", err)
	}
}

func (h *Handlers) writeError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(ErrorResponse{Error: message}); err != nil {
		log.Printf("encode error json: %v", err)
	}
}

func (h *Handlers) writeErr(w http.ResponseWriter, err error) {
	status, message := apperr.HTTPStatus(err)
	if status >= 500 {
		log.Printf("request error: %v", err)
	}
	h.writeError(w, message, status)
}
