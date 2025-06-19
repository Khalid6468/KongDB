# KongDB Technical Architecture Document

## 1. Project Structure

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
│   │   ├── bloom_filter/        # Bloom filter implementation
│   │   ├── btree/              # B+ tree implementation
│   │   ├── heap/               # Heap file storage
│   │   └── wal/                # Write-ahead logging
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
│   └── fixtures/                # Test data and fixtures
├── docs/                       # Documentation
├── examples/                   # Usage examples
├── go.mod
├── go.sum
└── README.md
```

## 2. Distributed RDBMS Challenges & Solutions

### 2.1 Real-World Distributed RDBMS Examples

#### 2.1.1 CockroachDB
- **Architecture**: Multi-region, globally distributed SQL database
- **Key Features**: 
  - Strong consistency with serializable isolation
  - Automatic data distribution and rebalancing
  - Multi-active availability
  - Built on RocksDB with Raft consensus
- **Inspiration**: Global data distribution, automatic sharding

#### 2.1.2 TiDB
- **Architecture**: MySQL-compatible distributed database
- **Key Features**:
  - Horizontal scalability with MySQL compatibility
  - Strong consistency with distributed transactions
  - Automatic sharding with TiKV storage layer
- **Inspiration**: MySQL compatibility, layered architecture

#### 2.1.3 YugabyteDB
- **Architecture**: PostgreSQL-compatible distributed database
- **Key Features**:
  - PostgreSQL compatibility with distributed storage
  - Multi-region deployment
  - Built on DocDB (document store) with Raft
- **Inspiration**: PostgreSQL compatibility, document-based storage

#### 2.1.4 Amazon Aurora
- **Architecture**: Cloud-native distributed MySQL/PostgreSQL
- **Key Features**:
  - Storage-compute separation
  - Multi-AZ replication
  - Shared storage architecture
- **Inspiration**: Storage-compute separation, cloud-native design

### 2.2 Key Challenges in Distributed RDBMS

#### 2.2.1 ACID Compliance Across Nodes
**Challenge**: Maintaining ACID properties across distributed nodes
**Solution**: Two-phase commit with distributed consensus

```go
// internal/transaction/distributed_acid.go
type DistributedACIDManager struct {
    consensus    consensus.ConsensusManager
    coordinator  *TwoPhaseCoordinator
    participants map[NodeID]*Participant
}

type TwoPhaseCoordinator struct {
    txnID        TransactionID
    participants []NodeID
    phase        CommitPhase
    log          *CommitLog
}

type CommitPhase int
const (
    PHASE_PREPARE CommitPhase = iota
    PHASE_COMMIT
    PHASE_ABORT
)
```

#### 2.2.2 Distributed Joins and Aggregations
**Challenge**: Performing complex SQL operations across partitioned data
**Solution**: Distributed query execution with data movement

```go
// internal/executor/distributed_joins.go
type DistributedJoinExecutor struct {
    router       routing.QueryRouter
    network      network.RPCClient
    coordinator  *JoinCoordinator
}

type JoinStrategy int
const (
    JOIN_BROADCAST JoinStrategy = iota
    JOIN_SHUFFLE
    JOIN_COLOCATED
    JOIN_REPLICATED
)

func (dje *DistributedJoinExecutor) ExecuteJoin(left, right Table, strategy JoinStrategy) (ResultSet, error) {
    switch strategy {
    case JOIN_BROADCAST:
        return dje.broadcastJoin(left, right)
    case JOIN_SHUFFLE:
        return dje.shuffleJoin(left, right)
    case JOIN_COLOCATED:
        return dje.colocatedJoin(left, right)
    case JOIN_REPLICATED:
        return dje.replicatedJoin(left, right)
    }
}
```

#### 2.2.3 Global Indexes and Constraints
**Challenge**: Maintaining indexes and constraints across distributed data
**Solution**: Global index coordination with eventual consistency

```go
// internal/storage/global_index.go
type GlobalIndexManager struct {
    localIndexes  map[string]*LocalIndex
    globalIndexes map[string]*GlobalIndex
    coordinator   *IndexCoordinator
}

type GlobalIndex struct {
    Name       string
    Table      string
    Columns    []string
    Strategy   IndexStrategy
    Replicas   map[NodeID]*IndexReplica
}

type IndexStrategy int
const (
    INDEX_GLOBAL IndexStrategy = iota
    INDEX_LOCAL
    INDEX_HYBRID
)
```

## 3. Distributed-Ready Architecture (RDBMS Focused)

### 3.1 High-Level Distributed Architecture
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

### 3.2 Core Abstractions for Distributed RDBMS

#### 3.2.1 Storage Interface (RDBMS Focused)
```go
// internal/storage/interface.go
type StorageEngine interface {
    // Local operations
    Read(key Key) (Value, error)
    Write(key Key, value Value) error
    Delete(key Key) error
    Scan(start, end Key) (Iterator, error)
    
    // RDBMS-specific operations
    BeginTransaction() (Transaction, error)
    CommitTransaction(txn Transaction) error
    RollbackTransaction(txn Transaction) error
    
    // Distributed operations
    MultiRead(keys []Key) (map[Key]Value, error)
    MultiWrite(operations []WriteOperation) error
    BatchOperation(operations []Operation) error
    
    // RDBMS-specific distributed operations
    DistributedJoin(left, right Table, condition Expression) (ResultSet, error)
    DistributedAggregate(table Table, aggFunc AggregateFunction) (Value, error)
    GlobalIndexLookup(indexName string, key Key) ([]Key, error)
    
    // Consistency and replication
    GetVersion(key Key) (Version, error)
    SetReplicationFactor(factor int) error
    SyncToReplicas() error
    EnsureConsistency(level ConsistencyLevel) error
}

type WriteOperation struct {
    Key   Key
    Value Value
    Type  OperationType // INSERT, UPDATE, DELETE
    TxnID TransactionID
}

type Operation struct {
    Type      OperationType
    Key       Key
    Value     Value
    Timestamp Timestamp
    NodeID    NodeID
    TxnID     TransactionID
}
```

#### 3.2.2 Transaction Interface (RDBMS Focused)
```go
// internal/transaction/interface.go
type TransactionManager interface {
    // Local transaction operations
    Begin() (Transaction, error)
    Commit(txn Transaction) error
    Rollback(txn Transaction) error
    
    // RDBMS-specific operations
    SetIsolationLevel(level IsolationLevel) error
    GetSnapshot() (Snapshot, error)
    LockResource(resource ResourceID, lockType LockType) error
    UnlockResource(resource ResourceID) error
    
    // Distributed transaction operations
    BeginDistributed() (DistributedTransaction, error)
    Prepare(txn DistributedTransaction) error
    CommitDistributed(txn DistributedTransaction) error
    AbortDistributed(txn DistributedTransaction) error
    
    // RDBMS-specific distributed operations
    GlobalDeadlockDetection() ([]TransactionID, error)
    DistributedSnapshotIsolation() (Snapshot, error)
    CrossNodeConstraintValidation(constraints []Constraint) error
}

type DistributedTransaction struct {
    ID        TransactionID
    Nodes     []NodeID
    Status    TransactionStatus
    Timestamp Timestamp
    Coordinator NodeID
    IsolationLevel IsolationLevel
    Locks     map[ResourceID]LockType
}

type IsolationLevel int
const (
    READ_UNCOMMITTED IsolationLevel = iota
    READ_COMMITTED
    REPEATABLE_READ
    SERIALIZABLE
)
```

#### 3.2.3 Query Execution Interface (RDBMS Focused)
```go
// internal/executor/interface.go
type QueryExecutor interface {
    // Local execution
    ExecuteLocal(plan ExecutionPlan) (ResultSet, error)
    
    // RDBMS-specific local operations
    ExecuteJoin(left, right Table, condition Expression) (ResultSet, error)
    ExecuteAggregate(table Table, aggFunc AggregateFunction) (Value, error)
    ExecuteSort(table Table, columns []string, order Order) (ResultSet, error)
    
    // Distributed execution
    ExecuteDistributed(plan DistributedExecutionPlan) (ResultSet, error)
    ExecuteParallel(plans []ExecutionPlan) ([]ResultSet, error)
    
    // RDBMS-specific distributed operations
    ExecuteDistributedJoin(left, right Table, strategy JoinStrategy) (ResultSet, error)
    ExecuteDistributedAggregate(table Table, aggFunc AggregateFunction) (Value, error)
    ExecuteGlobalIndexScan(indexName string, predicate Expression) (ResultSet, error)
    
    // Query routing
    RouteQuery(query Query) ([]NodeID, error)
    OptimizeForDistribution(plan ExecutionPlan) (DistributedExecutionPlan, error)
    AnalyzeDataLocality(tables []Table) (DataLocalityMap, error)
}

type DistributedExecutionPlan struct {
    SubPlans    map[NodeID]ExecutionPlan
    Coordination CoordinationStrategy
    MergeStrategy MergeStrategy
    DataMovement []DataMovement
    JoinStrategy JoinStrategy
}

