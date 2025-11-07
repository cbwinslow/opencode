package sidebar

import (
	"fmt"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/opencode-ai/opencode/internal/tui/styles"
)

// SystemInfoWidget displays system information and statistics
type SystemInfoWidget struct {
	BaseWidget
	memStats    runtime.MemStats
	numGoroutines int
	lspConnections int
}

func NewSystemInfoWidget() Widget {
	return &SystemInfoWidget{
		BaseWidget: BaseWidget{
			title: "System Info",
		},
	}
}

func (w *SystemInfoWidget) Init() tea.Cmd {
	w.updateStats()
	return nil
}

func (w *SystemInfoWidget) Update(msg tea.Msg) (Widget, tea.Cmd) {
	// Update stats periodically
	w.updateStats()
	return w, nil
}

func (w *SystemInfoWidget) View() string {
	if w.collapsed {
		return ""
	}

	var lines []string
	
	// Memory usage
	memMB := float64(w.memStats.Alloc) / 1024 / 1024
	memLine := fmt.Sprintf("Memory: %.1f MB", memMB)
	lines = append(lines, styles.BaseStyle.Foreground(styles.Forground).Render(memLine))
	
	// Goroutines
	goroutinesLine := fmt.Sprintf("Goroutines: %d", w.numGoroutines)
	lines = append(lines, styles.BaseStyle.Foreground(styles.Forground).Render(goroutinesLine))
	
	// LSP connections
	if w.lspConnections > 0 {
		lspLine := fmt.Sprintf("LSP Servers: %d", w.lspConnections)
		lines = append(lines, styles.BaseStyle.Foreground(styles.PrimaryColor).Render(lspLine))
	}
	
	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	
	return styles.BaseStyle.
		Width(w.width).
		Render(content)
}

func (w *SystemInfoWidget) GetHeight() int {
	if w.collapsed {
		return 0
	}
	
	height := 2 // Memory + Goroutines
	if w.lspConnections > 0 {
		height++ // LSP connections
	}
	return height
}

func (w *SystemInfoWidget) updateStats() {
	runtime.ReadMemStats(&w.memStats)
	w.numGoroutines = runtime.NumGoroutine()
}

func (w *SystemInfoWidget) SetLSPConnections(count int) {
	w.lspConnections = count
}
