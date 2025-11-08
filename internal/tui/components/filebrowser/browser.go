package filebrowser

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/opencode-ai/opencode/internal/tui/styles"
)

// FileItem represents a file or directory in the tree
type FileItem struct {
	name    string
	path    string
	isDir   bool
	size    int64
}

// Implement list.Item interface
func (i FileItem) FilterValue() string { return i.name }
func (i FileItem) Title() string       { 
	if i.isDir {
		return "📁 " + i.name
	}
	return "📄 " + i.name
}

func (i FileItem) Description() string { 
	if i.isDir {
		return i.path
	}
	return i.path
}

// FileBrowser is a file tree browser component
type FileBrowser struct {
	list          list.Model
	currentPath   string
	width         int
	height        int
	selectedFile  string
}

// NewFileBrowser creates a new file browser
func NewFileBrowser(startPath string) *FileBrowser {
	items := []list.Item{}
	
	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 0, 0)
	l.Title = "File Browser"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	
	fb := &FileBrowser{
		list:        l,
		currentPath: startPath,
	}
	
	// Load initial directory
	_ = fb.loadDirectory(startPath)
	
	return fb
}

// loadDirectory loads files from a directory
func (m *FileBrowser) loadDirectory(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	
	var items []list.Item
	
	// Add parent directory entry if not at root
	if path != "/" && path != "" {
		items = append(items, FileItem{
			name:  "..",
			path:  filepath.Dir(path),
			isDir: true,
		})
	}
	
	// Sort directories first, then files
	var dirs []os.DirEntry
	var files []os.DirEntry
	
	for _, entry := range entries {
		// Skip hidden files
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		
		if entry.IsDir() {
			dirs = append(dirs, entry)
		} else {
			files = append(files, entry)
		}
	}
	
	// Sort each group alphabetically
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Name() < dirs[j].Name()
	})
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	
	// Add directories
	for _, dir := range dirs {
		info, _ := dir.Info()
		items = append(items, FileItem{
			name:  dir.Name(),
			path:  filepath.Join(path, dir.Name()),
			isDir: true,
			size:  info.Size(),
		})
	}
	
	// Add files
	for _, file := range files {
		info, _ := file.Info()
		items = append(items, FileItem{
			name:  file.Name(),
			path:  filepath.Join(path, file.Name()),
			isDir: false,
			size:  info.Size(),
		})
	}
	
	m.list.SetItems(items)
	m.currentPath = path
	m.list.Title = "File Browser: " + path
	
	return nil
}

// Init implements tea.Model
func (m *FileBrowser) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m *FileBrowser) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, nil
		case "enter":
			// Navigate into directory or select file
			if selected, ok := m.list.SelectedItem().(FileItem); ok {
				if selected.isDir {
					// Navigate into directory
					_ = m.loadDirectory(selected.path)
					return m, nil
				} else {
					// File selected
					m.selectedFile = selected.path
					return m, nil
				}
			}
		case "backspace":
			// Go to parent directory
			parent := filepath.Dir(m.currentPath)
			if parent != m.currentPath {
				_ = m.loadDirectory(parent)
			}
			return m, nil
		}
	}
	
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m *FileBrowser) View() string {
	helpStyle := lipgloss.NewStyle().Foreground(styles.ForgroundDim)
	help := helpStyle.Render("\nenter: open • backspace: parent • /: filter • q/esc: close")
	
	return m.list.View() + "\n" + help
}

// SetSize sets the size of the browser
func (m *FileBrowser) SetSize(width, height int) {
	m.width = width
	m.height = height
	
	// Leave room for help text
	m.list.SetSize(width, height-2)
}

// GetSelectedFile returns the currently selected file path
func (m *FileBrowser) GetSelectedFile() string {
	return m.selectedFile
}

// GetCurrentPath returns the current directory path
func (m *FileBrowser) GetCurrentPath() string {
	return m.currentPath
}

// SetCurrentPath sets the current directory and loads it
func (m *FileBrowser) SetCurrentPath(path string) error {
	return m.loadDirectory(path)
}
