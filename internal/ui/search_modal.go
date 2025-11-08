package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"starsearch/internal/types"
)

// SearchModal displays a search interface for finding text in documents
type SearchModal struct {
	visible       bool
	input         textinput.Model
	results       []types.SearchResult
	selectedIdx   int
	currentMatch  int
	width         int
	height        int
	document      *types.Document
	caseSensitive bool
}

// SearchSubmitMsg is sent when a search is submitted
type SearchSubmitMsg struct {
	Query         string
	CaseSensitive bool
}

// SearchNavigateMsg is sent when navigating between search results
type SearchNavigateMsg struct {
	Direction string // "next" or "prev"
}

// SearchCloseMsg is sent when the search modal is closed
type SearchCloseMsg struct{}

func NewSearchModal() *SearchModal {
	input := textinput.New()
	input.Placeholder = "Search in page..."
	input.Focus()
	input.Width = 40

	return &SearchModal{
		visible:       false,
		input:         input,
		results:       []types.SearchResult{},
		selectedIdx:   0,
		currentMatch:  -1,
		caseSensitive: false,
	}
}

func (m *SearchModal) Show(document *types.Document) tea.Cmd {
	m.visible = true
	m.document = document
	m.input.SetValue("")
	m.input.Focus()
	m.results = []types.SearchResult{}
	m.selectedIdx = 0
	m.currentMatch = -1
	return textinput.Blink
}

func (m *SearchModal) Hide() {
	m.visible = false
	m.input.Blur()
}

func (m *SearchModal) IsVisible() bool {
	return m.visible
}

func (m *SearchModal) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.input.Width = min(width-20, 60)
}

func (m *SearchModal) Update(msg tea.Msg) (*SearchModal, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
			m.Hide()
			return m, func() tea.Msg {
				return SearchCloseMsg{}
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			query := strings.TrimSpace(m.input.Value())
			if query != "" && m.document != nil {
				m.performSearch(query)
				if len(m.results) > 0 {
					m.currentMatch = 0
					return m, func() tea.Msg {
						return SearchSubmitMsg{
							Query:         query,
							CaseSensitive: m.caseSensitive,
						}
					}
				}
			}
			return m, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))):
			m.caseSensitive = !m.caseSensitive
			return m, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("n"))):
			if len(m.results) > 0 {
				m.currentMatch = (m.currentMatch + 1) % len(m.results)
				return m, func() tea.Msg {
					return SearchNavigateMsg{Direction: "next"}
				}
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("N"))):
			if len(m.results) > 0 {
				m.currentMatch = m.currentMatch - 1
				if m.currentMatch < 0 {
					m.currentMatch = len(m.results) - 1
				}
				return m, func() tea.Msg {
					return SearchNavigateMsg{Direction: "prev"}
				}
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("j", "down"))):
			if len(m.results) > 0 {
				m.selectedIdx = (m.selectedIdx + 1) % len(m.results)
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("k", "up"))):
			if len(m.results) > 0 {
				m.selectedIdx = m.selectedIdx - 1
				if m.selectedIdx < 0 {
					m.selectedIdx = len(m.results) - 1
				}
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter", "tab"))):
			if len(m.results) > 0 {
				m.currentMatch = m.selectedIdx
				return m, func() tea.Msg {
					return SearchNavigateMsg{Direction: "goto"}
				}
			}
		}
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *SearchModal) performSearch(query string) {
	m.results = []types.SearchResult{}

	if m.document == nil {
		return
	}

	searchText := query
	if !m.caseSensitive {
		searchText = strings.ToLower(query)
	}

	for lineIdx, line := range m.document.Lines {
		text := line.Text
		if !m.caseSensitive {
			text = strings.ToLower(text)
		}

		// Find all occurrences in this line
		start := 0
		for {
			idx := strings.Index(text[start:], searchText)
			if idx == -1 {
				break
			}

			absStart := start + idx
			absEnd := absStart + len(query)

			result := types.SearchResult{
				Line:     lineIdx,
				Start:    absStart,
				End:      absEnd,
				Text:     line.Text[absStart:absEnd],
				Selected: false,
			}

			m.results = append(m.results, result)
			start = absStart + 1
		}
	}
}

func (m *SearchModal) GetCurrentResult() *types.SearchResult {
	if m.currentMatch >= 0 && m.currentMatch < len(m.results) {
		return &m.results[m.currentMatch]
	}
	return nil
}

func (m *SearchModal) GetResults() []types.SearchResult {
	return m.results
}

func (m *SearchModal) View() string {
	if !m.visible {
		return ""
	}

	var b strings.Builder

	// Calculate modal dimensions
	modalWidth := min(m.width-4, 80)
	if modalWidth < 40 {
		modalWidth = 40
	}

	modalHeight := min(m.height-4, 20)
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

	inputStyle := lipgloss.NewStyle().
		Width(modalWidth - 4)

	caseStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Width(modalWidth - 4).
		MarginBottom(1)

	selectedStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("12")).
		Foreground(lipgloss.Color("0")).
		Bold(true).
		Width(modalWidth - 8)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Width(modalWidth - 8)

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
	b.WriteString(titleStyle.Render("Search in Page"))
	b.WriteString("\n")

	// Search input
	b.WriteString(inputStyle.Render(m.input.View()))
	b.WriteString("\n")

	// Case sensitive indicator
	caseText := "Case Sensitive: "
	if m.caseSensitive {
		caseText += "ON (Ctrl+C to toggle)"
	} else {
		caseText += "OFF (Ctrl+C to toggle)"
	}
	b.WriteString(caseStyle.Render(caseText))
	b.WriteString("\n")

	// Results
	if len(m.results) > 0 {
		b.WriteString(fmt.Sprintf("Found %d matches:\n", len(m.results)))

		// Show visible results
		visibleResults := modalHeight - 10
		if visibleResults < 1 {
			visibleResults = 1
		}

		startIdx := 0
		if m.selectedIdx >= visibleResults {
			startIdx = m.selectedIdx - visibleResults + 1
		}

		endIdx := startIdx + visibleResults
		if endIdx > len(m.results) {
			endIdx = len(m.results)
		}

		for i := startIdx; i < endIdx; i++ {
			result := m.results[i]

			// Get line text with context
			lineText := ""
			if result.Line < len(m.document.Lines) {
				lineText = m.document.Lines[result.Line].Text
				// Truncate if too long
				if len(lineText) > 50 {
					lineText = lineText[:47] + "..."
				}
			}

			matchText := fmt.Sprintf("Line %d: %s", result.Line+1, lineText)

			prefix := "  "
			if i == m.selectedIdx {
				prefix = "▶ "
			}
			if i == m.currentMatch {
				prefix = "● "
			}

			if i == m.selectedIdx {
				b.WriteString(selectedStyle.Render(prefix + matchText))
			} else {
				b.WriteString(normalStyle.Render(prefix + matchText))
			}
			b.WriteString("\n")
		}
	} else if m.input.Value() != "" {
		b.WriteString("No matches found")
	} else {
		b.WriteString("Enter search text above")
	}

	// Help text
	helpText := "j/k: move • enter: goto • n/N: next/prev • Ctrl+C: case toggle • esc: close"
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
