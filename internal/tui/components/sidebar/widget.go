package sidebar

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Widget is the base interface for all sidebar widgets
type Widget interface {
	// Init initializes the widget
	Init() tea.Cmd

	// Update handles messages and returns updated model
	Update(msg tea.Msg) (Widget, tea.Cmd)

	// View renders the widget
	View() string

	// SetSize sets the width and height for the widget
	SetSize(width, height int)

	// GetHeight returns the current height of the widget
	GetHeight() int

	// IsCollapsed returns whether the widget is collapsed
	IsCollapsed() bool

	// ToggleCollapse toggles the collapsed state
	ToggleCollapse()

	// Title returns the widget's title
	Title() string
}

// BaseWidget provides common functionality for all widgets
type BaseWidget struct {
	width     int
	height    int
	collapsed bool
	title     string
}

func (w *BaseWidget) SetSize(width, height int) {
	w.width = width
	w.height = height
}

func (w *BaseWidget) IsCollapsed() bool {
	return w.collapsed
}

func (w *BaseWidget) ToggleCollapse() {
	w.collapsed = !w.collapsed
}

func (w *BaseWidget) Title() string {
	return w.title
}

func (w *BaseWidget) GetWidth() int {
	return w.width
}