type DataMovement struct {
    SourceNode NodeID
    TargetNode NodeID
    Data       []Key
    Operation  MovementOperation
}
```

## 4. Network Layer Design (RDBMS Optimized)

### 4.1 RPC Framework (SQL-Aware)
```go
// internal/network/rpc.go
type RPCClient interface {
    Call(method string, args interface{}, reply interface{}) error
    CallAsync(method string, args interface{}) (Future, error)
    Broadcast(method string, args interface{}) error
    StreamCall(method string, args interface{}) (Stream, error)
}

type RPCServer interface {
    Register(service interface{}) error
    Start(address string) error
    Stop() error
    SetCompression(compression CompressionType) error
}

// RDBMS-specific RPC services
type SQLService interface {
    ExecuteQuery(query string, params []Value) (ResultSet, error)
    ExecuteTransaction(operations []Operation) error
    GetTableSchema(tableName string) (Schema, error)
    CreateIndex(indexName string, tableName string, columns []string) error
}

type StorageService interface {
    Read(key Key) (Value, error)
    Write(key Key, value Value) error
    Delete(key Key) error
    Scan(start, end Key) (Iterator, error)
    BeginTransaction() (TransactionID, error)
    CommitTransaction(txnID TransactionID) error
}
```

### 4.2 Cluster Management (RDBMS Aware)
```go
// internal/cluster/manager.go
type ClusterManager interface {
    // Node management
    AddNode(node NodeInfo) error
    RemoveNode(nodeID NodeID) error
    GetNode(nodeID NodeID) (NodeInfo, error)
    ListNodes() ([]NodeInfo, error)
    
    // Health monitoring
    HealthCheck(nodeID NodeID) (HealthStatus, error)
    MonitorCluster() (ClusterStatus, error)
    
    // RDBMS-specific cluster operations
    RebalanceTables(tables []string) error
    MigratePartition(partitionID PartitionID, fromNode, toNode NodeID) error
    UpdateTableDistribution(tableName string, strategy DistributionStrategy) error
    
    // Load balancing
    GetOptimalNode(operation OperationType) (NodeID, error)
    RebalanceCluster() error
    GetTableLocality(tableName string) (map[NodeID]float64, error)
}

type NodeInfo struct {
    ID       NodeID
    Address  string
    Role     NodeRole // COORDINATOR, EXECUTOR, STORAGE
    Status   NodeStatus
    Metadata NodeMetadata
    Tables   []string // Tables stored on this node
    Load     float64  // Current load
}

type DistributionStrategy int
const (
    DISTRIBUTION_HASH DistributionStrategy = iota
    DISTRIBUTION_RANGE
    DISTRIBUTION_REPLICATED
    DISTRIBUTION_PARTITIONED
)
```

### 4.3 Query Routing (SQL-Aware)
```go
// internal/routing/router.go
type QueryRouter interface {
    // Route queries to appropriate nodes
    RouteQuery(query Query) ([]NodeID, error)
    RouteRead(key Key) (NodeID, error)
    RouteWrite(operation WriteOperation) ([]NodeID, error)
    
    // RDBMS-specific routing
    RouteTableScan(tableName string, predicate Expression) ([]NodeID, error)
    RouteJoin(leftTable, rightTable string, condition Expression) ([]NodeID, error)
    RouteIndexLookup(indexName string, key Key) ([]NodeID, error)
    
    // Load balancing
    GetLoadBalancedNode(operationType OperationType) (NodeID, error)
    UpdateNodeLoad(nodeID NodeID, load float64) error
    GetTableDistribution(tableName string) (map[NodeID]float64, error)
}

type RoutingStrategy interface {
    Route(query Query, nodes []NodeInfo) ([]NodeID, error)
    Name() string
}

type HashRouting struct{}
type RangeRouting struct{}
type LoadBasedRouting struct{}
type LocalityAwareRouting struct{}
```

## 5. Updated Core Components Design (RDBMS Focused)

### 5.1 SQL Parser Architecture (Unchanged - Stateless)
```go
// internal/parser/parser.go
type Parser struct {
    lexer    *Lexer
    current  Token
    errors   []ParseError
}

// Parser remains stateless and distributed-friendly
// Future: Support distributed query hints like /*+ DISTRIBUTE(table) */
```

### 5.2 Query Optimizer Architecture (Enhanced for Distributed RDBMS)
```go
// internal/optimizer/optimizer.go
type Optimizer struct {
    catalog    *catalog.Catalog
    statistics *Statistics
    rules      []OptimizationRule
    router     routing.QueryRouter
    cluster    cluster.ClusterManager
    locality   *DataLocalityAnalyzer
}

type OptimizationRule interface {
    Apply(plan Plan) Plan
    Name() string
    // New method for distributed optimization
    ApplyDistributed(plan DistributedPlan) DistributedPlan
}

type DistributedPlan struct {
    LocalPlans    map[NodeID]Plan
    Coordination  CoordinationStrategy
    DataMovement  []DataMovement
    MergeStrategy MergeStrategy
    JoinStrategy  JoinStrategy
    LocalityCost  float64
    NetworkCost   float64
}

// RDBMS-specific optimization rules
type DistributedJoinOptimizationRule struct{}
type LocalityOptimizationRule struct{}
type GlobalIndexOptimizationRule struct{}
type PartitionPruningRule struct{}
```

### 5.3 Execution Planner Architecture (Enhanced for Distributed RDBMS)
```go
// internal/planner/planner.go
type Planner struct {
    catalog *catalog.Catalog
    router  routing.QueryRouter
    cluster cluster.ClusterManager
    locality *DataLocalityAnalyzer
}

func (p *Planner) CreateDistributedPlan(stmt Statement) (DistributedExecutionPlan, error) {
    // Create local plan first
    localPlan, err := p.CreatePlan(stmt)
    if err != nil {
        return DistributedExecutionPlan{}, err
    }
    
    // Optimize for distribution
    return p.optimizeForDistribution(localPlan)
}

func (p *Planner) optimizeForDistribution(plan Plan) (DistributedExecutionPlan, error) {
    // Analyze data locality
    locality := p.locality.Analyze(plan)
    
    // Determine optimal node placement
    nodePlacement := p.determineNodePlacement(plan, locality)
    
    // Plan data movement if needed
    dataMovement := p.planDataMovement(plan, nodePlacement)
    
    // Create coordination strategy
    coordination := p.createCoordinationStrategy(plan, nodePlacement)
    
    // Choose join strategy for distributed joins
    joinStrategy := p.selectJoinStrategy(plan, locality)
    
    return DistributedExecutionPlan{
        SubPlans:      nodePlacement,
        Coordination:  coordination,
        DataMovement:  dataMovement,
        JoinStrategy:  joinStrategy,
        LocalityCost:  locality.Cost,
    }, nil
}

// RDBMS-specific planning methods
func (p *Planner) planDistributedJoin(left, right Table, condition Expression) (DistributedExecutionPlan, error) {
    // Analyze join condition for distribution
    // Choose optimal join strategy (broadcast, shuffle, colocated)
    // Plan data movement for non-colocated joins
}

func (p *Planner) planGlobalIndexScan(indexName string, predicate Expression) (DistributedExecutionPlan, error) {
    // Route index lookup to appropriate nodes
    // Plan parallel index scans
    // Merge results from multiple nodes
}
```

### 5.4 Storage Engine Architecture (Enhanced for Distributed RDBMS)
```go
// internal/storage/engine.go
type StorageEngine struct {
    localStorage LocalStorageEngine
    network      network.RPCClient
    cluster      cluster.ClusterManager
    router       routing.QueryRouter
    consensus    consensus.ConsensusManager
    globalIndex  *GlobalIndexManager
}

func (se *StorageEngine) Read(key Key) (Value, error) {
    // Check if key is local
    if se.isLocalKey(key) {
        return se.localStorage.Read(key)
    }
    
    // Route to appropriate node
    nodeID, err := se.router.RouteRead(key)
    if err != nil {
        return nil, err
    }
    
    return se.network.Call(nodeID, "Read", key, &value)
}

func (se *StorageEngine) MultiRead(keys []Key) (map[Key]Value, error) {
    // Group keys by node
    keyGroups := se.groupKeysByNode(keys)
    
    // Execute parallel reads
    results := make(map[Key]Value)
    var wg sync.WaitGroup
    var mu sync.Mutex
    
    for nodeID, nodeKeys := range keyGroups {
        wg.Add(1)
        go func(nid NodeID, keys []Key) {
            defer wg.Done()
            nodeResults, err := se.readFromNode(nid, keys)
            if err != nil {
                // Handle error
                return
            }
            
            mu.Lock()
            for k, v := range nodeResults {
                results[k] = v
            }
            mu.Unlock()
        }(nodeID, nodeKeys)
    }
    
    wg.Wait()
    return results, nil
}

