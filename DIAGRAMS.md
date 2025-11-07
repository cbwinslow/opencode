# OpenCode Visual Architecture Diagrams

## High-Level System Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         OpenCode Application                         │
│                                                                      │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │                    User Interface Layer                         │ │
│  │  ┌──────────────────────────────────────────────────────────┐  │ │
│  │  │         Terminal UI (Bubble Tea)                          │  │ │
│  │  │  ┌────────┐  ┌────────┐  ┌────────┐  ┌────────────────┐ │  │ │
│  │  │  │ Chat   │  │ Logs   │  │ Help   │  │   Dialogs      │ │  │ │
│  │  │  │ Page   │  │ Page   │  │ Dialog │  │   & Overlays   │ │  │ │
│  │  │  └────────┘  └────────┘  └────────┘  └────────────────┘ │  │ │
│  │  └──────────────────────────────────────────────────────────┘  │ │
│  └────────────────────────────────────────────────────────────────┘ │
│                                                                      │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │                   Application Services Layer                    │ │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌──────────┐ │ │
│  │  │  Session   │  │  Message   │  │  History   │  │Permission│ │ │
│  │  │  Service   │  │  Service   │  │  Service   │  │ Service  │ │ │
│  │  └────────────┘  └────────────┘  └────────────┘  └──────────┘ │ │
│  └────────────────────────────────────────────────────────────────┘ │
│                                                                      │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │                    AI Agent Layer                               │ │
│  │  ┌──────────────────────────────────────────────────────────┐  │ │
│  │  │              Agent Orchestrator                           │  │ │
│  │  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌─────────┐  │  │ │
│  │  │  │  Coder   │  │   Task   │  │  Title   │  │  Tool   │  │  │ │
│  │  │  │  Agent   │  │  Agent   │  │  Agent   │  │ System  │  │  │ │
│  │  │  └──────────┘  └──────────┘  └──────────┘  └─────────┘  │  │ │
│  │  └──────────────────────────────────────────────────────────┘  │ │
│  └────────────────────────────────────────────────────────────────┘ │
│                                                                      │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │                 Infrastructure Layer                            │ │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────┐  │ │
│  │  │ Database │  │   LSP    │  │  Logging │  │    PubSub    │  │ │
│  │  │ (SQLite) │  │ Clients  │  │  System  │  │  Events      │  │ │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────────┘  │ │
│  └────────────────────────────────────────────────────────────────┘ │
│                                                                      │
└────────────────────────────────────────┬───────────────────────────┘
                                         │
                    ┌────────────────────┼────────────────────┐
                    │                    │                    │
                    ▼                    ▼                    ▼
         ┌──────────────────┐ ┌──────────────────┐ ┌──────────────────┐
         │   LLM Providers  │ │   MCP Servers    │ │  File System     │
         │  ┌────────────┐  │ │  ┌────────────┐  │ │  & Shell         │
         │  │  OpenAI    │  │ │  │   Custom   │  │ │                  │
         │  │  Anthropic │  │ │  │    Tools   │  │ │                  │
         │  │  Gemini    │  │ │  │            │  │ │                  │
         │  │  Groq      │  │ │  └────────────┘  │ │                  │
         │  │  Bedrock   │  │ │                  │ │                  │
         │  └────────────┘  │ │                  │ │                  │
         └──────────────────┘ └──────────────────┘ └──────────────────┘
```

## Module Dependency Graph

```
                          ┌──────────┐
                          │   main   │
                          └─────┬────┘
                                │
                          ┌─────▼────┐
                          │   cmd    │  (CLI Framework)
                          └─────┬────┘
                                │
                          ┌─────▼────┐
                          │   app    │  (Service Container)
                          └─────┬────┘
                                │
              ┌─────────────────┼─────────────────┐
              │                 │                 │
        ┌─────▼────┐      ┌────▼─────┐     ┌────▼────┐
        │  agent   │      │ session  │     │   tui   │
        │  (LLM)   │      │ (State)  │     │  (View) │
        └─────┬────┘      └────┬─────┘     └────┬────┘
              │                │                 │
              │                │                 │
        ┌─────┴────────────────┴────────────────┘
        │
        ▼
