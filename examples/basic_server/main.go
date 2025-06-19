package main

import (
	"log"

	"github.com/khalid64/kongdb/internal/server"
	"github.com/khalid64/kongdb/pkg/utils"
	"go.uber.org/zap"
)

func main() {
	// Create a logger
	logger := utils.NewLogger(true) // Enable debug mode

	// Load configuration
	config, err := server.LoadConfig("configs/kongdb.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Create server
	srv, err := server.NewServer(config, logger)
	if err != nil {
		log.Fatal("Failed to create server:", err)
	}

	logger.Info("Starting KongDB server example...")
	logger.Info("Server will be available at", zap.String("address", "http://localhost:8080"))
	logger.Info("Health check available at", zap.String("endpoint", "http://localhost:8080/health"))

	// Start server (this will block)
	if err := srv.Start(); err != nil {
		log.Fatal("Server failed:", err)
	}
}
