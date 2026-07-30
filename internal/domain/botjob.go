package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	BotJobKindPostCreated    = "post_created"
	BotJobKindCommentCreated = "comment_created"

	BotJobStatusPending    = "pending"
	BotJobStatusProcessing = "processing"
	BotJobStatusDone       = "done"
	BotJobStatusFailed     = "failed"
)

// BotJob is an async unit of work for the AI critic worker.
type BotJob struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Kind        string
	TargetID    uuid.UUID
	Status      string
	Attempts    int
	LastError   string
	AvailableAt time.Time
}
