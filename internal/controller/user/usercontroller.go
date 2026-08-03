package user

import (
	"github.com/dakaii/vibegopher/internal/auth"
	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/dakaii/vibegopher/internal/repository/userrepo"
	"github.com/google/uuid"
)

type repository interface {
	GetExistingUser(username string) (*domain.User, error)
	GetByClerkUserID(clerkUserID string) (*domain.User, error)
	GetByID(id uuid.UUID) (*domain.User, error)
	CreateUser(user domain.User) (*domain.User, error)
	UpsertClerkUser(profile domain.ClerkProfile) (*domain.User, error)
}

type Controller struct {
	service         repository
	clerkValidator  auth.ClerkTokenValidator
}

func InitController(userRepo *userrepo.UserRepo) *Controller {
	return &Controller{
		service:        userRepo,
		clerkValidator: auth.NewClerkTokenValidator(),
	}
}

// WithClerkValidator overrides the Clerk token validator (tests).
func (c *Controller) WithClerkValidator(v auth.ClerkTokenValidator) *Controller {
	c.clerkValidator = v
	return c
}
