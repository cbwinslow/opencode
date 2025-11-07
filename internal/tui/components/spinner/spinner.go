package spinner

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/opencode-ai/opencode/internal/tui/styles"
)

// LoadingSpinner is a component that shows a loading indicator with a message
type LoadingSpinner struct {
	spinner spinner.Model
	message string
	active  bool
}

// NewLoadingSpinner creates a new loading spinner
func NewLoadingSpinner() *LoadingSpinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(styles.PrimaryColor)
	
	return &LoadingSpinner{
		spinner: s,
		message: "Loading...",
		active:  false,
	}
}

// SetMessage sets the loading message
func (m *LoadingSpinner) SetMessage(msg string) {
	m.message = msg
}

// Start starts the spinner
func (m *LoadingSpinner) Start() tea.Cmd {
	m.active = true
	return m.spinner.Tick
}

// Stop stops the spinner
func (m *LoadingSpinner) Stop() {
	m.active = false
}

// IsActive returns whether the spinner is active
func (m *LoadingSpinner) IsActive() bool {
	return m.active
}

// Init implements tea.Model
func (m *LoadingSpinner) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m *LoadingSpinner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.active {
		return m, nil
	}
	
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m *LoadingSpinner) View() string {
	if !m.active {
		return ""
	}
	
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.spinner.View(),
		" ",
		m.message,
	)
}

// InlineSpinner is a simple inline spinner without a message
type InlineSpinner struct {
	spinner spinner.Model
}

// NewInlineSpinner creates a new inline spinner
func NewInlineSpinner() *InlineSpinner {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(styles.PrimaryColor)
	
	return &InlineSpinner{
		spinner: s,
	}
}

// Init implements tea.Model
func (m *InlineSpinner) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update implements tea.Model
func (m *InlineSpinner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m *InlineSpinner) View() string {
	return m.spinner.View()
}

// ProgressIndicator shows progress with a message and percentage
type ProgressIndicator struct {
	current  int
	total    int
	message  string
	spinner  spinner.Model
	showSpinner bool
}

// NewProgressIndicator creates a new progress indicator
func NewProgressIndicator(total int) *ProgressIndicator {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(styles.PrimaryColor)
	
	return &ProgressIndicator{
		current:     0,
		total:       total,
		message:     "",
		spinner:     s,
		showSpinner: true,
	}
}

// SetProgress updates the current progress
func (m *ProgressIndicator) SetProgress(current int, message string) {
	m.current = current
	m.message = message
}

// Increment increments the progress by 1
func (m *ProgressIndicator) Increment(message string) {
	m.current++
	m.message = message
}

// SetTotal sets the total number of items
func (m *ProgressIndicator) SetTotal(total int) {
	m.total = total
}

// GetPercentage returns the completion percentage
func (m *ProgressIndicator) GetPercentage() int {
	if m.total == 0 {
		return 0
	}
	return (m.current * 100) / m.total
}

// IsComplete returns whether the progress is complete
func (m *ProgressIndicator) IsComplete() bool {
	return m.current >= m.total
}

// Init implements tea.Model
func (m *ProgressIndicator) Init() tea.Cmd {
	if m.showSpinner {
		return m.spinner.Tick
	}
	return nil
}

// Update implements tea.Model
func (m *ProgressIndicator) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.showSpinner || m.IsComplete() {
		return m, nil
	}
	
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m *ProgressIndicator) View() string {
	percentage := m.GetPercentage()
	
	progressBar := renderProgressBar(percentage, 20)
	percentText := lipgloss.NewStyle().Render(string(rune(percentage)) + "%")
	
	var spinnerView string
	if m.showSpinner && !m.IsComplete() {
		spinnerView = m.spinner.View() + " "
	}
	
	percentStr := lipgloss.NewStyle().
		Foreground(styles.PrimaryColor).
		Bold(true).
		Render(spinnerView + progressBar + " " + percentText)
	
	if m.message != "" {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			percentStr,
			" ",
			styles.BaseStyle.Foreground(styles.Forground).Render(m.message),
		)
	}
	
	return percentStr
}

// renderProgressBar renders a simple text-based progress bar
func renderProgressBar(percentage int, width int) string {
	filled := (percentage * width) / 100
	empty := width - filled
	
	bar := ""
	for i := 0; i < filled; i++ {
		bar += "█"
	}
	for i := 0; i < empty; i++ {
		bar += "░"
	}
	
	return lipgloss.NewStyle().
		Foreground(styles.PrimaryColor).
		Render(bar)
}
