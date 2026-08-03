package domain

import (
	"time"

	"github.com/google/uuid"
)

// BotUserID is the fixed UUID for the vibe_critic AI account (seeded by migration).
var BotUserID = uuid.MustParse("a0000000-0000-4000-8000-000000000001")

const BotUsername = "vibe_critic"

// User is the public user model. Password is never serialized to JSON/JWT.
type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Username    string    `json:"username"`
	Email       string    `json:"email,omitempty"`
	ClerkUserID string    `json:"-"`
	GoogleSub   string    `json:"-"` // legacy; unused by auth
	Password    string    `json:"-"`
	IsBot       bool      `json:"is_bot"`
}
