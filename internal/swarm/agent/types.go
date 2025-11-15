package agent

import (
	"context"
	"time"
)

// AgentType defines the specialization of an agent
type AgentType string

const (
	// Core agent types
	AgentTypeCoordinator    AgentType = "coordinator"     // Orchestrates swarm activities
	AgentTypeMonitor        AgentType = "monitor"         // Monitors logs and system state
	AgentTypeAnalyzer       AgentType = "analyzer"        // Analyzes data and patterns
	AgentTypeExecutor       AgentType = "executor"        // Executes tasks and actions
	AgentTypeMemory         AgentType = "memory"          // Manages memory systems
	AgentTypeLearning       AgentType = "learning"        // Learns from successes/failures
	AgentTypeDocumentation  AgentType = "documentation"   // Manages documentation
	AgentTypeTesting        AgentType = "testing"         // Runs and validates tests
	AgentTypeErrorHandler   AgentType = "error_handler"   // Handles errors and recovery
	AgentTypeHealthChecker  AgentType = "health_checker"  // Monitors agent and system health
)

// AgentStatus represents the current state of an agent
type AgentStatus string

const (
	AgentStatusIdle       AgentStatus = "idle"
	AgentStatusBusy       AgentStatus = "busy"
	AgentStatusError      AgentStatus = "error"
	AgentStatusStopped    AgentStatus = "stopped"
	AgentStatusStarting   AgentStatus = "starting"
)

// Agent represents a specialized AI agent in the swarm
type Agent interface {
	// Lifecycle
	Start(ctx context.Context) error
	Stop() error
	GetStatus() AgentStatus
	
	// Identity
	GetID() string
	GetType() AgentType
	GetCapabilities() []string
	
	// Task execution
	ExecuteTask(ctx context.Context, task Task) (*TaskResult, error)
	CanHandleTask(task Task) bool
	
	// Communication
	SendMessage(msg Message) error
	ReceiveMessages() <-chan Message
	
	// Health and metrics
	GetHealthScore() float64
	GetMetrics() AgentMetrics
}

// Task represents work to be done by an agent
type Task struct {
	ID          string
	Type        string
	Priority    int
	Description string
	Input       map[string]interface{}
	CreatedAt   time.Time
	Deadline    *time.Time
	RetryCount  int
	MaxRetries  int
}

// TaskResult contains the outcome of a task execution
type TaskResult struct {
	TaskID      string
	Success     bool
	Output      map[string]interface{}
	Error       error
	ExecutionTime time.Duration
	AgentID     string
	CompletedAt time.Time
	Metadata    map[string]interface{}
}

// Message represents communication between agents
type Message struct {
	ID        string
	From      string
	To        string // Empty for broadcast
	Type      MessageType
	Content   interface{}
	Timestamp time.Time
	ReplyTo   string
}

// MessageType defines different message categories
type MessageType string

const (
	MessageTypeTaskRequest    MessageType = "task_request"
	MessageTypeTaskResponse   MessageType = "task_response"
	MessageTypeStatusUpdate   MessageType = "status_update"
	MessageTypeVoteRequest    MessageType = "vote_request"
	MessageTypeVoteResponse   MessageType = "vote_response"
	MessageTypeHealthCheck    MessageType = "health_check"
	MessageTypeBroadcast      MessageType = "broadcast"
	MessageTypeLogEntry       MessageType = "log_entry"
	MessageTypeMemoryUpdate   MessageType = "memory_update"
)

// AgentMetrics contains performance and operational metrics
type AgentMetrics struct {
	TasksCompleted    int
	TasksFailed       int
	AverageTaskTime   time.Duration
	LastActivityTime  time.Time
	MessagesReceived  int
	MessagesSent      int
	ErrorCount        int
	UptimeSeconds     int64
	CPUUsage          float64
	MemoryUsage       int64
}

// AgentConfig contains configuration for an agent
type AgentConfig struct {
	ID              string
	Type            AgentType
	ProviderType    string // "openrouter", "ollama", "lmstudio", "huggingface", "jan"
	Model           string
	MaxConcurrency  int
	HealthCheckInterval time.Duration
	MessageBufferSize   int
	EnableLearning  bool
	Capabilities    []string
	CustomConfig    map[string]interface{}
}

// SwarmConfig contains configuration for the entire swarm
type SwarmConfig struct {
	Name               string
	Agents             []AgentConfig
	VotingThreshold    float64 // Percentage for democratic decisions
	MaxConcurrentTasks int
	HealthCheckInterval time.Duration
	EnableMemory       bool
	EnableLearning     bool
	EnableSelfHealing  bool
	LogLevel           string
}
