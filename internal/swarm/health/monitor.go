package health

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// HealthStatus represents the health state of a component
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusDegraded  HealthStatus = "degraded"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusCritical  HealthStatus = "critical"
)

// HealthCheck represents a health check result
type HealthCheck struct {
	ComponentID   string
	Status        HealthStatus
	Score         float64 // 0.0 to 1.0
	Message       string
	Details       map[string]interface{}
	Timestamp     time.Time
	ResponseTime  time.Duration
}

// HealthMonitor monitors system health and triggers recovery
type HealthMonitor struct {
	checks        map[string]*HealthCheck
	mu            sync.RWMutex
	checkInterval time.Duration
	alertThreshold float64
	
	// Recovery strategies
	recoveryStrategies map[string]RecoveryStrategy
	
	// Event channels
	alertChan   chan HealthAlert
	recoveryChan chan RecoveryAction
	
	ctx        context.Context
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
}

// HealthAlert represents a health alert
type HealthAlert struct {
	ComponentID string
	Status      HealthStatus
	Check       HealthCheck
	Severity    AlertSeverity
	Timestamp   time.Time
}

// AlertSeverity defines alert importance
type AlertSeverity string

const (
	AlertSeverityInfo     AlertSeverity = "info"
	AlertSeverityWarning  AlertSeverity = "warning"
	AlertSeverityError    AlertSeverity = "error"
	AlertSeverityCritical AlertSeverity = "critical"
)

// RecoveryAction represents an action to recover from an issue
type RecoveryAction struct {
	ComponentID string
	ActionType  RecoveryActionType
	Parameters  map[string]interface{}
	Timestamp   time.Time
}

// RecoveryActionType defines types of recovery actions
type RecoveryActionType string

const (
	RecoveryActionRestart   RecoveryActionType = "restart"
	RecoveryActionReset     RecoveryActionType = "reset"
	RecoveryActionReload    RecoveryActionType = "reload"
	RecoveryActionScale     RecoveryActionType = "scale"
	RecoveryActionFallback  RecoveryActionType = "fallback"
	RecoveryActionIsolate   RecoveryActionType = "isolate"
)

// RecoveryStrategy defines how to recover from failures
type RecoveryStrategy interface {
	CanRecover(check HealthCheck) bool
	Recover(ctx context.Context, check HealthCheck) error
	GetPriority() int
}

// HealthMonitorConfig configures the health monitor
type HealthMonitorConfig struct {
	CheckInterval  time.Duration
	AlertThreshold float64
	AlertBuffer    int
	RecoveryBuffer int
}

// NewHealthMonitor creates a new health monitor
func NewHealthMonitor(config HealthMonitorConfig) *HealthMonitor {
	if config.CheckInterval <= 0 {
		config.CheckInterval = 30 * time.Second
	}
	if config.AlertThreshold <= 0 {
		config.AlertThreshold = 0.5
	}
	if config.AlertBuffer <= 0 {
		config.AlertBuffer = 100
	}
	if config.RecoveryBuffer <= 0 {
		config.RecoveryBuffer = 100
	}
	
	ctx, cancel := context.WithCancel(context.Background())
	
	return &HealthMonitor{
		checks:             make(map[string]*HealthCheck),
		checkInterval:      config.CheckInterval,
		alertThreshold:     config.AlertThreshold,
		recoveryStrategies: make(map[string]RecoveryStrategy),
		alertChan:          make(chan HealthAlert, config.AlertBuffer),
		recoveryChan:       make(chan RecoveryAction, config.RecoveryBuffer),
		ctx:                ctx,
		cancelFunc:         cancel,
	}
}

// Start begins health monitoring
func (hm *HealthMonitor) Start() error {
	hm.wg.Add(2)
	go hm.monitorLoop()
	go hm.recoveryLoop()
	return nil
}

// Stop stops health monitoring
func (hm *HealthMonitor) Stop() error {
	hm.cancelFunc()
	hm.wg.Wait()
	close(hm.alertChan)
	close(hm.recoveryChan)
	return nil
}

// RegisterCheck adds a component to monitor
func (hm *HealthMonitor) RegisterCheck(componentID string) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	
	hm.checks[componentID] = &HealthCheck{
		ComponentID: componentID,
		Status:      HealthStatusHealthy,
		Score:       1.0,
		Timestamp:   time.Now(),
	}
}

// UpdateCheck updates a health check result
func (hm *HealthMonitor) UpdateCheck(check HealthCheck) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	
	check.Timestamp = time.Now()
	hm.checks[check.ComponentID] = &check
	
	// Trigger alert if unhealthy
	if check.Score < hm.alertThreshold {
		hm.triggerAlert(check)
	}
}

// GetCheck retrieves the latest health check for a component
func (hm *HealthMonitor) GetCheck(componentID string) (*HealthCheck, error) {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	
	check, exists := hm.checks[componentID]
	if !exists {
		return nil, fmt.Errorf("component not found: %s", componentID)
	}
	
	return check, nil
}

// GetAllChecks returns all health checks
func (hm *HealthMonitor) GetAllChecks() map[string]*HealthCheck {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	
	checks := make(map[string]*HealthCheck)
	for id, check := range hm.checks {
		checkCopy := *check
		checks[id] = &checkCopy
	}
	
	return checks
}

