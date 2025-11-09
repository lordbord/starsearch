package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"starsearch/internal/types"
)

// HistoryModal displays browsing history for viewing and navigation
type HistoryModal struct {
	visible      bool
	history      []types.HistoryEntry
	filtered     []types.HistoryEntry
	searchQuery  string
	selectedIdx  int
	width        int
	height       int
	scrollOffset int
}

// HistorySelectedMsg is sent when a history entry is selected to navigate to
type HistorySelectedMsg struct {
	URL string
}

func NewHistoryModal() *HistoryModal {
	return &HistoryModal{
		visible:      false,
		history:      []types.HistoryEntry{},
		filtered:     []types.HistoryEntry{},
		searchQuery:  "",
		selectedIdx:  0,
		scrollOffset: 0,
	}
}

func (m *HistoryModal) Show(history []types.HistoryEntry) {
	m.visible = true
	m.history = history
	m.searchQuery = ""
	m.filter()
	m.selectedIdx = 0
	m.scrollOffset = 0
}

func (m *HistoryModal) Hide() {
	m.visible = false
	m.searchQuery = ""
	m.filtered = []types.HistoryEntry{}
}

func (m *HistoryModal) IsVisible() bool {
	return m.visible
}

func (m *HistoryModal) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// filter filters history based on search query
func (m *HistoryModal) filter() {
	if m.searchQuery == "" {
		m.filtered = make([]types.HistoryEntry, len(m.history))
		copy(m.filtered, m.history)
		// Reverse to show newest first
		for i, j := 0, len(m.filtered)-1; i < j; i, j = i+1, j-1 {
			m.filtered[i], m.filtered[j] = m.filtered[j], m.filtered[i]
		}
		return
	}

	query := strings.ToLower(m.searchQuery)
	m.filtered = []types.HistoryEntry{}
	for _, entry := range m.history {
		if strings.Contains(strings.ToLower(entry.URL), query) ||
			strings.Contains(strings.ToLower(entry.Title), query) {
			m.filtered = append(m.filtered, entry)
		}
	}
	// Reverse to show newest first
	for i, j := 0, len(m.filtered)-1; i < j; i, j = i+1, j-1 {
		m.filtered[i], m.filtered[j] = m.filtered[j], m.filtered[i]
	}
}

