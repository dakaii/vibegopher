package repository

import (
	"github.com/dakaii/vibegopher/internal/repository/botjobrepo"
	"github.com/dakaii/vibegopher/internal/repository/commentrepo"
	"github.com/dakaii/vibegopher/internal/repository/postrepo"
	"github.com/dakaii/vibegopher/internal/repository/userrepo"
	"gorm.io/gorm"
)

type Repositories struct {
	UserRepo    *userrepo.UserRepo
	PostRepo    *postrepo.PostRepo
	CommentRepo *commentrepo.CommentRepo
	BotJobRepo  *botjobrepo.BotJobRepo
}

func InitRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		UserRepo:    userrepo.NewUserRepo(db),
		PostRepo:    postrepo.NewPostRepo(db),
		CommentRepo: commentrepo.NewCommentRepo(db),
		BotJobRepo:  botjobrepo.NewBotJobRepo(db),
	}
}
