package userrepo

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/dakaii/vibegopher/internal/envvar"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (repo *UserRepo) GetExistingUser(username string) (*domain.User, error) {
	var user UserEntity
	result := repo.db.Where("username = ? AND deleted_at IS NULL", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no user found with username: %s", username)
		}
		return nil, result.Error
	}
	return toDomain(user), nil
}

func (repo *UserRepo) GetByGoogleSub(googleSub string) (*domain.User, error) {
	var user UserEntity
	result := repo.db.Where("google_sub = ? AND deleted_at IS NULL", googleSub).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no user found with google_sub")
		}
		return nil, result.Error
	}
	return toDomain(user), nil
}

func (repo *UserRepo) GetByID(id uuid.UUID) (*domain.User, error) {
	var user UserEntity
	result := repo.db.Where("id = ? AND deleted_at IS NULL", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no user found with id")
		}
		return nil, result.Error
	}
	return toDomain(user), nil
}

func (repo *UserRepo) CreateUser(user domain.User) (*domain.User, error) {
	dbUser := UserEntity{
		Username:  user.Username,
		Email:     nullString(user.Email),
		GoogleSub: nullString(user.GoogleSub),
		IsBot:     user.IsBot,
	}
	if user.Password != "" {
		hashedPass, err := HashPassword(user.Password)
		if err != nil {
			return nil, err
		}
		dbUser.Password = &hashedPass
	}

	result := repo.db.Create(&dbUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomain(dbUser), nil
}

// UpsertGoogleUser finds by google_sub or creates a new Google-backed user.
func (repo *UserRepo) UpsertGoogleUser(profile domain.GoogleProfile) (*domain.User, error) {
	if existing, err := repo.GetByGoogleSub(profile.Sub); err == nil {
		return existing, nil
	}

	username, err := repo.uniqueUsernameFromProfile(profile)
	if err != nil {
		return nil, err
	}

	return repo.CreateUser(domain.User{
		Username:  username,
		Email:     profile.Email,
		GoogleSub: profile.Sub,
	})
}

func (repo *UserRepo) uniqueUsernameFromProfile(profile domain.GoogleProfile) (string, error) {
	base := sanitizeUsername(profile.GivenName)
	if base == "" {
		base = sanitizeUsername(strings.Split(profile.Email, "@")[0])
	}
	if base == "" {
		base = "user"
	}
	if len(base) < 3 {
		base = base + "user"
	}
	if len(base) > 40 {
		base = base[:40]
	}

	candidate := base
	for i := 0; i < 20; i++ {
		_, err := repo.GetExistingUser(candidate)
		if err != nil {
			// not found → available
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s%d", base, i+1)
	}
	return fmt.Sprintf("user_%s", profile.Sub[len(profile.Sub)-8:]), nil
}

func sanitizeUsername(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func HashPassword(password string) (string, error) {
	cost := envvar.HashCost()
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

type UserEntity struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Username  string         `gorm:"unique;index;not null"`
	Password  *string        `gorm:"type:varchar(1000)"`
	GoogleSub *string        `gorm:"column:google_sub;uniqueIndex"`
	Email     *string        `gorm:"type:varchar(255)"`
	IsBot     bool           `gorm:"not null;default:false"`
}

func (user *UserEntity) BeforeCreate(tx *gorm.DB) (err error) {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	return
}

func (UserEntity) TableName() string {
	return "users"
}

func toDomain(user UserEntity) *domain.User {
	out := &domain.User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Username:  user.Username,
		IsBot:     user.IsBot,
	}
	if user.Password != nil {
		out.Password = *user.Password
	}
	if user.GoogleSub != nil {
		out.GoogleSub = *user.GoogleSub
	}
	if user.Email != nil {
		out.Email = *user.Email
	}
	return out
}

func nullString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
