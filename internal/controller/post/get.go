package post

import (
	"time"

	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/google/uuid"
)

// GetPostByID retrieves a post by its ID
func (c *Controller) GetPostByID(id uuid.UUID) (*domain.Post, error) {
	post, err := c.service.GetPostByID(id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// GetAllPosts retrieves all posts
func (c *Controller) GetAllPosts() ([]domain.Post, error) {
	return c.ListPosts(100, time.Time{}, uuid.Nil)
}

// ListPosts retrieves a page of posts (newest first) with an optional cursor.
func (c *Controller) ListPosts(limit int, beforeCreatedAt time.Time, beforeID uuid.UUID) ([]domain.Post, error) {
	return c.service.ListPosts(limit, beforeCreatedAt, beforeID)
}

// GetPostsByUserID retrieves all posts by a specific user
func (c *Controller) GetPostsByUserID(userID uuid.UUID) ([]domain.Post, error) {
	posts, err := c.service.GetPostsByUserID(userID)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
