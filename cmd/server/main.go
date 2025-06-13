package main

import (
	"fmt"
	"net/http"

	"github.com/dakaii/vibegopher/internal/api"
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

	port := envvar.Port()
	fmt.Println("REST API server is started at: http://localhost:" + port + "/")
	fmt.Println("Available endpoints:")
	fmt.Println("  POST /api/signup")
	fmt.Println("  POST /api/login")
	fmt.Println("  GET  /api/me")

	http.ListenAndServe(":"+port, router)
}
