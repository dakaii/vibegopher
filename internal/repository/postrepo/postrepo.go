package postrepo

import (
	"errors"
	"fmt"
	"time"

	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PostRepo handles post database operations
type PostRepo struct {
	db *gorm.DB
}

// NewPostRepo creates a new post repository
func NewPostRepo(db *gorm.DB) *PostRepo {
	return &PostRepo{
		db: db,
	}
}

// CreatePost creates a new post in the database
func (repo *PostRepo) CreatePost(post domain.Post) (*domain.Post, error) {
	dbPost := PostEntity{
		Content: post.Content,
		UserID:  post.UserID,
	}

	result := repo.db.Create(&dbPost)
	if result.Error != nil {
		return nil, result.Error
	}

	return repo.GetPostByID(dbPost.ID)
}

// GetPostByID fetches a post by ID with user information
func (repo *PostRepo) GetPostByID(id uuid.UUID) (*domain.Post, error) {
	var post PostEntity
	result := repo.db.Preload("User").Where("id = ?", id).First(&post)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no post found with ID: %s", id)
		}
		return nil, result.Error
	}

	return &domain.Post{
		ID:        post.ID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Content:   post.Content,
		UserID:    post.UserID,
		User: domain.User{
			ID:        post.User.ID,
			Username:  post.User.Username,
			CreatedAt: post.User.CreatedAt,
			UpdatedAt: post.User.UpdatedAt,
		},
	}, nil
}

// GetAllPosts fetches all posts with user information, ordered by creation date
func (repo *PostRepo) GetAllPosts() ([]domain.Post, error) {
	var posts []PostEntity
	result := repo.db.Preload("User").Order("created_at DESC").Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}

	domainPosts := make([]domain.Post, len(posts))
	for i, post := range posts {
		domainPosts[i] = domain.Post{
			ID:        post.ID,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
			Content:   post.Content,
			UserID:    post.UserID,
			User: domain.User{
				ID:        post.User.ID,
				Username:  post.User.Username,
				CreatedAt: post.User.CreatedAt,
				UpdatedAt: post.User.UpdatedAt,
			},
		}
	}

	return domainPosts, nil
}

// GetPostsByUserID fetches all posts by a specific user
func (repo *PostRepo) GetPostsByUserID(userID uuid.UUID) ([]domain.Post, error) {
	var posts []PostEntity
	result := repo.db.Preload("User").Where("user_id = ?", userID).Order("created_at DESC").Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}

	domainPosts := make([]domain.Post, len(posts))
	for i, post := range posts {
		domainPosts[i] = domain.Post{
			ID:        post.ID,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
			Content:   post.Content,
			UserID:    post.UserID,
			User: domain.User{
				ID:        post.User.ID,
				Username:  post.User.Username,
				CreatedAt: post.User.CreatedAt,
				UpdatedAt: post.User.UpdatedAt,
			},
		}
	}

	return domainPosts, nil
}

// UpdatePost updates an existing post
func (repo *PostRepo) UpdatePost(id uuid.UUID, updates domain.Post) (*domain.Post, error) {
	result := repo.db.Model(&PostEntity{}).Where("id = ?", id).Updates(map[string]interface{}{
		"content":    updates.Content,
		"updated_at": time.Now(),
	})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("no post found with ID: %s", id)
	}

	return repo.GetPostByID(id)
}

// DeletePost deletes a post by ID
func (repo *PostRepo) DeletePost(id uuid.UUID) error {
	result := repo.db.Delete(&PostEntity{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no post found with ID: %s", id)
	}
	return nil
}

// PostEntity represents the post entity in the database
type PostEntity struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Content   string         `gorm:"type:text;not null"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null"`
	User      UserEntity     `gorm:"foreignKey:UserID;references:ID"`
}

// UserEntity embedded for preloading
type UserEntity struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string `gorm:"unique;index;not null"`
}

// BeforeCreate sets a UUID for the post
func (post *PostEntity) BeforeCreate(tx *gorm.DB) (err error) {
	post.ID = uuid.New()
	return
}

// TableName overrides the table name
func (post *PostEntity) TableName() string {
	return "posts"
}

// TableName for UserEntity when used in post context
func (user *UserEntity) TableName() string {
	return "users"
}
