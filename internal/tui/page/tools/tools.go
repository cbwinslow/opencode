package tools

import (
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/opencode-ai/opencode/internal/config"
	"github.com/opencode-ai/opencode/internal/tui/components/filebrowser"
	"github.com/opencode-ai/opencode/internal/tui/components/markdown"
	"github.com/opencode-ai/opencode/internal/tui/components/ssh"
	"github.com/opencode-ai/opencode/internal/tui/styles"
)

// ToolType represents different tool views
type ToolType int

const (
	ToolNone ToolType = iota
	ToolMarkdownViewer
	ToolSSHKeys
	ToolFileBrowser
)

// ToolsPage is a page that showcases various tools and utilities
type ToolsPage struct {
	width  int
	height int
	
	// Current tool being displayed
	currentTool ToolType
	
	// Tool components
	markdownViewer *markdown.MarkdownViewer
	sshViewer      *ssh.SSHKeyViewer
	fileBrowser    *filebrowser.FileBrowser
}

// NewToolsPage creates a new tools page
func NewToolsPage() *ToolsPage {
	workingDir := config.WorkingDirectory()
	
	return &ToolsPage{
		currentTool:    ToolNone,
		markdownViewer: markdown.NewMarkdownViewer(),
		sshViewer:      ssh.NewSSHKeyViewer(),
		fileBrowser:    filebrowser.NewFileBrowser(workingDir),
	}
}

// Init implements tea.Model
func (m *ToolsPage) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m *ToolsPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle tool-specific keys
		if m.currentTool != ToolNone {
			switch m.currentTool {
			case ToolMarkdownViewer:
				_, cmd := m.markdownViewer.Update(msg)
				cmds = append(cmds, cmd)
			case ToolSSHKeys:
				_, cmd := m.sshViewer.Update(msg)
				cmds = append(cmds, cmd)
			case ToolFileBrowser:
				_, cmd := m.fileBrowser.Update(msg)
				cmds = append(cmds, cmd)
			}
			
			// Escape key to return to menu
			if msg.String() == "esc" || msg.String() == "q" {
				m.currentTool = ToolNone
			}
			return m, tea.Batch(cmds...)
		}
		
		// Main menu keys
		switch msg.String() {
		case "1":
			m.currentTool = ToolMarkdownViewer
			// Load README as example
			readmePath := config.WorkingDirectory() + "/README.md"
			if content, err := os.ReadFile(readmePath); err == nil {
				_ = m.markdownViewer.SetContent(string(content))
			} else {
				_ = m.markdownViewer.SetContent("# Markdown Viewer\n\nNo README.md found in the current directory.\n\nThis viewer uses Glamour to render markdown beautifully in the terminal.")
			}
		case "2":
			m.currentTool = ToolSSHKeys
			_ = m.sshViewer.LoadKeys()
		case "3":
			m.currentTool = ToolFileBrowser
		case "q", "esc":
			// Return to previous page would be handled by parent
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// Update component sizes
		m.markdownViewer.SetSize(msg.Width, msg.Height)
		m.sshViewer.SetSize(msg.Width, msg.Height)
		m.fileBrowser.SetSize(msg.Width, msg.Height)
	}
	
	return m, tea.Batch(cmds...)
}

// View implements tea.Model
func (m *ToolsPage) View() string {
	// Show tool-specific view if a tool is active
	if m.currentTool != ToolNone {
		switch m.currentTool {
		case ToolMarkdownViewer:
			return m.markdownViewer.View()
		case ToolSSHKeys:
			return m.sshViewer.View()
		case ToolFileBrowser:
			return m.fileBrowser.View()
		}
	}
	
	// Show main menu
	return m.renderMenu()
}

// renderMenu renders the tools menu
func (m *ToolsPage) renderMenu() string {
	title := styles.BaseStyle.
		Bold(true).
		Foreground(styles.PrimaryColor).
		Render("🔧 OpenCode Tools")
	
	subtitle := styles.BaseStyle.
		Foreground(styles.ForgroundDim).
		Render("Enhanced features powered by Charm Bracelet")
	
	menuItems := []string{
		"1. 📖 Markdown Viewer - View README and markdown files with beautiful rendering",
		"2. 🔑 SSH Keys - View your SSH keys and configuration",
		"3. 📂 File Browser - Navigate project files with an interactive browser",
	}
	
	var styledItems []string
	for _, item := range menuItems {
		styledItems = append(styledItems, styles.BaseStyle.
			Foreground(styles.Forground).
			Render("  "+item))
	}
	
	help := styles.BaseStyle.
		Foreground(styles.ForgroundDim).
		Render("\nPress 1-3 to select a tool • q/esc to return")
	
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		title,
		subtitle,
		"",
		"",
		lipgloss.JoinVertical(lipgloss.Left, styledItems...),
		"",
		help,
	)
	
	// Center the content
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

// SetSize implements layout.Sizeable
func (m *ToolsPage) SetSize(width, height int) tea.Cmd {
	m.width = width
	m.height = height
	
	m.markdownViewer.SetSize(width, height)
	m.sshViewer.SetSize(width, height)
	m.fileBrowser.SetSize(width, height)
	
	return nil
}

// GetSize returns the current size
func (m *ToolsPage) GetSize() (int, int) {
	return m.width, m.height
}

// BindingKeys implements layout.Bindings
func (m *ToolsPage) BindingKeys() []key.Binding {
	if m.currentTool != ToolNone {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("esc", "q"),
				key.WithHelp("esc/q", "return to menu"),
			),
		}
	}
	
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("1", "2", "3"),
			key.WithHelp("1-3", "select tool"),
		),
		key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q/esc", "return"),
		),
	}
}
