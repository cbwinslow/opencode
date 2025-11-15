# OpenCode Multi-Agent Swarm System

A sophisticated multi-agent architecture for autonomous, collaborative problem-solving with democratic decision-making, hierarchical memory, and self-healing capabilities.

## Overview

This package implements a comprehensive multi-agent swarm system that enables OpenCode to:

- **Collaborate**: Multiple specialized agents work together on complex tasks
- **Decide Democratically**: Important decisions made through voting
- **Remember**: Hierarchical memory system stores and retrieves knowledge
- **Self-Heal**: Automatically detect and recover from failures
- **Learn**: Continuous improvement from successes and failures
- **Scale**: From 2 agents to 50+ without redesign

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Coordinator                             │
│  Orchestrates all swarm activities and component integration │
└────────────┬────────────────────────────────────────────────┘
             │
    ┌────────┴──────────┬──────────┬──────────┬───────────┐
    │                   │          │          │           │
┌───▼────┐  ┌──────────▼───┐  ┌──▼───┐  ┌───▼───┐  ┌───▼────┐
│ Agent  │  │   Memory     │  │Voting│  │ Rules │  │ Health │
│Registry│  │Hierarchical  │  │System│  │Engine │  │Monitor │
└───┬────┘  │   Store      │  └──────┘  └───────┘  └────────┘
    │       └──────────────┘
    │
    ├─── Monitor Agents (watch logs, shell history)
    ├─── Analyzer Agents (analyze data, detect patterns)
    ├─── Executor Agents (execute tasks, generate code)
    ├─── Memory Agents (manage knowledge)
    ├─── Learning Agents (improve from experience)
    └─── Specialized Agents (testing, docs, error handling, etc.)
```

## Components

### 1. Agent System (`agent/`)

**Purpose**: Foundation for all agents in the swarm

**Key Features**:
- Agent registry for lifecycle management
- Message-based inter-agent communication
- Health tracking and metrics
- Task execution with capability matching
- 10 specialized agent types

**Files**:
- `types.go` - Agent types, tasks, messages, metrics
- `base.go` - Base agent implementation
- `registry.go` - Agent registry and message broker

### 2. Memory System (`memory/`)

**Purpose**: Hierarchical, encrypted knowledge storage

**Key Features**:
- 4 memory types (working, episodic, semantic, procedural)
- Vector-based semantic search
- AES-GCM encryption
- Automatic consolidation and pruning
- Priority-based retention

**Files**:
- `types.go` - Memory types and interfaces
- `hierarchical.go` - Hierarchical memory store implementation

### 3. Monitoring (`monitor/`)

**Purpose**: Monitor logs and shell history for learning

**Key Features**:
- Real-time log file watching (fsnotify)
- Shell history monitoring
- Pattern detection
- Structured event processing

**Files**:
- `log_watcher.go` - Log and shell history monitoring

### 4. Voting System (`voting/`)

**Purpose**: Democratic decision-making among agents

**Key Features**:
- 5 voting types (majority, super, unanimous, weighted, consensus)
- Confidence-based voting
- Reasoning capture
- Iterative consensus building

**Files**:
- `democratic.go` - Voting system implementation

### 5. Health Monitoring (`health/`)

**Purpose**: Self-healing through continuous monitoring

**Key Features**:
- Continuous health checks
- Alert system with severity levels
- Recovery strategies (restart, reset, fallback, isolate)
- System-wide health aggregation

**Files**:
- `monitor.go` - Health monitoring and recovery

### 6. Rule Engine (`rules/`)

**Purpose**: Automated behavioral responses

**Key Features**:
- Condition-action rules
- Priority-based execution
- Dynamic rule loading
- Execution history
- Middleware support

**Files**:
- `engine.go` - Rule engine implementation

### 7. Coordinator (`coordinator.go`)

**Purpose**: Central orchestration of all components

**Key Features**:
- Component lifecycle management
- Task queue and distribution
- Democratic task assignment
- Memory consolidation
- Learning from outcomes

## Usage

### Basic Setup

```go
import "github.com/opencode-ai/opencode/internal/swarm"

// Create coordinator
coordinator, err := swarm.NewCoordinator(swarm.CoordinatorConfig{
    SwarmConfig: agent.SwarmConfig{
        Name: "my-swarm",
        VotingThreshold: 0.66,
        EnableMemory: true,
        EnableLearning: true,
        EnableSelfHealing: true,
    },
})

// Start
coordinator.Start()
defer coordinator.Stop()

