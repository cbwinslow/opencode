package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/opencode-ai/opencode/internal/swarm"
	"github.com/opencode-ai/opencode/internal/swarm/agent"
	"github.com/opencode-ai/opencode/internal/swarm/health"
	"github.com/opencode-ai/opencode/internal/swarm/memory"
	"github.com/opencode-ai/opencode/internal/swarm/rules"
	"github.com/opencode-ai/opencode/internal/swarm/voting"
)

func main() {
	fmt.Println("OpenCode Multi-Agent Swarm Example")
	fmt.Println("===================================\n")

	// Example 1: Basic Coordinator Setup
	basicExample()

	// Example 2: Memory System
	memoryExample()

	// Example 3: Democratic Voting
	votingExample()

	// Example 4: Health Monitoring
	healthExample()

	// Example 5: Rule Engine
	ruleEngineExample()

	// Example 6: Complete Workflow
	completeWorkflowExample()
}

// Example 1: Basic Coordinator Setup
func basicExample() {
	fmt.Println("=== Example 1: Basic Coordinator Setup ===\n")

	// Create coordinator configuration
	config := swarm.CoordinatorConfig{
		SwarmConfig: agent.SwarmConfig{
			Name:               "example-swarm",
			VotingThreshold:    0.66,
			MaxConcurrentTasks: 10,
			EnableMemory:       true,
			EnableLearning:     true,
			EnableSelfHealing:  true,
		},
		MemoryConfig: memory.HierarchicalMemoryConfig{
			MaxMemories:           10000,
			ConsolidationInterval: 1 * time.Hour,
			PruneOlderThan:        30 * 24 * time.Hour,
		},
		HealthConfig: health.HealthMonitorConfig{
			CheckInterval:  30 * time.Second,
			AlertThreshold: 0.5,
		},
		LogPaths:      []string{"/var/log/opencode/*.log"},
		ShellHistory:  os.Getenv("HOME") + "/.bash_history",
		TaskQueueSize: 1000,
	}

	// Create coordinator
	coordinator, err := swarm.NewCoordinator(config)
	if err != nil {
		log.Fatalf("Failed to create coordinator: %v", err)
	}

	// Start coordinator
	if err := coordinator.Start(); err != nil {
		log.Fatalf("Failed to start coordinator: %v", err)
	}
	defer coordinator.Stop()

	fmt.Println("✓ Coordinator started successfully")
	fmt.Println("✓ Components initialized: Registry, Memory, Voting, Rules, Health")
	fmt.Println()

	// Get system status
	status := coordinator.GetSystemStatus()
	fmt.Printf("System Status:\n")
	fmt.Printf("  Running: %v\n", status.Running)
	fmt.Printf("  Queued Tasks: %d\n", status.QueuedTasks)
	fmt.Printf("  System Health: %s (%.2f)\n", 
		status.SystemHealth.OverallStatus, 
		status.SystemHealth.OverallScore)
	fmt.Println()
}

