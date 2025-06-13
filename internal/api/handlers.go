package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/dakaii/vibegopher/internal/auth"
	"github.com/dakaii/vibegopher/internal/controller"
	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

// Post handlers

// CreatePost handles creating a new post
func (h *Handlers) CreatePost(w http.ResponseWriter, r *http.Request) {
	user, err := h.getCurrentUser(r)
	if err != nil {
		h.writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req domain.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	post := domain.Post{
		Content: req.Content,
		UserID:  user.ID,
	}

	createdPost, err := h.controllers.PostController.CreatePost(post)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.writeJSON(w, createdPost, http.StatusCreated)
}

// GetAllPosts handles retrieving all posts
func (h *Handlers) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.controllers.PostController.GetAllPosts()
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, posts, http.StatusOK)
}

// GetPostByID handles retrieving a post by ID
func (h *Handlers) GetPostByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.writeError(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.controllers.PostController.GetPostByID(id)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	h.writeJSON(w, post, http.StatusOK)
}

// GetPostsByUserID handles retrieving posts by user ID
func (h *Handlers) GetPostsByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userId"]

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.writeError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	posts, err := h.controllers.PostController.GetPostsByUserID(userID)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, posts, http.StatusOK)
}

// UpdatePost handles updating a post
func (h *Handlers) UpdatePost(w http.ResponseWriter, r *http.Request) {
	user, err := h.getCurrentUser(r)
	if err != nil {
		h.writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.writeError(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Check if post belongs to the current user
	existingPost, err := h.controllers.PostController.GetPostByID(id)
	if err != nil {
		h.writeError(w, "Post not found", http.StatusNotFound)
		return
	}

	if existingPost.UserID != user.ID {
		h.writeError(w, "You can only update your own posts", http.StatusForbidden)
		return
	}

	var req domain.UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updates := domain.Post{
		Content: req.Content,
	}

	updatedPost, err := h.controllers.PostController.UpdatePost(id, updates)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.writeJSON(w, updatedPost, http.StatusOK)
}

// DeletePost handles deleting a post
func (h *Handlers) DeletePost(w http.ResponseWriter, r *http.Request) {
	user, err := h.getCurrentUser(r)
	if err != nil {
		h.writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.writeError(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Check if post belongs to the current user
	existingPost, err := h.controllers.PostController.GetPostByID(id)
	if err != nil {
		h.writeError(w, "Post not found", http.StatusNotFound)
		return
	}

	if existingPost.UserID != user.ID {
		h.writeError(w, "You can only delete your own posts", http.StatusForbidden)
		return
	}

	err = h.controllers.PostController.DeletePost(id)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Comment handlers

// CreateComment handles creating a new comment
func (h *Handlers) CreateComment(w http.ResponseWriter, r *http.Request) {
	user, err := h.getCurrentUser(r)
	if err != nil {
		h.writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req domain.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	comment := domain.Comment{
		Content: req.Content,
		UserID:  user.ID,
		PostID:  req.PostID,
	}

	createdComment, err := h.controllers.CommentController.CreateComment(comment)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.writeJSON(w, createdComment, http.StatusCreated)
}

// GetCommentsByPostID handles retrieving comments by post ID
func (h *Handlers) GetCommentsByPostID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postIDStr := vars["postId"]

	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		h.writeError(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	comments, err := h.controllers.CommentController.GetCommentsByPostID(postID)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, comments, http.StatusOK)
}

// GetCommentsByUserID handles retrieving comments by user ID
func (h *Handlers) GetCommentsByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userId"]

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.writeError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	comments, err := h.controllers.CommentController.GetCommentsByUserID(userID)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, comments, http.StatusOK)
}

// GetCommentByID handles retrieving a comment by ID
func (h *Handlers) GetCommentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.writeError(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	comment, err := h.controllers.CommentController.GetCommentByID(id)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	h.writeJSON(w, comment, http.StatusOK)
}

// UpdateComment handles updating a comment
func (h *Handlers) UpdateComment(w http.ResponseWriter, r *http.Request) {
	user, err := h.getCurrentUser(r)
	if err != nil {
		h.writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.writeError(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// Check if comment belongs to the current user
	existingComment, err := h.controllers.CommentController.GetCommentByID(id)
	if err != nil {
		h.writeError(w, "Comment not found", http.StatusNotFound)
		return
	}

	if existingComment.UserID != user.ID {
		h.writeError(w, "You can only update your own comments", http.StatusForbidden)
		return
	}

	var req domain.UpdateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updates := domain.Comment{
		Content: req.Content,
	}

	updatedComment, err := h.controllers.CommentController.UpdateComment(id, updates)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.writeJSON(w, updatedComment, http.StatusOK)
}

// DeleteComment handles deleting a comment
func (h *Handlers) DeleteComment(w http.ResponseWriter, r *http.Request) {
	user, err := h.getCurrentUser(r)
	if err != nil {
		h.writeError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.writeError(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// Check if comment belongs to the current user
	existingComment, err := h.controllers.CommentController.GetCommentByID(id)
	if err != nil {
		h.writeError(w, "Comment not found", http.StatusNotFound)
		return
	}

	if existingComment.UserID != user.ID {
		h.writeError(w, "You can only delete your own comments", http.StatusForbidden)
		return
	}

	err = h.controllers.CommentController.DeleteComment(id)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
