package commentrepo

import (
	"errors"
	"fmt"
	"time"

	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentRepo struct {
	db *gorm.DB
}

func NewCommentRepo(db *gorm.DB) *CommentRepo {
	return &CommentRepo{db: db}
}

func (repo *CommentRepo) CreateComment(comment domain.Comment) (*domain.Comment, error) {
	dbComment := CommentEntity{
		Content: comment.Content,
		UserID:  comment.UserID,
		PostID:  comment.PostID,
	}
	if err := repo.db.Create(&dbComment).Error; err != nil {
		return nil, err
	}
	return repo.GetCommentByID(dbComment.ID)
}

func (repo *CommentRepo) GetCommentByID(id uuid.UUID) (*domain.Comment, error) {
	var comment CommentEntity
	err := repo.db.Preload("User").Where("id = ?", id).First(&comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no comment found with ID: %s", id)
		}
		return nil, err
	}
	out := toDomain(comment)
	return &out, nil
}

func (repo *CommentRepo) GetCommentsByPostID(postID uuid.UUID) ([]domain.Comment, error) {
	var comments []CommentEntity
	if err := repo.db.Preload("User").Where("post_id = ?", postID).Order("created_at ASC").Limit(200).Find(&comments).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Comment, len(comments))
	for i, comment := range comments {
		out[i] = toDomain(comment)
	}
	return out, nil
}

func (repo *CommentRepo) GetCommentsByUserID(userID uuid.UUID) ([]domain.Comment, error) {
	var comments []CommentEntity
	if err := repo.db.Preload("User").Where("user_id = ?", userID).Order("created_at DESC").Limit(100).Find(&comments).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Comment, len(comments))
	for i, comment := range comments {
		out[i] = toDomain(comment)
	}
	return out, nil
}

func (repo *CommentRepo) UpdateComment(id uuid.UUID, updates domain.Comment) (*domain.Comment, error) {
	result := repo.db.Model(&CommentEntity{}).Where("id = ?", id).Updates(map[string]interface{}{
		"content":    updates.Content,
		"updated_at": time.Now(),
	})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("no comment found with ID: %s", id)
	}
	return repo.GetCommentByID(id)
}

func (repo *CommentRepo) DeleteComment(id uuid.UUID) error {
	result := repo.db.Delete(&CommentEntity{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no comment found with ID: %s", id)
	}
	return nil
}

type CommentEntity struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Content   string     `gorm:"type:text;not null"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null"`
	PostID    uuid.UUID  `gorm:"type:uuid;not null"`
	User      UserEntity `gorm:"foreignKey:UserID;references:ID"`
}

type UserEntity struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string `gorm:"unique;index;not null"`
	IsBot     bool   `gorm:"column:is_bot;not null;default:false"`
}

func (comment *CommentEntity) BeforeCreate(tx *gorm.DB) error {
	if comment.ID == uuid.Nil {
		comment.ID = uuid.New()
	}
	return nil
}

func (CommentEntity) TableName() string { return "comments" }
func (UserEntity) TableName() string    { return "users" }

func toDomain(comment CommentEntity) domain.Comment {
	return domain.Comment{
		ID:        comment.ID,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		Content:   comment.Content,
		UserID:    comment.UserID,
		PostID:    comment.PostID,
		User: domain.User{
			ID:        comment.User.ID,
			Username:  comment.User.Username,
			CreatedAt: comment.User.CreatedAt,
			UpdatedAt: comment.User.UpdatedAt,
			IsBot:     comment.User.IsBot,
		},
	}
}
