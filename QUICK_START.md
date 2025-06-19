# KongDB Quick Start Guide

Get KongDB up and running on GitHub in minutes!

## 🚀 Quick Setup

### 1. Run the Setup Script

```bash
# Make sure you're in the KongDB project directory
cd /path/to/KongDB

# Run the setup script
./scripts/setup-github.sh
```

The script will:
- ✅ Check your Go version (requires 1.21+)
- ✅ Initialize Git repository
- ✅ Update go.mod with your GitHub username
- ✅ Install dependencies
- ✅ Create initial commit
- ✅ Set up develop branch

### 2. Create GitHub Repository

1. Go to [GitHub New Repository](https://github.com/new)
2. Repository name: `kongdb`
3. Description: `A relational database in Go with distributed capabilities and bloom filter optimization`
4. Make it **Public** or **Private** (your choice)
5. **Don't** initialize with README, .gitignore, or license (we already have them)
6. Click "Create repository"

### 3. Push to GitHub

After creating the repository, run these commands (replace `yourusername` with your actual GitHub username):

```bash
git remote add origin https://github.com/yourusername/kongdb.git
git push -u origin main
git push -u origin develop
```

### 4. Enable GitHub Actions

1. Go to your repository on GitHub
2. Click "Actions" tab
3. Click "Enable Actions"
4. The CI/CD pipeline will automatically run on your next push

## 🛠️ Development Commands

Once set up, you can use these commands:

```bash
# Show all available commands
make help

# Install dependencies
make install

# Run all tests
make test

# Run tests with coverage
make test-coverage

# Build KongDB binaries
make build

# Format code
make format

# Run linting
make lint

# Run performance benchmarks
make benchmark

# Clean build artifacts
make clean
```

## 🐳 Docker Development

Run KongDB in a distributed cluster:

```bash
# Start the full cluster (5 nodes)
docker-compose up

# Access the coordinator at http://localhost:8080
# Access Grafana monitoring at http://localhost:3000 (admin/admin)

# Stop the cluster
docker-compose down
```

## 📁 Project Structure

```
kongdb/
├── cmd/                    # Application entry points
├── internal/              # Core implementation
│   ├── parser/           # SQL parsing
│   ├── optimizer/        # Query optimization
│   ├── storage/          # Storage engine
│   │   └── bloom_filter/ # Bloom filter implementation
│   ├── transaction/      # Transaction management
│   └── ...
├── test/                 # Tests
├── docs/                 # Documentation
├── scripts/              # Setup scripts
├── .github/              # GitHub Actions
├── Dockerfile            # Container configuration
├── docker-compose.yml    # Multi-node setup
├── Makefile              # Development commands
└── README.md             # Project overview
```

## 🧪 Testing Strategy

KongDB follows Test-Driven Development (TDD):

```bash
# Run unit tests
make test-quick

# Run integration tests
make test-integration

# Run performance tests
make test-performance

# Run with race detector
make test-race

# Check test coverage
make test-coverage
```

## 🔧 Configuration

Environment variables for KongDB:

```bash
# Node configuration
KONGDB_NODE_ID=coordinator
KONGDB_NODE_ROLE=coordinator
KONGDB_DATA_DIR=/app/data
KONGDB_LOG_LEVEL=debug

# Cluster configuration
KONGDB_CLUSTER_MODE=true
KONGDB_BOOTSTRAP=true
KONGDB_COORDINATOR_HOST=localhost
KONGDB_COORDINATOR_PORT=8080
```

## 📊 Monitoring

KongDB includes built-in monitoring:

- **Health checks**: `http://localhost:8080/health`
- **Metrics**: `http://localhost:8080/metrics`
- **Grafana dashboard**: `http://localhost:3000` (when using Docker)

## 🚨 Troubleshooting

### Common Issues

1. **Go version too old**
   ```bash
   # Install Go 1.21+
   brew install go  # macOS
   sudo apt install golang-go  # Ubuntu
   ```

2. **Git not installed**
   ```bash
   # Install Git
   brew install git  # macOS
   sudo apt install git  # Ubuntu
   ```

3. **Docker not running**
   ```bash
   # Start Docker Desktop or Docker daemon
   sudo systemctl start docker  # Linux
   ```

4. **Port conflicts**
   ```bash
   # Check what's using port 8080
   lsof -i :8080
   
   # Kill process or change port in docker-compose.yml
   ```

### Getting Help

- 📖 Read the [Technical Architecture](TECHNICAL_ARCHITECTURE.md)
- 📋 Check the [Development Roadmap](DEVELOPMENT_ROADMAP.md)
- 🤝 See [Contributing Guidelines](CONTRIBUTING.md)
- 🐛 Report issues on GitHub
- 💬 Start discussions on GitHub

## 🎯 Next Steps

1. **Start with Week 1**: Follow the [Development Roadmap](DEVELOPMENT_ROADMAP.md)
2. **Write your first test**: Follow TDD principles
3. **Build the storage engine**: Start with basic CRUD operations
4. **Add bloom filters**: Implement performance optimizations
5. **Contribute**: Submit pull requests and help improve KongDB

## 🏆 Success Metrics

Your KongDB project will be successful when:

- ✅ All tests pass with 90%+ coverage
- ✅ Bloom filters improve query performance by 50%+
- ✅ ACID properties are maintained
- ✅ Distributed capabilities work correctly
- ✅ Documentation is comprehensive
- ✅ Code follows TDD principles

---

**Happy coding! 🚀**

Remember: Start with tests, build incrementally, and have fun building your own database! 