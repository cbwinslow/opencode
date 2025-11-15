package swarm

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/opencode-ai/opencode/internal/swarm/agent"
	"github.com/opencode-ai/opencode/internal/swarm/health"
	"github.com/opencode-ai/opencode/internal/swarm/memory"
	"github.com/opencode-ai/opencode/internal/swarm/monitor"
	"github.com/opencode-ai/opencode/internal/swarm/rules"
	"github.com/opencode-ai/opencode/internal/swarm/voting"
)

// Coordinator manages the entire multi-agent swarm system
type Coordinator struct {
	config agent.SwarmConfig
	
	// Core components
	registry      *agent.Registry
	memoryStore   memory.MemoryStore
	votingSystem  *voting.DemocraticVotingSystem
	ruleEngine    *rules.RuleEngine
	healthMonitor *health.HealthMonitor
	
	// Monitoring
	logWatcher     *monitor.LogWatcher
	historyWatcher *monitor.ShellHistoryWatcher
	
	// Task management
	taskQueue     chan agent.Task
	taskResults   chan *agent.TaskResult
	
	// Lifecycle
	ctx        context.Context
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
	mu         sync.Mutex
	running    bool
}

// CoordinatorConfig contains configuration for the coordinator
type CoordinatorConfig struct {
	SwarmConfig    agent.SwarmConfig
	MemoryConfig   memory.HierarchicalMemoryConfig
	HealthConfig   health.HealthMonitorConfig
	LogPaths       []string
	ShellHistory   string
	TaskQueueSize  int
}

// NewCoordinator creates a new swarm coordinator
func NewCoordinator(config CoordinatorConfig) (*Coordinator, error) {
	ctx, cancel := context.WithCancel(context.Background())
	
	if config.TaskQueueSize <= 0 {
		config.TaskQueueSize = 1000
	}
	
	// Initialize components
	registry := agent.NewRegistry()
	memoryStore := memory.NewHierarchicalMemoryStore(config.MemoryConfig)
	votingSystem := voting.NewDemocraticVotingSystem()
	ruleEngine := rules.NewRuleEngine(rules.RuleEngineConfig{
		MaxHistory:    10000,
		EnableHistory: true,
		ParallelExec:  true,
	})
	healthMonitor := health.NewHealthMonitor(config.HealthConfig)
	
	// Initialize monitoring
	var logWatcher *monitor.LogWatcher
	var historyWatcher *monitor.ShellHistoryWatcher
	var err error
	
	if len(config.LogPaths) > 0 {
		logWatcher, err = monitor.NewLogWatcher(monitor.LogWatcherConfig{
			Paths:      config.LogPaths,
			BufferSize: 1000,
		})
		if err != nil {
			cancel()
			return nil, fmt.Errorf("failed to create log watcher: %w", err)
		}
	}
	
	if config.ShellHistory != "" {
		historyWatcher, err = monitor.NewShellHistoryWatcher(config.ShellHistory, 100)
		if err != nil {
			cancel()
			return nil, fmt.Errorf("failed to create history watcher: %w", err)
		}
	}
	
	coordinator := &Coordinator{
		config:         config.SwarmConfig,
		registry:       registry,
		memoryStore:    memoryStore,
		votingSystem:   votingSystem,
		ruleEngine:     ruleEngine,
		healthMonitor:  healthMonitor,
		logWatcher:     logWatcher,
		historyWatcher: historyWatcher,
		taskQueue:      make(chan agent.Task, config.TaskQueueSize),
		taskResults:    make(chan *agent.TaskResult, config.TaskQueueSize),
		ctx:            ctx,
		cancelFunc:     cancel,
	}
	
	return coordinator, nil
}

