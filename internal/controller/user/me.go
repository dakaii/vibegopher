package user

import (
	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/google/uuid"
)

func (c *Controller) GetByID(id uuid.UUID) (*domain.User, error) {
	return c.service.GetByID(id)
}
