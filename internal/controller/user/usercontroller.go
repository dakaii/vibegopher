package user

import (
	"github.com/dakaii/vibegopher/internal/auth"
	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/dakaii/vibegopher/internal/repository/userrepo"
	"github.com/google/uuid"
)

type repository interface {
	GetExistingUser(username string) (*domain.User, error)
	GetByGoogleSub(googleSub string) (*domain.User, error)
	GetByID(id uuid.UUID) (*domain.User, error)
	CreateUser(user domain.User) (*domain.User, error)
	UpsertGoogleUser(profile domain.GoogleProfile) (*domain.User, error)
}

type Controller struct {
	service         repository
	googleValidator auth.GoogleTokenValidator
}

func InitController(userRepo *userrepo.UserRepo) *Controller {
	return &Controller{
		service:         userRepo,
		googleValidator: auth.NewGoogleTokenValidator(),
	}
}

// WithGoogleValidator overrides the Google token validator (tests).
func (c *Controller) WithGoogleValidator(v auth.GoogleTokenValidator) *Controller {
	c.googleValidator = v
	return c
}
