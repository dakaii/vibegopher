package botjobrepo

import (
	"errors"
	"time"

	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BotJobRepo struct {
	db *gorm.DB
}

func NewBotJobRepo(db *gorm.DB) *BotJobRepo {
	return &BotJobRepo{db: db}
}

func (r *BotJobRepo) Enqueue(kind string, targetID uuid.UUID) (*domain.BotJob, error) {
	entity := BotJobEntity{
		Kind:        kind,
		TargetID:    targetID,
		Status:      domain.BotJobStatusPending,
		AvailableAt: time.Now().UTC(),
	}
	if err := r.db.Create(&entity).Error; err != nil {
		return nil, err
	}
	return toDomain(entity), nil
}

// ClaimNext claims the next available pending job using FOR UPDATE SKIP LOCKED.
func (r *BotJobRepo) ClaimNext() (*domain.BotJob, error) {
	var entity BotJobEntity
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("status = ? AND available_at <= ?", domain.BotJobStatusPending, time.Now().UTC()).
			Order("available_at ASC").
			First(&entity).Error
		if err != nil {
			return err
		}
		entity.Status = domain.BotJobStatusProcessing
		entity.Attempts++
		entity.UpdatedAt = time.Now().UTC()
		return tx.Save(&entity).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomain(entity), nil
}

func (r *BotJobRepo) MarkDone(id uuid.UUID) error {
	return r.db.Model(&BotJobEntity{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     domain.BotJobStatusDone,
		"updated_at": time.Now().UTC(),
		"last_error": "",
	}).Error
}

func (r *BotJobRepo) MarkFailed(id uuid.UUID, jobErr string, retry bool) error {
	status := domain.BotJobStatusFailed
	availableAt := time.Now().UTC()
	if retry {
		status = domain.BotJobStatusPending
		availableAt = time.Now().UTC().Add(30 * time.Second)
	}
	return r.db.Model(&BotJobEntity{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       status,
		"last_error":   jobErr,
		"available_at": availableAt,
		"updated_at":   time.Now().UTC(),
	}).Error
}

type BotJobEntity struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Kind        string    `gorm:"type:varchar(32);not null"`
	TargetID    uuid.UUID `gorm:"type:uuid;not null"`
	Status      string    `gorm:"type:varchar(32);not null"`
	Attempts    int       `gorm:"not null;default:0"`
	LastError   string    `gorm:"type:text"`
	AvailableAt time.Time `gorm:"not null"`
}

func (BotJobEntity) TableName() string { return "bot_jobs" }

func (e *BotJobEntity) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

func toDomain(e BotJobEntity) *domain.BotJob {
	return &domain.BotJob{
		ID:          e.ID,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		Kind:        e.Kind,
		TargetID:    e.TargetID,
		Status:      e.Status,
		Attempts:    e.Attempts,
		LastError:   e.LastError,
		AvailableAt: e.AvailableAt,
	}
}