// RDBMS-specific distributed operations
func (se *StorageEngine) DistributedJoin(left, right Table, condition Expression) (ResultSet, error) {
    // Analyze join condition
    // Choose join strategy based on data locality
    // Execute distributed join
    // Merge results
}

func (se *StorageEngine) GlobalIndexLookup(indexName string, key Key) ([]Key, error) {
    // Route to index coordinator
    // Perform parallel index lookups
    // Merge results
}
```

### 5.5 Transaction Management (Enhanced for Distributed RDBMS)
```go
// internal/transaction/manager.go
type TransactionManager struct {
    localTxnMgr LocalTransactionManager
    network      network.RPCClient
    cluster      cluster.ClusterManager
    coordinator  *DistributedCoordinator
    consensus    consensus.ConsensusManager
}

type DistributedCoordinator struct {
    nodeID       NodeID
    participants map[TransactionID][]NodeID
    twoPhaseLog  map[TransactionID]*TwoPhaseLog
    deadlockDetector *GlobalDeadlockDetector
}

func (tm *TransactionManager) BeginDistributed() (DistributedTransaction, error) {
    txn := DistributedTransaction{
        ID:          generateTransactionID(),
        Coordinator: tm.coordinator.nodeID,
        Status:      TRANSACTION_ACTIVE,
        Timestamp:   getCurrentTimestamp(),
        IsolationLevel: SERIALIZABLE, // Default for distributed
    }
    
    return txn, nil
}

func (tm *TransactionManager) Prepare(txn DistributedTransaction) error {
    // Phase 1: Send prepare to all participants
    prepareResults := make(map[NodeID]bool)
    
    for _, nodeID := range txn.Nodes {
        if nodeID == tm.coordinator.nodeID {
            // Local prepare
            err := tm.localTxnMgr.Prepare(txn.ID)
            prepareResults[nodeID] = err == nil
        } else {
            // Remote prepare
            var result bool
            err := tm.network.Call(nodeID, "Prepare", txn.ID, &result)
            prepareResults[nodeID] = err == nil && result
        }
    }
    
    // Check if all participants are ready
    allReady := true
    for _, ready := range prepareResults {
        if !ready {
            allReady = false
            break
        }
    }
    
    if allReady {
        txn.Status = TRANSACTION_PREPARED
        return nil
    } else {
        // Abort transaction
        return tm.abortTransaction(txn)
    }
}

// RDBMS-specific distributed transaction methods
func (tm *TransactionManager) GlobalDeadlockDetection() ([]TransactionID, error) {
    // Build global wait-for graph
    // Detect cycles across nodes
    // Return transactions to abort
}

func (tm *TransactionManager) DistributedSnapshotIsolation() (Snapshot, error) {
    // Create consistent snapshot across all nodes
    // Use vector clocks for causality tracking
}
```

## 6. Data Distribution Strategies (RDBMS Specific)

### 6.1 Partitioning (RDBMS Aware)
```go
// internal/storage/partitioning.go
type PartitioningStrategy interface {
    GetPartition(key Key) (PartitionID, error)
    GetPartitionRange(partitionID PartitionID) (KeyRange, error)
    RebalancePartitions() error
    GetTablePartitions(tableName string) ([]PartitionID, error)
}

type HashPartitioning struct {
    numPartitions int
    hashFunction  func(Key) uint64
    tableName     string
}

type RangePartitioning struct {
    ranges []KeyRange
    tableName string
}

type KeyRange struct {
    Start Key
    End   Key
    NodeID NodeID
}

// RDBMS-specific partitioning
type TablePartitioning struct {
    TableName string
    Strategy  PartitioningStrategy
    Columns   []string // Partitioning columns
    Partitions map[PartitionID]*Partition
}

type Partition struct {
    ID       PartitionID
    TableName string
    KeyRange  KeyRange
    NodeID    NodeID
    Replicas  []NodeID
    Status    PartitionStatus
}
```

### 6.2 Replication (RDBMS Focused)
```go
// internal/storage/replication.go
type ReplicationManager interface {
    // Replication control
    SetReplicationFactor(factor int) error
    AddReplica(partitionID PartitionID, nodeID NodeID) error
    RemoveReplica(partitionID PartitionID, nodeID NodeID) error
    
    // RDBMS-specific replication
    ReplicateTransaction(txnID TransactionID, operations []Operation) error
    SyncReplicas(partitionID PartitionID) error
    GetConsistentRead(key Key) (Value, error)
    HandleReplicaFailure(failedNode NodeID) error
}

type ReplicationStrategy interface {
    SelectReplica(key Key, operation OperationType) (NodeID, error)
    HandleReplicaFailure(failedNode NodeID) error
    EnsureQuorumRead(partitionID PartitionID) ([]NodeID, error)
}

type QuorumReplication struct {
    replicationFactor int
    readQuorum       int
    writeQuorum      int
}

type LeaderReplication struct {
    leader NodeID
    followers []NodeID
}
```

## 7. Consistency Models (RDBMS Specific)

### 7.1 Consistency Levels
```go
// internal/consistency/levels.go
type ConsistencyLevel int

const (
    CONSISTENCY_EVENTUAL ConsistencyLevel = iota
    CONSISTENCY_READ_COMMITTED
    CONSISTENCY_REPEATABLE_READ
    CONSISTENCY_SERIALIZABLE
    CONSISTENCY_STRONG
)

type ConsistencyManager interface {
    Read(key Key, level ConsistencyLevel) (Value, error)
    Write(key Key, value Value, level ConsistencyLevel) error
    EnsureConsistency(level ConsistencyLevel) error
    GetConsistentSnapshot() (Snapshot, error)
}

// RDBMS-specific consistency
type SerializableConsistency struct {
    vectorClocks map[NodeID]Timestamp
    transactionLog map[TransactionID]*TransactionLog
}

type TransactionLog struct {
    TxnID TransactionID
    Operations []Operation
    Timestamp Timestamp
    VectorClock VectorClock
}
```

### 7.2 Vector Clocks (RDBMS Enhanced)
```go
// internal/consistency/vector_clock.go
type VectorClock map[NodeID]Timestamp

type VersionedValue struct {
    Value       Value
    VectorClock VectorClock
    Timestamp   Timestamp
    TransactionID TransactionID
}

func (vc VectorClock) Compare(other VectorClock) ClockComparison {
    // Implement vector clock comparison
    // Handle concurrent modifications
    // Detect causality violations
}

// RDBMS-specific versioning
type MVCCManager struct {
    versions map[Key][]VersionedValue
    garbageCollector *GarbageCollector
}

func (mvcc *MVCCManager) GetVersion(key Key, timestamp Timestamp) (Value, error) {
    // Find appropriate version based on timestamp
    // Handle concurrent access
    // Maintain consistency
}
```

## 8. Configuration for Distributed RDBMS

### 8.1 Distributed Configuration
```go
// internal/config/distributed.go
type DistributedConfig struct {
    // Cluster configuration
    Cluster struct {
        NodeID       NodeID
        Address      string
        Peers        []string
        Bootstrap    bool
        DataCenter   string
        Rack         string
    }
    
    // Network configuration
    Network struct {
        RPCAddress     string
        HeartbeatInterval time.Duration
        Timeout         time.Duration
        Compression     CompressionType
    }
    
    // Storage configuration
    Storage struct {
        DataDir        string
        ReplicationFactor int
        PartitioningStrategy string
        ConsistencyLevel ConsistencyLevel
    }
    
    // RDBMS-specific configuration
    RDBMS struct {
        IsolationLevel IsolationLevel
        GlobalIndexes  bool
        DistributedJoins bool
        AutoPartitioning bool
        CrossNodeConstraints bool
    }
    
    // Consistency configuration
    Consistency struct {
        DefaultLevel ConsistencyLevel
        ReadRepair   bool
        QuorumRead   bool
        VectorClocks bool
    }
}
```

## 9. Migration Path (RDBMS Specific)

### 9.1 Single Node to Distributed RDBMS
```go
// internal/migration/migrator.go
type MigrationManager interface {
    // Start distributed mode
    EnableDistribution(config DistributedConfig) error
    
    // Migrate existing data
    MigrateData(partitioningStrategy PartitioningStrategy) error
    
    // RDBMS-specific migration
    MigrateTables(tables []string) error
    CreateGlobalIndexes() error
    UpdateConstraints() error
    
    // Add new nodes
    AddNodeToCluster(nodeInfo NodeInfo) error
    
    // Rebalance data
    RebalanceData() error
    RebalanceTables(tables []string) error
}

