package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AddressBar represents the URL input bar
type AddressBar struct {
	input    textinput.Model
	focused  bool
	width    int
}

// NewAddressBar creates a new address bar
func NewAddressBar() *AddressBar {
	ti := textinput.New()
	ti.Placeholder = "gemini://..."
	ti.CharLimit = 1024
	ti.Width = 50

	return &AddressBar{
		input:   ti,
		focused: false,
		width:   80,
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
			switch msg.String() {
			case "enter":
				// Return the URL and unfocus
				url := a.input.Value()
				a.focused = false
				a.input.Blur()
				if url != "" {
					return a, func() tea.Msg { return NavigateMsg{URL: url} }
				}
				return a, nil
			case "esc":
				// Cancel input
				a.focused = false
				a.input.Blur()
				return a, nil
			}
		}
	}

	// Update the input if focused
	if a.focused {
		a.input, cmd = a.input.Update(msg)
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

	return style.Width(a.width).Render(a.input.View())
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

// NavigateMsg is sent when the user wants to navigate to a URL
type NavigateMsg struct {
	URL string
}
