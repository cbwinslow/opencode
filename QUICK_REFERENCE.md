# OpenCode Framework Quick Reference

## 🏗️ Application Framework

### Technology Stack
- **Language**: Go 1.24+
- **CLI**: Cobra (command handling)
- **TUI**: Bubble Tea (terminal UI)
- **Database**: SQLite with WAL mode
- **Config**: Viper (multi-source configuration)
- **Protocol**: MCP (Model Context Protocol)

### Architecture Pattern
```
Clean Architecture with:
├── Presentation Layer (TUI)
├── Application Layer (Services)
├── Domain Layer (Business Logic)
└── Infrastructure Layer (DB, APIs, Logging)
```

## ✅ Is OpenCode Modular?

**YES! Highly modular with 12+ independent modules:**

### Core Modules
| Module | Purpose | Dependencies |
|--------|---------|--------------|
| **cmd/** | CLI commands | Cobra |
| **app/** | Service orchestration | All services |
| **config/** | Configuration mgmt | Viper |
| **llm/** | AI integration | Providers, Tools |
| **tui/** | User interface | Bubble Tea |
| **db/** | Data persistence | SQLite |
| **lsp/** | Language servers | LSP protocol |
| **session/** | State management | Database |
| **message/** | Chat history | Database |
| **permission/** | Security | PubSub |
| **pubsub/** | Event system | None |
| **logging/** | Observability | None |

### Module Independence
- ✅ Clear interfaces between modules
- ✅ Dependency injection pattern
- ✅ Event-driven communication (PubSub)
- ✅ Swappable implementations
- ✅ Independent testing possible

## 🤖 AI Agent Inference System

### What Drives It?

**Multi-layer orchestration system:**

1. **Configuration** (`config/`)
   - Model selection (Claude, GPT-4, Gemini, etc.)
   - Token limits and reasoning settings
   - Provider API keys

2. **Provider Abstraction** (`llm/provider/`)
   - OpenAI, Anthropic, Google, Groq, AWS Bedrock
   - Streaming response support
   - Unified interface across all providers

3. **Prompt Engineering** (`llm/prompt/`)
   - System prompts per agent type
   - Context from project files
   - Provider-specific adaptations

4. **Tool System** (`llm/tools/`)
   - File operations (ls, view, write, edit)
   - Search (grep, glob)
   - Code intelligence (diagnostics)
   - Execution (bash)
   - Delegation (agent)

5. **Message History** (`message/`)
   - Full conversation context
   - SQLite persistence
   - Real-time updates

6. **Event Streaming** (`llm/agent/`)
   - Channel-based async processing
   - Tool call detection and execution
   - Permission checking
   - Loop until complete

### Agent Flow
```
User Input → Agent.Run() → Load History → Provider.StreamResponse()
    ↓
Events: Content, Thinking, ToolCalls → Permission Check
    ↓
Execute Tools → Create Tool Results → Loop (send back to LLM)
    ↓
Final Response → Save to DB → Update UI
```

### Supported Models (14+)

**OpenAI**: GPT-4.1, GPT-4o, O1, O3, O4
**Anthropic**: Claude 3.5/3.7 (Sonnet, Haiku)
**Google**: Gemini 2.0, 2.5 Flash
**Groq**: Llama 4, QWEN, Deepseek
**AWS**: Claude via Bedrock

## 🔌 Extension Points

### 1. Add New Tool
```go
// Implement BaseTool interface
type MyTool struct { ... }

func (t *MyTool) Info() ToolInfo { ... }
func (t *MyTool) Run(ctx, params) (ToolResponse, error) { ... }

// Register in agent.CoderAgentTools()
```

### 2. Add New Provider
```go
// Implement Provider interface
type MyProvider struct { ... }

func (p *MyProvider) StreamResponse(...) <-chan ProviderEvent { ... }
func (p *MyProvider) Model() models.Model { ... }
```

### 3. Add MCP Server
```json
{
  "mcpServers": {
    "my-server": {
      "type": "stdio",
      "command": "/path/to/server"
    }
  }
}
```

### 4. Add LSP Language
```json
{
  "lsp": {
    "python": {
      "command": "pylsp"
    }
  }
}
```

## 💾 Database Schema

### Tables
```sql
sessions (id, title, tokens, cost, timestamps)
messages (id, session_id, role, parts, model)
files (id, session_id, path, content, version)
```

### Features
- Foreign key constraints
- Automatic triggers
- Token/cost tracking
- File version history

## 📊 Data Flow

```
User Input (TUI)
    ↓
Agent Service (Async)
    ↓
Provider (LLM API) ← Message History (DB)
    ↓
Stream Events → Update DB → PubSub Events
    ↓
TUI Updates (Real-time)
```

## 🔮 RAG Integration Possibilities

### ✅ Can We Set Up RAG with OpenCode?

**YES! Multiple approaches available:**

### Approach 1: Built-in SQLite-vec
- ✅ Single database file
- ✅ No external dependencies
- ✅ Good for small-medium projects
- ⚠️ Limited scalability

### Approach 2: Embedded Vector DB
- ✅ Better performance
- ✅ Qdrant/Weaviate embedded
- ⚠️ Additional dependency

### Approach 3: External Vector DB
- ✅ Best performance
- ✅ Cloud or self-hosted
- ⚠️ Requires external service

### Approach 4: MCP RAG Server
- ✅ Leverages existing MCP
- ✅ Pluggable architecture
- ✅ User's choice of RAG solution

### Recommended: Hybrid
**Built-in SQLite-vec + MCP extension support**

### What RAG Would Add

1. **Semantic Code Search**
   - "Find authentication implementation"
   - "Show similar error handling"

2. **Large Codebase Support**
   - Index unlimited code size
   - Only send relevant chunks to LLM

3. **Cost Efficiency**
   - Reduce input token costs
   - Faster inference

4. **Better Accuracy**
   - Focused context
   - Project-aware responses

5. **Documentation Access**
   - Search wikis, READMEs
   - Find usage examples

### RAG Architecture
```
User Query
    ↓
Embed Query (OpenAI/Cohere/Local)
    ↓
Vector Search (SQLite-vec/Qdrant)
    ↓
Retrieve Top-K Chunks
    ↓
Augment LLM Prompt
    ↓
Generate Response with Context
```

### Implementation Phases

**Phase 1**: Foundation (embeddings, vector store)
**Phase 2**: Indexing (file crawler, incremental updates)
**Phase 3**: Retrieval (semantic search, ranking)
**Phase 4**: Agent integration (RAG tools)
**Phase 5**: UX polish (status, visualization)

### Storage Requirements
- 10MB code → 10K chunks → 60MB vectors
- Query: <100ms
- Index: <1 minute for 100K lines

## 🚀 Quick Start Examples

### Configuration
```json
{
  "agents": {
    "coder": {
      "model": "claude-3.7-sonnet",
      "maxTokens": 5000
    }
  },
  "providers": {
    "anthropic": {
      "apiKey": "sk-..."
    }
  }
}
```

### Add Tool
```go
tool := &MyTool{}
agent.RegisterTool(tool)
```

### Search Code
```go
results, _ := vectorStore.Search(ctx, "authentication flow", 10)
```

## 📚 Key Files Reference

| File | Purpose |
|------|---------|
| `main.go` | Entry point |
| `cmd/root.go` | CLI setup |
| `internal/app/app.go` | Service container |
| `internal/llm/agent/agent.go` | Agent orchestration |
| `internal/llm/provider/provider.go` | LLM abstraction |
| `internal/config/config.go` | Configuration |
| `internal/db/db.go` | Database operations |
| `internal/tui/tui.go` | UI root |

## 🎯 Summary

### Framework
✅ **Modular**: 12+ independent modules
✅ **Clean Architecture**: Clear separation
✅ **Event-Driven**: PubSub pattern
✅ **Extensible**: Interfaces and protocols

### AI Inference
✅ **Multi-Provider**: 5+ LLM providers
✅ **Tool-Based**: Extensible capabilities
✅ **Streaming**: Real-time responses
✅ **Persistent**: SQLite storage

### RAG Potential
✅ **Feasible**: Multiple integration paths
✅ **Valuable**: Large codebase support
✅ **Modular**: Fits existing architecture
✅ **Recommended**: Hybrid approach

---

**For detailed documentation, see:**
- `ARCHITECTURE.md` - Full architecture analysis
- `RAG_INTEGRATION.md` - RAG integration guide
- `README.md` - User documentation