┌───────────────────────────────────────────────────────┐
│          Shared Infrastructure Modules                │
│  ┌──────┐  ┌────────┐  ┌─────────┐  ┌───────────┐   │
│  │  db  │  │ config │  │ logging │  │  pubsub   │   │
│  └──────┘  └────────┘  └─────────┘  └───────────┘   │
│  ┌──────┐  ┌────────┐  ┌─────────┐                  │
│  │ lsp  │  │message │  │permission│                  │
│  └──────┘  └────────┘  └─────────┘                  │
└───────────────────────────────────────────────────────┘
```

## AI Agent Inference Flow

```
┌─────────────┐
│ User Input  │
└──────┬──────┘
       │
       ▼
┌────────────────────────┐
│ agent.Run()            │
│ - Create user message  │
│ - Load history from DB │
└──────┬─────────────────┘
       │
       ▼
┌────────────────────────────┐
│ provider.StreamResponse()  │
│ - Send to LLM API          │
│ - Include message history  │
│ - Include system prompt    │
│ - Include available tools  │
└──────┬─────────────────────┘
       │
       ▼
┌──────────────────────────────────────────────────┐
│          Event Stream Processing                  │
│                                                   │
│  ┌─────────────┐  ┌─────────────┐  ┌──────────┐ │
│  │ ContentDelta│  │ThinkingDelta│  │ToolUse   │ │
│  │    ▼        │  │    ▼        │  │   Start  │ │
│  │  Update UI  │  │  Update UI  │  │    ▼     │ │
│  │  Save to DB │  │  Save to DB │  │  Store   │ │
│  └─────────────┘  └─────────────┘  └──────────┘ │
└──────────────────────────┬───────────────────────┘
                           │
              ┌────────────┴────────────┐
              │                         │
              │ Complete Event?         │
              │                         │
         No   │                         │ Yes
    ┌─────────▼─────────┐               │
    │ Any Tool Calls?   │               │
    │                   │               │
    │  No       Yes     │               │
    │  │         │      │               │
    │  │    ┌────▼─────────────┐        │
    │  │    │ Permission Check │        │
    │  │    └────┬─────────────┘        │
    │  │         │                      │
    │  │         ▼                      │
    │  │    ┌─────────────────┐         │
    │  │    │ Execute Tools   │         │
    │  │    │ - bash          │         │
    │  │    │ - file ops      │         │
    │  │    │ - diagnostics   │         │
    │  │    │ - agent         │         │
    │  │    └────┬────────────┘         │
    │  │         │                      │
    │  │         ▼                      │
    │  │    ┌─────────────────┐         │
    │  │    │ Create Tool     │         │
    │  │    │ Result Message  │         │
    │  │    └────┬────────────┘         │
    │  │         │                      │
    │  │         │ Loop back to         │
    │  │         │ provider.Stream...   │
    │  │         └──────────────────┐   │
    │  │                            │   │
    │  └────────────────────────────┼───┼──┐
    │                               │   │  │
    │                               │   │  │
    └───────────────────────────────┴───┴──┘
                                           │
                                           ▼
                                    ┌─────────────┐
                                    │   Return    │
                                    │  Response   │
                                    │  to User    │
                                    └─────────────┘
```

## Tool System Architecture

```
┌──────────────────────────────────────────────────────┐
│                  Agent with Tools                    │
└────────────────────┬─────────────────────────────────┘
                     │
        ┌────────────┼────────────┬────────────┐
        │            │            │            │
        ▼            ▼            ▼            ▼
┌─────────────┐ ┌─────────┐ ┌─────────┐ ┌──────────┐
│ Built-in    │ │  LSP    │ │  MCP    │ │  Agent   │
│   Tools     │ │  Tools  │ │  Tools  │ │   Tool   │
└─────┬───────┘ └────┬────┘ └────┬────┘ └────┬─────┘
      │              │           │            │
      │              │           │            │
      ▼              ▼           ▼            ▼

