package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// BaseAgent provides common functionality for all agent implementations
type BaseAgent struct {
	id           string
	agentType    AgentType
	status       AgentStatus
	capabilities []string
	config       AgentConfig
	
	// Communication
	incomingMessages chan Message
	outgoingMessages chan Message
	
	// Metrics and health
	metrics      AgentMetrics
	metricsMutex sync.RWMutex
	healthScore  float64
	startTime    time.Time
	
	// Lifecycle
	ctx        context.Context
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
	statusMutex sync.RWMutex
}

// NewBaseAgent creates a new base agent
func NewBaseAgent(config AgentConfig) *BaseAgent {
	if config.ID == "" {
		config.ID = uuid.New().String()
	}
	
	if config.MessageBufferSize <= 0 {
		config.MessageBufferSize = 100
	}
	
	return &BaseAgent{
		id:               config.ID,
		agentType:        config.Type,
		status:           AgentStatusStopped,
		capabilities:     config.Capabilities,
		config:           config,
		incomingMessages: make(chan Message, config.MessageBufferSize),
		outgoingMessages: make(chan Message, config.MessageBufferSize),
		healthScore:      1.0,
		metrics: AgentMetrics{
			TasksCompleted:   0,
			TasksFailed:      0,
			AverageTaskTime:  0,
			LastActivityTime: time.Now(),
			MessagesReceived: 0,
			MessagesSent:     0,
			ErrorCount:       0,
			UptimeSeconds:    0,
		},
	}
}

// Start begins the agent's lifecycle
func (a *BaseAgent) Start(ctx context.Context) error {
	a.statusMutex.Lock()
	defer a.statusMutex.Unlock()
	
	if a.status != AgentStatusStopped {
		return fmt.Errorf("agent %s is already running", a.id)
	}
	
	a.status = AgentStatusStarting
	a.ctx, a.cancelFunc = context.WithCancel(ctx)
	a.startTime = time.Now()
	
	// Start message processing
	a.wg.Add(1)
	go a.processMessages()
	
	// Start health monitoring if configured
	if a.config.HealthCheckInterval > 0 {
		a.wg.Add(1)
		go a.monitorHealth()
	}
	
	a.status = AgentStatusIdle
	return nil
}

// Stop terminates the agent
func (a *BaseAgent) Stop() error {
	a.statusMutex.Lock()
	status := a.status
	a.status = AgentStatusStopped
	a.statusMutex.Unlock()
	
	if status == AgentStatusStopped {
		return nil
	}
	
	if a.cancelFunc != nil {
		a.cancelFunc()
	}
	
	// Close message channels
	close(a.incomingMessages)
	
	// Wait for goroutines to finish
	a.wg.Wait()
	
	return nil
}

// GetStatus returns the current agent status
func (a *BaseAgent) GetStatus() AgentStatus {
	a.statusMutex.RLock()
	defer a.statusMutex.RUnlock()
	return a.status
}

// SetStatus updates the agent status
func (a *BaseAgent) SetStatus(status AgentStatus) {
	a.statusMutex.Lock()
	defer a.statusMutex.Unlock()
	a.status = status
}

// GetID returns the agent's unique identifier
func (a *BaseAgent) GetID() string {
	return a.id
}

// GetType returns the agent's type
func (a *BaseAgent) GetType() AgentType {
	return a.agentType
}

// GetCapabilities returns the agent's capabilities
func (a *BaseAgent) GetCapabilities() []string {
	return a.capabilities
}

// SendMessage sends a message from this agent
func (a *BaseAgent) SendMessage(msg Message) error {
	if msg.From == "" {
		msg.From = a.id
	}
	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}
	
	select {
	case a.outgoingMessages <- msg:
		a.incrementMessagesSent()
		return nil
	case <-a.ctx.Done():
		return fmt.Errorf("agent context cancelled")
	default:
		return fmt.Errorf("outgoing message buffer full")
	}
}

// ReceiveMessages returns the channel for incoming messages
func (a *BaseAgent) ReceiveMessages() <-chan Message {
	return a.incomingMessages
}

