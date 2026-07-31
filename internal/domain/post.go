package domain

import (
	"time"

	"github.com/google/uuid"
)

// Post is a feed item.
type Post struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Content   string    `json:"content"`
	UserID    uuid.UUID `json:"user_id"`
	User      User      `json:"user"`
	Comments  []Comment `json:"comments,omitempty"`
}

type CreatePostRequest struct {
	Content string `json:"content"`
}

type UpdatePostRequest struct {
	Content string `json:"content"`
}
