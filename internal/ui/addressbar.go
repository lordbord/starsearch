package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AddressBar represents the URL input bar
type AddressBar struct {
	input       textinput.Model
	focused     bool
	width       int
	suggestions *Suggestions
}

// NewAddressBar creates a new address bar
func NewAddressBar() *AddressBar {
	ti := textinput.New()
	ti.Placeholder = "gemini://..."
	ti.CharLimit = 1024
	ti.Width = 50

	return &AddressBar{
		input:       ti,
		focused:     false,
		width:       80,
		suggestions: NewSuggestions(),
	}
}

// Init initializes the address bar
func (a *AddressBar) Init() tea.Cmd {
	return nil
}

// Update handles address bar updates
func (a *AddressBar) Update(msg tea.Msg) (*AddressBar, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if a.focused {
			// Handle suggestion navigation first
			if a.suggestions.IsVisible() {
				var suggestionCmd tea.Cmd
				a.suggestions, suggestionCmd = a.suggestions.Update(msg)
				if suggestionCmd != nil {
					// Check if a suggestion was selected
					if selectedMsg, ok := suggestionCmd().(SuggestionSelectedMsg); ok {
						a.input.SetValue(selectedMsg.URL)
						a.suggestions.Hide()
						a.focused = false
						a.input.Blur()
						return a, func() tea.Msg { return NavigateMsg{URL: selectedMsg.URL} }
					}
					return a, suggestionCmd
				}
				// If suggestions handled the key, don't process further
				if msg.String() == "up" || msg.String() == "down" || msg.String() == "ctrl+p" || msg.String() == "ctrl+n" || msg.String() == "tab" {
					return a, nil
				}
			}

			switch msg.String() {
			case "enter":
				// Return the URL and unfocus
				url := a.input.Value()
				a.focused = false
				a.input.Blur()
				a.suggestions.Hide()
				if url != "" {
					return a, func() tea.Msg { return NavigateMsg{URL: url} }
				}
				return a, nil
			case "esc":
				// Cancel input
				a.focused = false
				a.input.Blur()
				a.suggestions.Hide()
				return a, nil
			case "ctrl+c":
				// Clear the address bar text
				a.input.SetValue("")
				a.suggestions.Hide()
				return a, nil
			}
		}

	case tea.MouseMsg:
		if a.focused && msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
			// Calculate position within the text input
			// Border takes 1 char on left, padding takes 1 char = 2 chars offset
			// The address bar is rendered with border and padding, so adjust coordinates
			borderOffset := 2 // 1 for border + 1 for padding

			// Create an adjusted mouse message with coordinates relative to the text input
			adjustedMsg := tea.MouseMsg{
				X:      msg.X - borderOffset,
				Y:      0, // Text input is single line
				Type:   msg.Type,
				Action: msg.Action,
				Button: msg.Button,
			}

			// Only pass through if click is within the input area
			if adjustedMsg.X >= 0 && adjustedMsg.X < a.input.Width {
				a.input, cmd = a.input.Update(adjustedMsg)
			}
			// Return early to avoid passing the unadjusted message to the input
			return a, cmd
		}
	}

	// Update the input if focused
	if a.focused {
		oldValue := a.input.Value()
		a.input, cmd = a.input.Update(msg)
		newValue := a.input.Value()
		
		// If value changed, update suggestions (will be handled by app)
		if oldValue != newValue {
			// Suggestions will be updated by the app based on new value
		}
	}

	return a, cmd
}

// View renders the address bar
func (a *AddressBar) View() string {
	var style lipgloss.Style

	if a.focused {
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("12")).
			Padding(0, 1)
	} else {
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("7")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("8")).
			Padding(0, 1)
	}

	// Adjust input width to available space
	a.input.Width = a.width - 4 // Account for border and padding

	addressBarView := style.Width(a.width).Render(a.input.View())
	
	// Add suggestions if visible
	if a.suggestions.IsVisible() {
		a.suggestions.SetWidth(a.width)
		suggestionsView := a.suggestions.View()
		if suggestionsView != "" {
			return addressBarView + "\n" + suggestionsView
		}
	}
	
	return addressBarView
}

// Focus sets focus on the address bar
func (a *AddressBar) Focus() tea.Cmd {
	a.focused = true
	return a.input.Focus()
}

// Blur removes focus from the address bar
func (a *AddressBar) Blur() {
	a.focused = false
	a.input.Blur()
}

// SetValue sets the address bar value
func (a *AddressBar) SetValue(url string) {
	a.input.SetValue(url)
}

// Value returns the current address bar value
func (a *AddressBar) Value() string {
	return a.input.Value()
}

// SetWidth sets the address bar width
func (a *AddressBar) SetWidth(width int) {
	a.width = width
}

// IsFocused returns whether the address bar is focused
func (a *AddressBar) IsFocused() bool {
	return a.focused
}

// UpdateSuggestions updates the suggestions based on query
func (a *AddressBar) UpdateSuggestions(suggestions []Suggestion) {
	if a.focused && len(suggestions) > 0 {
		a.suggestions.Show(suggestions)
	} else {
		a.suggestions.Hide()
	}
}

// GetSuggestions returns the suggestions component
func (a *AddressBar) GetSuggestions() *Suggestions {
	return a.suggestions
}

// NavigateMsg is sent when the user wants to navigate to a URL
type NavigateMsg struct {
	URL string
}
