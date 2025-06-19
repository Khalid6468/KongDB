package node

import (
	"context"
	"time"

	"github.com/khalid64/kongdb/pkg/utils"
	"go.uber.org/zap"
)

// Config represents the node configuration
type Config struct {
	NodeID            string        `yaml:"node_id"`
	ClusterAddress    string        `yaml:"cluster_address"`
	DataDir           string        `yaml:"data_dir"`
	ShutdownTimeout   time.Duration `yaml:"shutdown_timeout"`
	HeartbeatInterval time.Duration `yaml:"heartbeat_interval"`
	ElectionTimeout   time.Duration `yaml:"election_timeout"`
}

// Node represents a KongDB node in distributed mode
type Node struct {
	config *Config
	logger *utils.Logger
	// TODO: Add consensus, storage, and network components
}

// LoadConfig loads node configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	// For now, return default config
	// TODO: Implement actual config loading from file
	return &Config{
		NodeID:            "",
		ClusterAddress:    "",
		DataDir:           "./data",
		ShutdownTimeout:   30 * time.Second,
		HeartbeatInterval: 100 * time.Millisecond,
		ElectionTimeout:   1000 * time.Millisecond,
	}, nil
}

// NewNode creates a new node instance
func NewNode(config *Config, logger *utils.Logger) (*Node, error) {
	n := &Node{
		config: config,
		logger: logger,
	}

	// TODO: Initialize consensus, storage, and network components

	return n, nil
}

// Start starts the node
func (n *Node) Start() error {
	n.logger.Info("Node starting", zap.String("node_id", n.config.NodeID), zap.String("cluster", n.config.ClusterAddress))

	// TODO: Start consensus, storage, and network components

	return nil
}

// Shutdown gracefully shuts down the node
func (n *Node) Shutdown(ctx context.Context) error {
	n.logger.Info("Node shutting down")

	// TODO: Shutdown consensus, storage, and network components

	return nil
}

// GetStatus returns the current node status
func (n *Node) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"node_id":   n.config.NodeID,
		"cluster":   n.config.ClusterAddress,
		"status":    "running", // TODO: Implement actual status
		"timestamp": time.Now().Format(time.RFC3339),
	}
}
