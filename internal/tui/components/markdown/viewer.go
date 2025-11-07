package markdown

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/opencode-ai/opencode/internal/tui/styles"
)

// MarkdownViewer is a component that renders markdown content with Glamour
type MarkdownViewer struct {
	viewport viewport.Model
	content  string
	width    int
	height   int
	renderer *glamour.TermRenderer
}

// NewMarkdownViewer creates a new markdown viewer
func NewMarkdownViewer() *MarkdownViewer {
	// Create a glamour renderer with a dark theme
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)

	return &MarkdownViewer{
		viewport: viewport.New(80, 20),
		renderer: renderer,
	}
}

// SetContent sets the markdown content to be rendered
func (m *MarkdownViewer) SetContent(content string) error {
	m.content = content
	
	// Render the markdown content
	rendered, err := m.renderer.Render(content)
	if err != nil {
		return err
	}
	
	m.viewport.SetContent(rendered)
	return nil
}

// Init implements tea.Model
func (m *MarkdownViewer) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m *MarkdownViewer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			// Close the viewer
			return m, nil
		}
	}
	
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m *MarkdownViewer) View() string {
	title := styles.BaseStyle.
		Bold(true).
		Foreground(styles.PrimaryColor).
		Render("Markdown Preview")
	
	help := styles.BaseStyle.
		Foreground(styles.ForgroundDim).
		Render("↑/↓: scroll • q/esc: close")
	
	header := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		help,
		"",
	)
	
	return lipgloss.JoinVertical(
		lipgloss.Top,
		header,
		m.viewport.View(),
	)
}

// SetSize sets the size of the viewer
func (m *MarkdownViewer) SetSize(width, height int) {
	m.width = width
	m.height = height
	
	// Update viewport size (subtract height of header)
	viewportHeight := height - 4
	if viewportHeight < 1 {
		viewportHeight = 1
	}
	
	m.viewport.Width = width
	m.viewport.Height = viewportHeight
	
	// Re-render with new width if we have content
	if m.content != "" {
		// Update renderer word wrap
		m.renderer, _ = glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(width-4),
		)
		
		// Re-render content
		rendered, err := m.renderer.Render(m.content)
		if err == nil {
			m.viewport.SetContent(rendered)
		}
	}
}

// GetContent returns the raw markdown content
func (m *MarkdownViewer) GetContent() string {
	return m.content
}

// MarkdownPreviewMsg is a message to show markdown preview
type MarkdownPreviewMsg struct {
	Content string
	Title   string
}

// RenderMarkdown is a helper function to quickly render markdown to a string
func RenderMarkdown(content string, width int) (string, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return "", err
	}
	
	return renderer.Render(content)
}

// RenderMarkdownFile renders markdown from a file path
func RenderMarkdownFile(filePath string, width int) (string, error) {
	// Read the file content would go here
	// For now, return an error as we need file reading implementation
	return "", nil
}

// TruncateMarkdown truncates markdown content to a certain number of lines
func TruncateMarkdown(content string, maxLines int) string {
	lines := strings.Split(content, "\n")
	if len(lines) <= maxLines {
		return content
	}
	
	truncated := strings.Join(lines[:maxLines], "\n")
	return truncated + "\n\n_[Content truncated...]_"
}