// Example 2: Memory System Usage
func memoryExample() {
	fmt.Println("=== Example 2: Memory System ===\n")

	// Create memory store
	memStore := memory.NewHierarchicalMemoryStore(memory.HierarchicalMemoryConfig{
		MaxMemories:           1000,
		ConsolidationInterval: 1 * time.Hour,
		PruneOlderThan:        7 * 24 * time.Hour,
	})

	// Store different types of memories

	// 1. Working memory (current task)
	workingMem := memory.Memory{
		Type:     memory.MemoryTypeWorking,
		Content:  "Currently analyzing code quality issues",
		Tags:     []string{"task", "active", "code-analysis"},
		Priority: memory.PriorityHigh,
	}
	memStore.Store(workingMem)
	fmt.Println("✓ Stored working memory")

	// 2. Episodic memory (event)
	episodicMem := memory.Memory{
		Type: memory.MemoryTypeEpisodic,
		Content: map[string]interface{}{
			"event":       "error_detected",
			"severity":    "high",
			"component":   "parser",
			"description": "Syntax error in file main.go:42",
		},
		Tags:     []string{"error", "parser", "syntax"},
		Priority: memory.PriorityHigh,
		Metadata: map[string]interface{}{
			"file": "main.go",
			"line": 42,
		},
	}
	memStore.Store(episodicMem)
	fmt.Println("✓ Stored episodic memory (error event)")

	// 3. Semantic memory (knowledge)
	semanticMem := memory.Memory{
		Type:     memory.MemoryTypeSemantic,
		Content:  "Go syntax requires semicolons or newlines to separate statements",
		Tags:     []string{"knowledge", "go", "syntax"},
		Priority: memory.PriorityNormal,
	}
	memStore.Store(semanticMem)
	fmt.Println("✓ Stored semantic memory (knowledge)")

	// 4. Procedural memory (how-to)
	proceduralMem := memory.Memory{
		Type: memory.MemoryTypeProcedural,
		Content: map[string]interface{}{
			"skill":       "fix_syntax_error",
			"steps": []string{
				"1. Identify the error location",
				"2. Check for missing semicolons or braces",
				"3. Verify proper statement termination",
				"4. Run syntax checker",
			},
		},
		Tags:     []string{"procedure", "debugging", "syntax"},
		Priority: memory.PriorityNormal,
	}
	memStore.Store(proceduralMem)
	fmt.Println("✓ Stored procedural memory (debugging steps)")

	// Query memories
	fmt.Println("\nQuerying memories:")

	query := memory.MemoryQuery{
		Tags:  []string{"error"},
		Limit: 10,
	}
	results, _ := memStore.Query(query)
	fmt.Printf("✓ Found %d memories tagged with 'error'\n", len(results))

	// Get statistics
	stats := memStore.GetStats()
	fmt.Printf("\nMemory Statistics:\n")
	fmt.Printf("  Total Memories: %d\n", stats.TotalMemories)
	for memType, count := range stats.MemoriesByType {
		fmt.Printf("  %s: %d\n", memType, count)
	}
	fmt.Println()
}

// Example 3: Democratic Voting
func votingExample() {
	fmt.Println("=== Example 3: Democratic Voting ===\n")

	votingSystem := voting.NewDemocraticVotingSystem()

	// Create a proposal
	proposal := voting.VoteProposal{
		Description: "Should we refactor the authentication module?",
		Options:     []string{"yes", "no", "defer"},
		Context: map[string]interface{}{
			"module":       "auth",
			"complexity":   "high",
			"risk":         "medium",
			"estimated_time": "2 days",
		},
		Deadline: time.Now().Add(1 * time.Minute),
	}

	// Create vote session
	session, err := votingSystem.CreateVoteSession(
		proposal,
		voting.VoteTypeMajority,
		3, // minimum 3 voters
		nil,
	)
	if err != nil {
		log.Printf("Failed to create vote session: %v", err)
		return
	}

	fmt.Printf("Created vote session: %s\n", session.ID)
	fmt.Printf("Proposal: %s\n", proposal.Description)
	fmt.Println()

	// Simulate agents voting
	votes := []voting.Vote{
		{
			AgentID:    "monitor-agent-1",
			Decision:   true,
			Confidence: 0.85,
			Reasoning:  "Code quality metrics show high technical debt in auth module",
		},
		{
			AgentID:    "analyzer-agent-1",
			Decision:   true,
			Confidence: 0.90,
			Reasoning:  "Security analysis indicates potential vulnerabilities",
		},
		{
			AgentID:    "executor-agent-1",
			Decision:   false,
			Confidence: 0.70,
			Reasoning:  "Current sprint is already overloaded, defer to next sprint",
		},
	}

	fmt.Println("Agents casting votes:")
	for _, vote := range votes {
		err := votingSystem.CastVote(session.ID, vote)
		if err != nil {
			log.Printf("Failed to cast vote: %v", err)
			continue
		}
		decision := "No"
		if vote.Decision {
			decision = "Yes"
		}
		fmt.Printf("  %s: %s (confidence: %.2f)\n", vote.AgentID, decision, vote.Confidence)
		fmt.Printf("    Reasoning: %s\n", vote.Reasoning)
	}
	fmt.Println()

	// Wait for result
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := votingSystem.WaitForResult(ctx, session.ID)
	if err != nil {
		log.Printf("Failed to get result: %v", err)
		return
	}

	// Display result
	fmt.Println("Vote Result:")
	fmt.Printf("  Decision: ")
	if result.Decision {
		fmt.Println("APPROVED")
	} else {
		fmt.Println("REJECTED")
	}
	fmt.Printf("  Yes: %d, No: %d\n", result.YesVotes, result.NoVotes)
	fmt.Printf("  Approval: %.1f%%\n", result.YesPercentage*100)
	fmt.Printf("  Average Confidence: %.2f\n", result.Confidence)
	fmt.Println()
}

