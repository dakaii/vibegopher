package domain

import (
	"time"

	"github.com/google/uuid"
)

// Comment is a reply on a post.
type Comment struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Content   string    `json:"content"`
	UserID    uuid.UUID `json:"user_id"`
	PostID    uuid.UUID `json:"post_id"`
	User      User      `json:"user"`
	Post      Post      `json:"post,omitempty"`
}

type CreateCommentRequest struct {
	Content string    `json:"content"`
	PostID  uuid.UUID `json:"post_id"`
}

type UpdateCommentRequest struct {
	Content string `json:"content"`
}
