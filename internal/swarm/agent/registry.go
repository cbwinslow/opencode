package agent

import (
	"context"
	"fmt"
	"sync"
)

// Registry manages all agents in the swarm
type Registry struct {
	agents      map[string]Agent
	agentsByType map[AgentType][]Agent
	mu          sync.RWMutex
	
	// Message routing
	messageBroker *MessageBroker
}

// NewRegistry creates a new agent registry
func NewRegistry() *Registry {
	return &Registry{
		agents:        make(map[string]Agent),
		agentsByType:  make(map[AgentType][]Agent),
		messageBroker: NewMessageBroker(),
	}
}

// RegisterAgent adds an agent to the registry
func (r *Registry) RegisterAgent(agent Agent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	id := agent.GetID()
	if _, exists := r.agents[id]; exists {
		return fmt.Errorf("agent with ID %s already registered", id)
	}
	
	r.agents[id] = agent
	
	agentType := agent.GetType()
	r.agentsByType[agentType] = append(r.agentsByType[agentType], agent)
	
	// Subscribe agent to message broker
	r.messageBroker.Subscribe(id, agent.ReceiveMessages())
	
	return nil
}

// UnregisterAgent removes an agent from the registry
func (r *Registry) UnregisterAgent(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	agent, exists := r.agents[id]
	if !exists {
		return fmt.Errorf("agent with ID %s not found", id)
	}
	
	// Remove from type map
	agentType := agent.GetType()
	agents := r.agentsByType[agentType]
	for i, a := range agents {
		if a.GetID() == id {
			r.agentsByType[agentType] = append(agents[:i], agents[i+1:]...)
			break
		}
	}
	
	delete(r.agents, id)
	r.messageBroker.Unsubscribe(id)
	
	return nil
}

// GetAgent retrieves an agent by ID
func (r *Registry) GetAgent(id string) (Agent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	agent, exists := r.agents[id]
	if !exists {
		return nil, fmt.Errorf("agent with ID %s not found", id)
	}
	
	return agent, nil
}

// GetAgentsByType retrieves all agents of a specific type
func (r *Registry) GetAgentsByType(agentType AgentType) []Agent {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	// Return a copy to avoid concurrent modification issues
	agents := make([]Agent, len(r.agentsByType[agentType]))
	copy(agents, r.agentsByType[agentType])
	
	return agents
}

// GetAllAgents returns all registered agents
func (r *Registry) GetAllAgents() []Agent {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	agents := make([]Agent, 0, len(r.agents))
	for _, agent := range r.agents {
		agents = append(agents, agent)
	}
	
	return agents
}

// FindAgentsForTask finds suitable agents for a task
func (r *Registry) FindAgentsForTask(task Task) []Agent {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var suitable []Agent
	for _, agent := range r.agents {
		if agent.GetStatus() == AgentStatusIdle && agent.CanHandleTask(task) {
			suitable = append(suitable, agent)
		}
	}
	
	return suitable
}

// BroadcastMessage sends a message to all agents
func (r *Registry) BroadcastMessage(msg Message) error {
	return r.messageBroker.Broadcast(msg)
}

// SendMessage sends a message to a specific agent
func (r *Registry) SendMessage(toID string, msg Message) error {
	msg.To = toID
	return r.messageBroker.Send(msg)
}

// StartAll starts all registered agents
func (r *Registry) StartAll(ctx context.Context) error {
	r.mu.RLock()
	agents := make([]Agent, 0, len(r.agents))
	for _, agent := range r.agents {
		agents = append(agents, agent)
	}
	r.mu.RUnlock()
	
	for _, agent := range agents {
		if err := agent.Start(ctx); err != nil {
			return fmt.Errorf("failed to start agent %s: %w", agent.GetID(), err)
		}
	}
	
	return nil
}

// StopAll stops all registered agents
func (r *Registry) StopAll() error {
	r.mu.RLock()
	agents := make([]Agent, 0, len(r.agents))
	for _, agent := range r.agents {
		agents = append(agents, agent)
	}
	r.mu.RUnlock()
	
	var lastErr error
	for _, agent := range agents {
		if err := agent.Stop(); err != nil {
			lastErr = err
		}
	}
	
	return lastErr
}

// GetHealthStatus returns health information for all agents
func (r *Registry) GetHealthStatus() map[string]AgentHealth {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	status := make(map[string]AgentHealth)
	for id, agent := range r.agents {
		status[id] = AgentHealth{
			ID:          id,
			Type:        agent.GetType(),
			Status:      agent.GetStatus(),
			HealthScore: agent.GetHealthScore(),
			Metrics:     agent.GetMetrics(),
		}
	}
	
	return status
}

// AgentHealth represents the health status of an agent
type AgentHealth struct {
	ID          string
	Type        AgentType
	Status      AgentStatus
	HealthScore float64
	Metrics     AgentMetrics
}

// MessageBroker handles message routing between agents
type MessageBroker struct {
	subscribers map[string]<-chan Message
	mu          sync.RWMutex
}

// NewMessageBroker creates a new message broker
func NewMessageBroker() *MessageBroker {
	return &MessageBroker{
		subscribers: make(map[string]<-chan Message),
	}
}

// Subscribe registers an agent's message channel
func (mb *MessageBroker) Subscribe(agentID string, msgChan <-chan Message) {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	mb.subscribers[agentID] = msgChan
}

// Unsubscribe removes an agent's message channel
func (mb *MessageBroker) Unsubscribe(agentID string) {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	delete(mb.subscribers, agentID)
}

// Send routes a message to a specific agent
func (mb *MessageBroker) Send(msg Message) error {
	mb.mu.RLock()
	defer mb.mu.RUnlock()
	
	if msg.To == "" {
		return fmt.Errorf("message must have a recipient")
	}
	
	// In a real implementation, this would route to the agent's input channel
	// For now, this is a placeholder
	return nil
}

// Broadcast sends a message to all subscribed agents
func (mb *MessageBroker) Broadcast(msg Message) error {
	mb.mu.RLock()
	defer mb.mu.RUnlock()
	
	// In a real implementation, this would send to all agents
	// For now, this is a placeholder
	return nil
}