// RegisterRecoveryStrategy adds a recovery strategy
func (hm *HealthMonitor) RegisterRecoveryStrategy(componentID string, strategy RecoveryStrategy) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	hm.recoveryStrategies[componentID] = strategy
}

// Alerts returns the alert channel
func (hm *HealthMonitor) Alerts() <-chan HealthAlert {
	return hm.alertChan
}

// RecoveryActions returns the recovery action channel
func (hm *HealthMonitor) RecoveryActions() <-chan RecoveryAction {
	return hm.recoveryChan
}

// monitorLoop periodically checks health
func (hm *HealthMonitor) monitorLoop() {
	defer hm.wg.Done()
	
	ticker := time.NewTicker(hm.checkInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			hm.performHealthChecks()
		case <-hm.ctx.Done():
			return
		}
	}
}

// performHealthChecks checks all registered components
func (hm *HealthMonitor) performHealthChecks() {
	hm.mu.RLock()
	checks := make([]*HealthCheck, 0, len(hm.checks))
	for _, check := range hm.checks {
		checks = append(checks, check)
	}
	hm.mu.RUnlock()
	
	for _, check := range checks {
		// Check if stale (no updates in 2x interval)
		if time.Since(check.Timestamp) > 2*hm.checkInterval {
			check.Status = HealthStatusUnhealthy
			check.Score = 0.3
			check.Message = "Component not responding"
			hm.UpdateCheck(*check)
		}
	}
}

// recoveryLoop handles recovery actions
func (hm *HealthMonitor) recoveryLoop() {
	defer hm.wg.Done()
	
	for {
		select {
		case alert := <-hm.alertChan:
			hm.handleAlert(alert)
		case <-hm.ctx.Done():
			return
		}
	}
}

// triggerAlert creates and sends an alert
func (hm *HealthMonitor) triggerAlert(check HealthCheck) {
	severity := hm.calculateSeverity(check)
	
	alert := HealthAlert{
		ComponentID: check.ComponentID,
		Status:      check.Status,
		Check:       check,
		Severity:    severity,
		Timestamp:   time.Now(),
	}
	
	select {
	case hm.alertChan <- alert:
	default:
		// Alert buffer full, skip
	}
}

// handleAlert processes an alert and initiates recovery
func (hm *HealthMonitor) handleAlert(alert HealthAlert) {
	hm.mu.RLock()
	strategy, hasStrategy := hm.recoveryStrategies[alert.ComponentID]
	hm.mu.RUnlock()
	
	if !hasStrategy {
		return
	}
	
	if strategy.CanRecover(alert.Check) {
		// Attempt recovery
		ctx, cancel := context.WithTimeout(hm.ctx, 30*time.Second)
		defer cancel()
		
		if err := strategy.Recover(ctx, alert.Check); err != nil {
			// Recovery failed, escalate
		} else {
			// Recovery successful
			action := RecoveryAction{
				ComponentID: alert.ComponentID,
				ActionType:  RecoveryActionRestart,
				Timestamp:   time.Now(),
			}
			
			select {
			case hm.recoveryChan <- action:
			default:
			}
		}
	}
}

// calculateSeverity determines alert severity based on health check
func (hm *HealthMonitor) calculateSeverity(check HealthCheck) AlertSeverity {
	if check.Score < 0.2 {
		return AlertSeverityCritical
	} else if check.Score < 0.4 {
		return AlertSeverityError
	} else if check.Score < 0.6 {
		return AlertSeverityWarning
	}
	return AlertSeverityInfo
}

// GetSystemHealth returns overall system health
func (hm *HealthMonitor) GetSystemHealth() SystemHealth {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	
	var totalScore float64
	statusCounts := make(map[HealthStatus]int)
	
	for _, check := range hm.checks {
		totalScore += check.Score
		statusCounts[check.Status]++
	}
	
	avgScore := 0.0
	if len(hm.checks) > 0 {
		avgScore = totalScore / float64(len(hm.checks))
	}
	
	overallStatus := HealthStatusHealthy
	if avgScore < 0.3 {
		overallStatus = HealthStatusCritical
	} else if avgScore < 0.5 {
		overallStatus = HealthStatusUnhealthy
	} else if avgScore < 0.8 {
		overallStatus = HealthStatusDegraded
	}
	
	return SystemHealth{
		OverallStatus:    overallStatus,
		OverallScore:     avgScore,
		ComponentCount:   len(hm.checks),
		HealthyCount:     statusCounts[HealthStatusHealthy],
		DegradedCount:    statusCounts[HealthStatusDegraded],
		UnhealthyCount:   statusCounts[HealthStatusUnhealthy],
		CriticalCount:    statusCounts[HealthStatusCritical],
		LastUpdated:      time.Now(),
	}
}

// SystemHealth represents overall system health
type SystemHealth struct {
	OverallStatus   HealthStatus
	OverallScore    float64
	ComponentCount  int
	HealthyCount    int
	DegradedCount   int
	UnhealthyCount  int
	CriticalCount   int
	LastUpdated     time.Time
}
