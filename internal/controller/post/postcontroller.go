package post

import (
	"time"

	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/dakaii/vibegopher/internal/repository/postrepo"
	"github.com/google/uuid"
)

// Repository interface allows us to easily swap out the actual implementation
type repository interface {
	CreatePost(post domain.Post) (*domain.Post, error)
	GetPostByID(id uuid.UUID) (*domain.Post, error)
	GetAllPosts() ([]domain.Post, error)
	ListPosts(limit int, beforeCreatedAt time.Time, beforeID uuid.UUID) ([]domain.Post, error)
	GetPostsByUserID(userID uuid.UUID) ([]domain.Post, error)
	UpdatePost(id uuid.UUID, updates domain.Post) (*domain.Post, error)
	DeletePost(id uuid.UUID) error
}

// Controller contains the service as an injectable dependency
type Controller struct {
	service repository
}

// InitController initializes the post controller
func InitController(postRepo *postrepo.PostRepo) *Controller {
	return &Controller{
		service: postRepo,
	}
}
