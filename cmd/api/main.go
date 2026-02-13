package main

import (
	"log"

	"upisettle/internal/config"
	"upisettle/internal/logger"
	"upisettle/internal/storage"
	httpserver "upisettle/internal/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logg := logger.New(cfg.Env)

	db, err := storage.NewDB(cfg, logg)
	if err != nil {
		logg.Fatalf("failed to connect to database: %v", err)
	}

	server := httpserver.NewServer(cfg, logg, db)

	if err := server.Run(); err != nil {
		logg.Fatalf("server exited with error: %v", err)
	}
}

