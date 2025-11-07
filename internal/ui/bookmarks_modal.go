package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"starsearch/internal/types"
)

// BookmarksModal displays a list of bookmarks for viewing and management
type BookmarksModal struct {
	visible      bool
	bookmarks    []types.Bookmark
	selectedIdx  int
	width        int
	height       int
	scrollOffset int
}

// BookmarkSelectedMsg is sent when a bookmark is selected to navigate to
type BookmarkSelectedMsg struct {
	URL string
}

// BookmarkDeleteMsg is sent when a bookmark should be deleted
type BookmarkDeleteMsg struct {
	URL string
}

func NewBookmarksModal() *BookmarksModal {
	return &BookmarksModal{
		visible:      false,
		bookmarks:    []types.Bookmark{},
		selectedIdx:  0,
		scrollOffset: 0,
	}
}

func (m *BookmarksModal) Show(bookmarks []types.Bookmark) {
	m.visible = true
	m.bookmarks = bookmarks
	m.selectedIdx = 0
	m.scrollOffset = 0
}

func (m *BookmarksModal) Hide() {
	m.visible = false
}

func (m *BookmarksModal) IsVisible() bool {
	return m.visible
}

func (m *BookmarksModal) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *BookmarksModal) Update(msg tea.Msg) (*BookmarksModal, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc", "q", "b"))):
			m.Hide()
			return m, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("j", "down"))):
			if m.selectedIdx < len(m.bookmarks)-1 {
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
			if len(m.bookmarks) > 0 {
				m.selectedIdx = len(m.bookmarks) - 1
				m.adjustScroll()
			}
			return m, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if m.selectedIdx < len(m.bookmarks) {
				url := m.bookmarks[m.selectedIdx].URL
				m.Hide()
				return m, func() tea.Msg {
					return BookmarkSelectedMsg{URL: url}
				}
			}
			return m, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("d", "delete"))):
			if m.selectedIdx < len(m.bookmarks) {
				url := m.bookmarks[m.selectedIdx].URL
				return m, func() tea.Msg {
					return BookmarkDeleteMsg{URL: url}
				}
			}
			return m, nil
		}

	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft && len(m.bookmarks) > 0 {
			// Calculate modal position and dimensions
			modalWidth := m.width - 4
			if modalWidth < 40 {
				modalWidth = 40
			}
			if modalWidth > 100 {
				modalWidth = 100
			}

			modalHeight := m.height - 4
			if modalHeight < 10 {
				modalHeight = 10
			}

			// Calculate visible height for bookmarks
			visibleHeight := modalHeight - 8
			if visibleHeight < 1 {
				visibleHeight = 1
			}

			// Calculate top padding (modal is centered)
			// Approximate content height - border (2) + padding (2) + title (2) + help (2) + bookmarks
			endIdx := m.scrollOffset + visibleHeight
			if endIdx > len(m.bookmarks) {
				endIdx = len(m.bookmarks)
			}
			visibleBookmarks := endIdx - m.scrollOffset
			contentHeight := 2 + 2 + 2 + 2 + (visibleBookmarks * 2)
			if m.scrollOffset > 0 {
				contentHeight++ // scroll indicator
			}
			if endIdx < len(m.bookmarks) {
				contentHeight++ // scroll indicator
			}

			topPadding := (m.height - contentHeight) / 2
			if topPadding < 0 {
				topPadding = 0
			}

			// Calculate where bookmarks start
			// topPadding + border (1) + padding (1) + title (2) + potential scroll indicator
			bookmarksStartY := topPadding + 1 + 1 + 2
			if m.scrollOffset > 0 {
				bookmarksStartY++ // scroll indicator above
			}

			// Check if click is within bookmark area
			if msg.Y >= bookmarksStartY {
				// Calculate which bookmark was clicked (each bookmark is 2 lines tall)
				relativeY := msg.Y - bookmarksStartY
				clickedIdx := m.scrollOffset + (relativeY / 2)

				// Check if the clicked index is valid
				if clickedIdx >= 0 && clickedIdx < len(m.bookmarks) {
					url := m.bookmarks[clickedIdx].URL
					m.Hide()
					return m, func() tea.Msg {
						return BookmarkSelectedMsg{URL: url}
					}
				}
			}
		}
	}

	return m, nil
}

