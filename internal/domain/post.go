package domain

import (
	"time"

	"github.com/google/uuid"
)

// Post represents a post in the social media platform
type Post struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Comments  []Comment `json:"comments,omitempty" gorm:"foreignKey:PostID;references:ID"`
}

// CreatePostRequest represents the request structure for creating a post
type CreatePostRequest struct {
	Content string `json:"content" validate:"required,min=1,max=280"`
}

// UpdatePostRequest represents the request structure for updating a post
type UpdatePostRequest struct {
	Content string `json:"content" validate:"required,min=1,max=280"`
}
