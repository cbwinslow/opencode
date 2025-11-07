package sidebar

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/opencode-ai/opencode/internal/tui/styles"
)

// ProgressWidget displays current AI operations and progress
type ProgressWidget struct {
	BaseWidget
	isBusy      bool
	currentTask string
	progress    float64
}

func NewProgressWidget() Widget {
	return &ProgressWidget{
		BaseWidget: BaseWidget{
			title: "Progress",
		},
	}
}

func (w *ProgressWidget) Init() tea.Cmd {
	return nil
}

func (w *ProgressWidget) Update(msg tea.Msg) (Widget, tea.Cmd) {
	// TODO: Handle progress update messages
	return w, nil
}

func (w *ProgressWidget) View() string {
	if w.collapsed {
		return ""
	}

	content := ""
	if w.isBusy {
		status := styles.BaseStyle.Foreground(styles.PrimaryColor).Render("● Active")
		if w.currentTask != "" {
			task := styles.BaseStyle.Foreground(styles.Forground).Render(fmt.Sprintf("\n  %s", w.currentTask))
			content = lipgloss.JoinVertical(lipgloss.Left, status, task)
		} else {
			content = status
		}
		
		if w.progress > 0 && w.progress < 1 {
			progressBar := renderProgressBar(w.width-4, w.progress)
			content = lipgloss.JoinVertical(lipgloss.Left, content, progressBar)
		}
	} else {
		content = styles.BaseStyle.Foreground(styles.ForgroundDim).Render("○ Idle")
	}

	return styles.BaseStyle.
		Width(w.width).
		Render(content)
}

func (w *ProgressWidget) GetHeight() int {
	if w.collapsed {
		return 0
	}
	if w.isBusy && w.currentTask != "" {
		if w.progress > 0 && w.progress < 1 {
			return 3 // Status + task + progress bar
		}
		return 2 // Status + task
	}
	return 1 // Just status
}

func (w *ProgressWidget) SetBusy(busy bool, task string) {
	w.isBusy = busy
	w.currentTask = task
}

func (w *ProgressWidget) SetProgress(progress float64) {
	w.progress = progress
}

func renderProgressBar(width int, progress float64) string {
	if width < 4 {
		return ""
	}
	
	filled := int(float64(width-2) * progress)
	empty := width - 2 - filled
	
	bar := "["
	for i := 0; i < filled; i++ {
		bar += "="
	}
	for i := 0; i < empty; i++ {
		bar += " "
	}
	bar += "]"
	
	return styles.BaseStyle.Foreground(styles.PrimaryColor).Render(bar)
}
