package monitor

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// LogEntry represents a parsed log entry
type LogEntry struct {
	Timestamp time.Time
	Level     string
	Source    string
	Message   string
	Fields    map[string]interface{}
}

// LogWatcher monitors log files for changes
type LogWatcher struct {
	paths       []string
	watcher     *fsnotify.Watcher
	entries     chan LogEntry
	ctx         context.Context
	cancelFunc  context.CancelFunc
	wg          sync.WaitGroup
	fileOffsets map[string]int64
	mu          sync.Mutex
}

// LogWatcherConfig configures the log watcher
type LogWatcherConfig struct {
	Paths       []string
	BufferSize  int
	ParseFormat string // "json", "logfmt", "plain"
}

// NewLogWatcher creates a new log watcher
func NewLogWatcher(config LogWatcherConfig) (*LogWatcher, error) {
	if config.BufferSize <= 0 {
		config.BufferSize = 1000
	}
	
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}
	
	ctx, cancel := context.WithCancel(context.Background())
	
	lw := &LogWatcher{
		paths:       config.Paths,
		watcher:     watcher,
		entries:     make(chan LogEntry, config.BufferSize),
		ctx:         ctx,
		cancelFunc:  cancel,
		fileOffsets: make(map[string]int64),
	}
	
	return lw, nil
}

// Start begins monitoring log files
func (lw *LogWatcher) Start() error {
	// Add all paths to the watcher
	for _, path := range lw.paths {
		// Expand glob patterns
		matches, err := filepath.Glob(path)
		if err != nil {
			return fmt.Errorf("invalid path pattern %s: %w", path, err)
		}
		
		for _, match := range matches {
			if err := lw.addFile(match); err != nil {
				return err
			}
		}
		
		// Watch directory for new files matching pattern
		dir := filepath.Dir(path)
		if err := lw.watcher.Add(dir); err != nil {
			return fmt.Errorf("failed to watch directory %s: %w", dir, err)
		}
	}
	
	// Start the event processing loop
	lw.wg.Add(1)
	go lw.processEvents()
	
	return nil
}

// Stop stops the log watcher
func (lw *LogWatcher) Stop() error {
	lw.cancelFunc()
	lw.wg.Wait()
	
	if err := lw.watcher.Close(); err != nil {
		return err
	}
	
	close(lw.entries)
	return nil
}

// Entries returns the channel of log entries
func (lw *LogWatcher) Entries() <-chan LogEntry {
	return lw.entries
}

// addFile starts monitoring a specific file
func (lw *LogWatcher) addFile(path string) error {
	lw.mu.Lock()
	defer lw.mu.Unlock()
	
	// Get current file size to start reading from the end
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat file %s: %w", path, err)
	}
	
	lw.fileOffsets[path] = info.Size()
	
	if err := lw.watcher.Add(path); err != nil {
		return fmt.Errorf("failed to watch file %s: %w", path, err)
	}
	
	return nil
}

// processEvents handles file system events
func (lw *LogWatcher) processEvents() {
	defer lw.wg.Done()
	
	for {
		select {
		case event, ok := <-lw.watcher.Events:
			if !ok {
				return
			}
			
			if event.Op&fsnotify.Write == fsnotify.Write {
				lw.handleFileWrite(event.Name)
			} else if event.Op&fsnotify.Create == fsnotify.Create {
				lw.handleFileCreate(event.Name)
			}
			
		case err, ok := <-lw.watcher.Errors:
			if !ok {
				return
			}
			// Log error but continue watching
			_ = err
			
		case <-lw.ctx.Done():
			return
		}
	}
}

// handleFileWrite processes new data written to a file
func (lw *LogWatcher) handleFileWrite(path string) {
	lw.mu.Lock()
	offset, exists := lw.fileOffsets[path]
	lw.mu.Unlock()
	
	if !exists {
		return
	}
	
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	
	// Seek to last known position
	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		return
	}
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		entry := lw.parseLine(line, path)
		
		select {
		case lw.entries <- entry:
		case <-lw.ctx.Done():
			return
		}
	}
	
	// Update offset
	newOffset, _ := file.Seek(0, io.SeekCurrent)
	lw.mu.Lock()
	lw.fileOffsets[path] = newOffset
	lw.mu.Unlock()
}

// handleFileCreate handles newly created files
func (lw *LogWatcher) handleFileCreate(path string) {
	// Check if this file matches any of our patterns
	for _, pattern := range lw.paths {
		matched, err := filepath.Match(pattern, path)
		if err == nil && matched {
			_ = lw.addFile(path)
			break
		}
	}
}

// parseLine parses a log line into a LogEntry
func (lw *LogWatcher) parseLine(line string, source string) LogEntry {
	// Basic parsing - could be enhanced with structured log parsing
	return LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Source:    source,
		Message:   line,
		Fields:    make(map[string]interface{}),
	}
}

// ShellHistoryWatcher monitors shell history
type ShellHistoryWatcher struct {
	historyFile string
	entries     chan string
	ctx         context.Context
	cancelFunc  context.CancelFunc
	wg          sync.WaitGroup
	lastOffset  int64
	mu          sync.Mutex
}

// NewShellHistoryWatcher creates a new shell history watcher
func NewShellHistoryWatcher(historyFile string, bufferSize int) (*ShellHistoryWatcher, error) {
	if bufferSize <= 0 {
		bufferSize = 100
	}
	
	ctx, cancel := context.WithCancel(context.Background())
	
	// Get initial file size
	info, err := os.Stat(historyFile)
	var offset int64
	if err == nil {
		offset = info.Size()
	}
	
	return &ShellHistoryWatcher{
		historyFile: historyFile,
		entries:     make(chan string, bufferSize),
		ctx:         ctx,
		cancelFunc:  cancel,
		lastOffset:  offset,
	}, nil
}

// Start begins monitoring shell history
func (shw *ShellHistoryWatcher) Start() error {
	shw.wg.Add(1)
	go shw.monitor()
	return nil
}

// Stop stops the shell history watcher
func (shw *ShellHistoryWatcher) Stop() error {
	shw.cancelFunc()
	shw.wg.Wait()
	close(shw.entries)
	return nil
}

// Entries returns the channel of history entries
func (shw *ShellHistoryWatcher) Entries() <-chan string {
	return shw.entries
}

// monitor periodically checks for new history entries
func (shw *ShellHistoryWatcher) monitor() {
	defer shw.wg.Done()
	
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			shw.readNewEntries()
		case <-shw.ctx.Done():
			return
		}
	}
}

// readNewEntries reads new entries from the history file
func (shw *ShellHistoryWatcher) readNewEntries() {
	shw.mu.Lock()
	defer shw.mu.Unlock()
	
	file, err := os.Open(shw.historyFile)
	if err != nil {
		return
	}
	defer file.Close()
	
	// Seek to last known position
	if _, err := file.Seek(shw.lastOffset, io.SeekStart); err != nil {
		return
	}
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			select {
			case shw.entries <- line:
			case <-shw.ctx.Done():
				return
			default:
				// Buffer full, skip
			}
		}
	}
	
	// Update offset
	newOffset, _ := file.Seek(0, io.SeekCurrent)
	shw.lastOffset = newOffset
}
