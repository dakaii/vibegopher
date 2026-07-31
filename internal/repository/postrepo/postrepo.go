package postrepo

import (
	"errors"
	"fmt"
	"time"

	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostRepo struct {
	db *gorm.DB
}

func NewPostRepo(db *gorm.DB) *PostRepo {
	return &PostRepo{db: db}
}

func (repo *PostRepo) CreatePost(post domain.Post) (*domain.Post, error) {
	dbPost := PostEntity{
		Content: post.Content,
		UserID:  post.UserID,
	}
	if err := repo.db.Create(&dbPost).Error; err != nil {
		return nil, err
	}
	return repo.GetPostByID(dbPost.ID)
}

func (repo *PostRepo) GetPostByID(id uuid.UUID) (*domain.Post, error) {
	var post PostEntity
	err := repo.db.Preload("User").Where("id = ?", id).First(&post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no post found with ID: %s", id)
		}
		return nil, err
	}
	out := toDomain(post)
	return &out, nil
}

func (repo *PostRepo) GetAllPosts() ([]domain.Post, error) {
	var posts []PostEntity
	if err := repo.db.Preload("User").Order("created_at DESC").Limit(100).Find(&posts).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Post, len(posts))
	for i, post := range posts {
		out[i] = toDomain(post)
	}
	return out, nil
}

func (repo *PostRepo) GetPostsByUserID(userID uuid.UUID) ([]domain.Post, error) {
	var posts []PostEntity
	if err := repo.db.Preload("User").Where("user_id = ?", userID).Order("created_at DESC").Limit(100).Find(&posts).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Post, len(posts))
	for i, post := range posts {
		out[i] = toDomain(post)
	}
	return out, nil
}

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

type PostEntity struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Content   string     `gorm:"type:text;not null"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null"`
	User      UserEntity `gorm:"foreignKey:UserID;references:ID"`
}

type UserEntity struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string `gorm:"unique;index;not null"`
	IsBot     bool   `gorm:"column:is_bot;not null;default:false"`
}

func (post *PostEntity) BeforeCreate(tx *gorm.DB) error {
	if post.ID == uuid.Nil {
		post.ID = uuid.New()
	}
	return nil
}

func (PostEntity) TableName() string { return "posts" }
func (UserEntity) TableName() string { return "users" }

func toDomain(post PostEntity) domain.Post {
	return domain.Post{
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
			IsBot:     post.User.IsBot,
		},
	}
}
