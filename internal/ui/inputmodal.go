package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// InputSubmitMsg is sent when the user submits input
type InputSubmitMsg struct {
	Input string
}

// InputCancelMsg is sent when the user cancels input
type InputCancelMsg struct{}

// InputModal displays a prompt and text input for user input
type InputModal struct {
	width     int
	height    int
	prompt    string
	input     textinput.Model
	sensitive bool // Whether this is sensitive input (masked)
}

// NewInputModal creates a new input modal
func NewInputModal() *InputModal {
	ti := textinput.New()
	ti.Placeholder = "Enter your response..."
	ti.CharLimit = 1024
	ti.Width = 60

	return &InputModal{
		input: ti,
	}
}

// SetSize sets the dimensions of the input modal
func (m *InputModal) SetSize(width, height int) {
	m.width = width
	m.height = height

	// Adjust input width based on modal width
	inputWidth := width - 20
	if inputWidth > 80 {
		inputWidth = 80
	}
	if inputWidth < 20 {
		inputWidth = 20
	}
	m.input.Width = inputWidth
}

// Show displays the input modal with a prompt
func (m *InputModal) Show(prompt string, sensitive bool) tea.Cmd {
	m.prompt = prompt
	m.sensitive = sensitive
	m.input.Reset()

	if sensitive {
		m.input.EchoMode = textinput.EchoPassword
		m.input.EchoCharacter = '•'
		m.input.Placeholder = "Enter sensitive input..."
	} else {
		m.input.EchoMode = textinput.EchoNormal
		m.input.Placeholder = "Enter your response..."
	}

	return m.input.Focus()
}

// Update handles input events
func (m *InputModal) Update(msg tea.Msg) (*InputModal, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Submit the input
			return m, func() tea.Msg {
				return InputSubmitMsg{Input: m.input.Value()}
			}
		case "esc", "ctrl+c":
			// Cancel input
			return m, func() tea.Msg {
				return InputCancelMsg{}
			}
		}
	}

	// Update the text input
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// View renders the input modal
func (m *InputModal) View() string {
	// Define styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		Background(lipgloss.Color("236")).
		Padding(0, 2).
		Width(m.width)

	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("11")).
		Bold(true).
		MarginBottom(1)

	containerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("12")).
		Padding(1, 2).
		Width(m.width - 4)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Italic(true).
		MarginTop(1)

	sensitiveWarningStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Italic(true).
		MarginBottom(1)

	// Build content
	var content strings.Builder

	title := "INPUT REQUIRED"
	if m.sensitive {
		title = "SENSITIVE INPUT REQUIRED"
	}
	content.WriteString(titleStyle.Render(title))
	content.WriteString("\n\n")

	// Show prompt
	if m.prompt != "" {
		content.WriteString(promptStyle.Render(m.prompt))
		content.WriteString("\n")
	}

	// Show warning for sensitive input
	if m.sensitive {
		content.WriteString(sensitiveWarningStyle.Render("⚠ Input will be masked"))
		content.WriteString("\n")
	}

	// Show input field
	content.WriteString(m.input.View())
	content.WriteString("\n")

	// Show help text
	content.WriteString(helpStyle.Render("Press Enter to submit • Esc to cancel"))

	return containerStyle.Render(content.String())
}

// IsFocused returns whether the input modal is currently focused
func (m *InputModal) IsFocused() bool {
	return m.input.Focused()
}
