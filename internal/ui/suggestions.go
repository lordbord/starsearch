package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"starsearch/internal/types"
)

// SuggestionType represents the type of a suggestion
type SuggestionType int

const (
	SuggestionHistory SuggestionType = iota
	SuggestionBookmark
	SuggestionURL
)

// Suggestion represents a single suggestion
type Suggestion struct {
	Text string
	URL  string
	Type SuggestionType
}

// Suggestions displays autocomplete suggestions
type Suggestions struct {
	suggestions []Suggestion
	selectedIdx int
	visible     bool
	width       int
	maxVisible  int
}

// NewSuggestions creates a new suggestions component
func NewSuggestions() *Suggestions {
	return &Suggestions{
		suggestions: []Suggestion{},
		selectedIdx: 0,
		visible:     false,
		maxVisible:  6,
	}
}

// Show displays suggestions
func (s *Suggestions) Show(suggestions []Suggestion) {
	s.suggestions = suggestions
	s.selectedIdx = 0
	s.visible = len(suggestions) > 0
	if s.selectedIdx >= len(s.suggestions) {
		s.selectedIdx = 0
	}
}

// Hide hides the suggestions
func (s *Suggestions) Hide() {
	s.visible = false
	s.suggestions = []Suggestion{}
	s.selectedIdx = 0
}

// IsVisible returns whether suggestions are visible
func (s *Suggestions) IsVisible() bool {
	return s.visible && len(s.suggestions) > 0
}

// SetWidth sets the width of the suggestions dropdown
func (s *Suggestions) SetWidth(width int) {
	s.width = width
}

// Update handles suggestions updates
func (s *Suggestions) Update(msg tea.Msg) (*Suggestions, tea.Cmd) {
	if !s.visible {
		return s, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "ctrl+p":
			if s.selectedIdx > 0 {
				s.selectedIdx--
			} else {
				s.selectedIdx = len(s.suggestions) - 1
			}
			return s, nil
		case "down", "ctrl+n", "tab":
			if s.selectedIdx < len(s.suggestions)-1 {
				s.selectedIdx++
			} else {
				s.selectedIdx = 0
			}
			return s, nil
		case "enter":
			if s.selectedIdx >= 0 && s.selectedIdx < len(s.suggestions) {
				return s, func() tea.Msg {
					return SuggestionSelectedMsg{
						URL: s.suggestions[s.selectedIdx].URL,
					}
				}
			}
		case "esc":
			s.Hide()
			return s, nil
		}
	}

	return s, nil
}

// View renders the suggestions dropdown
func (s *Suggestions) View() string {
	if !s.visible || len(s.suggestions) == 0 {
		return ""
	}

	var lines []string

	// Determine how many suggestions to show
	visibleCount := len(s.suggestions)
	if visibleCount > s.maxVisible {
		visibleCount = s.maxVisible
	}

	// Calculate start index for scrolling
	startIdx := s.selectedIdx - visibleCount/2
	if startIdx < 0 {
		startIdx = 0
	}
	if startIdx+visibleCount > len(s.suggestions) {
		startIdx = len(s.suggestions) - visibleCount
	}

	// Render visible suggestions
	for i := startIdx; i < startIdx+visibleCount && i < len(s.suggestions); i++ {
		suggestion := s.suggestions[i]
		isSelected := i == s.selectedIdx

		var prefix string
		var style lipgloss.Style

		switch suggestion.Type {
		case SuggestionHistory:
			prefix = "H "
			if isSelected {
				style = lipgloss.NewStyle().
					Foreground(lipgloss.Color("15")).
					Background(lipgloss.Color("12")).
					Width(s.width - 2).
					Padding(0, 1)
			} else {
				style = lipgloss.NewStyle().
					Foreground(lipgloss.Color("8")).
					Width(s.width - 2).
					Padding(0, 1)
			}
		case SuggestionBookmark:
			prefix = "★ "
			if isSelected {
				style = lipgloss.NewStyle().
					Foreground(lipgloss.Color("15")).
					Background(lipgloss.Color("12")).
					Width(s.width - 2).
					Padding(0, 1)
			} else {
				style = lipgloss.NewStyle().
					Foreground(lipgloss.Color("11")).
					Width(s.width - 2).
					Padding(0, 1)
			}
		default:
			prefix = "  "
			if isSelected {
				style = lipgloss.NewStyle().
					Foreground(lipgloss.Color("15")).
					Background(lipgloss.Color("12")).
					Width(s.width - 2).
					Padding(0, 1)
			} else {
				style = lipgloss.NewStyle().
					Foreground(lipgloss.Color("7")).
					Width(s.width - 2).
					Padding(0, 1)
			}
		}

		// Truncate text if too long
		text := suggestion.Text
		maxTextLen := s.width - len(prefix) - 4
		if len(text) > maxTextLen {
			text = text[:maxTextLen-3] + "..."
		}

		lines = append(lines, style.Render(prefix+text))
	}

	if len(lines) == 0 {
		return ""
	}

	// Wrap in border
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Width(s.width)

	return borderStyle.Render(strings.Join(lines, "\n"))
}

// GetSelected returns the currently selected suggestion
func (s *Suggestions) GetSelected() *Suggestion {
	if s.selectedIdx >= 0 && s.selectedIdx < len(s.suggestions) {
		return &s.suggestions[s.selectedIdx]
	}
	return nil
}

// SuggestionSelectedMsg is sent when a suggestion is selected
type SuggestionSelectedMsg struct {
	URL string
}

// FilterSuggestions filters suggestions based on query
func FilterSuggestions(query string, history []types.HistoryEntry, bookmarks []types.Bookmark) []Suggestion {
	query = strings.ToLower(query)
	suggestions := []Suggestion{}

	// Add matching history entries (most recent first)
	historyCount := 0
	for i := len(history) - 1; i >= 0 && historyCount < 5; i-- {
		entry := history[i]
		if strings.Contains(strings.ToLower(entry.URL), query) ||
			strings.Contains(strings.ToLower(entry.Title), query) {
			suggestions = append(suggestions, Suggestion{
				Text: entry.Title,
				URL:  entry.URL,
				Type: SuggestionHistory,
			})
			historyCount++
		}
	}

	// Add matching bookmarks
	bookmarkCount := 0
	for _, bookmark := range bookmarks {
		if bookmarkCount >= 3 {
			break
		}
		if strings.Contains(strings.ToLower(bookmark.URL), query) ||
			strings.Contains(strings.ToLower(bookmark.Title), query) {
			suggestions = append(suggestions, Suggestion{
				Text: bookmark.Title,
				URL:  bookmark.URL,
				Type: SuggestionBookmark,
			})
			bookmarkCount++
		}
	}

	return suggestions
}

