package table

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/opencode-ai/opencode/internal/tui/styles"
)

// DataTable is a wrapper around bubbles table with custom styling
type DataTable struct {
	table  table.Model
	width  int
	height int
}

// NewDataTable creates a new data table
func NewDataTable(columns []table.Column, rows []table.Row) *DataTable {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Custom styling
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styles.PrimaryColor).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(styles.Forground).
		Background(styles.PrimaryColor).
		Bold(false)
	
	t.SetStyles(s)

	return &DataTable{
		table: t,
	}
}

// SetRows updates the table rows
func (m *DataTable) SetRows(rows []table.Row) {
	m.table.SetRows(rows)
}

// SetColumns updates the table columns
func (m *DataTable) SetColumns(columns []table.Column) {
	m.table.SetColumns(columns)
}

// SelectedRow returns the currently selected row
func (m *DataTable) SelectedRow() table.Row {
	return m.table.SelectedRow()
}

// Init implements tea.Model
func (m *DataTable) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m *DataTable) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, nil
		}
	}
	
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m *DataTable) View() string {
	help := styles.BaseStyle.
		Foreground(styles.ForgroundDim).
		Render("↑/↓/j/k: navigate • enter: select • q/esc: close")
	
	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.table.View(),
		"",
		help,
	)
}

// SetSize sets the size of the table
func (m *DataTable) SetSize(width, height int) {
	m.width = width
	m.height = height
	
	// Update table size (leave room for help)
	tableHeight := height - 3
	if tableHeight < 1 {
		tableHeight = 1
	}
	
	m.table.SetWidth(width)
	m.table.SetHeight(tableHeight)
}

// Focus focuses the table
func (m *DataTable) Focus() {
	m.table.Focus()
}

// Blur removes focus from the table
func (m *DataTable) Blur() {
	m.table.Blur()
}

// Focused returns whether the table is focused
func (m *DataTable) Focused() bool {
	return m.table.Focused()
}
