package comment

import (
	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/dakaii/vibegopher/internal/repository/commentrepo"
	"github.com/google/uuid"
)

// Repository interface allows us to easily swap out the actual implementation
type repository interface {
	CreateComment(comment domain.Comment) (*domain.Comment, error)
	GetCommentByID(id uuid.UUID) (*domain.Comment, error)
	GetCommentsByPostID(postID uuid.UUID) ([]domain.Comment, error)
	GetCommentsByUserID(userID uuid.UUID) ([]domain.Comment, error)
	UpdateComment(id uuid.UUID, updates domain.Comment) (*domain.Comment, error)
	DeleteComment(id uuid.UUID) error
}

// Controller contains the service as an injectable dependency
type Controller struct {
	service repository
}

// InitController initializes the comment controller
func InitController(commentRepo *commentrepo.CommentRepo) *Controller {
	return &Controller{
		service: commentRepo,
	}
}
