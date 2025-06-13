package controller

import (
	"github.com/dakaii/vibegopher/internal/controller/comment"
	"github.com/dakaii/vibegopher/internal/controller/post"
	"github.com/dakaii/vibegopher/internal/controller/user"
	"github.com/dakaii/vibegopher/internal/repository"
)

// Controllers contains all the controllers
type Controllers struct {
	UserController    *user.Controller
	PostController    *post.Controller
	CommentController *comment.Controller
}

// InitControllers returns a new Controllers
func InitControllers(repositories *repository.Repositories) *Controllers {
	return &Controllers{
		UserController:    user.InitController(repositories.UserRepo),
		PostController:    post.InitController(repositories.PostRepo),
		CommentController: comment.InitController(repositories.CommentRepo),
	}
}
