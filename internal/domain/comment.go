package domain

import (
	"time"

	"github.com/google/uuid"
)

// Comment represents a comment on a post
type Comment struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	PostID    uuid.UUID `json:"post_id" gorm:"type:uuid;not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Post      Post      `json:"post,omitempty" gorm:"foreignKey:PostID;references:ID"`
}

// CreateCommentRequest represents the request structure for creating a comment
type CreateCommentRequest struct {
	Content string    `json:"content" validate:"required,min=1,max=280"`
	PostID  uuid.UUID `json:"post_id" validate:"required"`
}

// UpdateCommentRequest represents the request structure for updating a comment
type UpdateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=280"`
}
