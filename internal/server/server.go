package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/khalid64/kongdb/pkg/utils"
	"go.uber.org/zap"
)

// Config represents the server configuration
type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Storage StorageConfig `yaml:"storage"`
	Network NetworkConfig `yaml:"network"`
}

// ServerConfig contains server-specific configuration
type ServerConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"`
}

// StorageConfig contains storage-specific configuration
type StorageConfig struct {
	DataDir         string `yaml:"data_dir"`
	MaxFileSize     int64  `yaml:"max_file_size"`
	BloomFilterSize int    `yaml:"bloom_filter_size"`
	BloomFilterHash int    `yaml:"bloom_filter_hash"`
}

// NetworkConfig contains network-specific configuration
type NetworkConfig struct {
	EnableTLS      bool   `yaml:"enable_tls"`
	CertFile       string `yaml:"cert_file"`
	KeyFile        string `yaml:"key_file"`
	MaxConnections int    `yaml:"max_connections"`
}

// Server represents the KongDB server
type Server struct {
	config *Config
	logger *utils.Logger
	server *http.Server
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	// For now, return default config
	// TODO: Implement actual config loading from file
	return &Config{
		Server: ServerConfig{
			Host:            "localhost",
			Port:            6468,
			ShutdownTimeout: 30 * time.Second,
			ReadTimeout:     15 * time.Second,
			WriteTimeout:    15 * time.Second,
			IdleTimeout:     60 * time.Second,
		},
		Storage: StorageConfig{
			DataDir:         "./data",
			MaxFileSize:     1024 * 1024 * 1024, // 1GB
			BloomFilterSize: 1000000,
			BloomFilterHash: 7,
		},
		Network: NetworkConfig{
			EnableTLS:      false,
			MaxConnections: 1000,
		},
	}, nil
}

// NewServer creates a new server instance
func NewServer(config *Config, logger *utils.Logger) (*Server, error) {
	s := &Server{
		config: config,
		logger: logger,
	}

	// Create HTTP server
	s.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port),
		ReadTimeout:  config.Server.ReadTimeout,
		WriteTimeout: config.Server.WriteTimeout,
		IdleTimeout:  config.Server.IdleTimeout,
		Handler:      s.createHandler(),
	}

	return s, nil
}

// createHandler creates the HTTP handler for the server
func (s *Server) createHandler() http.Handler {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", s.healthHandler)

	// API endpoints (to be implemented)
	mux.HandleFunc("/api/v1/", s.apiHandler)

	return mux
}

// healthHandler handles health check requests
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
}

// apiHandler handles API requests
func (s *Server) apiHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement API routing
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"error":"API not implemented yet"}`))
}

// Start starts the server
func (s *Server) Start() error {
	s.logger.Info("Server starting", zap.String("address", s.server.Addr))
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Server shutting down")
	return s.server.Shutdown(ctx)
}
