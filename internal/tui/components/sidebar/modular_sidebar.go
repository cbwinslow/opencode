package sidebar

import (
	"context"
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/opencode-ai/opencode/internal/config"
	"github.com/opencode-ai/opencode/internal/diff"
	"github.com/opencode-ai/opencode/internal/history"
	"github.com/opencode-ai/opencode/internal/pubsub"
	"github.com/opencode-ai/opencode/internal/session"
	"github.com/opencode-ai/opencode/internal/tui/styles"
)

// ModularSidebar is an enhanced sidebar with collapsible widget sections
type ModularSidebar struct {
	width, height int
	session       session.Session
	history       history.Service
	
	// Core information
	modFiles map[string]struct {
		additions int
		removals  int
	}
	
	// Widgets
	widgets        []Widget
	progressWidget *ProgressWidget
	filesWidget    *FilesystemWidget
	systemWidget   *SystemInfoWidget
	
	// Collapsible sections
	showSession      bool
	showLSP          bool
	showModifiedFiles bool
}

func NewModularSidebar(session session.Session, history history.Service) tea.Model {
	// Create widgets
	progressWidget := NewProgressWidget().(*ProgressWidget)
	filesWidget := NewFilesystemWidget().(*FilesystemWidget)
	systemWidget := NewSystemInfoWidget().(*SystemInfoWidget)
	
	widgets := []Widget{
		progressWidget,
		filesWidget,
		systemWidget,
	}
	
	return &ModularSidebar{
		session:           session,
		history:           history,
		widgets:           widgets,
		progressWidget:    progressWidget,
		filesWidget:       filesWidget,
		systemWidget:      systemWidget,
		showSession:       true,
		showLSP:           true,
		showModifiedFiles: true,
	}
}

