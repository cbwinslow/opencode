# Swarm Configuration Guide

This guide explains how to configure and use the multi-agent swarm system in OpenCode.

## Configuration File Structure

Add the swarm configuration to your `.opencode.json`:

```json
{
  "swarm": {
    "enabled": true,
    "name": "opencode-swarm",
    "voting": {
      "threshold": 0.66,
      "type": "majority"
    },
    "maxConcurrentTasks": 10,
    "enableMemory": true,
    "enableLearning": true,
    "enableSelfHealing": true,
    "agents": [
      {
        "id": "monitor-1",
        "type": "monitor",
        "provider": "openrouter",
        "model": "anthropic/claude-3-haiku",
        "maxConcurrency": 3,
        "healthCheckInterval": "30s",
        "capabilities": ["log_analysis", "pattern_detection"]
      },
      {
        "id": "analyzer-1",
        "type": "analyzer",
        "provider": "ollama",
        "model": "llama3.2",
        "maxConcurrency": 2,
        "capabilities": ["data_analysis", "trend_detection"]
      },
      {
        "id": "executor-1",
        "type": "executor",
        "provider": "lmstudio",
        "model": "local-model",
        "maxConcurrency": 5,
        "capabilities": ["task_execution", "code_generation"]
      }
    ],
    "memory": {
      "maxMemories": 10000,
      "consolidationInterval": "1h",
      "pruneOlderThan": "720h",
      "enableEncryption": true
    },
    "monitoring": {
      "logPaths": [
        "/var/log/opencode/*.log",
        "~/.opencode/logs/*.log"
      ],
      "shellHistory": "~/.bash_history",
      "checkInterval": "30s"
    },
    "health": {
      "alertThreshold": 0.5,
      "checkInterval": "30s"
    }
  }
}
```

## Provider-Specific Configuration

### OpenRouter

Free and paid models available through a single API.

```json
{
  "agents": [
    {
      "provider": "openrouter",
      "model": "anthropic/claude-3-haiku",
      "config": {
        "apiKey": "${OPENROUTER_API_KEY}",
        "baseURL": "https://openrouter.ai/api/v1"
      }
    }
  ]
}
```

**Recommended Free Models**:
- `meta-llama/llama-3.2-3b-instruct:free`
- `mistralai/mistral-7b-instruct:free`
- `google/gemma-2-9b-it:free`

### Ollama

Local execution with privacy.

```json
{
  "agents": [
    {
      "provider": "ollama",
      "model": "llama3.2",
      "config": {
        "baseURL": "http://localhost:11434"
      }
    }
  ]
}
```

**Setup**:
```bash
# Install Ollama
curl -fsSL https://ollama.com/install.sh | sh

# Pull models
ollama pull llama3.2
ollama pull codellama
ollama pull mistral
```

### LM Studio

Local models with GUI management.

```json
{
  "agents": [
    {
      "provider": "lmstudio",
      "model": "local-model",
      "config": {
        "baseURL": "http://localhost:1234/v1"
      }
    }
  ]
}
```

**Setup**:
1. Download LM Studio from https://lmstudio.ai/
2. Load models through the UI
3. Start local server
4. Use in OpenCode

### Hugging Face

Access to thousands of models.

```json
{
  "agents": [
    {
      "provider": "huggingface",
      "model": "meta-llama/Llama-3.2-3B-Instruct",
      "config": {
        "apiKey": "${HUGGINGFACE_API_KEY}",
        "useInferenceAPI": true
      }
    }
  ]
}
```

### Jan

Privacy-focused local execution.

```json
{
  "agents": [
    {
      "provider": "jan",
      "model": "llama3.2",
      "config": {
        "baseURL": "http://localhost:1337/v1"
      }
    }
  ]
}
```

## Agent Type Configurations

### Monitor Agent

Watches logs and system state.

```json
{
  "id": "monitor-1",
  "type": "monitor",
  "provider": "openrouter",
  "model": "anthropic/claude-3-haiku",
  "capabilities": [
    "log_analysis",
    "pattern_detection",
    "anomaly_detection"
  ],
  "config": {
    "watchPaths": [
      "/var/log/*.log",
      "~/.opencode/logs/*.log"
    ],
    "alertOnKeywords": ["ERROR", "FATAL", "CRITICAL"]
  }
}
```

### Analyzer Agent

Analyzes data and identifies patterns.

```json
{
  "id": "analyzer-1",
  "type": "analyzer",
  "provider": "ollama",
  "model": "mistral",
  "capabilities": [
    "data_analysis",
    "trend_detection",
    "performance_analysis"
  ],
  "config": {
    "analysisInterval": "5m",
    "metricsToTrack": [
      "error_rate",
      "response_time",
      "success_rate"
    ]
  }
}
```

