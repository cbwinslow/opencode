package page

type PageID string

const (
	ChatPage  PageID = "chat"
	LogsPage  PageID = "logs"
	ToolsPage PageID = "tools"
)

// PageChangeMsg is used to change the current page
type PageChangeMsg struct {
	ID PageID
}