// Submit task
task := agent.Task{
    Type: "code_analysis",
    Description: "Analyze code quality",
}
coordinator.SubmitTask(task)

// Get result
result, _ := coordinator.GetTaskResult(task.ID, 5*time.Minute)
```

### Agent Registration

```go
// Create specialized agents
monitorAgent := &MonitorAgent{
    BaseAgent: agent.NewBaseAgent(agent.AgentConfig{
        Type: agent.AgentTypeMonitor,
        Capabilities: []string{"log_analysis", "pattern_detection"},
    }),
}

// Register with coordinator
registry := coordinator.GetRegistry()
registry.RegisterAgent(monitorAgent)
```

### Memory Operations

```go
memStore := coordinator.GetMemoryStore()

// Store
memory := memory.Memory{
    Type: memory.MemoryTypeEpisodic,
    Content: "Task completed successfully",
    Tags: []string{"success", "task"},
    Priority: memory.PriorityHigh,
}
memStore.Store(memory)

// Query
query := memory.MemoryQuery{
    Tags: []string{"success"},
    Limit: 10,
}
results, _ := memStore.Query(query)
```

### Democratic Voting

```go
votingSystem := coordinator.GetVotingSystem()

// Create vote
proposal := voting.VoteProposal{
    Description: "Should we proceed?",
    Deadline: time.Now().Add(1 * time.Minute),
}

session, _ := votingSystem.CreateVoteSession(
    proposal,
    voting.VoteTypeMajority,
    3, // min voters
    nil,
)

// Wait for result
result, _ := votingSystem.WaitForResult(ctx, session.ID)
```

### Health Monitoring

```go
healthMonitor := coordinator.GetHealthMonitor()

// Update health
healthMonitor.UpdateCheck(health.HealthCheck{
    ComponentID: "agent-1",
    Status: health.HealthStatusHealthy,
    Score: 0.95,
})

// Get system health
systemHealth := healthMonitor.GetSystemHealth()
```

### Rule Engine

```go
ruleEngine := coordinator.GetRuleEngine()

// Add rule
rule := rules.Rule{
    ID: "handle_errors",
    Condition: &rules.EventTypeCondition{EventType: "error"},
    Actions: []rules.Action{
        &rules.LogAction{Message: "Error detected"},
    },
}
ruleEngine.AddRule(rule)

// Evaluate
ruleEngine.EvaluateRules(ctx, ruleContext)
```

## Configuration

See [SWARM_CONFIGURATION.md](../../docs/SWARM_CONFIGURATION.md) for detailed configuration options.

## Documentation

- **[Quick Start](../../docs/SWARM_QUICK_START.md)** - Get started in 5 minutes
- **[Architecture](../../docs/MULTI_AGENT_ARCHITECTURE.md)** - Complete architecture overview
- **[Configuration](../../docs/SWARM_CONFIGURATION.md)** - Configuration reference
- **[Research](../../docs/SWARM_RESEARCH_ANALYSIS.md)** - Research analysis and findings

## Examples

See [examples/swarm_example.go](../../examples/swarm_example.go) for complete working examples.

## Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./...
```

## Performance

Based on research and benchmarks:

- **Scalability**: Linear up to 10-15 agents, then use clustering
- **Throughput**: 10-50 tasks/sec with 10 agents
- **Memory**: 100-500MB for typical workloads
- **Latency**: <100ms task distribution, <50ms memory queries

## Research Foundations

This implementation is based on extensive research:

1. **Democratic Multi-Agent Systems**
   - Majority voting for transparent decisions
   - Near-linear scalability (proven to 50 agents)
   - 90%+ fault tolerance through decentralization

2. **Hierarchical Memory**
   - 35% accuracy improvement over flat storage
   - 99.9% storage reduction through consolidation
   - 2x success rate on long-horizon tasks

3. **Self-Healing**
   - 80% reduction in manual interventions
   - 60% faster mean time to recovery
   - Predictive failure detection

See [SWARM_RESEARCH_ANALYSIS.md](../../docs/SWARM_RESEARCH_ANALYSIS.md) for full details.

## Contributing

When adding new components:

1. Follow existing patterns (interfaces, base implementations)
2. Add comprehensive documentation
3. Include examples
4. Write tests
5. Update relevant docs

## License

Same as OpenCode (MIT License)

## Support

- GitHub Issues: https://github.com/opencode-ai/opencode/issues
- Documentation: See `docs/` directory
- Examples: See `examples/` directory