// Start initializes and starts the swarm
func (c *Coordinator) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if c.running {
		return fmt.Errorf("coordinator already running")
	}
	
	// Start health monitor
	if err := c.healthMonitor.Start(); err != nil {
		return fmt.Errorf("failed to start health monitor: %w", err)
	}
	
	// Start monitoring
	if c.logWatcher != nil {
		if err := c.logWatcher.Start(); err != nil {
			return fmt.Errorf("failed to start log watcher: %w", err)
		}
		
		// Process log entries
		c.wg.Add(1)
		go c.processLogEntries()
	}
	
	if c.historyWatcher != nil {
		if err := c.historyWatcher.Start(); err != nil {
			return fmt.Errorf("failed to start history watcher: %w", err)
		}
		
		// Process history entries
		c.wg.Add(1)
		go c.processHistoryEntries()
	}
	
	// Start task processing
	c.wg.Add(1)
	go c.processTaskQueue()
	
	// Start result processing
	c.wg.Add(1)
	go c.processTaskResults()
	
	// Start agents
	if err := c.registry.StartAll(c.ctx); err != nil {
		return fmt.Errorf("failed to start agents: %w", err)
	}
	
	// Load default rules
	if err := c.loadDefaultRules(); err != nil {
		return fmt.Errorf("failed to load rules: %w", err)
	}
	
	c.running = true
	return nil
}

// Stop gracefully shuts down the swarm
func (c *Coordinator) Stop() error {
	c.mu.Lock()
	if !c.running {
		c.mu.Unlock()
		return nil
	}
	c.running = false
	c.mu.Unlock()
	
	// Stop components
	c.cancelFunc()
	
	// Stop agents
	if err := c.registry.StopAll(); err != nil {
		return err
	}
	
	// Stop monitoring
	if c.logWatcher != nil {
		_ = c.logWatcher.Stop()
	}
	if c.historyWatcher != nil {
		_ = c.historyWatcher.Stop()
	}
	
	// Stop health monitor
	_ = c.healthMonitor.Stop()
	
	// Wait for goroutines
	c.wg.Wait()
	
	// Close channels
	close(c.taskQueue)
	close(c.taskResults)
	
	return nil
}

// SubmitTask adds a task to the queue
func (c *Coordinator) SubmitTask(task agent.Task) error {
	select {
	case c.taskQueue <- task:
		return nil
	case <-c.ctx.Done():
		return fmt.Errorf("coordinator stopped")
	default:
		return fmt.Errorf("task queue full")
	}
}

