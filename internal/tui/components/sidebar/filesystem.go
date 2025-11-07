package sidebar

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/opencode-ai/opencode/internal/config"
	"github.com/opencode-ai/opencode/internal/tui/styles"
)

// FilesystemWidget displays a file browser for the project
type FilesystemWidget struct {
	BaseWidget
	rootPath     string
	currentPath  string
	files        []fileEntry
	maxFiles     int
	showHidden   bool
}

type fileEntry struct {
	name  string
	path  string
	isDir bool
}

func NewFilesystemWidget() Widget {
	return &FilesystemWidget{
		BaseWidget: BaseWidget{
			title: "Filesystem",
		},
		rootPath:   config.WorkingDirectory(),
		currentPath: config.WorkingDirectory(),
		maxFiles:   10,
		showHidden: false,
	}
}

func (w *FilesystemWidget) Init() tea.Cmd {
	w.loadDirectory()
	return nil
}

func (w *FilesystemWidget) Update(msg tea.Msg) (Widget, tea.Cmd) {
	return w, nil
}

func (w *FilesystemWidget) View() string {
	if w.collapsed {
		return ""
	}

	// Show current directory relative to root
	relPath, _ := filepath.Rel(w.rootPath, w.currentPath)
	if relPath == "." {
		relPath = "/"
	} else {
		relPath = "/" + relPath
	}
	
	header := styles.BaseStyle.
		Foreground(styles.ForgroundDim).
		Render(relPath)
	
	var fileViews []string
	displayCount := w.maxFiles
	if len(w.files) < displayCount {
		displayCount = len(w.files)
	}
	
	for i := 0; i < displayCount; i++ {
		entry := w.files[i]
		icon := "  "
		color := styles.Forground
		
		if entry.isDir {
			icon = "📁"
			color = styles.PrimaryColor
		} else {
			icon = "📄"
		}
		
		name := entry.name
		if len(name) > w.width-6 {
			name = name[:w.width-9] + "..."
		}
		
		fileView := styles.BaseStyle.
			Foreground(color).
			Render(fmt.Sprintf("%s %s", icon, name))
		fileViews = append(fileViews, fileView)
	}
	
	if len(w.files) > displayCount {
		more := styles.BaseStyle.
			Foreground(styles.ForgroundDim).
			Render(fmt.Sprintf("  ... and %d more", len(w.files)-displayCount))
		fileViews = append(fileViews, more)
	}
	
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		lipgloss.JoinVertical(lipgloss.Left, fileViews...),
	)
	
	return styles.BaseStyle.
		Width(w.width).
		Render(content)
}

func (w *FilesystemWidget) GetHeight() int {
	if w.collapsed {
		return 0
	}
	
	displayCount := w.maxFiles
	if len(w.files) < displayCount {
		displayCount = len(w.files)
	}
	
	height := 1 // header
	height += displayCount
	if len(w.files) > displayCount {
		height++ // "... and X more" line
	}
	
	return height
}

func (w *FilesystemWidget) loadDirectory() {
	w.files = []fileEntry{}
	
	entries, err := os.ReadDir(w.currentPath)
	if err != nil {
		return
	}
	
	for _, entry := range entries {
		// Skip hidden files unless showHidden is true
		if !w.showHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		
		// Skip common directories that are not useful
		if entry.IsDir() && (entry.Name() == "node_modules" || 
			entry.Name() == ".git" || 
			entry.Name() == "vendor" ||
			entry.Name() == "dist" ||
			entry.Name() == "build") {
			continue
		}
		
		w.files = append(w.files, fileEntry{
			name:  entry.Name(),
			path:  filepath.Join(w.currentPath, entry.Name()),
			isDir: entry.IsDir(),
		})
	}
	
	// Sort: directories first, then files, both alphabetically
	sort.Slice(w.files, func(i, j int) bool {
		if w.files[i].isDir != w.files[j].isDir {
			return w.files[i].isDir
		}
		return w.files[i].name < w.files[j].name
	})
}

func (w *FilesystemWidget) ToggleHidden() {
	w.showHidden = !w.showHidden
	w.loadDirectory()
}
