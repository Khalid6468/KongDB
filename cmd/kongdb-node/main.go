package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/khalid64/kongdb/internal/node"
	"github.com/khalid64/kongdb/pkg/utils"
	"go.uber.org/zap"
)

var (
	configPath = flag.String("config", "configs/node.yaml", "Path to node configuration file")
	nodeID     = flag.String("node-id", "", "Node ID (required)")
	cluster    = flag.String("cluster", "", "Cluster address (required)")
	debug      = flag.Bool("debug", false, "Enable debug mode")
)

func main() {
	flag.Parse()

	// Validate required flags
	if *nodeID == "" {
		log.Fatal("--node-id is required")
	}
	if *cluster == "" {
		log.Fatal("--cluster is required")
	}

	// Initialize logger
	logger := utils.NewLogger(*debug)
	logger.Info("Starting KongDB node...", zap.String("node_id", *nodeID))

	// Load node configuration
	config, err := node.LoadConfig(*configPath)
	if err != nil {
		logger.Error("Failed to load node configuration", zap.Error(err))
		os.Exit(1)
	}

	// Override config with command line flags
	config.NodeID = *nodeID
	config.ClusterAddress = *cluster

	// Create node instance
	n, err := node.NewNode(config, logger)
	if err != nil {
		logger.Error("Failed to create node", zap.Error(err))
		os.Exit(1)
	}

	// Start node in a goroutine
	go func() {
		logger.Info("Starting node", zap.String("node_id", config.NodeID), zap.String("cluster", config.ClusterAddress))
		if err := n.Start(); err != nil {
			logger.Error("Node failed to start", zap.Error(err))
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the node
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down node...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
	defer cancel()

	if err := n.Shutdown(ctx); err != nil {
		logger.Error("Node forced to shutdown", zap.Error(err))
	}

	logger.Info("Node exited")
}