File Operations    Code Intel    External      Sub-Agent
┌────────────┐    ┌─────────┐   ┌─────────┐  ┌────────┐
│ ls         │    │diagnose │   │ custom  │  │ task   │
│ view       │    │ (errors)│   │ tools   │  │ agent  │
│ write      │    └─────────┘   │ via     │  └────────┘
│ edit       │                  │ MCP     │
│ patch      │                  │ servers │
│ grep       │                  └─────────┘
│ glob       │
└────────────┘

Execution         Network
┌────────────┐    ┌─────────┐
│ bash       │    │ fetch   │
└────────────┘    │sourcegrph│
                  └─────────┘
```

## RAG Integration Architecture (Proposed)

```
┌─────────────────────────────────────────────────────────────────┐
│                    OpenCode with RAG                            │
└────────────────────────────┬────────────────────────────────────┘
                             │
                ┌────────────┴────────────┐
                │                         │
                ▼                         ▼
    ┌──────────────────────┐    ┌──────────────────────┐
    │    Built-in RAG      │    │   MCP RAG Servers    │
    │   (SQLite-vec)       │    │     (Optional)       │
    └───────┬──────────────┘    └──────────┬───────────┘
            │                              │
    ┌───────┴──────────┐          ┌────────┴─────────┐
    │                  │          │                  │
    ▼                  ▼          ▼                  ▼
┌────────┐      ┌──────────┐  ┌──────┐      ┌──────────┐
│Indexer │      │ Vector   │  │Qdrant│      │ Weaviate │
│        │──────▶│  Store   │  │Server│      │  Server  │
│        │      │(SQLite)  │  └──────┘      └──────────┘
└───┬────┘      └────┬─────┘
    │                │
    │                │
    ▼                ▼
┌─────────────────────────────────────┐
│      Embedding Service              │
│  ┌──────────┐  ┌──────────────┐    │
│  │ OpenAI   │  │   Cohere     │    │
│  │ API      │  │   API        │    │
│  └──────────┘  └──────────────┘    │
│  ┌──────────┐                      │
│  │  Local   │                      │
│  │  ONNX    │                      │
│  └──────────┘                      │
└─────────────────────────────────────┘

RAG Query Flow:
──────────────

User Query
    │
    ▼
Generate Embedding
    │
    ▼
Search Vector Store
    │
    ▼
Retrieve Top-K Chunks
    │
    ▼
Augment LLM Prompt
    │
    ▼
Agent Inference
    │
    ▼
Response with Context
```

## Data Flow Through System

```
┌────────┐
│  User  │
└───┬────┘
    │ Types message in TUI
    ▼
┌─────────────────┐
│   TUI (View)    │
│  Bubble Tea     │
└───┬─────────────┘
    │ Send to agent
    ▼
┌─────────────────────────┐
│   Agent Service         │
│  - Validate input       │
│  - Check session busy   │
│  - Start async process  │
└───┬─────────────────────┘
    │
    ├────────────────────────┐
    │                        │
    ▼                        ▼
┌──────────────┐      ┌─────────────┐
│  Message     │      │  Session    │
│  Service     │      │  Service    │
│  (Create)    │      │  (Load)     │
└───┬──────────┘      └─────┬───────┘
    │                       │
    │  Save message         │  Get history
    ▼                       ▼
