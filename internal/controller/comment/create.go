package comment

import (
	"errors"
	"strings"

	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/google/uuid"
)

func (c *Controller) CreateComment(comment domain.Comment) (*domain.Comment, error) {
	comment.Content = strings.TrimSpace(comment.Content)
	if !isValidContent(comment.Content) {
		return nil, errors.New("content cannot be empty and must be 280 characters or less")
	}
	if comment.UserID == uuid.Nil {
		return nil, errors.New("user ID is required")
	}
	if comment.PostID == uuid.Nil {
		return nil, errors.New("post ID is required")
	}
	return c.service.CreateComment(comment)
}

func isValidContent(content string) bool {
	return len(content) > 0 && len(content) <= 280
}