func (m *HistoryModal) Update(msg tea.Msg) (*HistoryModal, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc", "ctrl+c", "ctrl+h"))):
			m.Hide()
			return m, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("j", "down"))):
			if m.selectedIdx < len(m.filtered)-1 {
				m.selectedIdx++
				m.adjustScroll()
			}
			return m, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("k", "up"))):
			if m.selectedIdx > 0 {
				m.selectedIdx--
				m.adjustScroll()
			}
			return m, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("g"))):
			m.selectedIdx = 0
			m.scrollOffset = 0
			return m, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("G"))):
			if len(m.filtered) > 0 {
				m.selectedIdx = len(m.filtered) - 1
				m.adjustScroll()
			}
			return m, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if m.selectedIdx < len(m.filtered) {
				url := m.filtered[m.selectedIdx].URL
				m.Hide()
				return m, func() tea.Msg {
					return HistorySelectedMsg{URL: url}
				}
			}
			return m, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("/"))):
			// Start search mode - for now, just clear search
			// In a full implementation, you'd want a search input field
			m.searchQuery = ""
			m.filter()
			m.selectedIdx = 0
			m.scrollOffset = 0
			return m, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("backspace"))):
			// Handle backspace to remove last character from search
			if len(m.searchQuery) > 0 {
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				m.filter()
				m.selectedIdx = 0
				m.scrollOffset = 0
			}
			return m, nil

		default:
			// Handle typing for search
			if len(msg.Runes) > 0 {
				m.searchQuery += string(msg.Runes)
				m.filter()
				m.selectedIdx = 0
				m.scrollOffset = 0
				return m, nil
			}
		}

	case tea.MouseMsg:
		if msg.Type == tea.MouseWheelUp {
			// Scroll up with mouse wheel
			if m.scrollOffset > 0 {
				m.scrollOffset--
				if m.selectedIdx > 0 {
					m.selectedIdx--
				}
			} else if m.selectedIdx > 0 {
				m.selectedIdx--
				m.adjustScroll()
			}
			return m, nil
		}

		if msg.Type == tea.MouseWheelDown {
			// Scroll down with mouse wheel
			modalHeight := m.height - 6
			if modalHeight < 10 {
				modalHeight = 10
			}
			visibleHeight := modalHeight - 6
			if visibleHeight < 1 {
				visibleHeight = 1
			}

			// Calculate visible entries (each entry is 3 lines)
			visibleEntries := visibleHeight / 3
			if visibleEntries < 1 {
				visibleEntries = 1
			}

			maxScroll := len(m.filtered) - visibleEntries
			if maxScroll < 0 {
				maxScroll = 0
			}

			if m.scrollOffset < maxScroll {
				m.scrollOffset++
				if m.selectedIdx < len(m.filtered)-1 {
					m.selectedIdx++
				}
			} else if m.selectedIdx < len(m.filtered)-1 {
				m.selectedIdx++
				m.adjustScroll()
			}
			return m, nil
		}

		if msg.Type == tea.MouseLeft && len(m.filtered) > 0 {
			// Similar mouse handling as bookmarks modal
			modalWidth := m.width - 6
			if modalWidth < 60 {
				modalWidth = 60
			}
			if modalWidth > m.width - 4 {
				modalWidth = m.width - 4
			}

			modalHeight := m.height - 6
			if modalHeight < 10 {
				modalHeight = 10
			}

			visibleHeight := modalHeight - 6
			if visibleHeight < 1 {
				visibleHeight = 1
			}

			contentWidth := modalWidth + 6
			leftPadding := (m.width - contentWidth) / 2
			if leftPadding < 0 {
				leftPadding = 0
			}

			// Calculate modal position
			contentHeight := modalHeight
			topPadding := (m.height - contentHeight) / 2
			if topPadding < 0 {
				topPadding = 0
			}

			modalTop := topPadding + 1 // Account for border
			modalLeft := leftPadding + 1 // Account for border

			// Check if click is within modal bounds
			if msg.X >= modalLeft && msg.X < modalLeft+modalWidth &&
				msg.Y >= modalTop && msg.Y < modalTop+modalHeight {
				// Click is in modal - calculate which entry was clicked
				// Account for title (1 line) and help text (1 line) and padding
				clickY := msg.Y - modalTop - 3
				if clickY >= 0 {
					// Each entry is 3 lines (title, URL, timestamp)
					clickedIdx := m.scrollOffset + (clickY / 3)
					if clickedIdx >= 0 && clickedIdx < len(m.filtered) {
						m.selectedIdx = clickedIdx
						m.adjustScroll()
					}
				}
			}
		}
	}

	return m, nil
}

func (m *HistoryModal) adjustScroll() {
	modalHeight := m.height - 6
	if modalHeight < 10 {
		modalHeight = 10
	}
	// Account for title (1 line), help text (1 line), and padding (4 lines total)
	visibleHeight := modalHeight - 6
	if visibleHeight < 1 {
		visibleHeight = 1
	}

	// Each entry takes 3 lines (title, URL, timestamp)
	// So we can show visibleHeight/3 entries
	visibleEntries := visibleHeight / 3
	if visibleEntries < 1 {
		visibleEntries = 1
	}

	// Scroll down if selected item is below visible area
	if m.selectedIdx >= m.scrollOffset+visibleEntries {
		m.scrollOffset = m.selectedIdx - visibleEntries + 1
	}

	// Scroll up if selected item is above visible area
	if m.selectedIdx < m.scrollOffset {
		m.scrollOffset = m.selectedIdx
	}

	// Ensure scroll offset doesn't go negative
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}

	// Ensure we don't scroll past the end
	maxScroll := len(m.filtered) - visibleEntries
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.scrollOffset > maxScroll {
		m.scrollOffset = maxScroll
	}
}