// GetTaskResult waits for a task result
func (c *Coordinator) GetTaskResult(taskID string, timeout time.Duration) (*agent.TaskResult, error) {
	ctx, cancel := context.WithTimeout(c.ctx, timeout)
	defer cancel()
	
	for {
		select {
		case result := <-c.taskResults:
			if result.TaskID == taskID {
				return result, nil
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout waiting for task result")
		}
	}
}

// processTaskQueue handles task distribution
func (c *Coordinator) processTaskQueue() {
	defer c.wg.Done()
	
	for {
		select {
		case task, ok := <-c.taskQueue:
			if !ok {
				return
			}
			
			// Find suitable agents
			agents := c.registry.FindAgentsForTask(task)
			
			if len(agents) == 0 {
				// No agents available, requeue or fail
				continue
			}
			
			// If multiple agents can handle it, use democratic voting
			if len(agents) > 1 && c.config.VotingThreshold > 0 {
				c.handleTaskWithVoting(task, agents)
			} else {
				// Assign to first available agent
				go c.executeTask(agents[0], task)
			}
			
		case <-c.ctx.Done():
			return
		}
	}
}

// executeTask executes a task on an agent
func (c *Coordinator) executeTask(ag agent.Agent, task agent.Task) {
	ctx, cancel := context.WithTimeout(c.ctx, 5*time.Minute)
	defer cancel()
	
	result, err := ag.ExecuteTask(ctx, task)
	if err != nil {
		result = &agent.TaskResult{
			TaskID:      task.ID,
			Success:     false,
			Error:       err,
			AgentID:     ag.GetID(),
			CompletedAt: time.Now(),
		}
	}
	
	// Store result in memory
	c.storeTaskResult(result)
	
	// Send result
	select {
	case c.taskResults <- result:
	case <-c.ctx.Done():
	}
}

// handleTaskWithVoting uses democratic voting for task decisions
func (c *Coordinator) handleTaskWithVoting(task agent.Task, agents []agent.Agent) {
	// Create a vote on how to handle the task
	proposal := voting.VoteProposal{
		Description: fmt.Sprintf("Should we execute task: %s", task.Description),
		Context: map[string]interface{}{
			"task": task,
		},
		Deadline: time.Now().Add(30 * time.Second),
	}
	
	session, err := c.votingSystem.CreateVoteSession(
		proposal,
		voting.VoteTypeMajority,
		len(agents),
		nil,
	)
	if err != nil {
		return
	}
	
	// Collect votes from agents (simplified - would need actual agent input)
	for _, ag := range agents {
		vote := voting.Vote{
			AgentID:    ag.GetID(),
			Decision:   ag.CanHandleTask(task),
			Confidence: ag.GetHealthScore(),
			Reasoning:  "Agent capability assessment",
		}
		_ = c.votingSystem.CastVote(session.ID, vote)
	}
	
	// Wait for result
	ctx, cancel := context.WithTimeout(c.ctx, 1*time.Minute)
	defer cancel()
	
	result, err := c.votingSystem.WaitForResult(ctx, session.ID)
	if err == nil && result.Decision {
		// Execute on the agent with highest confidence
		bestAgent := agents[0]
		c.executeTask(bestAgent, task)
	}
}

// processTaskResults handles task results
func (c *Coordinator) processTaskResults() {
	defer c.wg.Done()
	
	for {
		select {
		case result, ok := <-c.taskResults:
			if !ok {
				return
			}
			
			// Analyze and learn from results
			c.learnFromResult(result)
			
		case <-c.ctx.Done():
			return
		}
	}
}

// processLogEntries handles log monitoring
func (c *Coordinator) processLogEntries() {
	defer c.wg.Done()
	
	for {
		select {
		case entry, ok := <-c.logWatcher.Entries():
			if !ok {
				return
			}
			
			// Store in memory
			mem := memory.Memory{
				Type:     memory.MemoryTypeEpisodic,
				Content:  entry,
				Tags:     []string{"log", entry.Level},
				Priority: memory.PriorityNormal,
			}
			_ = c.memoryStore.Store(mem)
			
			// Evaluate rules
			ruleCtx := rules.RuleContext{
				EventType: "log_entry",
				EventData: map[string]interface{}{
					"level":   entry.Level,
					"message": entry.Message,
					"source":  entry.Source,
				},
				Timestamp: entry.Timestamp,
			}
			_ = c.ruleEngine.EvaluateRules(c.ctx, ruleCtx)
			
		case <-c.ctx.Done():
			return
		}
	}
}

// processHistoryEntries handles shell history monitoring
func (c *Coordinator) processHistoryEntries() {
	defer c.wg.Done()
	
	for {
		select {
		case entry, ok := <-c.historyWatcher.Entries():
			if !ok {
				return
			}
			
			// Store in memory
			mem := memory.Memory{
				Type:     memory.MemoryTypeEpisodic,
				Content:  entry,
				Tags:     []string{"shell", "command"},
				Priority: memory.PriorityNormal,
			}
			_ = c.memoryStore.Store(mem)
			
		case <-c.ctx.Done():
			return
		}
	}
}

// storeTaskResult stores task results in memory
func (c *Coordinator) storeTaskResult(result *agent.TaskResult) {
	tags := []string{"task", "result"}
	priority := memory.PriorityNormal
	
	if result.Success {
		tags = append(tags, "success")
		priority = memory.PriorityHigh
	} else {
		tags = append(tags, "failure")
		priority = memory.PriorityHigh // Learn from failures
	}
	
	mem := memory.Memory{
		Type:     memory.MemoryTypeProcedural,
		Content:  result,
		Tags:     tags,
		Priority: priority,
		Metadata: map[string]interface{}{
			"task_id":  result.TaskID,
			"agent_id": result.AgentID,
			"success":  result.Success,
		},
	}
	
	_ = c.memoryStore.Store(mem)
}

// learnFromResult analyzes task results for learning
func (c *Coordinator) learnFromResult(result *agent.TaskResult) {
	// Query similar past results
	query := memory.MemoryQuery{
		Type:  memory.MemoryTypeProcedural,
		Tags:  []string{"task", "result"},
		Limit: 10,
	}
	
	similar, _ := c.memoryStore.Query(query)
	
	// Analyze patterns (simplified)
	successRate := 0.0
	if len(similar) > 0 {
		successCount := 0
		for _, mem := range similar {
			if taskResult, ok := mem.Content.(*agent.TaskResult); ok {
				if taskResult.Success {
					successCount++
				}
			}
		}
		successRate = float64(successCount) / float64(len(similar))
	}
	
	// Update agent health based on performance
	if result.Success {
		// Positive reinforcement
	} else {
		// Negative feedback, may trigger recovery
		c.healthMonitor.UpdateCheck(health.HealthCheck{
			ComponentID: result.AgentID,
			Status:      health.HealthStatusDegraded,
			Score:       successRate,
			Message:     "Task execution failed",
		})
	}
}

// loadDefaultRules loads predefined behavior rules
func (c *Coordinator) loadDefaultRules() error {
	// Error handling rule
	errorRule := rules.Rule{
		ID:          "handle_errors",
		Name:        "Error Handler",
		Description: "Respond to error events",
		Priority:    100,
		Enabled:     true,
		Condition: &rules.EventTypeCondition{
			EventType: "error",
		},
		Actions: []rules.Action{
			&rules.LogAction{
				Message: "Error detected, initiating recovery",
			},
		},
		Tags: []string{"error", "recovery"},
	}
	
	if err := c.ruleEngine.AddRule(errorRule); err != nil {
		return err
	}
	
	// Log analysis rule
	logRule := rules.Rule{
		ID:          "analyze_logs",
		Name:        "Log Analyzer",
		Description: "Analyze log entries",
		Priority:    50,
		Enabled:     true,
		Condition:   &rules.AlwaysCondition{},
		Actions: []rules.Action{
			&rules.LogAction{
				Message: "Processing log entry",
			},
		},
		Tags: []string{"log", "analysis"},
	}
	
	return c.ruleEngine.AddRule(logRule)
}

// GetRegistry returns the agent registry
func (c *Coordinator) GetRegistry() *agent.Registry {
	return c.registry
}

// GetMemoryStore returns the memory store
func (c *Coordinator) GetMemoryStore() memory.MemoryStore {
	return c.memoryStore
}

// GetVotingSystem returns the voting system
func (c *Coordinator) GetVotingSystem() *voting.DemocraticVotingSystem {
	return c.votingSystem
}

// GetRuleEngine returns the rule engine
func (c *Coordinator) GetRuleEngine() *rules.RuleEngine {
	return c.ruleEngine
}

// GetHealthMonitor returns the health monitor
func (c *Coordinator) GetHealthMonitor() *health.HealthMonitor {
	return c.healthMonitor
}

// GetSystemStatus returns overall system status
func (c *Coordinator) GetSystemStatus() SystemStatus {
	return SystemStatus{
		Running:       c.running,
		AgentHealth:   c.registry.GetHealthStatus(),
		SystemHealth:  c.healthMonitor.GetSystemHealth(),
		MemoryStats:   c.memoryStore.GetStats(),
		ActiveSessions: len(c.votingSystem.GetActiveSessions()),
		QueuedTasks:   len(c.taskQueue),
	}
}

// SystemStatus represents the overall system status
type SystemStatus struct {
	Running        bool
	AgentHealth    map[string]agent.AgentHealth
	SystemHealth   health.SystemHealth
	MemoryStats    memory.MemoryStats
	ActiveSessions int
	QueuedTasks    int
}
