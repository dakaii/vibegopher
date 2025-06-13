package comment

import (
	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/google/uuid"
)

// GetCommentByID retrieves a comment by its ID
func (c *Controller) GetCommentByID(id uuid.UUID) (*domain.Comment, error) {
	comment, err := c.service.GetCommentByID(id)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// GetCommentsByPostID retrieves all comments for a specific post
func (c *Controller) GetCommentsByPostID(postID uuid.UUID) ([]domain.Comment, error) {
	comments, err := c.service.GetCommentsByPostID(postID)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// GetCommentsByUserID retrieves all comments by a specific user
func (c *Controller) GetCommentsByUserID(userID uuid.UUID) ([]domain.Comment, error) {
	comments, err := c.service.GetCommentsByUserID(userID)
	if err != nil {
		return nil, err
	}
	return comments, nil
}
