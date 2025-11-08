package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"starsearch/internal/types"
)

// TabBar displays and manages browser tabs
type TabBar struct {
	tabs         []types.Tab
	activeIdx    int
	width        int
	height       int
	scrollOffset int
}

// TabSwitchMsg is sent when user switches tabs
type TabSwitchMsg struct {
	Index int
}

// TabCloseMsg is sent when user closes a tab
type TabCloseMsg struct {
	Index int
}

// TabNewMsg is sent when user creates a new tab
type TabNewMsg struct{}

func NewTabBar() *TabBar {
	return &TabBar{
		tabs:         []types.Tab{},
		activeIdx:    0,
		scrollOffset: 0,
	}
}

func (t *TabBar) AddTab(url, title string) {
	tab := types.Tab{
		ID:       len(t.tabs),
		Title:    title,
		URL:      url,
		Document: nil,
		Scroll:   0,
	}

	t.tabs = append(t.tabs, tab)
	t.activeIdx = len(t.tabs) - 1
	t.adjustScroll()
}

func (t *TabBar) CloseTab(index int) {
	if index < 0 || index >= len(t.tabs) {
		return
	}

	// Remove tab
	t.tabs = append(t.tabs[:index], t.tabs[index+1:]...)

	// Adjust active index
	if t.activeIdx >= index {
		t.activeIdx--
	}
	if t.activeIdx < 0 && len(t.tabs) > 0 {
		t.activeIdx = 0
	}

	// Adjust tab IDs
	for i := range t.tabs {
		t.tabs[i].ID = i
	}

	t.adjustScroll()
}

func (t *TabBar) SwitchTab(index int) {
	if index >= 0 && index < len(t.tabs) {
		t.activeIdx = index
		t.adjustScroll()
	}
}

func (t *TabBar) GetActiveTab() *types.Tab {
	if t.activeIdx >= 0 && t.activeIdx < len(t.tabs) {
		return &t.tabs[t.activeIdx]
	}
	return nil
}

func (t *TabBar) GetActiveIndex() int {
	return t.activeIdx
}

func (t *TabBar) GetTabs() []types.Tab {
	return t.tabs
}

func (t *TabBar) UpdateTab(index int, url, title string, document *types.Document, scroll int) {
	if index >= 0 && index < len(t.tabs) {
		t.tabs[index].URL = url
		t.tabs[index].Title = title
		t.tabs[index].Document = document
		t.tabs[index].Scroll = scroll
	}
}

func (t *TabBar) SetSize(width, height int) {
	t.width = width
	t.height = height
	t.adjustScroll()
}

func (t *TabBar) adjustScroll() {
	// Calculate total width needed
	totalWidth := 0
	for _, tab := range t.tabs {
		tabWidth := t.calculateTabWidth(tab)
		totalWidth += tabWidth + 1 // +1 for separator
	}

	// Adjust scroll offset to keep active tab visible
	if totalWidth > t.width {
		activeTabX := 0
		for i := 0; i < t.activeIdx; i++ {
			activeTabX += t.calculateTabWidth(t.tabs[i]) + 1
		}

		activeTabWidth := t.calculateTabWidth(t.tabs[t.activeIdx])

		// Scroll right if active tab is beyond visible area
		if activeTabX+activeTabWidth > t.scrollOffset+t.width {
			t.scrollOffset = activeTabX + activeTabWidth - t.width + 2
		}

		// Scroll left if active tab is before visible area
		if activeTabX < t.scrollOffset {
			t.scrollOffset = activeTabX
		}
	} else {
		t.scrollOffset = 0
	}
}

func (t *TabBar) calculateTabWidth(tab types.Tab) int {
	// Minimum width for tab content
	minWidth := 10

	// Calculate actual width needed
	title := tab.Title
	if title == "" {
		title = "Untitled"
	}

	// Add icon and padding
	width := len(title) + 4 // 2 for icon, 2 for padding

	if width < minWidth {
		width = minWidth
	}

	return width
}

