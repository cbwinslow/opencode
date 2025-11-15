package voting

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// VoteType defines different voting mechanisms
type VoteType string

const (
	VoteTypeMajority  VoteType = "majority"   // Simple majority (>50%)
	VoteTypeSuper     VoteType = "super"      // Super majority (>66%)
	VoteTypeUnanimous VoteType = "unanimous"  // All agree
	VoteTypeWeighted  VoteType = "weighted"   // Weighted by agent expertise
	VoteTypeConsensus VoteType = "consensus"  // Iterative consensus building
)

// Vote represents a single vote
type Vote struct {
	AgentID   string
	Decision  bool   // true for yes, false for no
	Confidence float64 // 0.0 to 1.0
	Reasoning string
	Timestamp time.Time
}

// VoteProposal represents something being voted on
type VoteProposal struct {
	ID          string
	Description string
	ProposedBy  string
	Options     []string
	Context     map[string]interface{}
	CreatedAt   time.Time
	Deadline    time.Time
}

// VoteSession manages a voting process
type VoteSession struct {
	ID          string
	Proposal    VoteProposal
	VoteType    VoteType
	Votes       map[string]Vote // AgentID -> Vote
	mu          sync.RWMutex
	Completed   bool
	Result      *VoteResult
	MinVoters   int
	AgentWeights map[string]float64 // For weighted voting
}

// VoteResult contains the outcome of a vote
type VoteResult struct {
	Decision      bool
	YesVotes      int
	NoVotes       int
	TotalVotes    int
	YesPercentage float64
	Confidence    float64 // Average confidence
	Reasoning     []string
	CompletedAt   time.Time
}

// DemocraticVotingSystem coordinates voting among agents
type DemocraticVotingSystem struct {
	sessions map[string]*VoteSession
	mu       sync.RWMutex
}

// NewDemocraticVotingSystem creates a new voting system
func NewDemocraticVotingSystem() *DemocraticVotingSystem {
	return &DemocraticVotingSystem{
		sessions: make(map[string]*VoteSession),
	}
}

// CreateVoteSession initiates a new vote
func (dvs *DemocraticVotingSystem) CreateVoteSession(
	proposal VoteProposal,
	voteType VoteType,
	minVoters int,
	agentWeights map[string]float64,
) (*VoteSession, error) {
	dvs.mu.Lock()
	defer dvs.mu.Unlock()
	
	if proposal.ID == "" {
		proposal.ID = uuid.New().String()
	}
	
	if proposal.CreatedAt.IsZero() {
		proposal.CreatedAt = time.Now()
	}
	
	session := &VoteSession{
		ID:           uuid.New().String(),
		Proposal:     proposal,
		VoteType:     voteType,
		Votes:        make(map[string]Vote),
		MinVoters:    minVoters,
		AgentWeights: agentWeights,
	}
	
	dvs.sessions[session.ID] = session
	return session, nil
}

// CastVote records a vote in a session
func (dvs *DemocraticVotingSystem) CastVote(sessionID string, vote Vote) error {
	dvs.mu.RLock()
	session, exists := dvs.sessions[sessionID]
	dvs.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("vote session not found: %s", sessionID)
	}
	
	session.mu.Lock()
	defer session.mu.Unlock()
	
	if session.Completed {
		return fmt.Errorf("vote session already completed")
	}
	
	if time.Now().After(session.Proposal.Deadline) {
		return fmt.Errorf("vote deadline passed")
	}
	
	vote.Timestamp = time.Now()
	session.Votes[vote.AgentID] = vote
	
	// Check if we can finalize
	if len(session.Votes) >= session.MinVoters {
		dvs.finalizeVote(session)
	}
	
	return nil
}

// GetVoteResult retrieves the result of a vote session
func (dvs *DemocraticVotingSystem) GetVoteResult(sessionID string) (*VoteResult, error) {
	dvs.mu.RLock()
	session, exists := dvs.sessions[sessionID]
	dvs.mu.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("vote session not found: %s", sessionID)
	}
	
	session.mu.RLock()
	defer session.mu.RUnlock()
	
	if !session.Completed {
		return nil, fmt.Errorf("vote session not completed")
	}
	
	return session.Result, nil
}

