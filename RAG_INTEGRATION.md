# RAG Database Integration for OpenCode

## Table of Contents
1. [Overview](#overview)
2. [What is RAG?](#what-is-rag)
3. [Why RAG for OpenCode?](#why-rag-for-opencode)
4. [Current State Analysis](#current-state-analysis)
5. [RAG Integration Approaches](#rag-integration-approaches)
6. [Recommended Architecture](#recommended-architecture)
7. [Implementation Roadmap](#implementation-roadmap)
8. [Code Examples](#code-examples)
9. [Performance Considerations](#performance-considerations)
10. [Alternatives and Trade-offs](#alternatives-and-trade-offs)

## Overview

This document explores the integration of **RAG (Retrieval-Augmented Generation)** capabilities into OpenCode. RAG would enhance the AI agent's ability to access and reason over large codebases, documentation, and project-specific knowledge that exceeds the context window limitations of LLMs.

## What is RAG?

**Retrieval-Augmented Generation (RAG)** is a technique that enhances LLM responses by:

1. **Indexing**: Converting documents into vector embeddings and storing them in a searchable database
2. **Retrieval**: Finding relevant information based on semantic similarity to user queries
3. **Augmentation**: Injecting retrieved context into LLM prompts
4. **Generation**: LLM produces responses informed by the retrieved context

### RAG Flow

```
User Query
    │
    ▼
┌─────────────────────┐
│  Embed Query        │  Convert to vector using embedding model
└──────┬──────────────┘
       │
       ▼
┌─────────────────────┐
│  Search Vector DB   │  Find similar documents/chunks
└──────┬──────────────┘
       │
       ▼
┌─────────────────────┐
│  Retrieve Chunks    │  Get top-k most relevant
└──────┬──────────────┘
       │
       ▼
┌─────────────────────┐
│  Augment Prompt     │  Add context to LLM prompt
└──────┬──────────────┘
       │
       ▼
┌─────────────────────┐
│  LLM Generation     │  Generate response with context
└─────────────────────┘
```

## Why RAG for OpenCode?

### Current Limitations

1. **Context Window Limits**: Even with large context windows (200k+ tokens), full codebases don't fit
2. **Cost**: Sending large contexts to LLMs is expensive (input token costs)
3. **Relevance**: LLMs receive all or nothing; can't efficiently focus on relevant parts
4. **Real-time Updates**: Code changes require re-sending entire context

### RAG Benefits

✅ **Semantic Search**: Find relevant code by meaning, not just keywords
✅ **Scalability**: Handle codebases of any size
✅ **Cost Efficiency**: Only send relevant chunks to LLM
✅ **Better Accuracy**: Focused context improves response quality
✅ **Project Knowledge**: Index documentation, wikis, issue history
✅ **Multi-project**: Support multiple projects with isolated indices

### Use Cases for OpenCode + RAG

1. **Code Understanding**
   - "Explain the authentication flow in this project"
   - "Find all database access patterns"
   - "Show me error handling examples"

2. **Code Search**
   - "Where is user validation implemented?"
   - "Find components that use the cache system"
   - "Show similar implementations to this function"

3. **Documentation**
   - "What's the deployment process?"
   - "How do I configure the logging system?"
   - "Find examples of using the API client"

4. **Issue Resolution**
   - "Find related bug fixes"
   - "Show code that might cause this error"
   - "What changed in the authentication module?"

5. **Code Generation**
   - "Generate a controller similar to UserController"
   - "Create tests following project patterns"
   - "Add error handling like in other modules"

## Current State Analysis

### What OpenCode Has

✅ **File Operations**: `ls`, `view`, `grep`, `glob` tools
✅ **Code Intelligence**: LSP integration for diagnostics
✅ **Search**: Sourcegraph integration for public code
✅ **Context Loading**: Automatic loading of project instructions
✅ **Tool System**: Extensible architecture for new capabilities
✅ **Database**: SQLite for persistent storage
✅ **MCP Support**: External tool integration protocol

### What's Missing for RAG

❌ **Embedding Generation**: No vector embedding capability
❌ **Vector Database**: No similarity search infrastructure
❌ **Chunking Strategy**: No intelligent code/doc splitting
❌ **Index Management**: No indexing or re-indexing system
❌ **Embedding Models**: No local or API-based embedding service
❌ **Semantic Search**: No semantic retrieval mechanism

## RAG Integration Approaches

### Approach 1: Native SQLite with Vector Extension

**Technology**: SQLite + sqlite-vec extension

**Pros**:
- ✅ Reuses existing SQLite infrastructure
- ✅ Single database file, easy deployment
- ✅ No additional services required
- ✅ Good for small-to-medium projects

**Cons**:
- ❌ Limited vector search performance at scale
- ❌ Less optimized than dedicated vector DBs
- ❌ Requires CGO and custom sqlite3 build

**Example**:
```sql
-- Using sqlite-vec extension
CREATE VIRTUAL TABLE code_embeddings USING vec0(
    embedding FLOAT[1536],  -- OpenAI ada-002 dimensions
    content TEXT,
    file_path TEXT,
    chunk_id TEXT
);

-- Similarity search
SELECT 
    file_path, 
    content,
    distance 
FROM code_embeddings 
WHERE embedding MATCH ?
ORDER BY distance 
LIMIT 10;
```

### Approach 2: Embedded Vector Database

**Technology**: Qdrant embedded mode or Meilisearch

**Pros**:
- ✅ Purpose-built for vector search
- ✅ Better performance than SQLite extensions
- ✅ No external services required
- ✅ Richer filtering and querying

**Cons**:
- ❌ Additional dependency
- ❌ Separate data files
- ❌ Increased binary size

**Libraries**:
- [Qdrant Go Client](https://github.com/qdrant/go-client)
- [Weaviate embedded](https://weaviate.io/)

### Approach 3: External Vector Database

**Technology**: Qdrant, Weaviate, Pinecone, Chroma

**Pros**:
- ✅ Best performance and scalability
- ✅ Advanced features (hybrid search, reranking)
- ✅ Multiple project support
- ✅ Cloud or self-hosted options

**Cons**:
- ❌ Requires external service setup
- ❌ Network dependency
- ❌ More complex deployment
- ❌ May require authentication

### Approach 4: MCP-based RAG Server

**Technology**: External MCP server with RAG capabilities

**Pros**:
- ✅ Leverages existing MCP integration
- ✅ Pluggable architecture
- ✅ Can be implemented independently
- ✅ Users can choose their preferred RAG solution

**Cons**:
- ❌ Not built-in experience
- ❌ Requires user configuration
- ❌ External dependency management

## Recommended Architecture

### Hybrid Approach: Built-in + MCP Extension

Combine **Approach 1** (SQLite-vec for baseline) + **Approach 4** (MCP for advanced RAG).

### Architecture Overview

```
┌───────────────────────────────────────────────────────────┐
│                      OpenCode Core                         │
├───────────────────────────────────────────────────────────┤
│                                                            │
│  ┌─────────────────────────────────────────────────────┐  │
│  │              Built-in RAG (Basic)                   │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────┐  │  │
│  │  │  Embeddings  │→ │ SQLite-vec   │→ │  Search  │  │  │
│  │  │  (via API)   │  │  (Vectors)   │  │  Tool    │  │  │
│  │  └──────────────┘  └──────────────┘  └──────────┘  │  │
│  └─────────────────────────────────────────────────────┘  │
│                                                            │
│  ┌─────────────────────────────────────────────────────┐  │
│  │           MCP RAG Servers (Advanced)                │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────┐  │  │
│  │  │   Qdrant     │  │   Weaviate   │  │  Chroma  │  │  │
│  │  │   Server     │  │   Server     │  │  Server  │  │  │
│  │  └──────────────┘  └──────────────┘  └──────────┘  │  │
│  └─────────────────────────────────────────────────────┘  │
│                                                            │
└───────────────────────────────────────────────────────────┘

┌───────────────────────────────────────────────────────────┐
│                   Embedding Services                       │
├───────────────────────────────────────────────────────────┤
│  OpenAI API │ Cohere API │ Local (ONNX) │ Anthropic API  │
└───────────────────────────────────────────────────────────┘
```

### Key Components

#### 1. Embedding Service (`internal/rag/embeddings/`)

```go
type EmbeddingService interface {
    // Generate embeddings for text
    Embed(ctx context.Context, texts []string) ([][]float32, error)
    
    // Get embedding dimensions
    Dimensions() int
    
    // Get model name
    Model() string
}

// Implementations
type OpenAIEmbeddings struct { ... }
type CohereEmbeddings struct { ... }
type LocalEmbeddings struct { ... }  // ONNX runtime
```

#### 2. Vector Store (`internal/rag/vectorstore/`)

```go
type VectorStore interface {
    // Add documents to the store
    AddDocuments(ctx context.Context, docs []Document) error
    
    // Search for similar documents
    Search(ctx context.Context, query string, k int) ([]Document, error)
    
    // Search by vector
    SearchByVector(ctx context.Context, vector []float32, k int) ([]Document, error)
    
    // Delete documents by filter
    Delete(ctx context.Context, filter Filter) error
}

type Document struct {
    ID       string
    Content  string
    Metadata map[string]interface{}
    Vector   []float32
    Score    float32
}
```

#### 3. Chunking Strategy (`internal/rag/chunker/`)

```go
type Chunker interface {
    // Split content into chunks
    Chunk(content string, metadata map[string]interface{}) []Chunk
}

type Chunk struct {
    Content  string
    Metadata map[string]interface{}
    Start    int
    End      int
}

// Implementations
type SemanticChunker struct { ... }       // Split by meaning
type TokenChunker struct { ... }          // Fixed token size
type CodeChunker struct { ... }           // Language-aware splits
type MarkdownChunker struct { ... }       // Markdown sections
```

#### 4. Indexer (`internal/rag/indexer/`)

```go
type Indexer interface {
    // Index a project
    IndexProject(ctx context.Context, projectPath string) error
    
    // Update index for changed files
    UpdateFiles(ctx context.Context, files []string) error
    
    // Get index status
    Status(ctx context.Context) IndexStatus
}

type IndexStatus struct {
    DocumentCount int
    LastIndexed   time.Time
    IsIndexing    bool
}
```

#### 5. RAG Tools (`internal/llm/tools/rag.go`)

```go
// New tools for agents
type SemanticSearchTool struct { ... }
type CodeSearchTool struct { ... }
type DocumentSearchTool struct { ... }
```

### Database Schema Extension

```sql
-- Embedding metadata
CREATE TABLE embeddings (
    id TEXT PRIMARY KEY,
    session_id TEXT,
    file_path TEXT NOT NULL,
    chunk_index INTEGER NOT NULL,
    chunk_text TEXT NOT NULL,
    embedding_model TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    FOREIGN KEY (session_id) REFERENCES sessions (id) ON DELETE CASCADE
);

CREATE INDEX idx_embeddings_session_id ON embeddings (session_id);
CREATE INDEX idx_embeddings_file_path ON embeddings (file_path);

-- Vector table (using sqlite-vec)
CREATE VIRTUAL TABLE code_vectors USING vec0(
    embedding_id TEXT PRIMARY KEY,
    vector FLOAT[1536]
);

-- Index configuration
CREATE TABLE index_config (
    project_path TEXT PRIMARY KEY,
    embedding_model TEXT NOT NULL,
    chunk_size INTEGER NOT NULL,
    chunk_overlap INTEGER NOT NULL,
    last_indexed INTEGER NOT NULL,
    document_count INTEGER NOT NULL
);
```

### Configuration Extension

```json
{
  "rag": {
    "enabled": true,
    "embeddingService": "openai",
    "embeddingModel": "text-embedding-3-small",
    "chunkSize": 1000,
    "chunkOverlap": 200,
    "vectorStore": "sqlite-vec",
    "autoIndex": true,
    "indexPatterns": [
      "**/*.go",
      "**/*.js",
      "**/*.py",
      "**/*.md",
      "!**/node_modules/**",
      "!**/vendor/**"
    ],
    "retrievalK": 10
  },
  "providers": {
    "openai": {
      "apiKey": "sk-...",
      "disabled": false
    }
  }
}
```

## Implementation Roadmap

### Phase 1: Foundation (Weeks 1-2)

- [ ] Design vector store interface
- [ ] Implement embedding service abstraction
- [ ] Add OpenAI embeddings provider
- [ ] Create basic chunking strategies
- [ ] Set up SQLite-vec integration

### Phase 2: Indexing (Weeks 3-4)

- [ ] Implement file crawler and indexer
- [ ] Add incremental indexing (watch for changes)
- [ ] Create index management CLI commands
- [ ] Add progress tracking and status
- [ ] Implement caching and deduplication

### Phase 3: Retrieval (Weeks 5-6)

- [ ] Implement semantic search tool
- [ ] Add hybrid search (keyword + semantic)
- [ ] Create result ranking and reranking
- [ ] Add metadata filtering
- [ ] Implement caching for queries

### Phase 4: Agent Integration (Weeks 7-8)

- [ ] Register RAG tools with agent
- [ ] Update prompts to use RAG context
- [ ] Add automatic context augmentation
- [ ] Implement smart context selection
- [ ] Add RAG usage tracking and metrics

### Phase 5: UX and Polish (Weeks 9-10)

- [ ] Add TUI for index management
- [ ] Create status indicators
- [ ] Add search result visualization
- [ ] Implement cost estimation
- [ ] Add configuration validation

### Phase 6: Advanced Features (Future)

- [ ] Multi-modal embeddings (code + docs)
- [ ] Conversation memory RAG
- [ ] Custom embedding models
- [ ] MCP RAG server protocol
- [ ] Distributed vector stores

## Code Examples

### Example 1: Basic RAG Tool Usage

```go
// internal/llm/tools/semantic_search.go
package tools

import (
    "context"
    "encoding/json"
    "fmt"
    
    "github.com/opencode-ai/opencode/internal/rag"
)

type SemanticSearchTool struct {
    vectorStore rag.VectorStore
    embeddings  rag.EmbeddingService
}

func (t *SemanticSearchTool) Info() ToolInfo {
    return ToolInfo{
        Name:        "semantic_search",
        Description: "Search the codebase using semantic similarity",
        Parameters: map[string]any{
            "type": "object",
            "properties": map[string]any{
                "query": map[string]any{
                    "type":        "string",
                    "description": "Natural language search query",
                },
                "k": map[string]any{
                    "type":        "integer",
                    "description": "Number of results to return",
                    "default":     10,
                },
                "file_types": map[string]any{
                    "type":        "array",
                    "description": "Filter by file extensions",
                    "items": map[string]any{
                        "type": "string",
                    },
                },
            },
            "required": []string{"query"},
        },
    }
}

func (t *SemanticSearchTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
    var params struct {
        Query     string   `json:"query"`
        K         int      `json:"k"`
        FileTypes []string `json:"file_types"`
    }
    
    if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
        return NewTextErrorResponse("Invalid parameters"), err
    }
    
    if params.K <= 0 {
        params.K = 10
    }
    
    // Search vector store
    results, err := t.vectorStore.Search(ctx, params.Query, params.K)
    if err != nil {
        return NewTextErrorResponse(fmt.Sprintf("Search failed: %v", err)), err
    }
    
    // Format results
    var response strings.Builder
    response.WriteString(fmt.Sprintf("Found %d relevant code chunks:\n\n", len(results)))
    
    for i, result := range results {
        response.WriteString(fmt.Sprintf("## Result %d (Score: %.3f)\n", i+1, result.Score))
        response.WriteString(fmt.Sprintf("File: %s\n", result.Metadata["file_path"]))
        response.WriteString(fmt.Sprintf("```\n%s\n```\n\n", result.Content))
    }
    
    return NewTextResponse(response.String()), nil
}
```

### Example 2: Automatic Context Augmentation

```go
// internal/llm/agent/rag_augmentation.go
package agent

import (
    "context"
    "fmt"
    
    "github.com/opencode-ai/opencode/internal/rag"
)

func (a *agent) augmentWithRAG(ctx context.Context, userMessage string) (string, error) {
    cfg := config.Get()
    if !cfg.RAG.Enabled {
        return userMessage, nil
    }
    
    // Search for relevant context
    results, err := a.vectorStore.Search(ctx, userMessage, cfg.RAG.RetrievalK)
    if err != nil {
        logging.Warn("RAG search failed", "error", err)
        return userMessage, nil
    }
    
    if len(results) == 0 {
        return userMessage, nil
    }
    
    // Build augmented prompt
    var augmented strings.Builder
    augmented.WriteString("# Relevant Context from Codebase\n\n")
    
    for i, result := range results {
        augmented.WriteString(fmt.Sprintf("## Context %d\n", i+1))
        augmented.WriteString(fmt.Sprintf("Source: %s\n", result.Metadata["file_path"]))
        augmented.WriteString(fmt.Sprintf("```\n%s\n```\n\n", result.Content))
    }
    
    augmented.WriteString("# User Request\n\n")
    augmented.WriteString(userMessage)
    
    return augmented.String(), nil
}
```

### Example 3: Indexing Service

```go
// internal/rag/indexer/indexer.go
package indexer

import (
    "context"
    "path/filepath"
    
    "github.com/opencode-ai/opencode/internal/rag"
)

type Indexer struct {
    vectorStore rag.VectorStore
    embeddings  rag.EmbeddingService
    chunker     rag.Chunker
}

func (idx *Indexer) IndexProject(ctx context.Context, projectPath string) error {
    cfg := config.Get()
    
    // Walk directory tree
    var files []string
    err := filepath.WalkDir(projectPath, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        
        if d.IsDir() {
            return nil
        }
        
        // Check if file matches patterns
        if idx.shouldIndex(path, cfg.RAG.IndexPatterns) {
            files = append(files, path)
        }
        
        return nil
    })
    
    if err != nil {
        return fmt.Errorf("failed to walk directory: %w", err)
    }
    
    // Index each file
    for _, file := range files {
        if err := idx.indexFile(ctx, file); err != nil {
            logging.Warn("failed to index file", "file", file, "error", err)
            continue
        }
    }
    
    return nil
}

func (idx *Indexer) indexFile(ctx context.Context, filePath string) error {
    // Read file
    content, err := os.ReadFile(filePath)
    if err != nil {
        return err
    }
    
    // Chunk content
    chunks := idx.chunker.Chunk(string(content), map[string]interface{}{
        "file_path": filePath,
        "file_type": filepath.Ext(filePath),
    })
    
    // Generate embeddings
    texts := make([]string, len(chunks))
    for i, chunk := range chunks {
        texts[i] = chunk.Content
    }
    
    vectors, err := idx.embeddings.Embed(ctx, texts)
    if err != nil {
        return err
    }
    
    // Create documents
    docs := make([]rag.Document, len(chunks))
    for i, chunk := range chunks {
        docs[i] = rag.Document{
            ID:       fmt.Sprintf("%s:%d", filePath, i),
            Content:  chunk.Content,
            Metadata: chunk.Metadata,
            Vector:   vectors[i],
        }
    }
    
    // Add to vector store
    return idx.vectorStore.AddDocuments(ctx, docs)
}
```

## Performance Considerations

### Embedding Generation

| Service | Cost per 1M tokens | Speed | Dimensions |
|---------|-------------------|-------|------------|
| OpenAI text-embedding-3-small | $0.02 | Fast | 1536 |
| OpenAI text-embedding-3-large | $0.13 | Fast | 3072 |
| Cohere embed-english-v3.0 | $0.10 | Fast | 1024 |
| Local ONNX (all-MiniLM-L6-v2) | Free | Medium | 384 |

### Storage Requirements

Average codebase: **~10MB of code** → **~10,000 chunks** → **~15MB vectors**

```
Chunk size: 1000 tokens
Embedding dimensions: 1536 (OpenAI)
Storage per chunk: 1536 * 4 bytes = 6KB

10,000 chunks = 60MB of vector data
```

### Query Performance

- **SQLite-vec**: ~1-5ms for 10k vectors
- **Qdrant**: ~1-2ms for 100k+ vectors
- **Weaviate**: ~2-3ms for 1M+ vectors

### Optimization Strategies

1. **Batch embeddings**: Generate in batches of 100-1000
2. **Cache embeddings**: Don't re-embed unchanged files
3. **Incremental indexing**: Only update changed files
4. **Smart chunking**: Preserve semantic boundaries
5. **Metadata filters**: Pre-filter before vector search
6. **Reranking**: Use LLM to rerank top-k results

## Alternatives and Trade-offs

### Alternative 1: No RAG, Use Larger Context Windows

**Pros**:
- ✅ Simpler architecture
- ✅ No embedding costs
- ✅ No indexing overhead

**Cons**:
- ❌ Expensive for large codebases
- ❌ Slower inference
- ❌ Not scalable beyond ~500k tokens

### Alternative 2: Keyword-based Search Only

**Pros**:
- ✅ Fast and cheap
- ✅ No embedding required
- ✅ Exact matches

**Cons**:
- ❌ Misses semantic similarity
- ❌ Requires knowing exact terms
- ❌ Poor for natural language queries

### Alternative 3: Full-text + Semantic Hybrid

**Pros**:
- ✅ Best of both worlds
- ✅ Handles exact and fuzzy matches
- ✅ Most accurate retrieval

**Cons**:
- ❌ More complex implementation
- ❌ Higher computational cost
- ❌ Requires tuning weights

## Conclusion

RAG integration is **highly feasible** for OpenCode and would provide significant value:

### Recommended Next Steps

1. **Prototype** with SQLite-vec for basic functionality
2. **Implement** OpenAI embeddings service integration
3. **Create** semantic search tool for agents
4. **Test** with real codebases to validate approach
5. **Iterate** based on user feedback and performance

### Success Metrics

- ✅ Index 100k+ lines of code in <1 minute
- ✅ Search latency <100ms per query
- ✅ Relevance: Top-5 results contain answer >80% of time
- ✅ Cost: <$0.10 per 1000 queries
- ✅ Storage: <100MB for average project

RAG will transform OpenCode from a **context-limited** assistant to a **project-aware** coding companion that truly understands your entire codebase.
