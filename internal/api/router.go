package api

import (
	"net/http"

	"github.com/dakaii/vibegopher/internal/controller"
	"github.com/gorilla/mux"
)

// SetupRouter creates and configures the HTTP router
func SetupRouter(controllers *controller.Controllers) *mux.Router {
	handlers := NewHandlers(controllers)

	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Auth routes
	api.HandleFunc("/signup", handlers.Signup).Methods("POST")
	api.HandleFunc("/login", handlers.Login).Methods("POST")
	api.HandleFunc("/me", handlers.Me).Methods("GET")

	// Post routes
	api.HandleFunc("/posts", handlers.CreatePost).Methods("POST")
	api.HandleFunc("/posts", handlers.GetAllPosts).Methods("GET")
	api.HandleFunc("/posts/{id}", handlers.GetPostByID).Methods("GET")
	api.HandleFunc("/posts/user/{userId}", handlers.GetPostsByUserID).Methods("GET")
	api.HandleFunc("/posts/{id}", handlers.UpdatePost).Methods("PATCH")
	api.HandleFunc("/posts/{id}", handlers.DeletePost).Methods("DELETE")

	// Comment routes
	api.HandleFunc("/comments", handlers.CreateComment).Methods("POST")
	api.HandleFunc("/comments/post/{postId}", handlers.GetCommentsByPostID).Methods("GET")
	api.HandleFunc("/comments/user/{userId}", handlers.GetCommentsByUserID).Methods("GET")
	api.HandleFunc("/comments/{id}", handlers.GetCommentByID).Methods("GET")
	api.HandleFunc("/comments/{id}", handlers.UpdateComment).Methods("PATCH")
	api.HandleFunc("/comments/{id}", handlers.DeleteComment).Methods("DELETE")

	// Add CORS middleware
	r.Use(corsMiddleware)

	return r
}

// corsMiddleware adds CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
