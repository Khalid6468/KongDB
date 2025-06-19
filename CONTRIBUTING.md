# Contributing to KongDB

Thank you for your interest in contributing to KongDB! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Test-Driven Development](#test-driven-development)
- [Code Style](#code-style)
- [Pull Request Process](#pull-request-process)
- [Reporting Bugs](#reporting-bugs)
- [Feature Requests](#feature-requests)
- [Documentation](#documentation)

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code.

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- Basic understanding of database concepts
- Familiarity with Test-Driven Development (TDD)

### Setting Up Your Development Environment

1. **Fork the repository**
   ```bash
   # Fork on GitHub, then clone your fork
   git clone https://github.com/yourusername/kongdb.git
   cd kongdb
   ```

2. **Set up the upstream remote**
   ```bash
   git remote add upstream https://github.com/original-owner/kongdb.git
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Run tests to ensure everything works**
   ```bash
   go test -v ./...
   ```

## Development Workflow

### Branch Strategy

- `main`: Production-ready code
- `develop`: Integration branch for features
- `feature/*`: Feature branches
- `bugfix/*`: Bug fix branches
- `hotfix/*`: Critical bug fixes

### Creating a Feature Branch

```bash
# Create and switch to a new feature branch
git checkout -b feature/your-feature-name

# Make your changes following TDD principles
# Write tests first, then implementation

# Commit your changes
git add .
git commit -m "feat: add your feature description"

# Push to your fork
git push origin feature/your-feature-name
```

## Test-Driven Development

KongDB follows strict TDD principles. All contributions must follow the Red-Green-Refactor cycle:

### TDD Process

1. **RED**: Write a failing test that describes the desired behavior
2. **GREEN**: Write minimal code to make the test pass
3. **REFACTOR**: Improve the code while keeping tests green

### Testing Requirements

- **Unit Tests**: 90%+ coverage for all new code
- **Integration Tests**: Test component interactions
- **Performance Tests**: Benchmark critical operations
- **Documentation Tests**: Ensure examples work

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test files
go test -v ./internal/storage/...

# Run benchmarks
go test -bench=. ./test/performance/...

# Run race detector
go test -race ./internal/...
```

### Test Examples

```go
// Example of good TDD test
func TestBloomFilterInsertion(t *testing.T) {
    t.Run("insert and query", func(t *testing.T) {
        // Arrange
        bf := NewBloomFilter(BloomFilterConfig{
            FalsePositiveRate: 0.01,
            ExpectedElements:  1000,
        })
        
        // Act
        bf.Insert(Key("test-key"))
        
        // Assert
        assert.True(t, bf.MightContain(Key("test-key")))
    })
}
```

## Code Style

### Go Conventions

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Run `go vet` to check for common mistakes
- Use `staticcheck` for additional static analysis

### Naming Conventions

- **Packages**: Lowercase, single word
- **Functions**: MixedCaps or mixedCaps
- **Variables**: MixedCaps or mixedCaps
- **Constants**: MixedCaps
- **Interfaces**: Single method interfaces end with -er

### Code Organization

```
internal/
├── component/
│   ├── component.go      # Main implementation
│   ├── component_test.go # Unit tests
│   └── types.go          # Type definitions
```

### Documentation

- All exported functions must have godoc comments
- Include examples in tests
- Update README.md for user-facing changes

```go
// BloomFilter provides probabilistic membership testing
// with configurable false positive rates.
type BloomFilter struct {
    // ... fields
}

// Insert adds a key to the bloom filter.
// The key will be hashed using multiple hash functions
// and corresponding bits will be set in the filter.
func (bf *BloomFilter) Insert(key Key) {
    // ... implementation
}
```

## Pull Request Process

### Before Submitting

1. **Ensure tests pass**
   ```bash
   go test -v ./...
   go test -race ./internal/...
   go vet ./...
   ```

2. **Check coverage**
   ```bash
   go test -coverprofile=coverage.out ./...
   go tool cover -func=coverage.out
   ```

3. **Update documentation**
   - Update README.md if needed
   - Add/update inline documentation
   - Update examples

4. **Squash commits** (if needed)
   ```bash
   git rebase -i HEAD~n  # where n is the number of commits
   ```

### Pull Request Guidelines

1. **Title**: Use conventional commit format
   - `feat: add bloom filter optimization`
   - `fix: resolve race condition in transaction manager`
   - `docs: update README with new features`

2. **Description**: Include
   - What the PR does
   - Why it's needed
   - How it was tested
   - Any breaking changes

3. **Checklist**:
   - [ ] Tests pass with 90%+ coverage
   - [ ] Code follows TDD principles
   - [ ] Documentation updated
   - [ ] No breaking changes (or documented)
   - [ ] Performance impact considered

### Review Process

1. **Automated checks** must pass
2. **Code review** by maintainers
3. **TDD compliance** verified
4. **Performance impact** assessed
5. **Documentation** reviewed

## Reporting Bugs

### Bug Report Template

```markdown
**Bug Description**
Clear and concise description of the bug.

**Steps to Reproduce**
1. Go to '...'
2. Click on '...'
3. Execute '...'
4. See error

**Expected Behavior**
What you expected to happen.

**Actual Behavior**
What actually happened.

**Environment**
- OS: [e.g. macOS, Linux, Windows]
- Go Version: [e.g. 1.21.0]
- KongDB Version: [e.g. 0.1.0]

**Additional Context**
Any other context about the problem.
```

## Feature Requests

### Feature Request Template

```markdown
**Feature Description**
Clear and concise description of the feature.

**Use Case**
Why is this feature needed? What problem does it solve?

**Proposed Solution**
How would you like to see this implemented?

**Alternatives Considered**
Any alternative solutions you've considered.

**Additional Context**
Any other context or screenshots.
```

## Documentation

### Contributing to Documentation

1. **User Documentation**: Update README.md for user-facing changes
2. **API Documentation**: Add godoc comments for all exported functions
3. **Architecture Documentation**: Update technical docs for significant changes
4. **Examples**: Provide working examples in tests

### Documentation Standards

- Use clear, concise language
- Include code examples
- Keep documentation up-to-date with code
- Use proper markdown formatting

## Getting Help

- **Issues**: Use GitHub issues for bugs and feature requests
- **Discussions**: Use GitHub Discussions for questions and ideas
- **Code Review**: Ask questions in pull request comments

## Recognition

Contributors will be recognized in:
- README.md contributors section
- Release notes
- GitHub contributors page

Thank you for contributing to KongDB! Your contributions help make this project better for everyone. 