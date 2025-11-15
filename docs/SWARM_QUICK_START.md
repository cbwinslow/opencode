# Multi-Agent Swarm Quick Start Guide

## Overview

OpenCode's multi-agent swarm system enables autonomous, collaborative problem-solving through specialized AI agents. This guide will get you up and running quickly.

## What You Get

- **Democratic Decision-Making**: Agents vote on important decisions
- **Hierarchical Memory**: Stores and retrieves knowledge efficiently
- **Self-Healing**: Automatically detects and recovers from errors
- **Multi-Provider Support**: Use OpenRouter, Ollama, LM Studio, Hugging Face, or Jan
- **Continuous Learning**: Improves over time from successes and failures

## 5-Minute Setup

### 1. Basic Configuration

Add to your `.opencode.json`:

```json
{
  "swarm": {
    "enabled": true,
    "agents": [
      {
        "type": "monitor",
        "provider": "openrouter",
        "model": "anthropic/claude-3-haiku"
      },
      {
        "type": "analyzer",
        "provider": "ollama",
        "model": "llama3.2"
      },
      {
        "type": "executor",
        "provider": "ollama",
        "model": "codellama"
      }
    ]
  }
}
```

### 2. Install Local Models (Optional)

For privacy and cost savings, use Ollama:

```bash
# Install Ollama
curl -fsSL https://ollama.com/install.sh | sh

# Pull models
ollama pull llama3.2
ollama pull codellama
```

### 3. Set API Keys (If Using Cloud Providers)

```bash
export OPENROUTER_API_KEY="your-key-here"
```

### 4. Start OpenCode

```bash
opencode swarm start
```

That's it! Your multi-agent swarm is now running.

## Quick Examples

### Example 1: Run the Demo

```bash
cd examples
go run swarm_example.go
```

This demonstrates:
- Coordinator setup
- Memory storage and retrieval
- Democratic voting
- Health monitoring
- Rule engine

### Example 2: Submit a Task Programmatically

```go
import "github.com/opencode-ai/opencode/internal/swarm"

// Create coordinator
coordinator, _ := swarm.NewCoordinator(config)
coordinator.Start()
defer coordinator.Stop()

// Submit task
task := agent.Task{
    Type: "code_review",
    Description: "Review authentication module",
    Input: map[string]interface{}{
        "files": []string{"auth.go", "middleware.go"},
    },
}
coordinator.SubmitTask(task)

// Get result
result, _ := coordinator.GetTaskResult(task.ID, 5*time.Minute)
```

### Example 3: Store and Query Memory

```go
// Store a success pattern
memory := memory.Memory{
    Type: memory.MemoryTypeProcedural,
    Content: "Successfully fixed authentication bug by...",
    Tags: []string{"success", "auth", "bug-fix"},
    Priority: memory.PriorityHigh,
}
memoryStore.Store(memory)

// Query similar memories
query := memory.MemoryQuery{
    Tags: []string{"auth", "bug-fix"},
    Limit: 5,
}
similar, _ := memoryStore.Query(query)
```

### Example 4: Democratic Voting

```go
// Create proposal
proposal := voting.VoteProposal{
    Description: "Should we proceed with refactoring?",
    Deadline: time.Now().Add(1 * time.Minute),
}

// Start vote
session, _ := votingSystem.CreateVoteSession(
    proposal,
    voting.VoteTypeMajority,
    3, // min voters
    nil,
)

// Wait for result
result, _ := votingSystem.WaitForResult(ctx, session.ID)
fmt.Printf("Decision: %v (%.0f%% approval)\n", 
    result.Decision, result.YesPercentage*100)
```

## Common Use Cases

### 1. Code Analysis

Agents analyze code quality, security, and performance:

```json
{
  "agents": [
    {
      "type": "analyzer",
      "capabilities": ["code_quality", "security_scan", "performance_analysis"]
    }
  ]
}
```

### 2. Automated Testing

Agents write and run tests:

```json
{
  "agents": [
    {
      "type": "testing",
      "capabilities": ["test_generation", "test_execution", "coverage_analysis"]
    }
  ]
}
```

### 3. Documentation

Agents maintain documentation:

```json
{
  "agents": [
    {
      "type": "documentation",
      "capabilities": ["doc_generation", "api_docs", "code_comments"]
    }
  ]
}
```

### 4. Error Recovery

Agents detect and fix errors:

```json
{
  "agents": [
    {
      "type": "error_handler",
      "capabilities": ["error_detection", "root_cause_analysis", "auto_fix"]
    }
  ]
}
```

## Provider Comparison

