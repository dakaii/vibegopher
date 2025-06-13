package commentrepo

import (
	"errors"
	"fmt"
	"time"

	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CommentRepo handles comment database operations
type CommentRepo struct {
	db *gorm.DB
}

// NewCommentRepo creates a new comment repository
func NewCommentRepo(db *gorm.DB) *CommentRepo {
	return &CommentRepo{
		db: db,
	}
}

// CreateComment creates a new comment in the database
func (repo *CommentRepo) CreateComment(comment domain.Comment) (*domain.Comment, error) {
	dbComment := CommentEntity{
		Content: comment.Content,
		UserID:  comment.UserID,
		PostID:  comment.PostID,
	}

	result := repo.db.Create(&dbComment)
	if result.Error != nil {
		return nil, result.Error
	}

	return repo.GetCommentByID(dbComment.ID)
}

// GetCommentByID fetches a comment by ID with user information
func (repo *CommentRepo) GetCommentByID(id uuid.UUID) (*domain.Comment, error) {
	var comment CommentEntity
	result := repo.db.Preload("User").Where("id = ?", id).First(&comment)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no comment found with ID: %s", id)
		}
		return nil, result.Error
	}

	return &domain.Comment{
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
		},
	}, nil
}

// GetCommentsByPostID fetches all comments for a specific post
func (repo *CommentRepo) GetCommentsByPostID(postID uuid.UUID) ([]domain.Comment, error) {
	var comments []CommentEntity
	result := repo.db.Preload("User").Where("post_id = ?", postID).Order("created_at ASC").Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}

	domainComments := make([]domain.Comment, len(comments))
	for i, comment := range comments {
		domainComments[i] = domain.Comment{
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
			},
		}
	}

	return domainComments, nil
}

// GetCommentsByUserID fetches all comments by a specific user
func (repo *CommentRepo) GetCommentsByUserID(userID uuid.UUID) ([]domain.Comment, error) {
	var comments []CommentEntity
	result := repo.db.Preload("User").Where("user_id = ?", userID).Order("created_at DESC").Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}

	domainComments := make([]domain.Comment, len(comments))
	for i, comment := range comments {
		domainComments[i] = domain.Comment{
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
			},
		}
	}

	return domainComments, nil
}

// UpdateComment updates an existing comment
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

// DeleteComment deletes a comment by ID
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

// CommentEntity represents the comment entity in the database
type CommentEntity struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Content   string         `gorm:"type:text;not null"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null"`
	PostID    uuid.UUID      `gorm:"type:uuid;not null"`
	User      UserEntity     `gorm:"foreignKey:UserID;references:ID"`
}

// UserEntity embedded for preloading
type UserEntity struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string `gorm:"unique;index;not null"`
}

// BeforeCreate sets a UUID for the comment
func (comment *CommentEntity) BeforeCreate(tx *gorm.DB) (err error) {
	comment.ID = uuid.New()
	return
}

// TableName overrides the table name
func (comment *CommentEntity) TableName() string {
	return "comments"
}

// TableName for UserEntity when used in comment context
func (user *UserEntity) TableName() string {
	return "users"
}
