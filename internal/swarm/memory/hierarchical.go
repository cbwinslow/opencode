package memory

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/google/uuid"
)

// HierarchicalMemoryStore implements a hierarchical memory system
type HierarchicalMemoryStore struct {
	memories    map[string]*Memory
	hierarchy   *HierarchicalNode
	mu          sync.RWMutex
	encryptionKey []byte
	
	// Configuration
	maxMemories      int
	consolidationInterval time.Duration
	pruneOlderThan   time.Duration
}

// HierarchicalMemoryConfig configures the memory store
type HierarchicalMemoryConfig struct {
	MaxMemories           int
	ConsolidationInterval time.Duration
	PruneOlderThan        time.Duration
	EncryptionKey         []byte
}

// NewHierarchicalMemoryStore creates a new hierarchical memory store
func NewHierarchicalMemoryStore(config HierarchicalMemoryConfig) *HierarchicalMemoryStore {
	if config.MaxMemories <= 0 {
		config.MaxMemories = 10000
	}
	if config.ConsolidationInterval <= 0 {
		config.ConsolidationInterval = 1 * time.Hour
	}
	if config.PruneOlderThan <= 0 {
		config.PruneOlderThan = 30 * 24 * time.Hour // 30 days
	}
	
	return &HierarchicalMemoryStore{
		memories:              make(map[string]*Memory),
		hierarchy:             &HierarchicalNode{ID: "root", Type: MemoryTypeSemantic, Level: 0},
		maxMemories:           config.MaxMemories,
		consolidationInterval: config.ConsolidationInterval,
		pruneOlderThan:        config.PruneOlderThan,
		encryptionKey:         config.EncryptionKey,
	}
}

// Store adds a memory to the store
func (hms *HierarchicalMemoryStore) Store(memory Memory) error {
	hms.mu.Lock()
	defer hms.mu.Unlock()
	
	if memory.ID == "" {
		memory.ID = uuid.New().String()
	}
	
	if memory.CreatedAt.IsZero() {
		memory.CreatedAt = time.Now()
	}
	
	// Encrypt if requested
	if memory.Encrypted && hms.encryptionKey != nil {
		encrypted, err := hms.encrypt(memory.Content)
		if err != nil {
			return fmt.Errorf("encryption failed: %w", err)
		}
		memory.Content = encrypted
	}
	
	hms.memories[memory.ID] = &memory
	
	// Add to hierarchy
	hms.addToHierarchy(&memory)
	
	// Check if we need to prune
	if len(hms.memories) > hms.maxMemories {
		hms.pruneOldest()
	}
	
	return nil
}

// Retrieve gets a memory by ID
func (hms *HierarchicalMemoryStore) Retrieve(id string) (*Memory, error) {
	hms.mu.RLock()
	defer hms.mu.RUnlock()
	
	memory, exists := hms.memories[id]
	if !exists {
		return nil, fmt.Errorf("memory not found: %s", id)
	}
	
	// Update access statistics
	memory.AccessCount++
	memory.LastAccessed = time.Now()
	
	// Decrypt if needed
	if memory.Encrypted && hms.encryptionKey != nil {
		decrypted, err := hms.decrypt(memory.Content)
		if err != nil {
			return nil, fmt.Errorf("decryption failed: %w", err)
		}
		
		// Return a copy with decrypted content
		decryptedMemory := *memory
		decryptedMemory.Content = decrypted
		return &decryptedMemory, nil
	}
	
	return memory, nil
}

// Update modifies an existing memory
func (hms *HierarchicalMemoryStore) Update(id string, memory Memory) error {
	hms.mu.Lock()
	defer hms.mu.Unlock()
	
	if _, exists := hms.memories[id]; !exists {
		return fmt.Errorf("memory not found: %s", id)
	}
	
	memory.ID = id
	
	if memory.Encrypted && hms.encryptionKey != nil {
		encrypted, err := hms.encrypt(memory.Content)
		if err != nil {
			return fmt.Errorf("encryption failed: %w", err)
		}
		memory.Content = encrypted
	}
	
	hms.memories[id] = &memory
	return nil
}

// Delete removes a memory
func (hms *HierarchicalMemoryStore) Delete(id string) error {
	hms.mu.Lock()
	defer hms.mu.Unlock()
	
	delete(hms.memories, id)
	return nil
}

// Query searches for memories matching criteria
func (hms *HierarchicalMemoryStore) Query(query MemoryQuery) ([]Memory, error) {
	hms.mu.RLock()
	defer hms.mu.RUnlock()
	
	var results []Memory
	
	for _, memory := range hms.memories {
		if hms.matchesQuery(memory, query) {
			results = append(results, *memory)
			if len(results) >= query.Limit && query.Limit > 0 {
				break
			}
		}
	}
	
	return results, nil
}