| Provider | Cost | Privacy | Speed | Best For |
|----------|------|---------|-------|----------|
| **Ollama** | Free | ✅ High | Fast | Local development, privacy |
| **OpenRouter** | Low | ⚠️ Medium | Fast | Access to many models |
| **LM Studio** | Free | ✅ High | Medium | GUI management |
| **Hugging Face** | Varies | ⚠️ Medium | Varies | Model experimentation |
| **Jan** | Free | ✅ High | Medium | Desktop app users |

**Recommendation**: Start with Ollama for free, private local development.

## Monitoring

### View System Status

```bash
opencode swarm status
```

Output:
```
System Status:
  Running: true
  Agents: 3 (3 healthy, 0 degraded)
  Memory: 1,247 entries
  Active Votes: 0
  Queued Tasks: 2
  Overall Health: 0.98
```

### View Agent Health

```bash
opencode swarm agents health
```

Output:
```
Agent Health:
  monitor-1: Healthy (0.95)
  analyzer-1: Healthy (0.92)
  executor-1: Healthy (0.98)
```

### View Memory Stats

```bash
opencode swarm memory stats
```

Output:
```
Memory Statistics:
  Total: 1,247
  Working: 12
  Episodic: 423
  Semantic: 651
  Procedural: 161
```

## Performance Tips

1. **Start Small**: Begin with 2-3 agents, scale up as needed
2. **Use Local Models**: Ollama is fast and free
3. **Memory Pruning**: Set appropriate retention periods
4. **Health Checks**: Monitor regularly for issues
5. **Voting Timeout**: Adjust based on task complexity

## Troubleshooting

### Agents Not Starting

**Problem**: Agents fail to initialize
**Solution**: Check provider connectivity

```bash
# Test Ollama
curl http://localhost:11434/api/tags

# Test OpenRouter
curl -H "Authorization: Bearer $OPENROUTER_API_KEY" \
  https://openrouter.ai/api/v1/models
```

### High Memory Usage

**Problem**: Memory store growing too large
**Solution**: Adjust pruning settings

```json
{
  "memory": {
    "maxMemories": 5000,
    "pruneOlderThan": "7d"
  }
}
```

### Slow Performance

**Problem**: Tasks taking too long
**Solution**: Add more executor agents

```json
{
  "agents": [
    {"type": "executor", "id": "executor-1"},
    {"type": "executor", "id": "executor-2"},
    {"type": "executor", "id": "executor-3"}
  ]
}
```

### Vote Timeouts

**Problem**: Votes not completing
**Solution**: Increase timeout or reduce min voters

```json
{
  "voting": {
    "timeout": "2m",
    "minVoters": 2
  }
}
```

## Next Steps

1. **Read Full Documentation**:
   - [Architecture Guide](MULTI_AGENT_ARCHITECTURE.md)
   - [Configuration Reference](SWARM_CONFIGURATION.md)
   - [Research Analysis](SWARM_RESEARCH_ANALYSIS.md)

2. **Explore Examples**:
   - `examples/swarm_example.go` - Complete working examples

3. **Customize Agents**:
   - Add specialized agents for your use case
   - Configure custom capabilities
   - Define domain-specific rules

4. **Monitor and Optimize**:
   - Track agent performance
   - Adjust configurations based on metrics
   - Scale based on workload

## Getting Help

- **GitHub Issues**: Report bugs and request features
- **Documentation**: Comprehensive guides in `docs/`
- **Examples**: Working code in `examples/`
- **Discord**: Join the OpenCode community

## Key Concepts

### Agents

Specialized AI workers that perform specific tasks:
- **Monitor**: Watches logs and system state
- **Analyzer**: Analyzes data and patterns
- **Executor**: Executes tasks and generates code
- **Memory**: Manages knowledge storage
- **Learning**: Improves from experience

### Memory Types

- **Working**: Current context (temporary)
- **Episodic**: Events and experiences
- **Semantic**: Facts and knowledge
- **Procedural**: How-to skills

### Voting Types

- **Majority**: >50% agreement
- **Super**: >66% agreement
- **Unanimous**: 100% agreement
- **Weighted**: Expertise-based
- **Consensus**: Iterative building (>75%)

### Health States

- **Healthy**: Operating normally (0.8-1.0)
- **Degraded**: Minor issues (0.5-0.8)
- **Unhealthy**: Significant problems (0.3-0.5)
- **Critical**: Failing (<0.3)

## Summary

The multi-agent swarm system provides:
- ✅ Autonomous problem-solving
- ✅ Democratic decision-making
- ✅ Persistent knowledge storage
- ✅ Self-healing capabilities
- ✅ Continuous learning
- ✅ Multiple LLM provider support

Start with the basic configuration, experiment with different agents and providers, and scale up as your needs grow.

**Ready to build? Run the example now:**

```bash
cd examples
go run swarm_example.go
```
