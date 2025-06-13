package comment

import (
	"errors"

	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/google/uuid"
)

// UpdateComment updates an existing comment
func (c *Controller) UpdateComment(id uuid.UUID, updates domain.Comment) (*domain.Comment, error) {
	if !isValidContent(updates.Content) {
		return nil, errors.New("content cannot be empty and must be 280 characters or less")
	}

	updatedComment, err := c.service.UpdateComment(id, updates)
	if err != nil {
		return nil, err
	}

	return updatedComment, nil
}

// DeleteComment deletes a comment by ID
func (c *Controller) DeleteComment(id uuid.UUID) error {
	err := c.service.DeleteComment(id)
	if err != nil {
		return err
	}
	return nil
}
