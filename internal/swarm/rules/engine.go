package rules

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Rule defines a behavior rule for agents
type Rule struct {
	ID          string
	Name        string
	Description string
	Priority    int
	Enabled     bool
	Condition   Condition
	Actions     []Action
	Tags        []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Condition defines when a rule should fire
type Condition interface {
	Evaluate(ctx context.Context, context RuleContext) (bool, error)
	String() string
}

// Action defines what happens when a rule fires
type Action interface {
	Execute(ctx context.Context, context RuleContext) error
	String() string
}

// RuleContext contains data for rule evaluation
type RuleContext struct {
	AgentID    string
	EventType  string
	EventData  map[string]interface{}
	Timestamp  time.Time
	Metadata   map[string]interface{}
}

// RuleEngine manages and executes rules
type RuleEngine struct {
	rules      map[string]*Rule
	mu         sync.RWMutex
	middleware []RuleMiddleware
	
	// Rule execution history
	history    []RuleExecution
	historyMu  sync.RWMutex
	maxHistory int
}

// RuleExecution records rule execution
type RuleExecution struct {
	RuleID      string
	Context     RuleContext
	Fired       bool
	Success     bool
	Error       error
	Duration    time.Duration
	Timestamp   time.Time
}

// RuleMiddleware can intercept rule execution
type RuleMiddleware interface {
	Before(ctx context.Context, rule *Rule, ruleCtx RuleContext) error
	After(ctx context.Context, rule *Rule, ruleCtx RuleContext, err error) error
}

// RuleEngineConfig configures the rule engine
type RuleEngineConfig struct {
	MaxHistory     int
	EnableHistory  bool
	ParallelExec   bool
}

// NewRuleEngine creates a new rule engine
func NewRuleEngine(config RuleEngineConfig) *RuleEngine {
	if config.MaxHistory <= 0 {
		config.MaxHistory = 1000
	}
	
	return &RuleEngine{
		rules:      make(map[string]*Rule),
		middleware: make([]RuleMiddleware, 0),
		history:    make([]RuleExecution, 0),
		maxHistory: config.MaxHistory,
	}
}

// AddRule registers a new rule
func (re *RuleEngine) AddRule(rule Rule) error {
	re.mu.Lock()
	defer re.mu.Unlock()
	
	if rule.ID == "" {
		return fmt.Errorf("rule ID cannot be empty")
	}
	
	if rule.Condition == nil {
		return fmt.Errorf("rule must have a condition")
	}
	
	if len(rule.Actions) == 0 {
		return fmt.Errorf("rule must have at least one action")
	}
	
	rule.UpdatedAt = time.Now()
	if rule.CreatedAt.IsZero() {
		rule.CreatedAt = time.Now()
	}
	
	re.rules[rule.ID] = &rule
	return nil
}

// RemoveRule deletes a rule
func (re *RuleEngine) RemoveRule(ruleID string) error {
	re.mu.Lock()
	defer re.mu.Unlock()
	
	if _, exists := re.rules[ruleID]; !exists {
		return fmt.Errorf("rule not found: %s", ruleID)
	}
	
	delete(re.rules, ruleID)
	return nil
}

// UpdateRule modifies an existing rule
func (re *RuleEngine) UpdateRule(rule Rule) error {
	re.mu.Lock()
	defer re.mu.Unlock()
	
	if _, exists := re.rules[rule.ID]; !exists {
		return fmt.Errorf("rule not found: %s", rule.ID)
	}
	
	rule.UpdatedAt = time.Now()
	re.rules[rule.ID] = &rule
	return nil
}

// GetRule retrieves a rule by ID
func (re *RuleEngine) GetRule(ruleID string) (*Rule, error) {
	re.mu.RLock()
	defer re.mu.RUnlock()
	
	rule, exists := re.rules[ruleID]
	if !exists {
		return nil, fmt.Errorf("rule not found: %s", ruleID)
	}
	
	return rule, nil
}

// GetAllRules returns all rules
func (re *RuleEngine) GetAllRules() []*Rule {
	re.mu.RLock()
	defer re.mu.RUnlock()
	
	rules := make([]*Rule, 0, len(re.rules))
	for _, rule := range re.rules {
		rules = append(rules, rule)
	}
	
	return rules
}

// EvaluateRules evaluates all rules against a context
func (re *RuleEngine) EvaluateRules(ctx context.Context, ruleCtx RuleContext) error {
	re.mu.RLock()
	rules := make([]*Rule, 0, len(re.rules))
	for _, rule := range re.rules {
		if rule.Enabled {
			rules = append(rules, rule)
		}
	}
	re.mu.RUnlock()
	
	// Sort by priority (higher first)
	for i := 0; i < len(rules); i++ {
		for j := i + 1; j < len(rules); j++ {
			if rules[j].Priority > rules[i].Priority {
				rules[i], rules[j] = rules[j], rules[i]
			}
		}
	}
	
	// Evaluate each rule
	for _, rule := range rules {
		if err := re.evaluateRule(ctx, rule, ruleCtx); err != nil {
			// Log error but continue with other rules
			continue
		}
	}
	
	return nil
}

// evaluateRule evaluates a single rule
func (re *RuleEngine) evaluateRule(ctx context.Context, rule *Rule, ruleCtx RuleContext) error {
	startTime := time.Now()
	
	execution := RuleExecution{
		RuleID:    rule.ID,
		Context:   ruleCtx,
		Timestamp: startTime,
	}
	
	// Run middleware before
	for _, mw := range re.middleware {
		if err := mw.Before(ctx, rule, ruleCtx); err != nil {
			execution.Error = err
			re.recordExecution(execution)
			return err
		}
	}
	
	// Evaluate condition
	fired, err := rule.Condition.Evaluate(ctx, ruleCtx)
	if err != nil {
		execution.Error = err
		re.recordExecution(execution)
		return err
	}
	
	execution.Fired = fired
	
	if !fired {
		execution.Duration = time.Since(startTime)
		re.recordExecution(execution)
		return nil
	}
	
	// Execute actions
	for _, action := range rule.Actions {
		if err := action.Execute(ctx, ruleCtx); err != nil {
			execution.Error = err
			execution.Duration = time.Since(startTime)
			re.recordExecution(execution)
			
			// Run middleware after (with error)
			for _, mw := range re.middleware {
				_ = mw.After(ctx, rule, ruleCtx, err)
			}
			
			return err
		}
	}
	
	execution.Success = true
	execution.Duration = time.Since(startTime)
	re.recordExecution(execution)
	
	// Run middleware after (success)
	for _, mw := range re.middleware {
		_ = mw.After(ctx, rule, ruleCtx, nil)
	}
	
	return nil
}

// AddMiddleware adds middleware to the engine
func (re *RuleEngine) AddMiddleware(mw RuleMiddleware) {
	re.mu.Lock()
	defer re.mu.Unlock()
	re.middleware = append(re.middleware, mw)
}

// recordExecution saves rule execution history
func (re *RuleEngine) recordExecution(execution RuleExecution) {
	re.historyMu.Lock()
	defer re.historyMu.Unlock()
	
	re.history = append(re.history, execution)
	
	// Trim history if needed
	if len(re.history) > re.maxHistory {
		re.history = re.history[len(re.history)-re.maxHistory:]
	}
}

// GetHistory returns rule execution history
func (re *RuleEngine) GetHistory(limit int) []RuleExecution {
	re.historyMu.RLock()
	defer re.historyMu.RUnlock()
	
	if limit <= 0 || limit > len(re.history) {
		limit = len(re.history)
	}
	
	history := make([]RuleExecution, limit)
	copy(history, re.history[len(re.history)-limit:])
	
	return history
}

// Common condition implementations

// AlwaysCondition always evaluates to true
type AlwaysCondition struct{}

func (ac *AlwaysCondition) Evaluate(ctx context.Context, context RuleContext) (bool, error) {
	return true, nil
}

func (ac *AlwaysCondition) String() string {
	return "always"
}

// EventTypeCondition matches specific event types
type EventTypeCondition struct {
	EventType string
}

func (etc *EventTypeCondition) Evaluate(ctx context.Context, context RuleContext) (bool, error) {
	return context.EventType == etc.EventType, nil
}

func (etc *EventTypeCondition) String() string {
	return fmt.Sprintf("event_type == %s", etc.EventType)
}

// FieldCondition checks a field value
type FieldCondition struct {
	Field    string
	Operator string // "==", "!=", ">", "<", ">=", "<=", "contains"
	Value    interface{}
}

func (fc *FieldCondition) Evaluate(ctx context.Context, context RuleContext) (bool, error) {
	fieldValue, exists := context.EventData[fc.Field]
	if !exists {
		return false, nil
	}
	
	switch fc.Operator {
	case "==":
		return fieldValue == fc.Value, nil
	case "!=":
		return fieldValue != fc.Value, nil
	// Add more operators as needed
	default:
		return false, fmt.Errorf("unknown operator: %s", fc.Operator)
	}
}

func (fc *FieldCondition) String() string {
	return fmt.Sprintf("%s %s %v", fc.Field, fc.Operator, fc.Value)
}

// LogAction logs a message
type LogAction struct {
	Message string
}

func (la *LogAction) Execute(ctx context.Context, context RuleContext) error {
	// In a real implementation, this would use a proper logger
	fmt.Printf("[Rule Action] %s\n", la.Message)
	return nil
}

func (la *LogAction) String() string {
	return fmt.Sprintf("log: %s", la.Message)
}

// CallbackAction executes a callback function
type CallbackAction struct {
	Callback func(context.Context, RuleContext) error
}

func (ca *CallbackAction) Execute(ctx context.Context, context RuleContext) error {
	return ca.Callback(ctx, context)
}

func (ca *CallbackAction) String() string {
	return "callback"
}
