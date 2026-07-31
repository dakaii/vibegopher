package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dakaii/vibegopher/internal/critic"
	"github.com/dakaii/vibegopher/internal/database"
	"github.com/dakaii/vibegopher/internal/envvar"
	"github.com/dakaii/vibegopher/internal/repository"
)

func main() {
	if err := envvar.ValidateRuntimeConfig(); err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db := database.GetDatabase()
	repos := repository.InitRepositories(db)
	worker := critic.NewWorker(repos.BotJobRepo, repos.PostRepo, repos.CommentRepo)
	log.Println("starting standalone critic worker")
	worker.Run(ctx)
}
