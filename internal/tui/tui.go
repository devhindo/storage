package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/devhindo/storage/internal/drive"
)

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99")).Padding(0, 1)
	folderStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	fileStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	selectedStyle = lipgloss.NewStyle().Background(lipgloss.Color("237")).Bold(true)
	statusStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Padding(1, 0, 0, 1)
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)

// Model is the Bubble Tea model for the file browser.
type Model struct {
	srv     *drive.Service
	entries []drive.FileEntry
	cursor  int
	path    []breadcrumb // navigation stack
	loading bool
	err     error
	width   int
	height  int
}

type breadcrumb struct {
	id   string
	name string
}

// Messages
type fetchedMsg struct {
	entries []drive.FileEntry
}

type errMsg struct {
	err error
}

func (e errMsg) Error() string { return e.err.Error() }

// Run starts the TUI application.
func Run(srv *drive.Service) error {
	m := Model{
		srv:     srv,
		loading: true,
		path:    []breadcrumb{{id: "root", name: "My Drive"}},
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func (m Model) Init() tea.Cmd {
	return m.fetchFolder("root")
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case fetchedMsg:
		m.entries = msg.entries
		m.cursor = 0
		m.loading = false
		m.err = nil
		return m, nil

	case errMsg:
		m.err = msg.err
		m.loading = false
		return m, nil
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < len(m.entries)-1 {
			m.cursor++
		}

	case "enter", "l":
		if m.cursor < len(m.entries) {
			entry := m.entries[m.cursor]
			if entry.IsFolder {
				m.path = append(m.path, breadcrumb{id: entry.ID, name: entry.Name})
				m.loading = true
				m.entries = nil
				return m, m.fetchFolder(entry.ID)
			}
			// TODO: open/download file
		}

	case "backspace", "h":
		if len(m.path) > 1 {
			m.path = m.path[:len(m.path)-1]
			parent := m.path[len(m.path)-1]
			m.loading = true
			m.entries = nil
			return m, m.fetchFolder(parent.id)
		}
	}

	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	// Title / breadcrumb
	var pathParts []string
	for _, bc := range m.path {
		pathParts = append(pathParts, bc.name)
	}
	b.WriteString(titleStyle.Render("storage > "+strings.Join(pathParts, " / ")) + "\n\n")

	if m.loading {
		b.WriteString("  Loading...\n")
		return b.String()
	}

	if m.err != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("  Error: %v", m.err)) + "\n")
		return b.String()
	}

	if len(m.entries) == 0 {
		b.WriteString("  (empty folder)\n")
		return b.String()
	}

	// File list
	visibleLines := m.height - 5
	if visibleLines < 1 {
		visibleLines = 10
	}

	start := 0
	if m.cursor >= visibleLines {
		start = m.cursor - visibleLines + 1
	}

	for i := start; i < len(m.entries) && i < start+visibleLines; i++ {
		entry := m.entries[i]
		line := formatEntry(entry)

		if i == m.cursor {
			line = selectedStyle.Render("> " + line)
		} else {
			line = "  " + line
		}

		b.WriteString(line + "\n")
	}

	// Status bar
	status := fmt.Sprintf("%d items | j/k: navigate | enter/l: open | backspace/h: back | q: quit", len(m.entries))
	b.WriteString(statusStyle.Render(status))

	return b.String()
}

func formatEntry(e drive.FileEntry) string {
	if e.IsFolder {
		return folderStyle.Render("📁 " + e.Name)
	}
	size := formatSize(e.Size)
	return fileStyle.Render(fmt.Sprintf("   %s  %s", e.Name, size))
}

func formatSize(bytes int64) string {
	if bytes == 0 {
		return ""
	}
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func (m Model) fetchFolder(folderID string) tea.Cmd {
	return func() tea.Msg {
		entries, err := m.srv.ListFolder(folderID)
		if err != nil {
			return errMsg{err: err}
		}
		return fetchedMsg{entries: entries}
	}
}
