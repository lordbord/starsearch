package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// StatusBar displays status information at the bottom
type StatusBar struct {
	message       string
	url           string
	scrollPercent float64
	width         int
	isLoading     bool
	errorMsg      string
}

// NewStatusBar creates a new status bar
func NewStatusBar(width int) *StatusBar {
	return &StatusBar{
		width:   width,
		message: "Ready",
	}
}

// SetMessage sets the status message
func (s *StatusBar) SetMessage(msg string) {
	s.message = msg
	s.errorMsg = ""
}

// SetError sets an error message
func (s *StatusBar) SetError(err string) {
	s.errorMsg = err
	s.message = ""
}

// SetURL sets the current URL
func (s *StatusBar) SetURL(url string) {
	s.url = url
}

// SetScrollPercent sets the scroll percentage
func (s *StatusBar) SetScrollPercent(percent float64) {
	s.scrollPercent = percent
}

// SetLoading sets the loading state
func (s *StatusBar) SetLoading(loading bool) {
	s.isLoading = loading
}

// SetWidth sets the status bar width
func (s *StatusBar) SetWidth(width int) {
	s.width = width
}

// View renders the status bar
func (s *StatusBar) View() string {
	// Define styles
	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("237")).
		Padding(0, 1)

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("1")).
		Padding(0, 1).
		Bold(true)

	urlStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")).
		Background(lipgloss.Color("237"))

	scrollStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Background(lipgloss.Color("237"))

	// Build status line
	var leftSection string

	if s.errorMsg != "" {
		leftSection = errorStyle.Render(" ERROR: " + s.errorMsg + " ")
	} else if s.isLoading {
		leftSection = normalStyle.Render(" ⟳ Loading... ")
	} else {
		leftSection = normalStyle.Render(" " + s.message + " ")
	}

	// Middle section: URL (if available)
	middleSection := ""
	if s.url != "" {
		// Truncate URL if too long
		maxURLLen := s.width - lipgloss.Width(leftSection) - 20
		if maxURLLen < 20 {
			maxURLLen = 20
		}

		displayURL := s.url
		if len(displayURL) > maxURLLen {
			displayURL = displayURL[:maxURLLen-3] + "..."
		}

		middleSection = urlStyle.Render(" " + displayURL + " ")
	}

	// Right section: Scroll position
	scrollText := fmt.Sprintf("%.0f%%", s.scrollPercent*100)
	rightSection := scrollStyle.Render(" " + scrollText + " ")

	// Calculate spacing
	usedWidth := lipgloss.Width(leftSection) + lipgloss.Width(middleSection) + lipgloss.Width(rightSection)
	spacing := s.width - usedWidth

	if spacing < 0 {
		spacing = 0
	}

	spacer := lipgloss.NewStyle().
		Background(lipgloss.Color("237")).
		Width(spacing).
		Render("")

	// Combine sections
	statusLine := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftSection,
		middleSection,
		spacer,
		rightSection,
	)

	return statusLine
}
