package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dakaii/vibegopher/internal/bot"
	"github.com/dakaii/vibegopher/internal/database"
	"github.com/dakaii/vibegopher/internal/repository"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db := database.GetDatabase()
	repos := repository.InitRepositories(db)
	worker := bot.NewWorker(repos.BotJobRepo, repos.PostRepo, repos.CommentRepo)
	log.Println("starting standalone bot worker")
	worker.Run(ctx)
}
