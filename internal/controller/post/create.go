package post

import (
	"errors"
	"strings"

	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/google/uuid"
)

func (c *Controller) CreatePost(post domain.Post) (*domain.Post, error) {
	post.Content = strings.TrimSpace(post.Content)
	if !isValidContent(post.Content) {
		return nil, errors.New("content cannot be empty and must be 280 characters or less")
	}
	if post.UserID == uuid.Nil {
		return nil, errors.New("user ID is required")
	}
	return c.service.CreatePost(post)
}

func isValidContent(content string) bool {
	return len(content) > 0 && len(content) <= 280
}
