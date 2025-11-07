package page

type PageID string

// PageID constants are defined in individual page files:
// - ChatPage in chat.go
// - LogsPage in logs.go
// - ToolsPage in toolspage.go

// PageChangeMsg is used to change the current page
type PageChangeMsg struct {
	ID PageID
}
