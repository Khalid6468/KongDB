# KongDB - An experimental database engine written in Golang—my open lab for poking at the guts of storage engines, indexing, MVCC, and Raft-style replication. 

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Status](https://img.shields.io/badge/Status-Planning-orange.svg)]()

## Overview

KongDB is a relational database management system (RDBMS) implemented in Go, designed to provide a comprehensive learning experience in database internals. The system features a SQL query parser, query optimizer, execution planner, and a storage engine with ACID compliance.

## Key Features

### Core Database Features
- **SQL Parser**: Full SQL support with complex query parsing
- **Query Optimizer**: Rule-based and cost-based optimization
- **Execution Engine**: Efficient query execution with multiple join algorithms
- **Storage Engine**: B+ tree indexing with heap file storage
- **Transaction Management**: MVCC with ACID compliance
- **Write-Ahead Logging**: Crash recovery and durability

### Development Approach
- **Test-Driven Development**: All features built using Red-Green-Refactor cycle
- **90%+ Test Coverage**: Comprehensive unit, integration, and performance tests
- **Modular Architecture**: Clean interfaces for easy testing and extension
- **Performance Focused**: Regular benchmarking and optimization

### (Future) Distributed Features
- **Horizontal Scaling**: Multi-node cluster support
- **Data Partitioning**: Hash and range-based partitioning
- **Distributed Transactions**: Two-phase commit with consensus
- **Global Indexes**: Distributed index management
- **SQL Compatibility**: Maintain SQL standard across distributed nodes

## Project Goals

### Learning Objectives
- Understand database internals at a granular level
- Learn modern database design patterns
- Experience building complex systems with TDD
- Develop better understanding of distributed systems

## Architecture

### High-Level Architecture
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   SQL Parser    │───▶│ Query Optimizer │───▶│ Execution Plan  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Storage       │◀───│ Transaction     │◀───│ Query Executor  │
│   Engine        │    │ Manager         │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Future Distributed Architecture(Deferred)
```
┌─────────────────────────────────────────────────────────────────┐
│                        Client Layer                             │
├─────────────────────────────────────────────────────────────────┤
│                    SQL Router/Load Balancer                     │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │   Node 1    │  │   Node 2    │  │   Node N    │             │
│  │ ┌─────────┐ │  │ ┌─────────┐ │  │ ┌─────────┐ │             │
│  │ │SQL Layer │ │  │ │SQL Layer │ │  │ │SQL Layer │ │             │
│  │ └─────────┘ │  │ └─────────┘ │  │ └─────────┘ │             │
│  │ ┌─────────┐ │  │ ┌─────────┐ │  │ ┌─────────┐ │             │
│  │ │Executor │ │  │ │Executor │ │  │ │Executor │ │             │
│  │ └─────────┘ │  │ └─────────┘ │  │ └─────────┘ │             │
│  │ ┌─────────┐ │  │ ┌─────────┐ │  │ ┌─────────┐ │             │
│  │ │Storage  │ │  │ │Storage  │ │  │ │Storage  │ │             │
│  │ └─────────┘ │  │ └─────────┘ │  │ └─────────┘ │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
├─────────────────────────────────────────────────────────────────┤
│                    Consensus Layer (Raft)                       │
└─────────────────────────────────────────────────────────────────┘
```

## Test-Driven Development Approach

### TDD Principles
- **Red-Green-Refactor**: Every feature follows this cycle
- **Test First**: Write tests before implementation
- **Minimal Implementation**: Write smallest code to pass tests
- **Continuous Refactoring**: Improve code while keeping tests green

### Testing Strategy
- **Unit Tests**: 90%+ coverage for all components
- **Integration Tests**: Test component interactions
- **Performance Tests**: Benchmark critical operations
- **End-to-End Tests**: Test complete user workflows
- **Distributed Tests**: Test cluster scenarios

### Testing Examples

#### SQL Parser Testing
```go
func TestSelectStatementParsing(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected *SelectStatement
    }{
        {
            name:  "simple select",
            input: "SELECT id, name FROM users",
            expected: &SelectStatement{
                Fields: []Expression{
                    &ColumnExpression{Name: "id"},
                    &ColumnExpression{Name: "name"},
                },
                From: []TableReference{
                    &TableReference{Name: "users"},
                },
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            parser := NewParser()
            result, err := parser.Parse(tt.input)
            
            require.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

#### Storage Engine Testing
```go
func TestStorageEngineACID(t *testing.T) {
    t.Run("atomicity", func(t *testing.T) {
        storage := NewStorageEngine()
        
        txn, err := storage.BeginTransaction()
        require.NoError(t, err)
        
        storage.Write(Key("key1"), Value("value1"), txn)
        storage.Write(Key("key2"), Value("value2"), txn)
        
        // Rollback transaction
        storage.RollbackTransaction(txn)
        
        // Verify no data was persisted
        result1, _ := storage.Read(Key("key1"))
        result2, _ := storage.Read(Key("key2"))
        assert.Nil(t, result1)
        assert.Nil(t, result2)
    })
}
```

## Getting Started

### Prerequisites
- Go 1.21 or later
- Git
- Docker (for testing environments)

### Installation
```bash
# Clone the repository
git clone https://github.com/yourusername/kongdb.git
cd kongdb

# Install dependencies
go mod download

# Run tests (TDD approach)
go test -v ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Running KongDB
```bash
# Start KongDB server
go run cmd/kongdb/main.go

# Connect and run SQL
# (Client implementation coming soon)
```

### Running Tests
```bash
# Unit tests
go test -v ./internal/...

# Integration tests
go test -v ./test/integration/...

# Performance tests
go test -bench=. ./test/performance/...

# All tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

## Project Structure

```
kongdb/
├── cmd/
│   ├── kongdb/                 # Main application entry point
│   └── kongdb-node/            # Distributed node entry point
├── internal/
│   ├── parser/                 # SQL parsing and AST
│   ├── optimizer/              # Query optimization
│   ├── planner/                # Execution planning
│   ├── executor/               # Query execution
│   ├── storage/                # Storage engine (local)
│   ├── transaction/            # Transaction management
│   ├── catalog/                # Metadata management
│   ├── types/                  # Data types and schemas
│   ├── network/                # Network layer and RPC
│   ├── consensus/              # Distributed consensus (Raft)
│   ├── routing/                # Query routing and load balancing
│   ├── cluster/                # Cluster management
│   └── distributed/            # Distributed-specific components
├── pkg/
│   ├── sql/                    # Public SQL interface
│   ├── client/                 # Client library
│   ├── utils/                  # Utility functions
│   └── distributed/            # Distributed database interfaces
├── test/                       # Integration tests
│   ├── unit/                   # Unit tests
│   ├── integration/            # Integration tests
│   ├── performance/            # Performance tests
│   ├── distributed/            # Distributed tests
│   └── fixtures/               # Test data and fixtures
├── docs/                       # Documentation
├── examples/                   # Usage examples
├── go.mod
├── go.sum
└── README.md
```

## Development Workflow

### TDD Development Process
1. **Write Test First**: Define expected behavior through tests
2. **Run Test**: Verify it fails (RED)
3. **Write Minimal Code**: Implement just enough to pass test (GREEN)
4. **Refactor**: Improve code while keeping tests green (REFACTOR)
5. **Repeat**: Add more test cases and functionality

### Feature Development
1. **Acceptance Test**: Write high-level acceptance test
2. **Unit Tests**: Write unit tests for each component
3. **Integration Tests**: Test component interactions
4. **Performance Tests**: Ensure performance requirements are met
5. **Documentation**: Update documentation with examples

### Code Quality
- **90%+ Test Coverage**: All components must meet this threshold
- **Code Reviews**: All changes reviewed with TDD focus
- **Performance Benchmarks**: Regular performance testing
- **Documentation**: Comprehensive documentation with examples

## Testing Strategy

### Unit Testing
- Individual component testing with mocks
- Edge case coverage
- Error condition testing
- Performance regression testing

### Integration Testing
- End-to-end query execution
- Transaction consistency testing
- Concurrency testing
- Recovery testing

### Performance Testing
- Query performance benchmarks
- Throughput testing
- Memory usage analysis
- Stress testing

### Distributed Testing
- Multi-node cluster validation
- Node failure and recovery scenarios
- Network partition simulation
- Distributed consistency verification

## Performance Goals

### Single Node Performance
- **Query Response Time**: < 100ms for simple queries on small datasets
- **Throughput**: 1000+ TPS for basic operations
- **Memory Usage**: Efficient buffer pool management
- **Disk I/O**: Minimize random access patterns

### Future Distributed Performance
- **Linear Scalability**: Scale with node count
- **Sub-second Queries**: Cross-node query performance
- **Network Efficiency**: Optimize for network communication
- **Fault Tolerance**: Handle node failures gracefully

## Contributing

### Development Guidelines
1. **Follow TDD**: Always write tests first
2. **Maintain Coverage**: Keep 90%+ test coverage
3. **Code Quality**: Follow Go best practices
4. **Documentation**: Update docs with changes
5. **Performance**: Consider performance implications

### Testing Requirements
- All new features must have comprehensive tests
- Performance benchmarks for critical paths
- Integration tests for component interactions
- Documentation tests for examples

### Pull Request Process
1. Write tests for new features
2. Implement features following TDD
3. Ensure all tests pass
4. Update documentation
5. Submit pull request with test coverage

## Learning Resources

### Database Internals
- **Database System Concepts** by Silberschatz, Korth, and Sudarshan
- **Database Design and Implementation** by Edward Sciore
- **Transaction Processing** by Jim Gray and Andreas Reuter

### Distributed Systems
- **Designing Data-Intensive Applications** by Martin Kleppmann
- **Distributed Systems: Concepts and Design** by George Coulouris
- **Consensus on Transaction Commit** by Jim Gray and Leslie Lamport

### Test-Driven Development
- **Test-Driven Development by Example** by Kent Beck
- **Growing Object-Oriented Software, Guided by Tests** by Steve Freeman and Nat Pryce
- **Refactoring: Improving the Design of Existing Code** by Martin Fowler

### Go Programming
- **The Go Programming Language** by Alan Donovan and Brian Kernighan
- **Concurrency in Go** by Katherine Cox-Buday
- **Go in Action** by William Kennedy

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by modern database systems like PostgreSQL, MySQL, and CockroachDB
- Built with lessons from distributed systems research
- Test-driven development approach inspired by Kent Beck and the TDD community
- Go community for excellent tooling and libraries

## Contact

- **Project Maintainer**: [Khalid Shareef]
- **Email**: [khalimohammed5@gmail.com]
- **GitHub**: [https://github.com/Khalid6468/kongdb]
- **Issues**: [https://github.com/Khalid6468/kongdb/issues]

## Roadmap

See [DEVELOPMENT_ROADMAP.md](DEVELOPMENT_ROADMAP.md) for detailed development phases and milestones.

## Documentation

- [Product Requirements Document](PRD_KongDB.md) - Comprehensive project requirements
- [Technical Architecture](TECHNICAL_ARCHITECTURE.md) - Detailed technical design
- [TDD Strategy](TDD_STRATEGY.md) - Test-driven development approach
- [Development Roadmap](DEVELOPMENT_ROADMAP.md) - Step-by-step development plan

---

**KongDB**: Building a database system the right way - with comprehensive testing and a clear path to distributed capabilities. 