// VectorSearch performs similarity search using vectors
func (hms *HierarchicalMemoryStore) VectorSearch(vector []float64, limit int) ([]Memory, error) {
	hms.mu.RLock()
	defer hms.mu.RUnlock()
	
	// Calculate cosine similarity for all memories with vectors
	type scoredMemory struct {
		memory *Memory
		score  float64
	}
	
	var scored []scoredMemory
	for _, memory := range hms.memories {
		if len(memory.Vector) > 0 {
			similarity := cosineSimilarity(vector, memory.Vector)
			scored = append(scored, scoredMemory{memory, similarity})
		}
	}
	
	// Sort by score (descending)
	// Simple bubble sort for now
	for i := 0; i < len(scored); i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score > scored[i].score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}
	
	// Return top results
	var results []Memory
	for i := 0; i < len(scored) && i < limit; i++ {
		results = append(results, *scored[i].memory)
	}
	
	return results, nil
}

// Consolidate merges and organizes memories
func (hms *HierarchicalMemoryStore) Consolidate() error {
	hms.mu.Lock()
	defer hms.mu.Unlock()
	
	// Group similar episodic memories into semantic memories
	episodicMemories := make([]*Memory, 0)
	for _, memory := range hms.memories {
		if memory.Type == MemoryTypeEpisodic {
			episodicMemories = append(episodicMemories, memory)
		}
	}
	
	// Consolidate episodic memories (simplified version)
	// In a real implementation, this would use clustering or LLM summarization
	
	return nil
}

// Prune removes memories based on criteria
func (hms *HierarchicalMemoryStore) Prune(criteria PruneCriteria) error {
	hms.mu.Lock()
	defer hms.mu.Unlock()
	
	cutoffTime := time.Now().Add(-criteria.MaxAge)
	toDelete := make([]string, 0)
	
	for id, memory := range hms.memories {
		// Skip if it has a preserved tag
		if hasAnyTag(memory.Tags, criteria.PreserveTags) {
			continue
		}
		
		// Check criteria
		if memory.CreatedAt.Before(cutoffTime) ||
			memory.AccessCount < criteria.MinAccessCount {
			toDelete = append(toDelete, id)
		}
	}
	
	// Delete marked memories
	for _, id := range toDelete {
		delete(hms.memories, id)
	}
	
	return nil
}

// GetStats returns statistics about the memory store
func (hms *HierarchicalMemoryStore) GetStats() MemoryStats {
	hms.mu.RLock()
	defer hms.mu.RUnlock()
	
	stats := MemoryStats{
		TotalMemories:  len(hms.memories),
		MemoriesByType: make(map[MemoryType]int),
	}
	
	var totalAccess int
	var oldest, newest time.Time
	
	for _, memory := range hms.memories {
		stats.MemoriesByType[memory.Type]++
		totalAccess += memory.AccessCount
		
		if oldest.IsZero() || memory.CreatedAt.Before(oldest) {
			oldest = memory.CreatedAt
		}
		if newest.IsZero() || memory.CreatedAt.After(newest) {
			newest = memory.CreatedAt
		}
	}
	
	if len(hms.memories) > 0 {
		stats.AverageAccessCount = float64(totalAccess) / float64(len(hms.memories))
	}
	
	stats.OldestMemory = oldest
	stats.NewestMemory = newest
	
	return stats
}

// Helper methods

func (hms *HierarchicalMemoryStore) addToHierarchy(memory *Memory) {
	// Simplified hierarchy addition
	// In a real implementation, this would use semantic clustering
}

func (hms *HierarchicalMemoryStore) pruneOldest() {
	// Find and remove oldest, least accessed memories
	var oldest *Memory
	for _, memory := range hms.memories {
		if oldest == nil || memory.CreatedAt.Before(oldest.CreatedAt) {
			if memory.AccessCount == 0 {
				oldest = memory
			}
		}
	}
	
	if oldest != nil {
		delete(hms.memories, oldest.ID)
	}
}

func (hms *HierarchicalMemoryStore) matchesQuery(memory *Memory, query MemoryQuery) bool {
	if query.Type != "" && memory.Type != query.Type {
		return false
	}
	
	if memory.Priority < query.MinPriority {
		return false
	}
	
	if len(query.Tags) > 0 && !hasAnyTag(memory.Tags, query.Tags) {
		return false
	}
	
	if query.TimeRange != nil {
		if memory.CreatedAt.Before(query.TimeRange.Start) ||
			memory.CreatedAt.After(query.TimeRange.End) {
			return false
		}
	}
	
	return true
}

func (hms *HierarchicalMemoryStore) encrypt(data interface{}) ([]byte, error) {
	plaintext, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	
	block, err := aes.NewCipher(hms.encryptionKey)
	if err != nil {
		return nil, err
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (hms *HierarchicalMemoryStore) decrypt(data interface{}) (interface{}, error) {
	ciphertext, ok := data.([]byte)
	if !ok {
		return data, nil // Not encrypted
	}
	
	block, err := aes.NewCipher(hms.encryptionKey)
	if err != nil {
		return nil, err
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	
	var result interface{}
	if err := json.Unmarshal(plaintext, &result); err != nil {
		return nil, err
	}
	
	return result, nil
}

func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}
	
	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	
	if normA == 0 || normB == 0 {
		return 0
	}
	
	return dotProduct / (normA * normB)
}

func hasAnyTag(tags, searchTags []string) bool {
	for _, tag := range tags {
		for _, searchTag := range searchTags {
			if tag == searchTag {
				return true
			}
		}
	}
	return false
}
