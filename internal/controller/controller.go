package controller

import (
	"log"

	"github.com/dakaii/vibegopher/internal/controller/comment"
	"github.com/dakaii/vibegopher/internal/controller/post"
	"github.com/dakaii/vibegopher/internal/controller/user"
	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/dakaii/vibegopher/internal/repository"
	"github.com/dakaii/vibegopher/internal/repository/botjobrepo"
	"github.com/google/uuid"
)

type Controllers struct {
	UserController    *user.Controller
	PostController    *post.Controller
	CommentController *comment.Controller
	BotJobs           *botjobrepo.BotJobRepo
}

func InitControllers(repositories *repository.Repositories) *Controllers {
	return &Controllers{
		UserController:    user.InitController(repositories.UserRepo),
		PostController:    post.InitController(repositories.PostRepo),
		CommentController: comment.InitController(repositories.CommentRepo),
		BotJobs:           repositories.BotJobRepo,
	}
}

// EnqueueBotJob queues async critic work. Failures are logged and do not fail the user request.
func (c *Controllers) EnqueueBotJob(kind string, targetID uuid.UUID, authorID uuid.UUID) {
	if c.BotJobs == nil {
		return
	}
	if authorID == domain.BotUserID {
		return
	}
	if _, err := c.BotJobs.Enqueue(kind, targetID); err != nil {
		log.Printf("enqueue critic job %s/%s: %v", kind, targetID, err)
	}
}
