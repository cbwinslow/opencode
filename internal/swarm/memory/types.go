package memory

import (
	"time"
)

// MemoryType defines different types of memory
type MemoryType string

const (
	MemoryTypeWorking    MemoryType = "working"    // Short-term, current context
	MemoryTypeEpisodic   MemoryType = "episodic"   // Event-based memories
	MemoryTypeSemantic   MemoryType = "semantic"   // Factual knowledge
	MemoryTypeProcedural MemoryType = "procedural" // How-to knowledge
)

// MemoryPriority defines importance levels
type MemoryPriority int

const (
	PriorityLow MemoryPriority = iota
	PriorityNormal
	PriorityHigh
	PriorityCritical
)

// Memory represents a single memory unit
type Memory struct {
	ID          string
	Type        MemoryType
	Content     interface{}
	Metadata    map[string]interface{}
	Vector      []float64 // Embedding for semantic search
	Tags        []string
	Priority    MemoryPriority
	AccessCount int
	LastAccessed time.Time
	CreatedAt   time.Time
	ExpiresAt   *time.Time
	Encrypted   bool
	Parent      string // For hierarchical organization
	Children    []string
}

// MemoryQuery represents a query for memories
type MemoryQuery struct {
	Type         MemoryType
	Tags         []string
	SearchText   string
	Vector       []float64
	Limit        int
	MinPriority  MemoryPriority
	TimeRange    *TimeRange
	IncludeChildren bool
}

// TimeRange defines a time period
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// MemoryStore defines the interface for memory storage
type MemoryStore interface {
	// CRUD operations
	Store(memory Memory) error
	Retrieve(id string) (*Memory, error)
	Update(id string, memory Memory) error
	Delete(id string) error
	
	// Query operations
	Query(query MemoryQuery) ([]Memory, error)
	VectorSearch(vector []float64, limit int) ([]Memory, error)
	
	// Maintenance operations
	Consolidate() error
	Prune(criteria PruneCriteria) error
	
	// Statistics
	GetStats() MemoryStats
}

// PruneCriteria defines what memories to remove
type PruneCriteria struct {
	MaxAge         time.Duration
	MinAccessCount int
	MaxMemories    int
	PreserveTags   []string
}

// MemoryStats contains statistics about the memory store
type MemoryStats struct {
	TotalMemories      int
	MemoriesByType     map[MemoryType]int
	TotalSize          int64
	AverageAccessCount float64
	OldestMemory       time.Time
	NewestMemory       time.Time
}

// HierarchicalNode represents a node in the memory hierarchy
type HierarchicalNode struct {
	ID       string
	Type     MemoryType
	Summary  string
	Children []*HierarchicalNode
	Memories []Memory
	Level    int
}

// ConsolidationStrategy defines how memories are consolidated
type ConsolidationStrategy string

const (
	ConsolidationByTime     ConsolidationStrategy = "time"     // Group by time period
	ConsolidationByTopic    ConsolidationStrategy = "topic"    // Group by semantic similarity
	ConsolidationByActivity ConsolidationStrategy = "activity" // Group by activity type
)
