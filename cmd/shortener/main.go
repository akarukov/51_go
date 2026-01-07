package main

import (
	"github.com/akarukov/51_go.git/internal/config"
	handlers "github.com/akarukov/51_go.git/internal/handler"
	"github.com/akarukov/51_go.git/internal/repository"
	"github.com/akarukov/51_go.git/internal/service"
	"log"
)

func main() {
	if err := runServer(); err != nil {
		log.Fatal(err)
	}
}

func runServer() error {
	cfg := config.GetConfig()
	store := repository.NewStore()
	shortenerService := service.NewShortenerService(store)

	log.Printf("Start server at %s\n", cfg.ServerAddr)
	return handlers.Serve(cfg, shortenerService)
}