// Example 4: Health Monitoring
func healthExample() {
	fmt.Println("=== Example 4: Health Monitoring ===\n")

	healthMonitor := health.NewHealthMonitor(health.HealthMonitorConfig{
		CheckInterval:  10 * time.Second,
		AlertThreshold: 0.5,
	})

	if err := healthMonitor.Start(); err != nil {
		log.Printf("Failed to start health monitor: %v", err)
		return
	}
	defer healthMonitor.Stop()

	fmt.Println("✓ Health monitor started")

	// Register components
	components := []string{
		"parser-agent",
		"executor-agent",
		"memory-system",
	}

	for _, comp := range components {
		healthMonitor.RegisterCheck(comp)
		fmt.Printf("✓ Registered component: %s\n", comp)
	}
	fmt.Println()

	// Simulate health updates
	fmt.Println("Simulating health checks:")

	// Healthy component
	healthMonitor.UpdateCheck(health.HealthCheck{
		ComponentID: "parser-agent",
		Status:      health.HealthStatusHealthy,
		Score:       0.95,
		Message:     "Operating normally",
		Details: map[string]interface{}{
			"requests_processed": 1000,
			"errors":            5,
			"avg_response_time": "50ms",
		},
	})
	fmt.Println("  parser-agent: Healthy (0.95)")

	// Degraded component
	healthMonitor.UpdateCheck(health.HealthCheck{
		ComponentID: "executor-agent",
		Status:      health.HealthStatusDegraded,
		Score:       0.65,
		Message:     "Higher than normal error rate",
		Details: map[string]interface{}{
			"error_rate": 0.15,
			"threshold":  0.10,
		},
	})
	fmt.Println("  executor-agent: Degraded (0.65)")

	// Unhealthy component
	healthMonitor.UpdateCheck(health.HealthCheck{
		ComponentID: "memory-system",
		Status:      health.HealthStatusUnhealthy,
		Score:       0.30,
		Message:     "High memory usage and slow queries",
		Details: map[string]interface{}{
			"memory_usage": "85%",
			"query_time":   "500ms",
		},
	})
	fmt.Println("  memory-system: Unhealthy (0.30)")
	fmt.Println()

	// Get system health
	systemHealth := healthMonitor.GetSystemHealth()
	fmt.Println("Overall System Health:")
	fmt.Printf("  Status: %s\n", systemHealth.OverallStatus)
	fmt.Printf("  Score: %.2f\n", systemHealth.OverallScore)
	fmt.Printf("  Components: %d total\n", systemHealth.ComponentCount)
	fmt.Printf("    Healthy: %d\n", systemHealth.HealthyCount)
	fmt.Printf("    Degraded: %d\n", systemHealth.DegradedCount)
	fmt.Printf("    Unhealthy: %d\n", systemHealth.UnhealthyCount)
	fmt.Println()
}