### Executor Agent

Executes tasks and generates code.

```json
{
  "id": "executor-1",
  "type": "executor",
  "provider": "lmstudio",
  "model": "codellama-13b",
  "capabilities": [
    "task_execution",
    "code_generation",
    "test_writing"
  ],
  "config": {
    "maxRetries": 3,
    "timeout": "5m"
  }
}
```

### Memory Agent

Manages the memory system.

```json
{
  "id": "memory-1",
  "type": "memory",
  "provider": "openrouter",
  "model": "anthropic/claude-3-sonnet",
  "capabilities": [
    "memory_consolidation",
    "semantic_search",
    "knowledge_extraction"
  ],
  "config": {
    "consolidationSchedule": "0 2 * * *",
    "vectorDimensions": 1536
  }
}
```

### Learning Agent

Learns from outcomes.

```json
{
  "id": "learning-1",
  "type": "learning",
  "provider": "huggingface",
  "model": "meta-llama/Llama-3.2-3B-Instruct",
  "capabilities": [
    "pattern_learning",
    "strategy_optimization",
    "failure_analysis"
  ],
  "config": {
    "learningRate": 0.01,
    "minSamples": 10
  }
}
```

## Rule Configuration Examples

### Error Handling Rule

```json
{
  "rules": [
    {
      "id": "handle_critical_errors",
      "name": "Critical Error Handler",
      "priority": 100,
      "enabled": true,
      "condition": {
        "type": "field",
        "field": "level",
        "operator": "==",
        "value": "CRITICAL"
      },
      "actions": [
        {
          "type": "notify",
          "config": {
            "channels": ["slack", "email"],
            "message": "Critical error detected: {{message}}"
          }
        },
        {
          "type": "trigger_recovery",
          "config": {
            "strategy": "restart"
          }
        }
      ]
    }
  ]
}
```

### Memory Consolidation Rule

```json
{
  "rules": [
    {
      "id": "consolidate_memories",
      "name": "Memory Consolidator",
      "priority": 50,
      "enabled": true,
      "condition": {
        "type": "schedule",
        "cron": "0 2 * * *"
      },
      "actions": [
        {
          "type": "consolidate_memory",
          "config": {
            "strategy": "semantic_clustering",
            "minSimilarity": 0.85
          }
        }
      ]
    }
  ]
}
```

### Performance Monitoring Rule

```json
{
  "rules": [
    {
      "id": "monitor_performance",
      "name": "Performance Monitor",
      "priority": 75,
      "enabled": true,
      "condition": {
        "type": "field",
        "field": "response_time",
        "operator": ">",
        "value": 5000
      },
      "actions": [
        {
          "type": "log",
          "message": "Slow response detected: {{response_time}}ms"
        },
        {
          "type": "scale",
          "config": {
            "direction": "up",
            "amount": 1
          }
        }
      ]
    }
  ]
}
```

## Voting Configuration

### Majority Voting (Default)

Simple majority decision-making.

```json
{
  "voting": {
    "type": "majority",
    "threshold": 0.5,
    "minVoters": 3,
    "timeout": "30s"
  }
}
```

### Super Majority

Requires strong agreement.

```json
{
  "voting": {
    "type": "super",
    "threshold": 0.66,
    "minVoters": 5,
    "timeout": "1m"
  }
}
```

### Weighted Voting

Different agents have different voting power.

```json
{
  "voting": {
    "type": "weighted",
    "threshold": 0.5,
    "weights": {
      "analyzer-1": 2.0,
      "monitor-1": 1.5,
      "executor-1": 1.0
    }
  }
}
```

### Consensus Building

Iterative approach to agreement.

```json
{
  "voting": {
    "type": "consensus",
    "threshold": 0.75,
    "maxRounds": 3,
    "roundTimeout": "30s"
  }
}
```

## Memory Configuration

### Basic Memory Setup

```json
{
  "memory": {
    "maxMemories": 10000,
    "consolidationInterval": "1h",
    "pruneOlderThan": "720h"
  }
}
```

### Advanced Memory Setup

```json
{
  "memory": {
    "maxMemories": 50000,
    "consolidationInterval": "1h",
    "pruneOlderThan": "720h",
    "enableEncryption": true,
    "encryptionKey": "${MEMORY_ENCRYPTION_KEY}",
    "vectorDimensions": 1536,
    "embeddingProvider": "openai",
    "embeddingModel": "text-embedding-3-small",
    "consolidationStrategy": "semantic_clustering",
    "pruning": {
      "strategy": "least_recently_used",
      "preserveTags": ["critical", "knowledge", "success"],
      "minAccessCount": 2
    }
  }
}
```

