package api

import (
	"encoding/json"
	"net/http"

	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// CreatePost handles creating a new post
func (h *Handlers) CreatePost(w http.ResponseWriter, r *http.Request, user domain.User) {
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
		h.writeErr(w, err)
		return
	}

	h.controllers.EnqueueBotJob(domain.BotJobKindPostCreated, createdPost.ID, user.ID)
	h.writeJSON(w, createdPost, http.StatusCreated)
}

// GetAllPosts handles retrieving all posts
func (h *Handlers) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.controllers.PostController.GetAllPosts()
	if err != nil {
		h.writeErr(w, err)
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
		h.writeErr(w, err)
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
		h.writeErr(w, err)
		return
	}

	h.writeJSON(w, posts, http.StatusOK)
}

// UpdatePost handles updating a post
func (h *Handlers) UpdatePost(w http.ResponseWriter, r *http.Request, user domain.User) {
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
		h.writeErr(w, err)
		return
	}

	h.writeJSON(w, updatedPost, http.StatusOK)
}

// DeletePost handles deleting a post
func (h *Handlers) DeletePost(w http.ResponseWriter, r *http.Request, user domain.User) {
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
		h.writeErr(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