// RDBMS-specific migration steps
type TableMigration struct {
    TableName string
    SourceNode NodeID
    TargetNodes []NodeID
    Strategy MigrationStrategy
    Status MigrationStatus
}
```

## 10. Testing for Distributed RDBMS

### 10.1 Distributed Testing
```go
// test/distributed/cluster_test.go
func TestDistributedRDBMSQueryExecution(t *testing.T) {
    cluster := setupTestCluster(t, 3)
    defer cleanupTestCluster(t, cluster)
    
    // Test distributed SQL operations
    result, err := cluster.ExecuteQuery("SELECT u.name, o.order_id FROM users u JOIN orders o ON u.id = o.user_id WHERE u.age > 18")
    require.NoError(t, err)
    
    // Verify ACID properties across nodes
    cluster.VerifyACIDProperties()
    
    // Test distributed transactions
    err = cluster.ExecuteTransaction([]string{
        "INSERT INTO users (id, name, age) VALUES (1, 'Alice', 25)",
        "INSERT INTO orders (order_id, user_id, amount) VALUES (1001, 1, 150.00)",
    })
    require.NoError(t, err)
    
    // Verify consistency
    cluster.VerifyConsistency()
}

func TestDistributedJoinStrategies(t *testing.T) {
    cluster := setupTestCluster(t, 3)
    defer cleanupTestCluster(t, cluster)
    
    // Test different join strategies
    strategies := []JoinStrategy{JOIN_BROADCAST, JOIN_SHUFFLE, JOIN_COLOCATED}
    
    for _, strategy := range strategies {
        result, err := cluster.ExecuteJoinWithStrategy("users", "orders", "users.id = orders.user_id", strategy)
        require.NoError(t, err)
        require.NotEmpty(t, result)
    }
}

