package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"starsearch/internal/types"
)

// DownloadModal displays download progress and management
type DownloadModal struct {
	visible      bool
	downloads    []types.Download
	selectedIdx  int
	width        int
	height       int
	progress     progress.Model
	scrollOffset int
}

// DownloadStartMsg is sent when a download starts
type DownloadStartMsg struct {
	Download *types.Download
}

// DownloadProgressMsg is sent when download progress updates
type DownloadProgressMsg struct {
	ID         string
	Downloaded int64
}

// DownloadCompleteMsg is sent when a download completes
type DownloadCompleteMsg struct {
	ID string
}

// DownloadCancelMsg is sent when user cancels a download
type DownloadCancelMsg struct {
	ID string
}

// DownloadCloseMsg is sent when download modal is closed
type DownloadCloseMsg struct{}

func NewDownloadModal() *DownloadModal {
	prog := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)

	return &DownloadModal{
		visible:      false,
		downloads:    []types.Download{},
		selectedIdx:  0,
		scrollOffset: 0,
		progress:     prog,
	}
}

func (m *DownloadModal) Show(downloads []types.Download) tea.Cmd {
	m.visible = true
	m.downloads = downloads
	m.selectedIdx = 0
	m.scrollOffset = 0
	return nil
}

func (m *DownloadModal) Hide() {
	m.visible = false
}

func (m *DownloadModal) IsVisible() bool {
	return m.visible
}

func (m *DownloadModal) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.progress.Width = min(width-20, 40)
}

func (m *DownloadModal) Update(msg tea.Msg) (*DownloadModal, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc", "q", "d"))):
			m.Hide()
			return m, func() tea.Msg {
				return DownloadCloseMsg{}
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("j", "down"))):
			if m.selectedIdx < len(m.downloads)-1 {
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
			if len(m.downloads) > 0 {
				m.selectedIdx = len(m.downloads) - 1
				m.adjustScroll()
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("c", "delete"))):
			if m.selectedIdx < len(m.downloads) {
				download := m.downloads[m.selectedIdx]
				if download.Status == types.Downloading || download.Status == types.DownloadPending {
					return m, func() tea.Msg {
						return DownloadCancelMsg{ID: download.ID}
					}
				}
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
			if m.selectedIdx < len(m.downloads) {
				download := m.downloads[m.selectedIdx]
				if download.Status == types.DownloadFailed {
					// TODO: Implement retry functionality
					return m, nil
				}
			}
		}
	}

	var progressCmd tea.Cmd
	progressModel, progressCmd := m.progress.Update(msg)
	m.progress = progressModel.(progress.Model)
	if progressCmd != nil {
		cmd = progressCmd
	}
	return m, cmd
}

func (m *DownloadModal) adjustScroll() {
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

func (m *DownloadModal) formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func (m *DownloadModal) formatDuration(seconds int64) string {
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	} else if seconds < 3600 {
		return fmt.Sprintf("%dm %ds", seconds/60, seconds%60)
	} else {
		hours := seconds / 3600
		minutes := (seconds % 3600) / 60
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
}

func (m *DownloadModal) View() string {
	if !m.visible {
		return ""
	}

	var b strings.Builder

	// Calculate modal dimensions
	modalWidth := min(m.width-4, 80)
	if modalWidth < 60 {
		modalWidth = 60
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

	selectedStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("12")).
		Foreground(lipgloss.Color("0")).
		Bold(true).
		Width(modalWidth - 4)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Width(modalWidth - 4)

	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")).
		Width(modalWidth - 4)

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
	b.WriteString(titleStyle.Render(fmt.Sprintf("Downloads (%d active)", len(m.downloads))))
	b.WriteString("\n")

	if len(m.downloads) == 0 {
		b.WriteString(statusStyle.Render("No active downloads"))
		b.WriteString("\n")
	} else {
		// Calculate visible range
		visibleHeight := modalHeight - 8
		if visibleHeight < 1 {
			visibleHeight = 1
		}

		endIdx := m.scrollOffset + visibleHeight
		if endIdx > len(m.downloads) {
			endIdx = len(m.downloads)
		}

		// Show scroll indicator if needed
		if m.scrollOffset > 0 {
			b.WriteString(statusStyle.Render("▲ more above ▲"))
			b.WriteString("\n")
		}

		// Render visible downloads
		for i := m.scrollOffset; i < endIdx; i++ {
			download := m.downloads[i]

			// Format status
			statusText := ""
			switch download.Status {
			case types.DownloadPending:
				statusText = "Pending"
			case types.Downloading:
				statusText = "Downloading"
			case types.DownloadCompleted:
				statusText = "Completed"
			case types.DownloadFailed:
				statusText = "Failed: " + download.Error
			case types.DownloadCancelled:
				statusText = "Cancelled"
			}

			// Calculate progress
			percentage := float64(0)
			if download.Size > 0 {
				percentage = float64(download.Downloaded) / float64(download.Size) * 100
			}

			// Format file info
			fileInfo := fmt.Sprintf("%s (%s/%s)",
				download.Filename,
				m.formatBytes(download.Downloaded),
				m.formatBytes(download.Size))

			// Calculate speed and ETA if downloading
			speedText := ""
			etaText := ""
			if download.Status == types.Downloading && download.StartTime > 0 {
				elapsed := time.Now().Unix() - download.StartTime
				if elapsed > 0 {
					speed := float64(download.Downloaded) / float64(elapsed)
					speedText = fmt.Sprintf(" @ %s/s", m.formatBytes(int64(speed)))

					if download.Size > 0 && download.Downloaded > 0 {
						remaining := download.Size - download.Downloaded
						eta := int64(float64(remaining) / speed)
						etaText = fmt.Sprintf(" ETA: %s", m.formatDuration(eta))
					}
				}
			}

			// Build download line
			line := fmt.Sprintf("%s\n%s%s%s\n  %s%s",
				fileInfo,
				statusStyle.Render("["+statusText+"]"),
				speedText,
				etaText,
				m.progress.ViewAs(percentage/100),
				statusStyle.Render(fmt.Sprintf(" %.1f%%", percentage)))

			if i == m.selectedIdx {
				b.WriteString(selectedStyle.Render(line))
			} else {
				b.WriteString(normalStyle.Render(line))
			}
			b.WriteString("\n")
		}

		// Show scroll indicator if needed
		if endIdx < len(m.downloads) {
			b.WriteString(statusStyle.Render("▼ more below ▼"))
			b.WriteString("\n")
		}
	}

	// Help text
	helpText := "j/k: move • c: cancel • r: retry • esc/q/d: close"
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
