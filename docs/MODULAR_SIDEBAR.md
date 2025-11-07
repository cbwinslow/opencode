# Modular Sidebar Documentation

## Overview

The OpenCode TUI now features an enhanced modular sidebar system that provides comprehensive information about your coding session in an organized, collapsible interface.

## Features

### 1. Session Information
- **Title**: Current session name
- **Status**: Session state and activity
- **Toggle**: `Ctrl+T S`

### 2. LSP Configuration
- **Servers**: List of configured language servers
- **Commands**: Language server executable paths
- **Toggle**: `Ctrl+T L`

### 3. Modified Files
- **Files**: List of files changed in current session
- **Diff Stats**: Addition (+) and deletion (-) counts
- **Real-time Updates**: Automatically updates as files change
- **Toggle**: `Ctrl+T M`

### 4. Progress Widget
- **Status**: Shows whether AI is active or idle
- **Current Task**: Displays the current operation being performed
- **Progress Bar**: Visual progress indicator (when applicable)
- **Toggle**: `Ctrl+T P`

### 5. Filesystem Widget
- **Browser**: Navigate project directory structure
- **Icons**: Visual indicators for files (📄) and directories (📁)
- **Smart Filtering**: Hides common build artifacts (node_modules, .git, etc.)
- **Compact View**: Shows top 10 items with "... and X more" indicator
- **Toggle**: `Ctrl+T F`

### 6. System Info Widget
- **Memory**: Current memory usage in MB
- **Goroutines**: Number of active goroutines
- **LSP Connections**: Active language server connections
- **Toggle**: `Ctrl+T I`

## Keyboard Shortcuts

All sidebar sections can be toggled with keyboard shortcuts:

| Shortcut | Action |
|----------|--------|
| `Ctrl+T S` | Toggle Session section |
| `Ctrl+T L` | Toggle LSP Configuration section |
| `Ctrl+T M` | Toggle Modified Files section |
| `Ctrl+T P` | Toggle Progress widget |
| `Ctrl+T F` | Toggle Filesystem widget |
| `Ctrl+T I` | Toggle System Info widget |

## Visual Indicators

- **▼** - Section is expanded
- **▶** - Section is collapsed
- **●** - AI is active
- **○** - AI is idle
- **📁** - Directory
- **📄** - File

## Architecture

The sidebar is built with a modular widget system:

```
ModularSidebar
├── BaseWidget (interface)
│   ├── Init()
│   ├── Update()
│   ├── View()
│   ├── SetSize()
│   └── ToggleCollapse()
├── ProgressWidget
├── FilesystemWidget
└── SystemInfoWidget
```

Each widget is self-contained and can be:
- Independently collapsed/expanded
- Updated with its own data
- Sized dynamically based on content

## Implementation Details

### File Structure

```
internal/tui/components/sidebar/
├── widget.go           # Base widget interface
├── progress.go         # Progress tracking widget
├── filesystem.go       # File browser widget
├── system_info.go      # System stats widget
└── modular_sidebar.go  # Main sidebar component
```

### Integration

The modular sidebar is integrated into the chat page and can be toggled on/off:

```go
// In chat.go
chatPage{
    useModularSidebar: true, // Enable modular sidebar
    ...
}
```

## Future Enhancements

Potential additions to the sidebar system:

1. **Git Integration Widget**
   - Current branch
   - Uncommitted changes
   - Recent commits

2. **Tool Usage Widget**
   - List of tools used in session
   - Frequency and success rate

3. **Performance Metrics Widget**
   - Response times
   - Token usage
   - API call statistics

4. **Keyboard Shortcuts Widget**
   - Quick reference for all shortcuts
   - Context-sensitive help

5. **Custom Widgets**
   - Plugin system for user-defined widgets
   - Configuration-based widget ordering