func (t *TabBar) Update(msg tea.Msg) (*TabBar, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+tab"))):
			// Next tab
			if len(t.tabs) > 0 {
				nextIdx := (t.activeIdx + 1) % len(t.tabs)
				t.SwitchTab(nextIdx)
				return t, func() tea.Msg {
					return TabSwitchMsg{Index: nextIdx}
				}
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+shift+tab"))):
			// Previous tab
			if len(t.tabs) > 0 {
				prevIdx := t.activeIdx - 1
				if prevIdx < 0 {
					prevIdx = len(t.tabs) - 1
				}
				t.SwitchTab(prevIdx)
				return t, func() tea.Msg {
					return TabSwitchMsg{Index: prevIdx}
				}
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+w"))):
			// Close current tab
			if len(t.tabs) > 0 {
				closeIdx := t.activeIdx
				t.CloseTab(closeIdx)
				return t, func() tea.Msg {
					return TabCloseMsg{Index: closeIdx}
				}
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+t"))):
			// New tab
			t.AddTab("", "New Tab")
			return t, func() tea.Msg {
				return TabNewMsg{}
			}

		// Number keys 1-9 for direct tab switching
		case key.Matches(msg, key.NewBinding(key.WithKeys("1"))):
			if len(t.tabs) >= 1 {
				t.SwitchTab(0)
				return t, func() tea.Msg {
					return TabSwitchMsg{Index: 0}
				}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("2"))):
			if len(t.tabs) >= 2 {
				t.SwitchTab(1)
				return t, func() tea.Msg {
					return TabSwitchMsg{Index: 1}
				}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("3"))):
			if len(t.tabs) >= 3 {
				t.SwitchTab(2)
				return t, func() tea.Msg {
					return TabSwitchMsg{Index: 2}
				}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("4"))):
			if len(t.tabs) >= 4 {
				t.SwitchTab(3)
				return t, func() tea.Msg {
					return TabSwitchMsg{Index: 3}
				}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("5"))):
			if len(t.tabs) >= 5 {
				t.SwitchTab(4)
				return t, func() tea.Msg {
					return TabSwitchMsg{Index: 4}
				}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("6"))):
			if len(t.tabs) >= 6 {
				t.SwitchTab(5)
				return t, func() tea.Msg {
					return TabSwitchMsg{Index: 5}
				}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("7"))):
			if len(t.tabs) >= 7 {
				t.SwitchTab(6)
				return t, func() tea.Msg {
					return TabSwitchMsg{Index: 6}
				}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("8"))):
			if len(t.tabs) >= 8 {
				t.SwitchTab(7)
				return t, func() tea.Msg {
					return TabSwitchMsg{Index: 7}
				}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("9"))):
			if len(t.tabs) >= 9 {
				t.SwitchTab(8)
				return t, func() tea.Msg {
					return TabSwitchMsg{Index: 8}
				}
			}
		}

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
			// Calculate which tab was clicked
			x := msg.X
			currentX := -t.scrollOffset

			for i, tab := range t.tabs {
				tabWidth := t.calculateTabWidth(tab)
				if x >= currentX && x < currentX+tabWidth {
					t.SwitchTab(i)
					return t, func() tea.Msg {
						return TabSwitchMsg{Index: i}
					}
				}
				currentX += tabWidth + 1
			}
		}
	}

	return t, nil
}

func (t *TabBar) View() string {
	if len(t.tabs) == 0 {
		return ""
	}

	var b strings.Builder

	// Styles
	activeStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("12")).
		Foreground(lipgloss.Color("0")).
		Bold(true)

	inactiveStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("8")).
		Foreground(lipgloss.Color("15"))

	separatorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("7"))

	// Calculate visible range
	currentX := -t.scrollOffset
	visibleTabs := 0

	for i, tab := range t.tabs {
		tabWidth := t.calculateTabWidth(tab)

		// Check if tab is visible
		tabEndX := currentX + tabWidth
		if tabEndX > 0 && currentX < t.width {
			// Render tab
			title := tab.Title
			if title == "" {
				title = "Untitled"
			}

			// Truncate if too long
			maxTitleLen := tabWidth - 4 // Account for icon and padding
			if len(title) > maxTitleLen {
				title = title[:maxTitleLen-3] + "..."
			}

			// Add icon
			icon := "🌐"
			if i == t.activeIdx {
				icon = "🌍"
			}

			tabText := fmt.Sprintf(" %s %s ", icon, title)

			if i == t.activeIdx {
				b.WriteString(activeStyle.Render(tabText))
			} else {
				b.WriteString(inactiveStyle.Render(tabText))
			}

			visibleTabs++
		}

		currentX += tabWidth + 1 // +1 for separator

		// Add separator if not the last tab and if we're still in visible range
		if i < len(t.tabs)-1 && currentX-1 < t.width {
			b.WriteString(separatorStyle.Render("│"))
		}

		// Stop if we've exceeded width
		if currentX > t.width {
			break
		}
	}

	// Fill remaining space
	usedWidth := currentX
	if usedWidth < t.width {
		b.WriteString(strings.Repeat(" ", t.width-usedWidth))
	}

	return b.String()
}