func (m *HistoryModal) View() string {
	if !m.visible {
		return ""
	}

	// Use more of the available width (wider modal)
	modalWidth := m.width - 6
	if modalWidth < 60 {
		modalWidth = 60
	}
	if modalWidth > m.width - 4 {
		modalWidth = m.width - 4
	}

	// Use less height to make it shorter
	modalHeight := m.height - 6
	if modalHeight < 10 {
		modalHeight = 10
	}
	if modalHeight > m.height - 6 {
		modalHeight = m.height - 6
	}

	visibleHeight := modalHeight - 6 // Less padding for more content
	if visibleHeight < 1 {
		visibleHeight = 1
	}

	// Calculate how many entries we can show (each entry is 3 lines)
	visibleEntries := visibleHeight / 3
	if visibleEntries < 1 {
		visibleEntries = 1
	}

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("12")).
		Padding(0, 1).
		Width(modalWidth - 4)

	title := "History"
	if m.searchQuery != "" {
		title += fmt.Sprintf(" (filter: %s)", m.searchQuery)
	}
	titleText := titleStyle.Render(title)

	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Padding(0, 1).
		Width(modalWidth - 4)

	helpText := helpStyle.Render("Enter: Navigate | Esc/Ctrl+C: Close | /: Search | Mouse: Scroll")

	// History entries - show entries starting from scrollOffset
	var entries []string
	startIdx := m.scrollOffset
	endIdx := startIdx + visibleEntries
	if endIdx > len(m.filtered) {
		endIdx = len(m.filtered)
	}

	for i := startIdx; i < endIdx; i++ {
		entry := m.filtered[i]
		isSelected := i == m.selectedIdx

		// Format timestamp
		timestamp := time.Unix(entry.Timestamp, 0)
		timeStr := timestamp.Format("2006-01-02 15:04")

		var style lipgloss.Style
		if isSelected {
			style = lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				Background(lipgloss.Color("12")).
				Padding(0, 1).
				Width(modalWidth - 4)
		} else {
			style = lipgloss.NewStyle().
				Foreground(lipgloss.Color("7")).
				Padding(0, 1).
				Width(modalWidth - 4)
		}

		// Truncate title and URL if needed - use more width
		title := entry.Title
		maxTitleLen := modalWidth - 10
		if len(title) > maxTitleLen {
			title = title[:maxTitleLen-3] + "..."
		}
		url := entry.URL
		maxURLLen := modalWidth - 10
		if len(url) > maxURLLen {
			url = url[:maxURLLen-3] + "..."
		}

		entryText := fmt.Sprintf("%s\n  %s\n  %s", title, url, timeStr)
		entries = append(entries, style.Render(entryText))
	}

	if len(entries) == 0 {
		noResultsStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Padding(0, 1).
			Width(modalWidth - 4)
		entries = append(entries, noResultsStyle.Render("No history entries found"))
	}

	// Combine all parts
	content := strings.Join([]string{titleText, helpText, strings.Join(entries, "\n")}, "\n")

	// Wrap in border
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("12")).
		Padding(1, 2).
		Width(modalWidth).
		MaxHeight(modalHeight)

	modalContent := borderStyle.Render(content)

	// Center the modal
	contentLines := strings.Split(modalContent, "\n")
	contentHeight := len(contentLines)
	if contentHeight > modalHeight {
		contentHeight = modalHeight
		// Truncate if too tall
		modalContent = strings.Join(contentLines[:modalHeight], "\n")
	}
	
	contentWidth := modalWidth + 6 // Account for border and padding

	topPadding := (m.height - contentHeight) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	leftPadding := (m.width - contentWidth) / 2
	if leftPadding < 0 {
		leftPadding = 0
	}

	// Add padding
	result := strings.Repeat("\n", topPadding)
	for _, line := range strings.Split(modalContent, "\n") {
		result += strings.Repeat(" ", leftPadding) + line + "\n"
	}

	return result
}

