# OpenCode Architecture Documentation

## Table of Contents
1. [Overview](#overview)
2. [Application Framework](#application-framework)
3. [Modularity Analysis](#modularity-analysis)
4. [AI Agent Inference System](#ai-agent-inference-system)
5. [Data Flow and Components](#data-flow-and-components)
6. [Extension Points](#extension-points)

## Overview

OpenCode is a **terminal-based AI assistant** built in Go that provides intelligent coding assistance through a sophisticated modular architecture. The application follows clean architecture principles with clear separation of concerns and dependency injection patterns.

## Application Framework

### Core Technology Stack

```
┌─────────────────────────────────────────────────────────────┐
│                    OpenCode Application                      │
├─────────────────────────────────────────────────────────────┤
│  Language: Go 1.24+                                          │
│  CLI Framework: Cobra (spf13/cobra)                          │
│  TUI Framework: Bubble Tea (charmbracelet/bubbletea)        │
│  Database: SQLite with go-sqlite3                            │
│  Configuration: Viper (spf13/viper)                          │
│  Protocol: Model Context Protocol (MCP)                      │
└─────────────────────────────────────────────────────────────┘
```

### High-Level Architecture

```
┌──────────────┐
│   main.go    │  Entry point with panic recovery
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   cmd/       │  CLI command handling (Cobra)
└──────┬───────┘
       │
       ▼
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │   TUI    │  │   App    │  │  Config  │  │   LSP    │   │
│  │ (Views)  │  │ (Core)   │  │ (Settings)│  │(Language)│   │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘   │
└───────┼────────────┼─────────────┼──────────────┼──────────┘
        │            │              │              │
        ▼            ▼              ▼              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Domain Layer                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │ Agents   │  │ Sessions │  │ Messages │  │   Tools  │   │
│  │ (LLM)    │  │ (State)  │  │ (Chat)   │  │ (Actions)│   │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘   │
└───────┼────────────┼─────────────┼──────────────┼──────────┘
        │            │              │              │
        ▼            ▼              ▼              ▼
┌─────────────────────────────────────────────────────────────┐
│                 Infrastructure Layer                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │ Providers│  │ Database │  │  Logging │  │  PubSub  │   │
│  │ (LLM API)│  │ (SQLite) │  │ (Events) │  │ (Events) │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## Modularity Analysis

### ✅ Yes, OpenCode is Highly Modular!

The application is organized into well-defined, loosely-coupled modules:

### Module Breakdown

#### 1. **cmd/** - Command Line Interface
- **Purpose**: Entry point and CLI command handling
- **Dependencies**: Internal modules only
- **Key Files**: `root.go` (main command definition)

#### 2. **internal/app/** - Application Core
- **Purpose**: Main application orchestration and service coordination
- **Responsibilities**:
  - Service initialization and lifecycle management
  - LSP client management
  - Resource cleanup and shutdown handling
- **Key Components**:
  - `App` struct: Central service container
  - Session management
  - Message handling
  - File history tracking
  - Permission system

#### 3. **internal/config/** - Configuration Management
- **Purpose**: Centralized configuration handling
- **Features**:
  - Multi-source configuration (env vars, files, defaults)
  - Provider API key management
  - Agent model configuration
  - MCP server configuration
  - LSP configuration
- **Format**: JSON-based configuration files

#### 4. **internal/llm/** - Language Model Integration
**Sub-modules**:

##### a. **llm/agent/** - AI Agent Core
- `agent.go`: Main agent orchestration
- `agent-tool.go`: Sub-agent/task delegation
- `tools.go`: Tool registration for agents
- `mcp-tools.go`: MCP protocol tool integration

##### b. **llm/models/** - Model Definitions
- Model metadata (context windows, pricing, capabilities)
- Provider-to-model mappings
- Support for: OpenAI, Anthropic, Gemini, Groq, Bedrock

##### c. **llm/provider/** - LLM Provider Implementations
- Abstract provider interface
- Concrete implementations for each LLM provider
- Request/response handling
- Streaming support

##### d. **llm/tools/** - Tool System
- `tools.go`: Base tool interface
- File operations: `ls.go`, `file.go`, `edit.go`
- Code operations: `grep.go`, `diagnostics.go`
- External: `fetch.go`, `bash.go`

##### e. **llm/prompt/** - Prompt Engineering
- System prompts for different agents
- Provider-specific prompt adaptations

#### 5. **internal/tui/** - Terminal User Interface
**Sub-modules**:
- **tui/components/**: Reusable UI components
- **tui/layout/**: Screen layout management
- **tui/page/**: Individual page implementations
- **tui/styles/**: Consistent styling

#### 6. **internal/db/** - Data Persistence
- **Purpose**: Database operations and migrations
- **Technology**: SQLite with WAL mode
- **Schema**:
  - `sessions`: Conversation sessions
  - `messages`: Chat message history
  - `files`: File version tracking
- **Features**:
  - SQLC for type-safe queries
  - Goose for migrations
  - Foreign key constraints

#### 7. **internal/lsp/** - Language Server Protocol
- **Purpose**: IDE-like code intelligence
- **Features**:
  - Multi-language support
  - Diagnostics (errors, warnings)
  - File watching
  - Protocol implementation
- **Integration**: Exposes diagnostics to AI agents

#### 8. **internal/session/** - Session Management
- Session creation and retrieval
- Session state tracking
- Token usage and cost tracking

#### 9. **internal/message/** - Message Handling
- Message CRUD operations
- Content part management
- Tool call tracking
- PubSub event publishing

#### 10. **internal/permission/** - Security & Permissions
- Tool execution permission system
- User approval workflow
- Per-session permission persistence

#### 11. **internal/pubsub/** - Event System
- Event-driven communication between modules
- Decoupled component interactions
- Real-time UI updates

#### 12. **internal/logging/** - Observability
- Structured logging
- Persistent error logging
- Debug mode support
- Panic recovery

### Dependency Graph

```
┌────────┐
│  main  │
└───┬────┘
    │
    ▼
┌────────┐     ┌──────────┐
│  cmd   │────▶│   app    │
└────────┘     └────┬─────┘
                    │
        ┌───────────┼───────────┬────────────┬──────────┐
        ▼           ▼           ▼            ▼          ▼
    ┌───────┐  ┌────────┐  ┌─────────┐  ┌──────┐  ┌──────┐
    │ agent │  │session │  │ message │  │ lsp  │  │ tui  │
    └───┬───┘  └───┬────┘  └────┬────┘  └──┬───┘  └──┬───┘
        │          │            │           │         │
        ▼          ▼            ▼           ▼         ▼
    ┌──────────────────────────────────────────────────┐
    │           Shared Infrastructure                   │
    │  [db] [config] [logging] [pubsub] [permission]   │
    └──────────────────────────────────────────────────┘
```

## AI Agent Inference System

### What Drives AI Agent Inference?

The AI agent inference is driven by a **sophisticated multi-layer orchestration system**:

### 1. **Agent Architecture** (`internal/llm/agent/`)

```go
type Service interface {
    Run(ctx context.Context, sessionID string, content string) (<-chan AgentEvent, error)
    Cancel(sessionID string)
    IsSessionBusy(sessionID string) bool
    IsBusy() bool
}
```

**Key Components**:

#### Agent Types
- **Coder Agent**: Main coding assistant with full tool access
- **Task Agent**: Subtask delegation and execution
- **Title Agent**: Session title generation

#### Agent Flow

```
User Input
    │
    ▼
┌─────────────────────────────────────────────┐
│        Agent.Run(sessionID, content)        │
└─────────────────┬───────────────────────────┘
                  │
    ┌─────────────┴─────────────┐
    │  Generate Title (async)    │
    │  (First message only)      │
    └────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────┐
│      Create User Message in DB              │
└─────────────────┬───────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────┐
│   Load Message History from Database        │
└─────────────────┬───────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────┐
│        Provider.StreamResponse()            │
│  (Send to LLM: OpenAI, Claude, Gemini, etc) │
└─────────────────┬───────────────────────────┘
                  │
    ┌─────────────┴──────────────┐
    │  Streaming Event Processing │
    └─────────────┬────────────────┘
                  │
        ┌─────────┴──────────┐
        │   Event Types:     │
        │  - ContentDelta    │
        │  - ThinkingDelta   │
        │  - ToolUseStart    │
        │  - ToolUseStop     │
        │  - Complete        │
        │  - Error           │
        └─────────┬──────────┘
                  │
                  ▼
      ┌───────────────────────┐
      │  Tools Requested?     │
      └───────┬───────────────┘
              │
         Yes  │  No
    ┌─────────┴────┐
    │              ▼
    │      ┌──────────────┐
    │      │   Return     │
    │      │  Response    │
    │      └──────────────┘
    ▼
┌──────────────────────────┐
│  Permission Check        │
│  (User Approval)         │
└─────────┬────────────────┘
          │
          ▼
┌──────────────────────────┐
│  Execute Tools           │
│  - bash commands         │
│  - file operations       │
│  - code searches         │
│  - diagnostics           │
│  - etc.                  │
└─────────┬────────────────┘
          │
          ▼
┌──────────────────────────┐
│  Create Tool Result Msg  │
│  (Store in DB)           │
└─────────┬────────────────┘
          │
          ▼
┌──────────────────────────┐
│  Loop: Send Results      │
│  Back to LLM             │
│  (Continue conversation) │
└──────────────────────────┘
```

### 2. **Provider System** (`internal/llm/provider/`)

The provider layer abstracts different LLM APIs:

```go
type Provider interface {
    StreamResponse(ctx, messages, tools) <-chan ProviderEvent
    SendMessages(ctx, messages, tools) (Response, error)
    Model() models.Model
}
```

**Supported Providers**:
- **OpenAI**: GPT-4.1, GPT-4o, O1, O3, O4
- **Anthropic**: Claude 3.5/3.7 (Sonnet, Haiku, Opus)
- **Google**: Gemini 2.0, 2.5, 2.5 Flash
- **AWS Bedrock**: Claude via AWS
- **Groq**: Llama 4, QWEN, Deepseek

### 3. **Prompt Engineering** (`internal/llm/prompt/`)

System prompts are dynamically generated based on:
- Agent type (coder, task, title)
- Model provider (affects prompt format)
- Available tools and capabilities
- Context from configuration files

**Context Loading**: OpenCode automatically loads project-specific instructions from:
- `.github/copilot-instructions.md`
- `.cursorrules`
- `opencode.md`, `OPENCODE.md`
- `CLAUDE.md`

### 4. **Tool System** (`internal/llm/tools/`)

Tools extend agent capabilities:

```go
type BaseTool interface {
    Info() ToolInfo
    Run(ctx context.Context, params ToolCall) (ToolResponse, error)
}
```

**Available Tools**:

| Category | Tool | Description |
|----------|------|-------------|
| **File Ops** | `ls` | List directory contents |
| | `view` | Read file contents |
| | `write` | Write to files |
| | `edit` | Edit files with line ranges |
| | `patch` | Apply diff patches |
| **Search** | `grep` | Search in files |
| | `glob` | Find files by pattern |
| **Execution** | `bash` | Execute shell commands |
| **Network** | `fetch` | HTTP requests |
| | `sourcegraph` | Search public code |
| **Code Intel** | `diagnostics` | LSP diagnostics |
| **Delegation** | `agent` | Run sub-agents |

### 5. **MCP Integration** (Model Context Protocol)

OpenCode implements MCP for extensibility:

```json
{
  "mcpServers": {
    "custom-tool": {
      "type": "stdio",
      "command": "/path/to/mcp-server",
      "args": [],
      "env": []
    }
  }
}
```

**MCP Flow**:
1. MCP servers defined in config
2. Tools auto-discovered on startup
3. Tools exposed to AI agents
4. Permission system applies
5. Results returned to agent

### 6. **Inference Configuration**

Agents are configured per model:

```json
{
  "agents": {
    "coder": {
      "model": "claude-3.7-sonnet",
      "maxTokens": 5000,
      "reasoningEffort": "medium"
    }
  }
}
```

**Key Parameters**:
- **model**: Which LLM to use
- **maxTokens**: Maximum response length
- **reasoningEffort**: For reasoning models (low/medium/high)

### 7. **State Management**

The agent maintains state through:
- **Sessions**: Persistent conversation containers
- **Messages**: Full message history with roles (user/assistant/tool)
- **Tool Calls**: Tracked per message with results
- **Files**: Version tracking of modified files
- **Usage Tracking**: Token counts and cost estimation

## Data Flow and Components

### Request Flow

```
User Types in TUI
       │
       ▼
┌──────────────────┐
│  tui/page/chat   │  Bubble Tea message handling
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│  agent.Run()     │  Async channel-based execution
└────────┬─────────┘
         │
    ┌────┴──────────────────┐
    │                       │
    ▼                       ▼
┌─────────────┐    ┌──────────────────┐
│  Provider   │    │  Message Service │
│  (LLM API)  │◄───┤  (Load history)  │
└──────┬──────┘    └──────────────────┘
       │
       ▼
┌──────────────────┐
│  Stream Events   │  Channel-based streaming
└────────┬─────────┘
         │
    ┌────┴─────┬─────────┐
    │          │         │
    ▼          ▼         ▼
┌────────┐ ┌──────┐ ┌────────┐
│Content │ │Tools │ │Complete│
└───┬────┘ └──┬───┘ └───┬────┘
    │         │         │
    ▼         ▼         ▼
┌──────────────────────────┐
│  Message Service         │  Update DB
│  (Real-time save)        │
└────────┬─────────────────┘
         │
         ▼
┌──────────────────────────┐
│  PubSub Events           │
└────────┬─────────────────┘
         │
         ▼
┌──────────────────────────┐
│  TUI Updates             │  UI re-renders
└──────────────────────────┘
```

### Database Schema

```sql
-- Sessions: Conversation containers
sessions (
    id TEXT PRIMARY KEY,
    parent_session_id TEXT,
    title TEXT,
    message_count INTEGER,
    prompt_tokens INTEGER,
    completion_tokens INTEGER,
    cost REAL,
    updated_at INTEGER,
    created_at INTEGER
)

-- Messages: Chat history
messages (
    id TEXT PRIMARY KEY,
    session_id TEXT,
    role TEXT,  -- user, assistant, tool
    parts TEXT, -- JSON array of content parts
    model TEXT,
    created_at INTEGER,
    updated_at INTEGER,
    finished_at INTEGER
)

-- Files: Change tracking
files (
    id TEXT PRIMARY KEY,
    session_id TEXT,
    path TEXT,
    content TEXT,
    version TEXT,
    created_at INTEGER,
    updated_at INTEGER
)
```

## Extension Points

### Adding New Capabilities

#### 1. **New LLM Provider**
- Implement `provider.Provider` interface
- Add to `internal/llm/provider/`
- Register in models configuration

#### 2. **New Tool**
- Implement `tools.BaseTool` interface
- Add to `internal/llm/tools/`
- Register in `agent.CoderAgentTools()`

#### 3. **MCP Server**
- Create external MCP server
- Add to config `mcpServers`
- Auto-discovered on startup

#### 4. **LSP Language**
- Add language server command to config
- LSP client auto-initialized
- Diagnostics available to agents

#### 5. **New Agent Type**
- Define agent in `config.AgentName`
- Create prompt in `internal/llm/prompt/`
- Initialize in application

### Extensibility Features

✅ **Configuration-driven**: Most features configurable via JSON
✅ **Protocol-based**: MCP for external tool integration
✅ **Event-driven**: PubSub for loose coupling
✅ **Interface-based**: Easy to swap implementations
✅ **Tool-based**: AI capabilities extensible via tools

## Summary

OpenCode is a **modular, extensible, and well-architected** AI coding assistant with:

- ✅ **Clean Architecture**: Clear separation of concerns
- ✅ **Modular Design**: Independent, reusable components
- ✅ **Multiple LLM Support**: Provider abstraction layer
- ✅ **Extensible Tools**: Plugin-like tool system
- ✅ **MCP Integration**: External tool protocol
- ✅ **LSP Integration**: Language server support
- ✅ **Persistent State**: SQLite-backed storage
- ✅ **Event-Driven**: Real-time UI updates via PubSub

The AI agent inference is driven by a combination of:
1. **Configuration** (models, tokens, reasoning)
2. **Provider Implementation** (LLM API clients)
3. **Prompt Engineering** (system prompts, context)
4. **Tool System** (capabilities and actions)
5. **Message History** (conversation context)
6. **Streaming Events** (real-time processing)