func TestGlobalIndexes(t *testing.T) {
    cluster := setupTestCluster(t, 3)
    defer cleanupTestCluster(t, cluster)
    
    // Create global index
    err := cluster.CreateGlobalIndex("idx_users_email", "users", []string{"email"})
    require.NoError(t, err)
    
    // Test index lookup
    result, err := cluster.ExecuteQuery("SELECT * FROM users WHERE email = 'alice@example.com'")
    require.NoError(t, err)
    require.NotEmpty(t, result)
}
```

## 11. Test-Driven Development Strategy

### 11.1 TDD Principles for KongDB

#### 11.1.1 Red-Green-Refactor Cycle
```go
// Example TDD cycle for SQL Parser
func TestSelectStatementParsing(t *testing.T) {
    // RED: Write failing test first
    tests := []struct {
        name     string
        input    string
        expected *SelectStatement
    }{
        {
            name:  "simple select",
            input: "SELECT id, name FROM users",
            expected: &SelectStatement{
                Fields: []Expression{/* ... */},
                From:   []TableReference{/* ... */},
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            parser := NewParser()
            result, err := parser.Parse(tt.input)
            
            // Test should fail initially (RED)
            require.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}

// GREEN: Implement minimal code to pass test
func (p *Parser) Parse(input string) (Statement, error) {
    // Minimal implementation to make test pass
    return &SelectStatement{
        Fields: []Expression{},
        From:   []TableReference{},
    }, nil
}

// REFACTOR: Improve implementation while keeping tests green
func (p *Parser) Parse(input string) (Statement, error) {
    // Improved implementation with proper parsing
    lexer := NewLexer(input)
    tokens, err := lexer.Tokenize()
    if err != nil {
        return nil, err
    }
    
    return p.parseStatement(tokens)
}
```

#### 11.1.2 Test-First Development Guidelines
1. **Write Test First**: Always write the test before implementation
2. **Minimal Implementation**: Write the smallest code to make test pass
3. **Refactor Continuously**: Improve code while keeping tests green
4. **Test Coverage**: Aim for 90%+ test coverage
5. **Test Isolation**: Each test should be independent and isolated

### 11.2 Testing Strategy by Component

#### 11.2.1 SQL Parser Testing
```go
// internal/parser/parser_test.go
func TestParserTDD(t *testing.T) {
    // Test 1: Basic SELECT statement
    t.Run("parse basic select", func(t *testing.T) {
        input := "SELECT id FROM users"
        expected := &SelectStatement{
            Fields: []Expression{&ColumnExpression{Name: "id"}},
            From:   []TableReference{&TableReference{Name: "users"}},
        }
        
        parser := NewParser()
        result, err := parser.Parse(input)
        
        require.NoError(t, err)
        assert.Equal(t, expected, result)
    })
    
    // Test 2: SELECT with WHERE clause
    t.Run("parse select with where", func(t *testing.T) {
        input := "SELECT id FROM users WHERE age > 18"
        expected := &SelectStatement{
            Fields: []Expression{&ColumnExpression{Name: "id"}},
            From:   []TableReference{&TableReference{Name: "users"}},
            Where:  &BinaryExpression{
                Left:  &ColumnExpression{Name: "age"},
                Op:    ">",
                Right: &LiteralExpression{Value: 18},
            },
        }
        
        parser := NewParser()
        result, err := parser.Parse(input)
        
        require.NoError(t, err)
        assert.Equal(t, expected, result)
    })
    
    // Test 3: Error handling
    t.Run("parse invalid sql", func(t *testing.T) {
        input := "SELECT FROM" // Invalid SQL
        
        parser := NewParser()
        result, err := parser.Parse(input)
        
        require.Error(t, err)
        assert.Nil(t, result)
        assert.Contains(t, err.Error(), "syntax error")
    })
}
```

#### 11.2.2 Storage Engine Testing
```go
// internal/storage/storage_test.go
func TestStorageEngineTDD(t *testing.T) {
    // Test 1: Basic read/write operations
    t.Run("basic read write", func(t *testing.T) {
        storage := NewStorageEngine()
        key := Key("test-key")
        value := Value("test-value")
        
        // Write test
        err := storage.Write(key, value)
        require.NoError(t, err)
        
        // Read test
        result, err := storage.Read(key)
        require.NoError(t, err)
        assert.Equal(t, value, result)
    })
    
    // Test 2: Transaction operations
    t.Run("transaction operations", func(t *testing.T) {
        storage := NewStorageEngine()
        
        txn, err := storage.BeginTransaction()
        require.NoError(t, err)
        
        // Write within transaction
        err = storage.Write(Key("txn-key"), Value("txn-value"), txn)
        require.NoError(t, err)
        
        // Read within transaction
        result, err := storage.Read(Key("txn-key"), txn)
        require.NoError(t, err)
        assert.Equal(t, Value("txn-value"), result)
        
        // Commit transaction
        err = storage.CommitTransaction(txn)
        require.NoError(t, err)
        
        // Verify data is persisted
        result, err = storage.Read(Key("txn-key"))
        require.NoError(t, err)
        assert.Equal(t, Value("txn-value"), result)
    })
    
    // Test 3: ACID properties
    t.Run("acid properties", func(t *testing.T) {
        storage := NewStorageEngine()
        
        // Test Atomicity
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

#### 11.2.3 Query Optimizer Testing
```go
// internal/optimizer/optimizer_test.go
func TestOptimizerTDD(t *testing.T) {
    // Test 1: Predicate pushdown
    t.Run("predicate pushdown", func(t *testing.T) {
        inputPlan := &SelectPlan{
            Fields: []Expression{&ColumnExpression{Name: "id"}},
            From:   &TableScan{Table: "users"},
            Where:  &BinaryExpression{
                Left:  &ColumnExpression{Name: "age"},
                Op:    ">",
                Right: &LiteralExpression{Value: 18},
            },
        }
        
        optimizer := NewOptimizer()
        result := optimizer.Optimize(inputPlan)
        
        // Verify predicate was pushed down to table scan
        tableScan := result.(*SelectPlan).From.(*TableScan)
        assert.NotNil(t, tableScan.Predicate)
    })
    
    // Test 2: Join reordering
    t.Run("join reordering", func(t *testing.T) {
        inputPlan := &JoinPlan{
            Left:  &TableScan{Table: "large_table"},
            Right: &TableScan{Table: "small_table"},
            Condition: &BinaryExpression{
                Left:  &ColumnExpression{Name: "id"},
                Op:    "=",
                Right: &ColumnExpression{Name: "user_id"},
            },
        }
        
        optimizer := NewOptimizer()
        result := optimizer.Optimize(inputPlan)
        
        // Verify smaller table is on the right (for nested loop join)
        joinPlan := result.(*JoinPlan)
        assert.Equal(t, "small_table", joinPlan.Right.(*TableScan).Table)
    })
}
```

#### 11.2.4 Transaction Manager Testing
```go
// internal/transaction/manager_test.go
func TestTransactionManagerTDD(t *testing.T) {
    // Test 1: Basic transaction lifecycle
    t.Run("transaction lifecycle", func(t *testing.T) {
        txnMgr := NewTransactionManager()
        
        txn, err := txnMgr.Begin()
        require.NoError(t, err)
        assert.Equal(t, TRANSACTION_ACTIVE, txn.Status)
        
        err = txnMgr.Commit(txn)
        require.NoError(t, err)
        assert.Equal(t, TRANSACTION_COMMITTED, txn.Status)
    })
    
    // Test 2: MVCC isolation
    t.Run("mvcc isolation", func(t *testing.T) {
        txnMgr := NewTransactionManager()
        storage := NewStorageEngine()
        
        // Write initial value
        storage.Write(Key("test"), Value("initial"))
        
        // Start transaction 1
        txn1, _ := txnMgr.Begin()
        storage.Write(Key("test"), Value("updated"), txn1)
        
        // Start transaction 2 (should see initial value)
        txn2, _ := txnMgr.Begin()
        result, err := storage.Read(Key("test"), txn2)
        require.NoError(t, err)
        assert.Equal(t, Value("initial"), result)
        
        // Commit transaction 1
        txnMgr.Commit(txn1)
        
        // Transaction 2 should still see initial value
        result, err = storage.Read(Key("test"), txn2)
        require.NoError(t, err)
        assert.Equal(t, Value("initial"), result)
    })
    
    // Test 3: Deadlock detection
    t.Run("deadlock detection", func(t *testing.T) {
        txnMgr := NewTransactionManager()
        
        txn1, _ := txnMgr.Begin()
        txn2, _ := txnMgr.Begin()
        
        // Create deadlock: txn1 locks A, txn2 locks B
        txnMgr.LockResource(ResourceID("A"), LOCK_EXCLUSIVE, txn1)
        txnMgr.LockResource(ResourceID("B"), LOCK_EXCLUSIVE, txn2)
        
        // Try to create deadlock
        err1 := txnMgr.LockResource(ResourceID("B"), LOCK_EXCLUSIVE, txn1)
        err2 := txnMgr.LockResource(ResourceID("A"), LOCK_EXCLUSIVE, txn2)
        
        // One should succeed, one should be aborted
        assert.True(t, (err1 == nil && err2 != nil) || (err1 != nil && err2 == nil))
    })
}
```

### 11.3 Integration Testing Strategy

#### 11.3.1 End-to-End Query Testing
```go
// test/integration/query_test.go
func TestEndToEndQueryTDD(t *testing.T) {
    // Test 1: Complete query execution
    t.Run("complete query execution", func(t *testing.T) {
        db := setupTestDatabase(t)
        defer cleanupTestDatabase(t, db)
        
        // Create table
        _, err := db.Exec("CREATE TABLE users (id INT, name VARCHAR(50), age INT)")
        require.NoError(t, err)
        
        // Insert data
        _, err = db.Exec("INSERT INTO users VALUES (1, 'Alice', 25), (2, 'Bob', 30)")
        require.NoError(t, err)
        
        // Execute query
        rows, err := db.Query("SELECT name FROM users WHERE age > 20")
        require.NoError(t, err)
        defer rows.Close()
        
        // Verify results
        var names []string
        for rows.Next() {
            var name string
            err := rows.Scan(&name)
            require.NoError(t, err)
            names = append(names, name)
        }
        
        assert.ElementsMatch(t, []string{"Alice", "Bob"}, names)
    })
    
    // Test 2: Transaction consistency
    t.Run("transaction consistency", func(t *testing.T) {
        db := setupTestDatabase(t)
        defer cleanupTestDatabase(t, db)
        
        // Create table
        db.Exec("CREATE TABLE accounts (id INT, balance DECIMAL(10,2))")
        db.Exec("INSERT INTO accounts VALUES (1, 100.00), (2, 200.00)")
        
        // Start transaction
        tx, err := db.Begin()
        require.NoError(t, err)
        
        // Transfer money
        _, err = tx.Exec("UPDATE accounts SET balance = balance - 50 WHERE id = 1")
        require.NoError(t, err)
        _, err = tx.Exec("UPDATE accounts SET balance = balance + 50 WHERE id = 2")
        require.NoError(t, err)
        
        // Commit transaction
        err = tx.Commit()
        require.NoError(t, err)
        
        // Verify results
        var balance1, balance2 float64
        db.QueryRow("SELECT balance FROM accounts WHERE id = 1").Scan(&balance1)
        db.QueryRow("SELECT balance FROM accounts WHERE id = 2").Scan(&balance2)
        
        assert.Equal(t, 50.0, balance1)
        assert.Equal(t, 250.0, balance2)
    })
}
```

#### 11.3.2 Performance Testing
```go
// test/performance/benchmarks_test.go
func BenchmarkQueryExecution(b *testing.B) {
    db := setupBenchmarkDatabase(b)
    defer cleanupBenchmarkDatabase(b, db)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        rows, err := db.Query("SELECT * FROM users WHERE age > 18")
        if err != nil {
            b.Fatal(err)
        }
        rows.Close()
    }
}

func BenchmarkTransactionThroughput(b *testing.B) {
    db := setupBenchmarkDatabase(b)
    defer cleanupBenchmarkDatabase(b, db)
    
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            tx, err := db.Begin()
            if err != nil {
                b.Fatal(err)
            }
            
            _, err = tx.Exec("INSERT INTO users (name, age) VALUES (?, ?)", "test", 25)
            if err != nil {
                b.Fatal(err)
            }
            
            err = tx.Commit()
            if err != nil {
                b.Fatal(err)
            }
        }
    })
}
```

### 11.4 Distributed Testing Strategy

#### 11.4.1 Cluster Testing
```go
// test/distributed/cluster_test.go
func TestDistributedClusterTDD(t *testing.T) {
    // Test 1: Basic cluster operations
    t.Run("basic cluster operations", func(t *testing.T) {
        cluster := setupTestCluster(t, 3)
        defer cleanupTestCluster(t, cluster)
        
        // Test cluster health
        status := cluster.GetClusterStatus()
        assert.Equal(t, 3, len(status.Nodes))
        assert.Equal(t, CLUSTER_HEALTHY, status.Health)
    })
    
    // Test 2: Distributed query execution
    t.Run("distributed query execution", func(t *testing.T) {
        cluster := setupTestCluster(t, 3)
        defer cleanupTestCluster(t, cluster)
        
        // Create distributed table
        err := cluster.ExecuteQuery("CREATE TABLE users (id INT, name VARCHAR(50)) DISTRIBUTE BY HASH(id)")
        require.NoError(t, err)
        
        // Insert data
        err = cluster.ExecuteQuery("INSERT INTO users VALUES (1, 'Alice'), (2, 'Bob')")
        require.NoError(t, err)
        
        // Execute distributed query
        result, err := cluster.ExecuteQuery("SELECT * FROM users WHERE id = 1")
        require.NoError(t, err)
        assert.NotEmpty(t, result)
    })
    
    // Test 3: Node failure handling
    t.Run("node failure handling", func(t *testing.T) {
        cluster := setupTestCluster(t, 3)
        defer cleanupTestCluster(t, cluster)
        
        // Kill a node
        cluster.KillNode(1)
        
        // Verify cluster continues to operate
        result, err := cluster.ExecuteQuery("SELECT 1")
        require.NoError(t, err)
        assert.NotEmpty(t, result)
        
        // Restart node
        cluster.RestartNode(1)
        
        // Verify cluster is healthy again
        status := cluster.GetClusterStatus()
        assert.Equal(t, CLUSTER_HEALTHY, status.Health)
    })
}
```

### 11.5 Test Infrastructure

#### 11.5.1 Test Utilities
```go
// test/utils/test_helpers.go
type TestDatabase struct {
    Storage     StorageEngine
    Parser      Parser
    Optimizer   Optimizer
    Planner     Planner
    Executor    Executor
    TxnManager  TransactionManager
}

func setupTestDatabase(t *testing.T) *TestDatabase {
    storage := NewInMemoryStorage()
    parser := NewParser()
    optimizer := NewOptimizer()
    planner := NewPlanner()
    executor := NewExecutor()
    txnManager := NewTransactionManager()
    
    return &TestDatabase{
        Storage:    storage,
        Parser:     parser,
        Optimizer:  optimizer,
        Planner:    planner,
        Executor:   executor,
        TxnManager: txnManager,
    }
}

func cleanupTestDatabase(t *testing.T, db *TestDatabase) {
    // Cleanup resources
}

type TestCluster struct {
    Nodes []*TestNode
}

func setupTestCluster(t *testing.T, nodeCount int) *TestCluster {
    cluster := &TestCluster{}
    
    for i := 0; i < nodeCount; i++ {
        node := NewTestNode(fmt.Sprintf("node-%d", i))
        cluster.Nodes = append(cluster.Nodes, node)
    }
    
    // Start cluster
    err := cluster.Start()
    require.NoError(t, err)
    
    return cluster
}
```

#### 11.5.2 Mock Implementations
```go
// test/mocks/mock_storage.go
type MockStorageEngine struct {
    data map[Key]Value
    mu   sync.RWMutex
}

func NewMockStorageEngine() *MockStorageEngine {
    return &MockStorageEngine{
        data: make(map[Key]Value),
    }
}

func (m *MockStorageEngine) Read(key Key) (Value, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    if value, exists := m.data[key]; exists {
        return value, nil
    }
    return nil, errors.New("key not found")
}

func (m *MockStorageEngine) Write(key Key, value Value) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    m.data[key] = value
    return nil
}
```

## 12. TDD Development Workflow

### 12.1 Component Development Process
1. **Write Test First**: Define expected behavior through tests
2. **Run Test**: Verify it fails (RED)
3. **Write Minimal Code**: Implement just enough to pass test (GREEN)
4. **Refactor**: Improve code while keeping tests green (REFACTOR)
5. **Repeat**: Add more test cases and functionality

### 12.2 Feature Development Process
1. **Acceptance Test**: Write high-level acceptance test
2. **Unit Tests**: Write unit tests for each component
3. **Integration Tests**: Test component interactions
4. **Performance Tests**: Ensure performance requirements are met
5. **Documentation**: Update documentation with examples

### 12.3 Continuous Integration
```yaml
# .github/workflows/test.yml
name: Test-Driven Development

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.21
    
    - name: Run unit tests
      run: go test -v ./internal/...
    
    - name: Run integration tests
      run: go test -v ./test/integration/...
    
    - name: Run performance tests
      run: go test -bench=. ./test/performance/...
    
    - name: Check test coverage
      run: go test -coverprofile=coverage.out ./...
    
    - name: Upload coverage
      uses: codecov/codecov-action@v1
```

This TDD approach ensures that KongDB is built with comprehensive testing from the ground up, making it more reliable, maintainable, and easier to extend with distributed features. 

## 13. Bloom Filter Architecture

### 13.1 Core Bloom Filter Implementation

#### 13.1.1 Basic Bloom Filter
```go
// internal/storage/bloom_filter/bloom_filter.go
type BloomFilter struct {
    bitset    []bool
    hashFuncs []HashFunction
    size      int
    capacity  int
    count     int
    mu        sync.RWMutex
}

type HashFunction func(data []byte) uint64

type BloomFilterConfig struct {
    FalsePositiveRate float64  // e.g., 0.01 (1%)
    ExpectedElements  int      // Expected number of keys
    HashFunctions     int      // Number of hash functions
    OptimalSize       int      // Calculated optimal size
}

func NewBloomFilter(config BloomFilterConfig) *BloomFilter {
    size := calculateOptimalSize(config.ExpectedElements, config.FalsePositiveRate)
    hashFuncs := generateHashFunctions(config.HashFunctions)
    
    return &BloomFilter{
        bitset:    make([]bool, size),
        hashFuncs: hashFuncs,
        size:      size,
        capacity:  config.ExpectedElements,
        count:     0,
    }
}

func (bf *BloomFilter) Insert(key Key) {
    bf.mu.Lock()
    defer bf.mu.Unlock()
    
    keyBytes := key.Bytes()
    for _, hashFunc := range bf.hashFuncs {
        hash := hashFunc(keyBytes)
        index := hash % uint64(bf.size)
        bf.bitset[index] = true
    }
    bf.count++
}

func (bf *BloomFilter) MightContain(key Key) bool {
    bf.mu.RLock()
    defer bf.mu.RUnlock()
    
    keyBytes := key.Bytes()
    for _, hashFunc := range bf.hashFuncs {
        hash := hashFunc(keyBytes)
        index := hash % uint64(bf.size)
        if !bf.bitset[index] {
            return false // Definitely not in the set
        }
    }
    return true // Might be in the set (could be false positive)
}

func (bf *BloomFilter) Clear() {
    bf.mu.Lock()
    defer bf.mu.Unlock()
    
    for i := range bf.bitset {
        bf.bitset[i] = false
    }
    bf.count = 0
}

func (bf *BloomFilter) GetStatistics() BloomFilterStats {
    bf.mu.RLock()
    defer bf.mu.RUnlock()
    
    return BloomFilterStats{
        Size:              bf.size,
        Capacity:          bf.capacity,
        CurrentElements:   bf.count,
        FalsePositiveRate: bf.calculateFalsePositiveRate(),
        MemoryUsage:       bf.size / 8, // bits to bytes
    }
}
```

#### 13.1.2 Hash Function Generation
```go
// internal/storage/bloom_filter/hash_functions.go
type HashFunctionGenerator struct {
    seeds []uint64
}

func NewHashFunctionGenerator(numFunctions int) *HashFunctionGenerator {
    seeds := make([]uint64, numFunctions)
    for i := 0; i < numFunctions; i++ {
        seeds[i] = uint64(i*2654435761 + 1) // FNV-like seeds
    }
    
    return &HashFunctionGenerator{seeds: seeds}
}

func (hfg *HashFunctionGenerator) GenerateHashFunctions() []HashFunction {
    functions := make([]HashFunction, len(hfg.seeds))
    
    for i, seed := range hfg.seeds {
        functions[i] = func(data []byte) uint64 {
            return hfg.hashWithSeed(data, seed)
        }
    }
    
    return functions
}

func (hfg *HashFunctionGenerator) hashWithSeed(data []byte, seed uint64) uint64 {
    hash := seed
    for _, b := range data {
        hash = hash*1099511628211 + uint64(b) // FNV-1a
    }
    return hash
}
```

#### 13.1.3 Bloom Filter Statistics
```go
// internal/storage/bloom_filter/statistics.go
type BloomFilterStats struct {
    Size              int     `json:"size"`
    Capacity          int     `json:"capacity"`
    CurrentElements   int     `json:"current_elements"`
    FalsePositiveRate float64 `json:"false_positive_rate"`
    MemoryUsage       int     `json:"memory_usage_bytes"`
    HitRate           float64 `json:"hit_rate"`
    MissRate          float64 `json:"miss_rate"`
}

func (bf *BloomFilter) calculateFalsePositiveRate() float64 {
    if bf.count == 0 {
        return 0.0
    }
    
    // Theoretical false positive rate
    k := float64(len(bf.hashFuncs))
    m := float64(bf.size)
    n := float64(bf.count)
    
    return math.Pow(1-math.Exp(-k*n/m), k)
}
```

### 13.2 Bloom Filter Manager

#### 13.2.1 Centralized Management
```go
// internal/storage/bloom_filter/manager.go
type BloomFilterManager struct {
    filters    map[string]*BloomFilter
    configs    map[string]BloomFilterConfig
    statistics *BloomFilterStatistics
    mu         sync.RWMutex
}

type BloomFilterStatistics struct {
    TotalFilters     int                    `json:"total_filters"`
    TotalMemoryUsage int                    `json:"total_memory_usage"`
    FilterStats      map[string]BloomFilterStats `json:"filter_stats"`
    PerformanceStats BloomFilterPerformance `json:"performance_stats"`
}

type BloomFilterPerformance struct {
    TotalQueries     int64   `json:"total_queries"`
    PositiveQueries  int64   `json:"positive_queries"`
    NegativeQueries  int64   `json:"negative_queries"`
    FalsePositives   int64   `json:"false_positives"`
    AverageQueryTime float64 `json:"average_query_time_ns"`
}

func NewBloomFilterManager() *BloomFilterManager {
    return &BloomFilterManager{
        filters:    make(map[string]*BloomFilter),
        configs:    make(map[string]BloomFilterConfig),
        statistics: &BloomFilterStatistics{
            FilterStats: make(map[string]BloomFilterStats),
        },
    }
}

func (bfm *BloomFilterManager) CreateBloomFilter(name string, config BloomFilterConfig) error {
    bfm.mu.Lock()
    defer bfm.mu.Unlock()
    
    if _, exists := bfm.filters[name]; exists {
        return errors.New("bloom filter already exists")
    }
    
    bfm.filters[name] = NewBloomFilter(config)
    bfm.configs[name] = config
    
    return nil
}

func (bfm *BloomFilterManager) GetBloomFilter(name string) (*BloomFilter, error) {
    bfm.mu.RLock()
    defer bfm.mu.RUnlock()
    
    filter, exists := bfm.filters[name]
    if !exists {
        return nil, errors.New("bloom filter not found")
    }
    
    return filter, nil
}

func (bfm *BloomFilterManager) UpdateBloomFilter(name string, key Key) error {
    filter, err := bfm.GetBloomFilter(name)
    if err != nil {
        return err
    }
    
    filter.Insert(key)
    return nil
}

func (bfm *BloomFilterManager) GetStatistics() BloomFilterStatistics {
    bfm.mu.RLock()
    defer bfm.mu.RUnlock()
    
    stats := BloomFilterStatistics{
        TotalFilters: len(bfm.filters),
        FilterStats:  make(map[string]BloomFilterStats),
    }
    
    for name, filter := range bfm.filters {
        stats.FilterStats[name] = filter.GetStatistics()
        stats.TotalMemoryUsage += filter.GetStatistics().MemoryUsage
    }
    
    return stats
}
```

### 13.3 Storage Engine Integration

#### 13.3.1 Enhanced Storage Engine
```go
// internal/storage/engine.go
type StorageEngine struct {
    localStorage     LocalStorageEngine
    bloomFilterMgr   *bloom_filter.BloomFilterManager
    network          network.RPCClient
    cluster          cluster.ClusterManager
    router           routing.QueryRouter
    consensus        consensus.ConsensusManager
    globalIndex      *GlobalIndexManager
}

func NewStorageEngine() *StorageEngine {
    return &StorageEngine{
        localStorage:   NewLocalStorageEngine(),
        bloomFilterMgr: bloom_filter.NewBloomFilterManager(),
        // ... other initializations
    }
}

func (se *StorageEngine) Read(key Key) (Value, error) {
    // Check bloom filter first for fast rejection
    tableName := se.getTableNameForKey(key)
    bloomFilter, err := se.bloomFilterMgr.GetBloomFilter(tableName)
    if err == nil && !bloomFilter.MightContain(key) {
        return nil, errors.New("key not found")
    }
    
    // Check if key is local
    if se.isLocalKey(key) {
        return se.localStorage.Read(key)
    }
    
    // Route to appropriate node
    nodeID, err := se.router.RouteRead(key)
    if err != nil {
        return nil, err
    }
    
    return se.network.Call(nodeID, "Read", key, &value)
}

func (se *StorageEngine) Write(key Key, value Value) error {
    // Update bloom filter when writing
    tableName := se.getTableNameForKey(key)
    se.bloomFilterMgr.UpdateBloomFilter(tableName, key)
    
    // Check if key is local
    if se.isLocalKey(key) {
        return se.localStorage.Write(key, value)
    }
    
    // Route to appropriate node
    nodeID, err := se.router.RouteWrite(WriteOperation{Key: key, Value: value})
    if err != nil {
        return err
    }
    
    return se.network.Call(nodeID, "Write", WriteOperation{Key: key, Value: value}, &result)
}

func (se *StorageEngine) MultiRead(keys []Key) (map[Key]Value, error) {
    // Group keys by table and check bloom filters
    tableGroups := se.groupKeysByTable(keys)
    results := make(map[Key]Value)
    
    for tableName, tableKeys := range tableGroups {
        bloomFilter, err := se.bloomFilterMgr.GetBloomFilter(tableName)
        if err != nil {
            // No bloom filter for this table, proceed normally
            tableResults, err := se.readKeysFromTable(tableKeys)
            if err != nil {
                return nil, err
            }
            for k, v := range tableResults {
                results[k] = v
            }
            continue
        }
        
        // Filter keys using bloom filter
        filteredKeys := make([]Key, 0)
        for _, key := range tableKeys {
            if bloomFilter.MightContain(key) {
                filteredKeys = append(filteredKeys, key)
            }
        }
        
        // Read only the keys that might exist
        if len(filteredKeys) > 0 {
            tableResults, err := se.readKeysFromTable(filteredKeys)
            if err != nil {
                return nil, err
            }
            for k, v := range tableResults {
                results[k] = v
            }
        }
    }
    
    return results, nil
}
```

#### 13.3.2 Table-Specific Bloom Filters
```go
// internal/storage/bloom_filter/table_bloom_filter.go
type TableBloomFilter struct {
    tableName    string
    bloomFilter  *BloomFilter
    columnFilter map[string]*BloomFilter // Column-specific filters
    config       TableBloomFilterConfig
}

type TableBloomFilterConfig struct {
    EnableTableFilter    bool    `json:"enable_table_filter"`
    EnableColumnFilters  bool    `json:"enable_column_filters"`
    FalsePositiveRate    float64 `json:"false_positive_rate"`
    ExpectedTableSize    int     `json:"expected_table_size"`
    ColumnConfigs        map[string]ColumnBloomFilterConfig `json:"column_configs"`
}

type ColumnBloomFilterConfig struct {
    EnableFilter       bool    `json:"enable_filter"`
    FalsePositiveRate  float64 `json:"false_positive_rate"`
    ExpectedCardinality int    `json:"expected_cardinality"`
}

func NewTableBloomFilter(tableName string, config TableBloomFilterConfig) *TableBloomFilter {
    tbf := &TableBloomFilter{
        tableName:   tableName,
        columnFilter: make(map[string]*BloomFilter),
        config:      config,
    }
    
    if config.EnableTableFilter {
        tbf.bloomFilter = NewBloomFilter(BloomFilterConfig{
            FalsePositiveRate: config.FalsePositiveRate,
            ExpectedElements:  config.ExpectedTableSize,
            HashFunctions:     5, // Default
        })
    }
    
    if config.EnableColumnFilters {
        for columnName, columnConfig := range config.ColumnConfigs {
            if columnConfig.EnableFilter {
                tbf.columnFilter[columnName] = NewBloomFilter(BloomFilterConfig{
                    FalsePositiveRate: columnConfig.FalsePositiveRate,
                    ExpectedElements:  columnConfig.ExpectedCardinality,
                    HashFunctions:     5,
                })
            }
        }
    }
    
    return tbf
}

func (tbf *TableBloomFilter) InsertRow(row Row) {
    if tbf.bloomFilter != nil {
        // Insert primary key into table bloom filter
        primaryKey := row.GetPrimaryKey()
        tbf.bloomFilter.Insert(primaryKey)
    }
    
    // Insert column values into column bloom filters
    for columnName, bloomFilter := range tbf.columnFilter {
        if value, exists := row.GetValue(columnName); exists {
            bloomFilter.Insert(Key(value.String()))
        }
    }
}

func (tbf *TableBloomFilter) MightContainKey(key Key) bool {
    if tbf.bloomFilter == nil {
        return true // No filter, assume it might exist
    }
    return tbf.bloomFilter.MightContain(key)
}

func (tbf *TableBloomFilter) MightContainColumnValue(columnName string, value Value) bool {
    bloomFilter, exists := tbf.columnFilter[columnName]
    if !exists {
        return true // No filter for this column
    }
    return bloomFilter.MightContain(Key(value.String()))
}
```

### 13.4 Query Optimizer Integration

#### 13.4.1 Bloom Filter Optimizer
```go
// internal/optimizer/bloom_filter_optimizer.go
type BloomFilterOptimizer struct {
    storageEngine StorageEngine
    statistics    *BloomFilterStatistics
}

func (bfo *BloomFilterOptimizer) Optimize(plan Plan) Plan {
    switch p := plan.(type) {
    case *TableScanPlan:
        return bfo.optimizeTableScan(p)
    case *IndexScanPlan:
        return bfo.optimizeIndexScan(p)
    case *JoinPlan:
        return bfo.optimizeJoin(p)
    default:
        return plan
    }
}

func (bfo *BloomFilterOptimizer) optimizeTableScan(plan *TableScanPlan) Plan {
    // Check if we can use bloom filter to optimize the scan
    if bfo.canUseBloomFilter(plan.Table, plan.Predicate) {
        return &BloomFilterTableScanPlan{
            TableScanPlan: *plan,
            BloomFilter:   bfo.getBloomFilter(plan.Table),
        }
    }
    return plan
}

func (bfo *BloomFilterOptimizer) optimizeIndexScan(plan *IndexScanPlan) Plan {
    // Use bloom filter to check if index might contain the key
    if bfo.canUseBloomFilter(plan.IndexName, plan.Predicate) {
        return &BloomFilterIndexScanPlan{
            IndexScanPlan: *plan,
            BloomFilter:   bfo.getBloomFilter(plan.IndexName),
        }
    }
    return plan
}

func (bfo *BloomFilterOptimizer) optimizeJoin(plan *JoinPlan) Plan {
    // Use bloom filter for join optimization
    if bfo.canUseBloomFilterForJoin(plan) {
        return &BloomFilterJoinPlan{
            JoinPlan:     *plan,
            LeftFilter:   bfo.getBloomFilter(plan.LeftTable),
            RightFilter:  bfo.getBloomFilter(plan.RightTable),
        }
    }
    return plan
}

func (bfo *BloomFilterOptimizer) canUseBloomFilter(tableName string, predicate Expression) bool {
    // Check if predicate can be optimized with bloom filter
    if predicate == nil {
        return false
    }
    
    // Check if bloom filter exists for this table
    _, err := bfo.storageEngine.GetBloomFilter(tableName)
    if err != nil {
        return false
    }
    
    // Check if predicate is suitable for bloom filter optimization
    return bfo.isPredicateSuitable(predicate)
}

func (bfo *BloomFilterOptimizer) isPredicateSuitable(predicate Expression) bool {
    // Check if predicate can be evaluated using bloom filter
    switch p := predicate.(type) {
    case *EqualityExpression:
        return true
    case *InExpression:
        return true
    case *RangeExpression:
        // Bloom filters are less effective for ranges
        return false
    default:
        return false
    }
}
```

#### 13.4.2 Bloom Filter Execution Plans
```go
// internal/planner/bloom_filter_plans.go
type BloomFilterTableScanPlan struct {
    TableScanPlan
    BloomFilter *BloomFilter
}

func (bfp *BloomFilterTableScanPlan) Execute(ctx ExecutionContext) (ResultSet, error) {
    // Check bloom filter first
    if bfp.Predicate != nil {
        key := bfp.extractKeyFromPredicate(bfp.Predicate)
        if !bfp.BloomFilter.MightContain(key) {
            return EmptyResultSet{}, nil // No results
        }
    }
    
    // Proceed with normal table scan
    return bfp.TableScanPlan.Execute(ctx)
}

type BloomFilterJoinPlan struct {
    JoinPlan
    LeftFilter  *BloomFilter
    RightFilter *BloomFilter
}

func (bfp *BloomFilterJoinPlan) Execute(ctx ExecutionContext) (ResultSet, error) {
    // Use bloom filter to filter one side of the join
    if bfp.LeftFilter != nil {
        // Filter right table using left table's bloom filter
        filteredRight := bfp.filterWithBloomFilter(bfp.RightTable, bfp.LeftFilter, bfp.Condition)
        return bfp.performJoin(bfp.LeftTable, filteredRight, bfp.Condition)
    }
    
    // Fall back to normal join
    return bfp.JoinPlan.Execute(ctx)
}
```

### 13.5 Distributed Bloom Filters

#### 13.5.1 Distributed Bloom Filter Manager
```go
// internal/storage/bloom_filter/distributed_manager.go
type DistributedBloomFilterManager struct {
    localManager    *BloomFilterManager
    remoteFilters   map[NodeID]*RemoteBloomFilter
    syncStrategy    BloomFilterSyncStrategy
    network         network.RPCClient
    cluster         cluster.ClusterManager
}

type RemoteBloomFilter struct {
    NodeID      NodeID
    BloomFilter *BloomFilter
    LastSync    time.Time
    SyncStatus  SyncStatus
}

type BloomFilterSyncStrategy int
const (
    SYNC_ON_WRITE BloomFilterSyncStrategy = iota
    SYNC_PERIODIC
    SYNC_ON_DEMAND
)

func NewDistributedBloomFilterManager(network network.RPCClient, cluster cluster.ClusterManager) *DistributedBloomFilterManager {
    return &DistributedBloomFilterManager{
        localManager:  NewBloomFilterManager(),
        remoteFilters: make(map[NodeID]*RemoteBloomFilter),
        syncStrategy:  SYNC_PERIODIC,
        network:       network,
        cluster:       cluster,
    }
}

func (dbfm *DistributedBloomFilterManager) MightContain(key Key) bool {
    // Check local bloom filter first
    localResult := dbfm.localManager.MightContain(key)
    if !localResult {
        return false // Definitely not in local data
    }
    
    // Check remote bloom filters
    for nodeID, remoteFilter := range dbfm.remoteFilters {
        if remoteFilter.BloomFilter.MightContain(key) {
            // Key might exist on this remote node
            return true
        }
    }
    
    return localResult
}

func (dbfm *DistributedBloomFilterManager) SyncBloomFilters() error {
    nodes := dbfm.cluster.GetAllNodes()
    
    for _, node := range nodes {
        if node.ID == dbfm.cluster.GetLocalNodeID() {
            continue // Skip local node
        }
        
        // Sync bloom filters with remote node
        err := dbfm.syncWithNode(node.ID)
        if err != nil {
            // Log error but continue with other nodes
            log.Printf("Failed to sync bloom filter with node %s: %v", node.ID, err)
        }
    }
    
    return nil
}

func (dbfm *DistributedBloomFilterManager) syncWithNode(nodeID NodeID) error {
    // Get local bloom filter data
    localFilters := dbfm.localManager.GetAllFilters()
    
    // Send to remote node
    var response BloomFilterSyncResponse
    err := dbfm.network.Call(nodeID, "SyncBloomFilters", localFilters, &response)
    if err != nil {
        return err
    }
    
    // Update remote filters
    dbfm.remoteFilters[nodeID] = &RemoteBloomFilter{
        NodeID:      nodeID,
        BloomFilter: response.BloomFilter,
        LastSync:    time.Now(),
        SyncStatus:  SYNC_SUCCESS,
    }
    
    return nil
}
```

#### 13.5.2 Bloom Filter Serialization
```go
// internal/storage/bloom_filter/serialization.go
type BloomFilterData struct {
    Name           string    `json:"name"`
    Bitset         []bool    `json:"bitset"`
    HashFunctions  int       `json:"hash_functions"`
    Size           int       `json:"size"`
    Capacity       int       `json:"capacity"`
    Count          int       `json:"count"`
    LastUpdated    time.Time `json:"last_updated"`
}

func (bf *BloomFilter) Serialize() BloomFilterData {
    bf.mu.RLock()
    defer bf.mu.RUnlock()
    
    return BloomFilterData{
        Name:          bf.name,
        Bitset:        bf.bitset,
        HashFunctions: len(bf.hashFuncs),
        Size:          bf.size,
        Capacity:      bf.capacity,
        Count:         bf.count,
        LastUpdated:   time.Now(),
    }
}

func (bf *BloomFilter) Deserialize(data BloomFilterData) error {
    bf.mu.Lock()
    defer bf.mu.Unlock()
    
    bf.bitset = data.Bitset
    bf.size = data.Size
    bf.capacity = data.Capacity
    bf.count = data.Count
    
    // Regenerate hash functions
    generator := NewHashFunctionGenerator(data.HashFunctions)
    bf.hashFuncs = generator.GenerateHashFunctions()
    
    return nil
}
```

## 14. Bloom Filter Testing Strategy

### 14.1 Unit Testing
```go
// internal/storage/bloom_filter/bloom_filter_test.go
func TestBloomFilterBasicOperations(t *testing.T) {
    t.Run("insert and query", func(t *testing.T) {
        config := BloomFilterConfig{
            FalsePositiveRate: 0.01,
            ExpectedElements:  1000,
            HashFunctions:     5,
        }
        
        bf := NewBloomFilter(config)
        
        // Insert keys
        keys := []Key{"key1", "key2", "key3", "key4", "key5"}
        for _, key := range keys {
            bf.Insert(key)
        }
        
        // Query existing keys (should all return true)
        for _, key := range keys {
            assert.True(t, bf.MightContain(key), "Existing key should be found")
        }
        
        // Query non-existing keys
        nonExisting := []Key{"key6", "key7", "key8", "key9", "key10"}
        falsePositives := 0
        for _, key := range nonExisting {
            if bf.MightContain(key) {
                falsePositives++
            }
        }
        
        // False positive rate should be close to expected
        actualRate := float64(falsePositives) / float64(len(nonExisting))
        assert.Less(t, actualRate, 0.05, "False positive rate should be low")
    })
}

func TestBloomFilterPerformance(t *testing.T) {
    t.Run("performance comparison", func(t *testing.T) {
        config := BloomFilterConfig{
            FalsePositiveRate: 0.01,
            ExpectedElements:  100000,
            HashFunctions:     5,
        }
        
        bf := NewBloomFilter(config)
        
        // Insert many keys
        for i := 0; i < 100000; i++ {
            bf.Insert(Key(fmt.Sprintf("key%d", i)))
        }
        
        // Benchmark bloom filter lookup
        start := time.Now()
        for i := 0; i < 10000; i++ {
            bf.MightContain(Key(fmt.Sprintf("key%d", i)))
        }
        duration := time.Since(start)
        
        // Should be very fast (microseconds per lookup)
        avgTime := duration.Nanoseconds() / 10000
        assert.Less(t, avgTime, int64(1000), "Bloom filter lookup should be fast")
    })
}
```

### 14.2 Integration Testing
```go
// test/integration/bloom_filter_test.go
func TestBloomFilterStorageIntegration(t *testing.T) {
    t.Run("storage engine with bloom filter", func(t *testing.T) {
        storage := NewStorageEngine()
        
        // Create table with bloom filter
        tableName := "users"
        config := TableBloomFilterConfig{
            EnableTableFilter:   true,
            EnableColumnFilters: true,
            FalsePositiveRate:   0.01,
            ExpectedTableSize:   10000,
        }
        
        err := storage.CreateTableBloomFilter(tableName, config)
        require.NoError(t, err)
        
        // Insert data
        for i := 0; i < 1000; i++ {
            key := Key(fmt.Sprintf("user_%d", i))
            value := Value(fmt.Sprintf("user_data_%d", i))
            err := storage.Write(key, value)
            require.NoError(t, err)
        }
        
        // Test bloom filter effectiveness
        start := time.Now()
        for i := 0; i < 1000; i++ {
            key := Key(fmt.Sprintf("user_%d", i))
            _, err := storage.Read(key)
            require.NoError(t, err)
        }
        withBloomFilter := time.Since(start)
        
        // Compare with storage without bloom filter
        storageNoBF := NewStorageEngineWithoutBloomFilter()
        // ... same operations
        withoutBloomFilter := time.Since(start)
        
        // Bloom filter should improve performance
        assert.Less(t, withBloomFilter, withoutBloomFilter, "Bloom filter should improve performance")
    })
}
```

This comprehensive bloom filter integration provides significant performance benefits while maintaining the modular, testable architecture of KongDB. The implementation covers both single-node and distributed scenarios, making it valuable for the eventual distributed version of the database. 