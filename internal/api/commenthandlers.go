package api

import (
	"encoding/json"
	"net/http"

	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// CreateComment handles creating a new comment
func (h *Handlers) CreateComment(w http.ResponseWriter, r *http.Request, user domain.User) {
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

	h.controllers.EnqueueBotJob(domain.BotJobKindCommentCreated, createdComment.ID, user.ID)
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
func (h *Handlers) UpdateComment(w http.ResponseWriter, r *http.Request, user domain.User) {
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
func (h *Handlers) DeleteComment(w http.ResponseWriter, r *http.Request, user domain.User) {
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
