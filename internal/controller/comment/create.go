package comment

import (
	"errors"

	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/google/uuid"
)

// CreateComment creates a new comment
func (c *Controller) CreateComment(comment domain.Comment) (*domain.Comment, error) {
	if !isValidContent(comment.Content) {
		return nil, errors.New("content cannot be empty and must be 280 characters or less")
	}

	if comment.UserID == uuid.Nil {
		return nil, errors.New("user ID is required")
	}

	if comment.PostID == uuid.Nil {
		return nil, errors.New("post ID is required")
	}

	createdComment, err := c.service.CreateComment(comment)
	if err != nil {
		return nil, err
	}

	return createdComment, nil
}

func isValidContent(content string) bool {
	return len(content) > 0 && len(content) <= 280
}