┌────────────────────────────────────┐
│         Database (SQLite)          │
│  messages table  │  sessions table │
└────────────────────────────────────┘
    │                       │
    │                       │
    └───────────┬───────────┘
                │
                ▼
    ┌───────────────────────┐
    │ Provider              │
    │ (LLM API Client)      │
    │ - OpenAI              │
    │ - Anthropic           │
    │ - Gemini              │
    └───┬───────────────────┘
        │ HTTP/WebSocket
        ▼
    ┌───────────────────────┐
    │   External LLM API    │
    │   (Cloud Service)     │
    └───┬───────────────────┘
        │ Stream events
        ▼
    ┌───────────────────────┐
    │  Event Processing     │
    │  - Content deltas     │
    │  - Tool calls         │
    │  - Complete           │
    └───┬───────────────────┘
        │
        ├─────────┬──────────┐
        │         │          │
        ▼         ▼          ▼
    ┌───────┐ ┌──────┐  ┌────────┐
    │Update │ │PubSub│  │Execute │
    │  DB   │ │Events│  │ Tools  │
    └───────┘ └──┬───┘  └───┬────┘
                 │          │
                 │          │
                 └────┬─────┘
                      │
                      ▼
                 ┌─────────┐
                 │   TUI   │
                 │ Updates │
                 └─────────┘
                      │
                      ▼
                 ┌─────────┐
                 │  User   │
                 │  Sees   │
                 │Response │
                 └─────────┘
```

## Configuration Flow

```
┌──────────────────┐
│  Application     │
│  Startup         │
└────┬─────────────┘
     │
     ▼
┌─────────────────────────────────────────┐
│       Configuration Loading             │
│                                         │
│  1. Environment Variables               │
│     └─▶ ANTHROPIC_API_KEY              │
│     └─▶ OPENAI_API_KEY                 │
│     └─▶ GEMINI_API_KEY                 │
│                                         │
│  2. Global Config Files                 │
│     └─▶ ~/.opencode.json               │
│     └─▶ ~/.config/opencode/.opencode.json │
│                                         │
│  3. Local Config Files                  │
│     └─▶ ./.opencode.json (project)     │
│                                         │
│  4. Merge & Validate                    │
│     └─▶ Set defaults                   │
│     └─▶ Check models exist             │
│     └─▶ Verify API keys                │
│                                         │
└────┬────────────────────────────────────┘
     │
     ▼
┌─────────────────────────────────────────┐
│     Initialized Configuration           │
│                                         │
│  • Providers (with API keys)            │
│  • Agents (with models & tokens)        │
│  • MCP Servers                          │
│  • LSP Clients                          │
│  • Data directory                       │
│                                         │
└─────────────────────────────────────────┘
```

## Session Lifecycle

```
┌──────────────┐
│ User Opens   │
│ OpenCode     │
└──────┬───────┘
       │
       ▼
┌──────────────────┐
│ Load or Create   │
│ Session          │
└──────┬───────────┘
       │
       ├─────── New Session ────────┐
       │                            │
       ▼                            ▼
   Load Existing          Create New Session
   Session from DB        ├─▶ Generate UUID
       │                  ├─▶ Set timestamp
       │                  └─▶ Save to DB
       │                            │
       └────────┬───────────────────┘
                │
                ▼
        ┌───────────────┐
        │ User Interacts│
        │ - Send msgs   │
        │ - View history│
        │ - Switch sess │
        └───┬───────────┘
            │
            │ For each message:
            ▼
    ┌──────────────────┐
    │ Create Message   │
    │ ├─▶ User msg     │
    │ ├─▶ Agent resp   │
    │ └─▶ Tool results │
    └───┬──────────────┘
        │
        │ Each message updates:
        ▼
    ┌──────────────────────┐
    │ Update Session       │
    │ ├─▶ message_count++  │
    │ ├─▶ prompt_tokens += │
    │ ├─▶ completion_tokens│
    │ ├─▶ cost +=          │
    │ └─▶ updated_at       │
    └───┬──────────────────┘
        │
        │ First message:
        ├──────────────────┐
        │                  │
        ▼                  ▼
    Generate Title    Continue Chat
    (async)               │
        │                 │
        └────────┬────────┘
                 │
                 ▼
         ┌───────────────┐
         │ User Exits    │
         │ - Session     │
         │   persisted   │
         │ - Can resume  │
         │   later       │
         └───────────────┘
```
