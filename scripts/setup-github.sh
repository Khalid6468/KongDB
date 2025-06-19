#!/bin/bash

# KongDB GitHub Setup Script
# This script helps set up the Git repository and prepare it for GitHub

set -e

echo "🚀 Setting up KongDB for GitHub..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in the right directory
if [ ! -f "go.mod" ] || [ ! -f "README.md" ]; then
    print_error "This script must be run from the KongDB project root directory"
    exit 1
fi

# Check if Git is installed
if ! command -v git &> /dev/null; then
    print_error "Git is not installed. Please install Git first."
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go 1.21+ first."
    exit 1
fi

print_status "Checking Go version..."
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
if [[ $(echo "$GO_VERSION 1.21" | tr " " "\n" | sort -V | head -n 1) != "1.21" ]]; then
    print_error "Go version $GO_VERSION is too old. Please install Go 1.21+"
    exit 1
fi
print_success "Go version $GO_VERSION is compatible"

# Initialize Git repository if not already initialized
if [ ! -d ".git" ]; then
    print_status "Initializing Git repository..."
    git init
    print_success "Git repository initialized"
else
    print_status "Git repository already exists"
fi

# Update go.mod with your GitHub username
print_status "Please enter your GitHub username:"
read -r GITHUB_USERNAME

if [ -z "$GITHUB_USERNAME" ]; then
    print_error "GitHub username is required"
    exit 1
fi

# Update go.mod file
print_status "Updating go.mod with your GitHub username..."
sed -i.bak "s/github.com\/yourusername\/kongdb/github.com\/$GITHUB_USERNAME\/kongdb/g" go.mod
rm -f go.mod.bak
print_success "go.mod updated"

# Install dependencies
print_status "Installing Go dependencies..."
go mod download
go mod tidy
print_success "Dependencies installed"

# Create initial commit
print_status "Creating initial commit..."
git add .
git commit -m "feat: initial commit - KongDB relational database

- Complete project structure and documentation
- Technical architecture with distributed capabilities
- TDD-focused development roadmap
- Bloom filter integration for performance optimization
- Comprehensive testing strategy
- Docker and CI/CD setup

This is a relational database implementation in Go designed to be
neither too simple nor too complex, with eventual distributed
capabilities and advanced features like bloom filters for query
optimization."
print_success "Initial commit created"

# Create develop branch
print_status "Creating develop branch..."
git checkout -b develop
print_success "Develop branch created"

# Switch back to main
git checkout main

print_success "✅ KongDB repository setup complete!"

echo ""
echo "📋 Next steps:"
echo "1. Create a new repository on GitHub: https://github.com/new"
echo "2. Name it 'kongdb'"
echo "3. Make it public or private as you prefer"
echo "4. Don't initialize with README, .gitignore, or license (we already have them)"
echo ""
echo "5. After creating the repository, run these commands:"
echo "   git remote add origin https://github.com/$GITHUB_USERNAME/kongdb.git"
echo "   git push -u origin main"
echo "   git push -u origin develop"
echo ""
echo "6. Enable GitHub Actions in your repository settings"
echo "7. Set up branch protection rules for main and develop"
echo ""
echo "🎉 Your KongDB project is ready for GitHub!"
echo ""
echo "💡 Useful commands:"
echo "   make help          - Show available make commands"
echo "   make test          - Run all tests"
echo "   make build         - Build KongDB binaries"
echo "   make test-coverage - Run tests with coverage"
echo "   docker-compose up  - Run KongDB in distributed mode" 