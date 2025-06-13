package repository

import (
	"github.com/dakaii/vibegopher/internal/repository/commentrepo"
	"github.com/dakaii/vibegopher/internal/repository/postrepo"
	"github.com/dakaii/vibegopher/internal/repository/userrepo"
	"gorm.io/gorm"
)

// Repositories contains all the repo structs
type Repositories struct {
	UserRepo    *userrepo.UserRepo
	PostRepo    *postrepo.PostRepo
	CommentRepo *commentrepo.CommentRepo
}

// InitRepositories should be called in main.go
func InitRepositories(db *gorm.DB) *Repositories {
	userRepo := userrepo.NewUserRepo(db)
	postRepo := postrepo.NewPostRepo(db)
	commentRepo := commentrepo.NewCommentRepo(db)
	return &Repositories{
		UserRepo:    userRepo,
		PostRepo:    postRepo,
		CommentRepo: commentRepo,
	}
}
