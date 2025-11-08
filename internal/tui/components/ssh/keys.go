package ssh

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/opencode-ai/opencode/internal/tui/styles"
)

// SSHKeyInfo represents information about an SSH key
type SSHKeyInfo struct {
	Path        string
	Type        string
	Fingerprint string
	Comment     string
}

// SSHKeyViewer displays SSH keys found in the user's .ssh directory
type SSHKeyViewer struct {
	viewport viewport.Model
	keys     []SSHKeyInfo
	width    int
	height   int
}

// NewSSHKeyViewer creates a new SSH key viewer
func NewSSHKeyViewer() *SSHKeyViewer {
	return &SSHKeyViewer{
		viewport: viewport.New(80, 20),
		keys:     []SSHKeyInfo{},
	}
}

// LoadKeys scans the .ssh directory for SSH keys
func (m *SSHKeyViewer) LoadKeys() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	
	sshDir := filepath.Join(homeDir, ".ssh")
	
	// Check if .ssh directory exists
	if _, err := os.Stat(sshDir); os.IsNotExist(err) {
		return fmt.Errorf(".ssh directory not found")
	}
	
	// Common SSH key file patterns
	keyPatterns := []string{
		"id_rsa",
		"id_dsa",
		"id_ecdsa",
		"id_ed25519",
	}
	
	m.keys = []SSHKeyInfo{}
	
	// Look for public keys
	for _, pattern := range keyPatterns {
		pubKeyPath := filepath.Join(sshDir, pattern+".pub")
		if _, err := os.Stat(pubKeyPath); err == nil {
			// Read the public key file
			content, err := os.ReadFile(pubKeyPath)
			if err != nil {
				continue
			}
			
			// Parse the key info
			parts := strings.Fields(string(content))
			keyType := ""
			comment := ""
			
			if len(parts) >= 1 {
				keyType = parts[0]
			}
			if len(parts) >= 3 {
				comment = strings.Join(parts[2:], " ")
			}
			
			keyInfo := SSHKeyInfo{
				Path:    pubKeyPath,
				Type:    keyType,
				Comment: comment,
			}
			
			m.keys = append(m.keys, keyInfo)
		}
	}
	
	// Update viewport content
	m.updateContent()
	
	return nil
}

// updateContent updates the viewport with key information
func (m *SSHKeyViewer) updateContent() {
	if len(m.keys) == 0 {
		m.viewport.SetContent("No SSH keys found in ~/.ssh directory")
		return
	}
	
	var content strings.Builder
	
	for i, key := range m.keys {
		if i > 0 {
			content.WriteString("\n\n")
		}
		
		// Key header
		header := styles.BaseStyle.
			Bold(true).
			Foreground(styles.PrimaryColor).
			Render(fmt.Sprintf("Key %d: %s", i+1, filepath.Base(key.Path)))
		content.WriteString(header)
		content.WriteString("\n")
		
		// Key details
		typeLabel := styles.BaseStyle.Foreground(styles.ForgroundDim).Render("Type: ")
		typeValue := styles.BaseStyle.Render(key.Type)
		content.WriteString(typeLabel + typeValue + "\n")
		
		pathLabel := styles.BaseStyle.Foreground(styles.ForgroundDim).Render("Path: ")
		pathValue := styles.BaseStyle.Render(key.Path)
		content.WriteString(pathLabel + pathValue + "\n")
		
		if key.Comment != "" {
			commentLabel := styles.BaseStyle.Foreground(styles.ForgroundDim).Render("Comment: ")
			commentValue := styles.BaseStyle.Render(key.Comment)
			content.WriteString(commentLabel + commentValue + "\n")
		}
	}
	
	m.viewport.SetContent(content.String())
}

// Init implements tea.Model
func (m *SSHKeyViewer) Init() tea.Cmd {
	// Load keys on initialization
	_ = m.LoadKeys()
	return nil
}

// Update implements tea.Model
func (m *SSHKeyViewer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, nil
		case "r":
			// Reload keys
			_ = m.LoadKeys()
			return m, nil
		}
	}
	
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m *SSHKeyViewer) View() string {
	title := styles.BaseStyle.
		Bold(true).
		Foreground(styles.PrimaryColor).
		Render("SSH Keys")
	
	help := styles.BaseStyle.
		Foreground(styles.ForgroundDim).
		Render("↑/↓: scroll • r: reload • q/esc: close")
	
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
func (m *SSHKeyViewer) SetSize(width, height int) {
	m.width = width
	m.height = height
	
	// Update viewport size (subtract height of header)
	viewportHeight := height - 4
	if viewportHeight < 1 {
		viewportHeight = 1
	}
	
	m.viewport.Width = width
	m.viewport.Height = viewportHeight
	
	// Update content with new width
	m.updateContent()
}

// GetKeys returns the list of SSH keys
func (m *SSHKeyViewer) GetKeys() []SSHKeyInfo {
	return m.keys
}

// SSHKeyListItem implements list.Item for use in a list component
type SSHKeyListItem struct {
	info SSHKeyInfo
}

func (i SSHKeyListItem) FilterValue() string {
	return i.info.Path
}

func (i SSHKeyListItem) Title() string {
	return filepath.Base(i.info.Path)
}

func (i SSHKeyListItem) Description() string {
	return fmt.Sprintf("%s - %s", i.info.Type, i.info.Comment)
}

// NewSSHKeyListItems converts SSH keys to list items
func NewSSHKeyListItems(keys []SSHKeyInfo) []list.Item {
	items := make([]list.Item, len(keys))
	for i, key := range keys {
		items[i] = SSHKeyListItem{info: key}
	}
	return items
}