// GetHealthScore returns the agent's health score (0.0 to 1.0)
func (a *BaseAgent) GetHealthScore() float64 {
	a.metricsMutex.RLock()
	defer a.metricsMutex.RUnlock()
	return a.healthScore
}

// GetMetrics returns the agent's metrics
func (a *BaseAgent) GetMetrics() AgentMetrics {
	a.metricsMutex.RLock()
	defer a.metricsMutex.RUnlock()
	
	// Update uptime
	metrics := a.metrics
	metrics.UptimeSeconds = int64(time.Since(a.startTime).Seconds())
	
	return metrics
}

// processMessages handles incoming messages
func (a *BaseAgent) processMessages() {
	defer a.wg.Done()
	
	for {
		select {
		case msg, ok := <-a.incomingMessages:
			if !ok {
				return
			}
			a.incrementMessagesReceived()
			a.handleMessage(msg)
			
		case <-a.ctx.Done():
			return
		}
	}
}

// handleMessage processes a single message
func (a *BaseAgent) handleMessage(msg Message) {
	// Default implementation - can be overridden by specialized agents
	switch msg.Type {
	case MessageTypeHealthCheck:
		a.respondToHealthCheck(msg)
	default:
		// Specialized agents should override this
	}
}

// respondToHealthCheck sends health information
func (a *BaseAgent) respondToHealthCheck(msg Message) {
	response := Message{
		ID:        uuid.New().String(),
		From:      a.id,
		To:        msg.From,
		Type:      MessageTypeStatusUpdate,
		Content:   a.GetMetrics(),
		Timestamp: time.Now(),
		ReplyTo:   msg.ID,
	}
	
	_ = a.SendMessage(response)
}

// monitorHealth periodically checks agent health
func (a *BaseAgent) monitorHealth() {
	defer a.wg.Done()
	
	ticker := time.NewTicker(a.config.HealthCheckInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			a.updateHealthScore()
		case <-a.ctx.Done():
			return
		}
	}
}

// updateHealthScore calculates and updates the health score
func (a *BaseAgent) updateHealthScore() {
	a.metricsMutex.Lock()
	defer a.metricsMutex.Unlock()
	
	// Calculate health based on various factors
	errorRate := 0.0
	totalTasks := a.metrics.TasksCompleted + a.metrics.TasksFailed
	if totalTasks > 0 {
		errorRate = float64(a.metrics.TasksFailed) / float64(totalTasks)
	}
	
	// Health score: 1.0 = perfect, 0.0 = critical
	a.healthScore = 1.0 - errorRate
	
	// Additional factors could include:
	// - Response time
	// - Resource usage
	// - Recent activity
}

// Metric update helpers
func (a *BaseAgent) incrementTasksCompleted() {
	a.metricsMutex.Lock()
	defer a.metricsMutex.Unlock()
	a.metrics.TasksCompleted++
	a.metrics.LastActivityTime = time.Now()
}

func (a *BaseAgent) incrementTasksFailed() {
	a.metricsMutex.Lock()
	defer a.metricsMutex.Unlock()
	a.metrics.TasksFailed++
	a.metrics.LastActivityTime = time.Now()
}

func (a *BaseAgent) incrementMessagesReceived() {
	a.metricsMutex.Lock()
	defer a.metricsMutex.Unlock()
	a.metrics.MessagesReceived++
}

func (a *BaseAgent) incrementMessagesSent() {
	a.metricsMutex.Lock()
	defer a.metricsMutex.Unlock()
	a.metrics.MessagesSent++
}

func (a *BaseAgent) incrementErrorCount() {
	a.metricsMutex.Lock()
	defer a.metricsMutex.Unlock()
	a.metrics.ErrorCount++
}

func (a *BaseAgent) updateAverageTaskTime(duration time.Duration) {
	a.metricsMutex.Lock()
	defer a.metricsMutex.Unlock()
	
	// Calculate rolling average
	if a.metrics.AverageTaskTime == 0 {
		a.metrics.AverageTaskTime = duration
	} else {
		a.metrics.AverageTaskTime = (a.metrics.AverageTaskTime + duration) / 2
	}
}
