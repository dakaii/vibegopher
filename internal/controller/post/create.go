package post

import (
	"errors"

	"github.com/dakaii/vibegopher/internal/domain"
)

// CreatePost creates a new post
func (c *Controller) CreatePost(post domain.Post) (*domain.Post, error) {
	if !isValidContent(post.Content) {
		return nil, errors.New("content cannot be empty and must be 280 characters or less")
	}

	if post.UserID.String() == "00000000-0000-0000-0000-000000000000" {
		return nil, errors.New("user ID is required")
	}

	createdPost, err := c.service.CreatePost(post)
	if err != nil {
		return nil, err
	}

	return createdPost, nil
}

func isValidContent(content string) bool {
	return len(content) > 0 && len(content) <= 280
}
