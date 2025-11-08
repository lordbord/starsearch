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

// CertificateModal displays certificate information and management
type CertificateModal struct {
	visible      bool
	certificates []types.CertificateInfo
	selectedIdx  int
	width        int
	height       int
	scrollOffset int
}

// CertificateTrustMsg is sent when user trusts a certificate
type CertificateTrustMsg struct {
	Host string
}

// CertificateUntrustMsg is sent when user untrusts a certificate
type CertificateUntrustMsg struct {
	Host string
}

// CertificateCloseMsg is sent when certificate modal is closed
type CertificateCloseMsg struct{}

func NewCertificateModal() *CertificateModal {
	return &CertificateModal{
		visible:      false,
		certificates: []types.CertificateInfo{},
		selectedIdx:   0,
		scrollOffset:  0,
	}
}

func (m *CertificateModal) Show(certificates []types.CertificateInfo) tea.Cmd {
	m.visible = true
	m.certificates = certificates
	m.selectedIdx = 0
	m.scrollOffset = 0
	return nil
}

func (m *CertificateModal) Hide() {
	m.visible = false
}

func (m *CertificateModal) IsVisible() bool {
	return m.visible
}

func (m *CertificateModal) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *CertificateModal) Update(msg tea.Msg) (*CertificateModal, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc", "q", "c"))):
			m.Hide()
			return m, func() tea.Msg {
				return CertificateCloseMsg{}
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("j", "down"))):
			if m.selectedIdx < len(m.certificates)-1 {
				m.selectedIdx++
				m.adjustScroll()
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("k", "up"))):
			if m.selectedIdx > 0 {
				m.selectedIdx--
				m.adjustScroll()
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("g"))):
			m.selectedIdx = 0
			m.scrollOffset = 0

		case key.Matches(msg, key.NewBinding(key.WithKeys("G"))):
			if len(m.certificates) > 0 {
				m.selectedIdx = len(m.certificates) - 1
				m.adjustScroll()
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("t"))):
			if m.selectedIdx < len(m.certificates) {
				cert := m.certificates[m.selectedIdx]
				if !cert.Trusted {
					return m, func() tea.Msg {
						return CertificateTrustMsg{Host: cert.Host}
					}
				}
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("u", "delete"))):
			if m.selectedIdx < len(m.certificates) {
				cert := m.certificates[m.selectedIdx]
				if cert.Trusted {
					return m, func() tea.Msg {
						return CertificateUntrustMsg{Host: cert.Host}
					}
				}
			}
		}
	}

	return m, nil
}

func (m *CertificateModal) adjustScroll() {
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

func (m *CertificateModal) formatFingerprint(fp string) string {
	// Format SHA256 fingerprint in groups of 4 characters
	if len(fp) != 64 {
		return fp
	}

	var result strings.Builder
	for i := 0; i < len(fp); i += 4 {
		if i > 0 {
			result.WriteString(":")
		}
		end := i + 4
		if end > len(fp) {
			end = len(fp)
		}
		result.WriteString(fp[i:end])
	}
	return result.String()
}

func (m *CertificateModal) formatTime(t time.Time) string {
	if t.IsZero() {
		return "Unknown"
	}
	return t.Format("2006-01-02 15:04:05")
}

func (m *CertificateModal) View() string {
	if !m.visible {
		return ""
	}

	var b strings.Builder

	// Calculate modal dimensions
	modalWidth := min(m.width-4, 100)
	if modalWidth < 80 {
		modalWidth = 80
	}

	modalHeight := min(m.height-4, 25)
	if modalHeight < 15 {
		modalHeight = 15
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

	trustedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		Bold(true)

	untrustedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Bold(true)

	fieldStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Bold(true)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Width(modalWidth - 12)

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
	b.WriteString(titleStyle.Render(fmt.Sprintf("Certificate Manager (%d certificates)", len(m.certificates))))
	b.WriteString("\n")

	if len(m.certificates) == 0 {
		b.WriteString(normalStyle.Render("No certificates found"))
		b.WriteString("\n")
	} else {
		// Calculate visible range
		visibleHeight := modalHeight - 8
		if visibleHeight < 1 {
			visibleHeight = 1
		}

		endIdx := m.scrollOffset + visibleHeight
		if endIdx > len(m.certificates) {
			endIdx = len(m.certificates)
		}

		// Show scroll indicator if needed
		if m.scrollOffset > 0 {
			b.WriteString(normalStyle.Render("▲ more above ▲"))
			b.WriteString("\n")
		}

		// Render visible certificates
		for i := m.scrollOffset; i < endIdx; i++ {
			cert := m.certificates[i]

			// Trust status
			trustText := "UNTRUSTED"
			trustStyle := untrustedStyle
			if cert.Trusted {
				trustText = "TRUSTED"
				trustStyle = trustedStyle
			}

			// Build certificate info
			info := fmt.Sprintf(
				"%s\n\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s\n%s %s",
				trustStyle.Render("["+trustText+"]"),
				fieldStyle.Render("Host:"),
				valueStyle.Render(cert.Host),
				fieldStyle.Render("Subject:"),
				valueStyle.Render(cert.Subject),
				fieldStyle.Render("Issuer:"),
				valueStyle.Render(cert.Issuer),
				fieldStyle.Render("Fingerprint:"),
				valueStyle.Render(m.formatFingerprint(cert.Fingerprint)),
				fieldStyle.Render("Valid From:"),
				valueStyle.Render(m.formatTime(cert.NotBefore)),
				fieldStyle.Render("Valid Until:"),
				valueStyle.Render(m.formatTime(cert.NotAfter)),
			)

			if i == m.selectedIdx {
				b.WriteString(selectedStyle.Render(info))
			} else {
				b.WriteString(normalStyle.Render(info))
			}
			b.WriteString("\n")
		}

		// Show scroll indicator if needed
		if endIdx < len(m.certificates) {
			b.WriteString(normalStyle.Render("▼ more below ▼"))
			b.WriteString("\n")
		}
	}

	// Help text
	helpText := "j/k: move • t: trust • u: untrust • esc/q/c: close"
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