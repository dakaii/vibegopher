package api

import (
	"net/http"
	"strings"

	"github.com/dakaii/vibegopher/internal/controller"
	"github.com/dakaii/vibegopher/internal/envvar"
	"github.com/gorilla/mux"
)

func SetupRouter(controllers *controller.Controllers) *mux.Router {
	handlers := NewHandlers(controllers)

	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/auth/google", handlers.GoogleAuth).Methods("POST")
	if envvar.PasswordAuthEnabled() {
		api.HandleFunc("/signup", handlers.Signup).Methods("POST")
		api.HandleFunc("/login", handlers.Login).Methods("POST")
	}
	api.HandleFunc("/me", handlers.withAuth(handlers.Me)).Methods("GET")

	api.HandleFunc("/posts", handlers.withAuth(handlers.CreatePost)).Methods("POST")
	api.HandleFunc("/posts", handlers.GetAllPosts).Methods("GET")
	api.HandleFunc("/posts/{id}", handlers.GetPostByID).Methods("GET")
	api.HandleFunc("/posts/user/{userId}", handlers.GetPostsByUserID).Methods("GET")
	api.HandleFunc("/posts/{id}", handlers.withAuth(handlers.UpdatePost)).Methods("PATCH")
	api.HandleFunc("/posts/{id}", handlers.withAuth(handlers.DeletePost)).Methods("DELETE")

	api.HandleFunc("/comments", handlers.withAuth(handlers.CreateComment)).Methods("POST")
	api.HandleFunc("/comments/post/{postId}", handlers.GetCommentsByPostID).Methods("GET")
	api.HandleFunc("/comments/user/{userId}", handlers.GetCommentsByUserID).Methods("GET")
	api.HandleFunc("/comments/{id}", handlers.GetCommentByID).Methods("GET")
	api.HandleFunc("/comments/{id}", handlers.withAuth(handlers.UpdateComment)).Methods("PATCH")
	api.HandleFunc("/comments/{id}", handlers.withAuth(handlers.DeleteComment)).Methods("DELETE")

	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}).Methods("GET")

	r.Use(corsMiddleware)
	return r
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originCfg := envvar.CORSOrigin()
		if originCfg == "" {
			originCfg = "*"
		}
		reqOrigin := r.Header.Get("Origin")
		if originCfg == "*" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else if reqOrigin != "" {
			for _, allowed := range strings.Split(originCfg, ",") {
				if strings.TrimSpace(allowed) == reqOrigin {
					w.Header().Set("Access-Control-Allow-Origin", reqOrigin)
					w.Header().Set("Vary", "Origin")
					break
				}
			}
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
