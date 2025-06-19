package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/khalid64/kongdb/internal/server"
	"github.com/khalid64/kongdb/pkg/utils"
	"go.uber.org/zap"
)

var (
	configPath = flag.String("config", "configs/kongdb.yaml", "Path to configuration file")
	port       = flag.Int("port", 8080, "Server port")
	host       = flag.String("host", "localhost", "Server host")
	debug      = flag.Bool("debug", false, "Enable debug mode")
)

func main() {
	flag.Parse()

	// Initialize logger
	logger := utils.NewLogger(*debug)
	logger.Info("Starting KongDB server...")

	// Load configuration
	config, err := server.LoadConfig(*configPath)
	if err != nil {
		logger.Error("Failed to load configuration", zap.Error(err))
		os.Exit(1)
	}

	// Override config with command line flags
	if *port != 8080 {
		config.Server.Port = *port
	}
	if *host != "localhost" {
		config.Server.Host = *host
	}

	// Create server instance
	srv, err := server.NewServer(config, logger)
	if err != nil {
		logger.Error("Failed to create server", zap.Error(err))
		os.Exit(1)
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting server", zap.String("host", config.Server.Host), zap.Int("port", config.Server.Port))
		if err := srv.Start(); err != nil {
			logger.Error("Server failed to start", zap.Error(err))
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), config.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
