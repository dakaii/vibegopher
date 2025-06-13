package post

import (
	"errors"

	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/google/uuid"
)

// UpdatePost updates an existing post
func (c *Controller) UpdatePost(id uuid.UUID, updates domain.Post) (*domain.Post, error) {
	if !isValidContent(updates.Content) {
		return nil, errors.New("content cannot be empty and must be 280 characters or less")
	}

	updatedPost, err := c.service.UpdatePost(id, updates)
	if err != nil {
		return nil, err
	}

	return updatedPost, nil
}

// DeletePost deletes a post by ID
func (c *Controller) DeletePost(id uuid.UUID) error {
	err := c.service.DeletePost(id)
	if err != nil {
		return err
	}
	return nil
}