func (m *ModularSidebar) Init() tea.Cmd {
	cmds := []tea.Cmd{}
	
	// Initialize all widgets
	for _, widget := range m.widgets {
		cmd := widget.Init()
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	
	if m.history != nil {
		ctx := context.Background()
		// Subscribe to file events
		filesCh := m.history.Subscribe(ctx)

		// Initialize the modified files map
		m.modFiles = make(map[string]struct {
			additions int
			removals  int
		})

		// Load initial files and calculate diffs
		m.loadModifiedFiles(ctx)

		// Return a command that will send file events to the Update method
		cmds = append(cmds, func() tea.Msg {
			return <-filesCh
		})
	}
	
	return tea.Batch(cmds...)
}

func (m *ModularSidebar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	
	switch msg := msg.(type) {
	case pubsub.Event[session.Session]:
		if msg.Type == pubsub.UpdatedEvent {
			if m.session.ID == msg.Payload.ID {
				m.session = msg.Payload
			}
		}
	case pubsub.Event[history.File]:
		if msg.Payload.SessionID == m.session.ID {
			// Process the individual file change
			ctx := context.Background()
			m.processFileChanges(ctx, msg.Payload)

			// Return a command to continue receiving events
			return m, func() tea.Msg {
				ctx := context.Background()
				filesCh := m.history.Subscribe(ctx)
				return <-filesCh
			}
		}
	}
	
	// Update all widgets
	for i, widget := range m.widgets {
		updated, cmd := widget.Update(msg)
		m.widgets[i] = updated
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	
	return m, tea.Batch(cmds...)
}

func (m *ModularSidebar) View() string {
	sections := []string{
		m.renderHeader(),
		"",
	}
	
	// Session section
	if m.showSession {
		sections = append(sections, m.renderSection("Session", m.sessionContent()))
		sections = append(sections, "")
	}
	
	// LSP section
	if m.showLSP {
		sections = append(sections, m.renderSection("LSP Configuration", m.lspContent()))
		sections = append(sections, "")
	}
	
	// Modified Files section
	if m.showModifiedFiles {
		sections = append(sections, m.renderSection("Modified Files", m.modifiedFilesContent()))
		sections = append(sections, "")
	}
	
	// Widget sections
	for _, widget := range m.widgets {
		if !widget.IsCollapsed() {
			sections = append(sections, m.renderSection(widget.Title(), widget.View()))
			sections = append(sections, "")
		}
	}
	
	content := lipgloss.JoinVertical(lipgloss.Top, sections...)
	
	return styles.BaseStyle.
		Width(m.width).
		PaddingLeft(4).
		PaddingRight(2).
		Height(m.height - 1).
		Render(content)
}

func (m *ModularSidebar) renderHeader() string {
	logo := fmt.Sprintf("%s %s", styles.OpenCodeIcon, "OpenCode")
	version := styles.BaseStyle.Foreground(styles.ForgroundDim).Render("Sidebar")
	
	header := lipgloss.JoinHorizontal(
		lipgloss.Left,
		styles.BaseStyle.Bold(true).Render(logo),
		" ",
		version,
	)
	
	cwd := fmt.Sprintf("cwd: %s", config.WorkingDirectory())
	cwdLine := styles.BaseStyle.Foreground(styles.ForgroundDim).Render(cwd)
	
	return lipgloss.JoinVertical(
		lipgloss.Top,
		header,
		cwdLine,
	)
}

func (m *ModularSidebar) renderSection(title string, content string) string {
	titleStyle := styles.BaseStyle.
		Width(m.width).
		Foreground(styles.PrimaryColor).
		Bold(true)
	
	return lipgloss.JoinVertical(
		lipgloss.Top,
		titleStyle.Render(title),
		content,
	)
}

func (m *ModularSidebar) sessionContent() string {
	sessionKey := styles.BaseStyle.Foreground(styles.Forground).Render("Title")
	sessionValue := styles.BaseStyle.
		Foreground(styles.Forground).
		Render(fmt.Sprintf(": %s", m.session.Title))
	return lipgloss.JoinHorizontal(lipgloss.Left, sessionKey, sessionValue)
}

func (m *ModularSidebar) lspContent() string {
	cfg := config.Get()
	
	// Get LSP names and sort them
	var lspNames []string
	for name := range cfg.LSP {
		lspNames = append(lspNames, name)
	}
	sort.Strings(lspNames)
	
	if len(lspNames) == 0 {
		return styles.BaseStyle.Foreground(styles.ForgroundDim).Render("No LSP servers configured")
	}
	
	var lspViews []string
	for _, name := range lspNames {
		lsp := cfg.LSP[name]
		lspLine := styles.BaseStyle.Foreground(styles.Forground).Render(
			fmt.Sprintf("• %s (%s)", name, lsp.Command),
		)
		lspViews = append(lspViews, lspLine)
	}
	
	return lipgloss.JoinVertical(lipgloss.Left, lspViews...)
}

func (m *ModularSidebar) modifiedFilesContent() string {
	// If no modified files, show a placeholder message
	if m.modFiles == nil || len(m.modFiles) == 0 {
		return styles.BaseStyle.Foreground(styles.ForgroundDim).Render("No modified files")
	}
	
	// Sort file paths alphabetically
	var paths []string
	for path := range m.modFiles {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	
	// Create views for each file
	var fileViews []string
	for _, path := range paths {
		stats := m.modFiles[path]
		fileViews = append(fileViews, m.renderModifiedFile(path, stats.additions, stats.removals))
	}
	
	return lipgloss.JoinVertical(lipgloss.Left, fileViews...)
}

func (m *ModularSidebar) renderModifiedFile(filePath string, additions, removals int) string {
	stats := ""
	if additions > 0 && removals > 0 {
		addStr := styles.BaseStyle.Foreground(styles.Green).Render(fmt.Sprintf("+%d", additions))
		remStr := styles.BaseStyle.Foreground(styles.Red).Render(fmt.Sprintf("-%d", removals))
		stats = fmt.Sprintf(" [%s %s]", addStr, remStr)
	} else if additions > 0 {
		stats = fmt.Sprintf(" [%s]", styles.BaseStyle.Foreground(styles.Green).Render(fmt.Sprintf("+%d", additions)))
	} else if removals > 0 {
		stats = fmt.Sprintf(" [%s]", styles.BaseStyle.Foreground(styles.Red).Render(fmt.Sprintf("-%d", removals)))
	}
	
	filePathStr := styles.BaseStyle.Render(filePath)
	return filePathStr + stats
}

func (m *ModularSidebar) SetSize(width, height int) tea.Cmd {
	m.width = width
	m.height = height
	
	// Update widget sizes
	for _, widget := range m.widgets {
		widget.SetSize(width, 0) // Height will be calculated dynamically
	}
	
	return nil
}

func (m *ModularSidebar) GetSize() (int, int) {
	return m.width, m.height
}

// Toggle methods for sections
func (m *ModularSidebar) ToggleSession() {
	m.showSession = !m.showSession
}

func (m *ModularSidebar) ToggleLSP() {
	m.showLSP = !m.showLSP
}

func (m *ModularSidebar) ToggleModifiedFiles() {
	m.showModifiedFiles = !m.showModifiedFiles
}

// File tracking methods (from original sidebar)
func (m *ModularSidebar) loadModifiedFiles(ctx context.Context) {
	if m.history == nil || m.session.ID == "" {
		return
	}

	// Get all latest files for this session
	latestFiles, err := m.history.ListLatestSessionFiles(ctx, m.session.ID)
	if err != nil {
		return
	}

	// Get all files for this session
	allFiles, err := m.history.ListBySession(ctx, m.session.ID)
	if err != nil {
		return
	}

	// Clear the existing map
	m.modFiles = make(map[string]struct {
		additions int
		removals  int
	})

	// Process each latest file
	for _, file := range latestFiles {
		if file.Version == history.InitialVersion {
			continue
		}

		// Find the initial version
		var initialVersion history.File
		for _, v := range allFiles {
			if v.Path == file.Path && v.Version == history.InitialVersion {
				initialVersion = v
				break
			}
		}

		if initialVersion.ID == "" {
			continue
		}
		if initialVersion.Content == file.Content {
			continue
		}

		// Calculate diff
		_, additions, removals := diff.GenerateDiff(initialVersion.Content, file.Content, file.Path)

		if additions > 0 || removals > 0 {
			displayPath := file.Path
			workingDir := config.WorkingDirectory()
			displayPath = strings.TrimPrefix(displayPath, workingDir)
			displayPath = strings.TrimPrefix(displayPath, "/")

			m.modFiles[displayPath] = struct {
				additions int
				removals  int
			}{
				additions: additions,
				removals:  removals,
			}
		}
	}
}

func (m *ModularSidebar) processFileChanges(ctx context.Context, file history.File) {
	if file.Version == history.InitialVersion {
		return
	}

	initialVersion, err := m.findInitialVersion(ctx, file.Path)
	if err != nil || initialVersion.ID == "" {
		return
	}

	displayPath := getDisplayPath(file.Path)
	
	if initialVersion.Content == file.Content {
		delete(m.modFiles, displayPath)
		return
	}

	_, additions, removals := diff.GenerateDiff(initialVersion.Content, file.Content, file.Path)

	if additions > 0 || removals > 0 {
		m.modFiles[displayPath] = struct {
			additions int
			removals  int
		}{
			additions: additions,
			removals:  removals,
		}
	} else {
		delete(m.modFiles, displayPath)
	}
}

func (m *ModularSidebar) findInitialVersion(ctx context.Context, path string) (history.File, error) {
	fileVersions, err := m.history.ListBySession(ctx, m.session.ID)
	if err != nil {
		return history.File{}, err
	}

	for _, v := range fileVersions {
		if v.Path == path && v.Version == history.InitialVersion {
			return v, nil
		}
	}

	return history.File{}, fmt.Errorf("initial version not found")
}

func getDisplayPath(path string) string {
	workingDir := config.WorkingDirectory()
	displayPath := strings.TrimPrefix(path, workingDir)
	return strings.TrimPrefix(displayPath, "/")
}