// Example 5: Rule Engine
func ruleEngineExample() {
	fmt.Println("=== Example 5: Rule Engine ===\n")

	ruleEngine := rules.NewRuleEngine(rules.RuleEngineConfig{
		MaxHistory:    1000,
		EnableHistory: true,
	})

	// Define error handling rule
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
				Message: "Error detected, logging for analysis",
			},
			&rules.CallbackAction{
				Callback: func(ctx context.Context, ruleCtx rules.RuleContext) error {
					fmt.Printf("  → Triggered recovery for: %v\n", ruleCtx.EventData["message"])
					return nil
				},
			},
		},
		Tags: []string{"error", "recovery"},
	}

	if err := ruleEngine.AddRule(errorRule); err != nil {
		log.Printf("Failed to add rule: %v", err)
		return
	}
	fmt.Println("✓ Added error handling rule")

	// Define performance monitoring rule
	perfRule := rules.Rule{
		ID:          "monitor_performance",
		Name:        "Performance Monitor",
		Description: "Monitor performance metrics",
		Priority:    50,
		Enabled:     true,
		Condition: &rules.FieldCondition{
			Field:    "response_time",
			Operator: ">",
			Value:    1000.0, // > 1 second
		},
		Actions: []rules.Action{
			&rules.LogAction{
				Message: "Slow response detected",
			},
		},
		Tags: []string{"performance", "monitoring"},
	}

	if err := ruleEngine.AddRule(perfRule); err != nil {
		log.Printf("Failed to add rule: %v", err)
		return
	}
	fmt.Println("✓ Added performance monitoring rule")
	fmt.Println()

	// Simulate events
	fmt.Println("Simulating events:")

	// Error event
	errorCtx := rules.RuleContext{
		EventType: "error",
		EventData: map[string]interface{}{
			"message": "Failed to parse configuration file",
			"level":   "critical",
		},
		Timestamp: time.Now(),
	}
	fmt.Println("1. Error event:")
	ruleEngine.EvaluateRules(context.Background(), errorCtx)

	// Performance event
	perfCtx := rules.RuleContext{
		EventType: "performance",
		EventData: map[string]interface{}{
			"response_time": 1500.0,
			"endpoint":      "/api/analyze",
		},
		Timestamp: time.Now(),
	}
	fmt.Println("\n2. Performance event:")
	ruleEngine.EvaluateRules(context.Background(), perfCtx)

	// Get execution history
	history := ruleEngine.GetHistory(10)
	fmt.Printf("\nRule Execution History: %d events\n", len(history))
	for i, exec := range history {
		status := "Did not fire"
		if exec.Fired {
			status = "Fired"
			if exec.Success {
				status += " (success)"
			} else {
				status += " (failed)"
			}
		}
		fmt.Printf("  %d. Rule %s: %s\n", i+1, exec.RuleID, status)
	}
	fmt.Println()
}

// Example 6: Complete Workflow
func completeWorkflowExample() {
	fmt.Println("=== Example 6: Complete Workflow ===\n")

	fmt.Println("This example demonstrates a complete workflow:")
	fmt.Println("1. Initialize coordinator with all components")
	fmt.Println("2. Register specialized agents")
	fmt.Println("3. Submit a complex task")
	fmt.Println("4. Agents vote on approach")
	fmt.Println("5. Execute with monitoring")
	fmt.Println("6. Store results in memory")
	fmt.Println("7. Learn from outcome")
	fmt.Println()

	// For brevity, showing the concept
	fmt.Println("Workflow Steps:")
	fmt.Println("  ✓ Coordinator initialized")
	fmt.Println("  ✓ Monitor, Analyzer, and Executor agents registered")
	fmt.Println("  ✓ Task submitted: 'Refactor authentication module'")
	fmt.Println("  ✓ Agents vote: 2 yes, 1 no → Approved (66%)")
	fmt.Println("  ✓ Executor agent begins work")
	fmt.Println("  ✓ Health monitor tracks progress")
	fmt.Println("  ✓ Task completed successfully in 45 seconds")
	fmt.Println("  ✓ Results stored in memory (episodic + procedural)")
	fmt.Println("  ✓ Success pattern learned for future tasks")
	fmt.Println("  ✓ Memory consolidation scheduled")
	fmt.Println()

	fmt.Println("System Benefits:")
	fmt.Println("  • Democratic decision-making ensures consensus")
	fmt.Println("  • Health monitoring enables auto-recovery")
	fmt.Println("  • Memory system preserves knowledge")
	fmt.Println("  • Learning improves future performance")
	fmt.Println("  • Rule engine automates responses")
	fmt.Println()
}
