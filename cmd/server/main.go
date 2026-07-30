package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dakaii/vibegopher/internal/api"
	"github.com/dakaii/vibegopher/internal/bot"
	"github.com/dakaii/vibegopher/internal/controller"
	"github.com/dakaii/vibegopher/internal/database"
	"github.com/dakaii/vibegopher/internal/envvar"
	"github.com/dakaii/vibegopher/internal/repository"
)

func main() {
	db := database.GetDatabase()
	repos := repository.InitRepositories(db)
	controllers := controller.InitControllers(repos)
	router := api.SetupRouter(controllers)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if envvar.BotWorkerEnabled() {
		worker := bot.NewWorker(repos.BotJobRepo, repos.PostRepo, repos.CommentRepo)
		go worker.Run(ctx)
	}

	port := envvar.Port()
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		fmt.Println("REST API server is started at: http://localhost:" + port + "/")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = server.Shutdown(shutdownCtx)
}
