package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// HelpModal displays keyboard shortcuts and commands
type HelpModal struct {
	width  int
	height int
}

// NewHelpModal creates a new help modal
func NewHelpModal() *HelpModal {
	return &HelpModal{}
}

// SetSize sets the dimensions of the help modal
func (h *HelpModal) SetSize(width, height int) {
	h.width = width
	h.height = height
}

// View renders the help modal
func (h *HelpModal) View() string {
	// Define styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		Background(lipgloss.Color("236")).
		Padding(0, 2).
		Width(h.width)

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("11")).
		MarginTop(1)

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		Width(20)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7"))

	containerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("12")).
		Padding(1, 2).
		Width(h.width - 4)

	dismissStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Italic(true).
		MarginTop(1)

	// Build help content
	var content strings.Builder

	content.WriteString(titleStyle.Render("STARSEARCH KEYBOARD SHORTCUTS"))
	content.WriteString("\n\n")

	// Navigation commands
	content.WriteString(headerStyle.Render("Navigation"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render("Ctrl+L") + descStyle.Render("Focus address bar"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render("G") + descStyle.Render("Enter link number mode"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render("0-9") + descStyle.Render("Input link number (in link mode)"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render("Enter") + descStyle.Render("Navigate to link/URL"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render("R") + descStyle.Render("Reload current page"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render("H / ← / Alt+←") + descStyle.Render("Go back in history"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render("L / → / Alt+→") + descStyle.Render("Go forward in history"))
	content.WriteString("\n\n")

	// Scrolling commands
	content.WriteString(headerStyle.Render("Scrolling"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render("J / ↓") + descStyle.Render("Scroll down"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render("K / ↑") + descStyle.Render("Scroll up"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render("PgDown / Space") + descStyle.Render("Page down"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render("PgUp") + descStyle.Render("Page up"))
	content.WriteString("\n\n")

	// Other commands
	content.WriteString(headerStyle.Render("Other"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render("D") + descStyle.Render("Toggle bookmark"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render("?") + descStyle.Render("Show this help"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render("Esc") + descStyle.Render("Exit link mode / Close help"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render("Q / Ctrl+C") + descStyle.Render("Quit"))
	content.WriteString("\n")

	content.WriteString(dismissStyle.Render("\nPress Esc or Q to close this help"))

	return containerStyle.Render(content.String())
}