func (m *BookmarksModal) adjustScroll() {
	// Calculate visible area (leave space for header and help text)
	visibleHeight := m.height - 8
	if visibleHeight < 1 {
		visibleHeight = 1
	}

	// Scroll down if selected item is below visible area
	if m.selectedIdx >= m.scrollOffset+visibleHeight {
		m.scrollOffset = m.selectedIdx - visibleHeight + 1
	}

	// Scroll up if selected item is above visible area
	if m.selectedIdx < m.scrollOffset {
		m.scrollOffset = m.selectedIdx
	}
}

func (m *BookmarksModal) View() string {
	if !m.visible {
		return ""
	}

	var b strings.Builder

	// Calculate modal dimensions
	modalWidth := m.width - 4
	if modalWidth < 40 {
		modalWidth = 40
	}
	if modalWidth > 100 {
		modalWidth = 100
	}

	modalHeight := m.height - 4
	if modalHeight < 10 {
		modalHeight = 10
	}

	// Styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		Width(modalWidth).
		Align(lipgloss.Center).
		MarginBottom(1)

	selectedStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("12")).
		Foreground(lipgloss.Color("0")).
		Bold(true).
		Width(modalWidth - 4)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Width(modalWidth - 4)

	emptyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Italic(true).
		Width(modalWidth).
		Align(lipgloss.Center).
		MarginTop(2).
		MarginBottom(2)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")).
		Width(modalWidth).
		Align(lipgloss.Center).
		MarginTop(1)

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("12")).
		Padding(1, 2).
		Width(modalWidth)

	// Build content
	b.WriteString(titleStyle.Render(fmt.Sprintf("Bookmarks (%d)", len(m.bookmarks))))
	b.WriteString("\n")

	if len(m.bookmarks) == 0 {
		b.WriteString(emptyStyle.Render("No bookmarks yet"))
		b.WriteString("\n")
		b.WriteString(emptyStyle.Render("Press 'd' on any page to add a bookmark"))
		b.WriteString("\n")
	} else {
		// Calculate visible range
		visibleHeight := modalHeight - 8
		if visibleHeight < 1 {
			visibleHeight = 1
		}

		endIdx := m.scrollOffset + visibleHeight
		if endIdx > len(m.bookmarks) {
			endIdx = len(m.bookmarks)
		}

		// Show scroll indicator if needed
		if m.scrollOffset > 0 {
			b.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("8")).
				Width(modalWidth - 4).
				Align(lipgloss.Center).
				Render("▲ more above ▲"))
			b.WriteString("\n")
		}

		// Render visible bookmarks
		for i := m.scrollOffset; i < endIdx; i++ {
			bookmark := m.bookmarks[i]

			// Truncate title if too long
			title := bookmark.Title
			if title == "" {
				title = "Untitled"
			}
			maxTitleLen := modalWidth - 10
			if len(title) > maxTitleLen {
				title = title[:maxTitleLen-3] + "..."
			}

			// Truncate URL if too long
			url := bookmark.URL
			maxURLLen := modalWidth - 10
			if len(url) > maxURLLen {
				url = url[:maxURLLen-3] + "..."
			}

			line := fmt.Sprintf("%s\n  %s", title, url)

			if i == m.selectedIdx {
				b.WriteString(selectedStyle.Render(line))
			} else {
				b.WriteString(normalStyle.Render(line))
			}
			b.WriteString("\n")
		}

		// Show scroll indicator if needed
		if endIdx < len(m.bookmarks) {
			b.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("8")).
				Width(modalWidth - 4).
				Align(lipgloss.Center).
				Render("▼ more below ▼"))
			b.WriteString("\n")
		}
	}

	// Help text
	helpText := "j/k: move • enter: open • d: delete • esc/q/b: close"
	b.WriteString(helpStyle.Render(helpText))

	// Wrap in border
	content := borderStyle.Render(b.String())

	// Center the modal
	contentHeight := strings.Count(content, "\n") + 1
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
	for _, line := range strings.Split(content, "\n") {
		result += strings.Repeat(" ", leftPadding) + line + "\n"
	}

	return result
}