## Health Monitoring Configuration

### Basic Health Setup

```json
{
  "health": {
    "checkInterval": "30s",
    "alertThreshold": 0.5
  }
}
```

### Advanced Health Setup

```json
{
  "health": {
    "checkInterval": "30s",
    "alertThreshold": 0.5,
    "alerts": {
      "channels": ["slack", "email", "webhook"],
      "webhook": {
        "url": "https://your-webhook.com/alert",
        "method": "POST"
      }
    },
    "recovery": {
      "autoRecover": true,
      "strategies": [
        {
          "type": "restart",
          "maxAttempts": 3,
          "backoff": "exponential"
        },
        {
          "type": "reset",
          "condition": "score < 0.3"
        },
        {
          "type": "fallback",
          "fallbackAgent": "backup-executor"
        }
      ]
    }
  }
}
```

## Environment Variables

```bash
# API Keys
export OPENROUTER_API_KEY="your-key"
export HUGGINGFACE_API_KEY="your-key"

# Memory encryption
export MEMORY_ENCRYPTION_KEY="your-32-byte-key"

# Monitoring
export SWARM_LOG_LEVEL="info"
export SWARM_DEBUG="false"

# Providers
export OLLAMA_HOST="http://localhost:11434"
export LMSTUDIO_HOST="http://localhost:1234"
export JAN_HOST="http://localhost:1337"
```

## CLI Commands

```bash
# Start swarm
opencode swarm start

# Stop swarm
opencode swarm stop

# Status
opencode swarm status

# List agents
opencode swarm agents list

# View agent health
opencode swarm agents health

# View memory stats
opencode swarm memory stats

# View active votes
opencode swarm voting active

# Query memories
opencode swarm memory query --tags "success,task" --limit 10

# View rules
opencode swarm rules list

# Add rule
opencode swarm rules add --file rule.json

# View logs
opencode swarm logs --follow

# Health check
opencode swarm health
```

## Programmatic Usage

### Go API

```go
import "github.com/opencode-ai/opencode/internal/swarm"

// Create coordinator
coordinator, err := swarm.NewCoordinator(swarm.CoordinatorConfig{
    SwarmConfig: agent.SwarmConfig{
        Name: "my-swarm",
        VotingThreshold: 0.66,
        EnableMemory: true,
        EnableLearning: true,
    },
})

// Start
coordinator.Start()
defer coordinator.Stop()

// Submit task
task := agent.Task{
    Type: "code_analysis",
    Description: "Analyze this code",
    Input: map[string]interface{}{
        "code": sourceCode,
    },
}
coordinator.SubmitTask(task)

// Get result
result, err := coordinator.GetTaskResult(task.ID, 5*time.Minute)
```

## Best Practices

1. **Start Small**: Begin with 2-3 agents and scale up
2. **Provider Mix**: Use local models for privacy, cloud for power
3. **Memory Management**: Set appropriate pruning intervals
4. **Health Monitoring**: Configure alerts for critical components
5. **Rule Testing**: Test rules thoroughly before enabling
6. **Voting Strategy**: Match strategy to decision importance
7. **Resource Limits**: Set maxConcurrency to prevent overload
8. **Logging**: Enable comprehensive logging during setup
9. **Backups**: Regularly backup memory stores
10. **Monitoring**: Set up dashboards for system health

## Troubleshooting

### Agent Not Starting

Check provider connectivity:
```bash
# Ollama
curl http://localhost:11434/api/tags

# LM Studio
curl http://localhost:1234/v1/models

# OpenRouter
curl -H "Authorization: Bearer $OPENROUTER_API_KEY" \
  https://openrouter.ai/api/v1/models
```

### High Memory Usage

Adjust memory configuration:
```json
{
  "memory": {
    "maxMemories": 5000,
    "pruneOlderThan": "168h"
  }
}
```

### Slow Performance

Scale up agents:
```json
{
  "maxConcurrentTasks": 20,
  "agents": [
    // Add more executor agents
  ]
}
```

### Vote Timeout

Increase voting timeout:
```json
{
  "voting": {
    "timeout": "2m"
  }
}
```

## Support

For issues and questions:
- GitHub Issues: https://github.com/opencode-ai/opencode/issues
- Documentation: https://docs.opencode.ai/swarm
- Discord: https://discord.gg/opencode