// WaitForResult blocks until a vote is completed or times out
func (dvs *DemocraticVotingSystem) WaitForResult(
	ctx context.Context,
	sessionID string,
) (*VoteResult, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			result, err := dvs.GetVoteResult(sessionID)
			if err == nil {
				return result, nil
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// finalizeVote calculates and stores the vote result
func (dvs *DemocraticVotingSystem) finalizeVote(session *VoteSession) {
	var yesCount, noCount int
	var totalConfidence float64
	var reasoning []string
	
	switch session.VoteType {
	case VoteTypeWeighted:
		yesCount, noCount = dvs.calculateWeightedVotes(session)
	default:
		for _, vote := range session.Votes {
			if vote.Decision {
				yesCount++
			} else {
				noCount++
			}
			totalConfidence += vote.Confidence
			if vote.Reasoning != "" {
				reasoning = append(reasoning, vote.Reasoning)
			}
		}
	}
	
	totalVotes := yesCount + noCount
	yesPercentage := 0.0
	if totalVotes > 0 {
		yesPercentage = float64(yesCount) / float64(totalVotes)
	}
	
	avgConfidence := 0.0
	if len(session.Votes) > 0 {
		avgConfidence = totalConfidence / float64(len(session.Votes))
	}
	
	decision := dvs.determineDecision(session.VoteType, yesPercentage, yesCount, totalVotes)
	
	session.Result = &VoteResult{
		Decision:      decision,
		YesVotes:      yesCount,
		NoVotes:       noCount,
		TotalVotes:    totalVotes,
		YesPercentage: yesPercentage,
		Confidence:    avgConfidence,
		Reasoning:     reasoning,
		CompletedAt:   time.Now(),
	}
	
	session.Completed = true
}

// calculateWeightedVotes calculates votes with agent weights
func (dvs *DemocraticVotingSystem) calculateWeightedVotes(session *VoteSession) (int, int) {
	var yesWeight, noWeight float64
	
	for agentID, vote := range session.Votes {
		weight := 1.0
		if w, exists := session.AgentWeights[agentID]; exists {
			weight = w
		}
		
		if vote.Decision {
			yesWeight += weight
		} else {
			noWeight += weight
		}
	}
	
	return int(yesWeight), int(noWeight)
}

// determineDecision applies voting rules to determine outcome
func (dvs *DemocraticVotingSystem) determineDecision(
	voteType VoteType,
	yesPercentage float64,
	yesCount, totalVotes int,
) bool {
	switch voteType {
	case VoteTypeMajority:
		return yesPercentage > 0.5
	case VoteTypeSuper:
		return yesPercentage > 0.66
	case VoteTypeUnanimous:
		return yesCount == totalVotes && totalVotes > 0
	case VoteTypeWeighted:
		return yesPercentage > 0.5
	case VoteTypeConsensus:
		// Consensus requires high agreement
		return yesPercentage > 0.75
	default:
		return yesPercentage > 0.5
	}
}

// GetActiveSessions returns all active voting sessions
func (dvs *DemocraticVotingSystem) GetActiveSessions() []*VoteSession {
	dvs.mu.RLock()
	defer dvs.mu.RUnlock()
	
	var active []*VoteSession
	for _, session := range dvs.sessions {
		session.mu.RLock()
		if !session.Completed {
			active = append(active, session)
		}
		session.mu.RUnlock()
	}
	
	return active
}

// CleanupCompletedSessions removes old completed sessions
func (dvs *DemocraticVotingSystem) CleanupCompletedSessions(olderThan time.Duration) {
	dvs.mu.Lock()
	defer dvs.mu.Unlock()
	
	cutoff := time.Now().Add(-olderThan)
	toDelete := make([]string, 0)
	
	for id, session := range dvs.sessions {
		session.mu.RLock()
		if session.Completed && session.Result.CompletedAt.Before(cutoff) {
			toDelete = append(toDelete, id)
		}
		session.mu.RUnlock()
	}
	
	for _, id := range toDelete {
		delete(dvs.sessions, id)
	}
}

// ConsensusBuilder helps build consensus through iterative voting
type ConsensusBuilder struct {
	maxRounds      int
	currentRound   int
	votingSystems  *DemocraticVotingSystem
	proposal       VoteProposal
	roundResults   []*VoteResult
}

// NewConsensusBuilder creates a new consensus builder
func NewConsensusBuilder(
	votingSystem *DemocraticVotingSystem,
	proposal VoteProposal,
	maxRounds int,
) *ConsensusBuilder {
	return &ConsensusBuilder{
		maxRounds:     maxRounds,
		votingSystems: votingSystem,
		proposal:      proposal,
		roundResults:  make([]*VoteResult, 0),
	}
}

// RunConsensusRound executes one round of consensus building
func (cb *ConsensusBuilder) RunConsensusRound(ctx context.Context, minVoters int) (*VoteResult, bool, error) {
	cb.currentRound++
	
	if cb.currentRound > cb.maxRounds {
		return nil, false, fmt.Errorf("max rounds exceeded")
	}
	
	session, err := cb.votingSystems.CreateVoteSession(
		cb.proposal,
		VoteTypeConsensus,
		minVoters,
		nil,
	)
	if err != nil {
		return nil, false, err
	}
	
	// Wait for votes (with timeout)
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	
	result, err := cb.votingSystems.WaitForResult(ctx, session.ID)
	if err != nil {
		return nil, false, err
	}
	
	cb.roundResults = append(cb.roundResults, result)
	
	// Check if consensus reached (75% agreement with high confidence)
	consensusReached := result.YesPercentage > 0.75 && result.Confidence > 0.7
	
	return result, consensusReached, nil
}
