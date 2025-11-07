package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"starsearch/internal/types"
)

// ContentViewport displays Gemini document content
type ContentViewport struct {
	viewport     viewport.Model
	document     *types.Document
	width        int
	height       int
	yPosition    int // Y position of viewport in the screen layout
	selectedLink int // Currently selected link for keyboard navigation
	lineMapping  map[int]int // Maps rendered line number to document line index
}

// NewContentViewport creates a new content viewport
func NewContentViewport(width, height int) *ContentViewport {
	vp := viewport.New(width, height)
	vp.MouseWheelEnabled = true

	return &ContentViewport{
		viewport:     vp,
		width:        width,
		height:       height,
		selectedLink: -1,
	}
}

// Init initializes the viewport
func (c *ContentViewport) Init() tea.Cmd {
	return nil
}

// Update handles viewport updates
func (c *ContentViewport) Update(msg tea.Msg) (*ContentViewport, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.MouseMsg:
		// Handle mouse clicks on links
		if msg.Button == tea.MouseButtonLeft && msg.Action == tea.MouseActionPress && c.document != nil {
			// Calculate which line was clicked
			// Subtract viewport's Y position to convert from screen coordinates to viewport coordinates
			viewportY := msg.Y - c.yPosition
			if viewportY >= 0 {
				// Calculate the rendered line number (accounting for scroll offset)
				renderedLineNum := c.viewport.YOffset + viewportY

				// Use line mapping to find the corresponding document line index
				if docLineIdx, ok := c.lineMapping[renderedLineNum]; ok {
					if docLineIdx >= 0 && docLineIdx < len(c.document.Lines) {
						line := c.document.Lines[docLineIdx]
						if line.Type == types.LineLink && line.URL != "" {
							return c, func() tea.Msg { return NavigateMsg{URL: line.URL} }
						}
					}
				}
			}
		}
	}

	c.viewport, cmd = c.viewport.Update(msg)
	return c, cmd
}

// View renders the viewport
func (c *ContentViewport) View() string {
	return c.viewport.View()
}

// SetDocument sets the document to display
func (c *ContentViewport) SetDocument(doc *types.Document) {
	c.document = doc
	c.selectedLink = -1
	c.viewport.YOffset = 0 // Reset scroll to top

	// Render the document
	content := c.renderDocument()
	c.viewport.SetContent(content)
}

// SetSize sets the viewport size
func (c *ContentViewport) SetSize(width, height int) {
	c.width = width
	c.height = height
	c.viewport.Width = width
	c.viewport.Height = height

	// Re-render document if present
	if c.document != nil {
		content := c.renderDocument()
		c.viewport.SetContent(content)
	}
}

// SetYPosition sets the viewport's Y position in the screen layout
func (c *ContentViewport) SetYPosition(y int) {
	c.yPosition = y
}

// renderDocument renders a Gemini document to styled text
func (c *ContentViewport) renderDocument() string {
	if c.document == nil {
		return "No document loaded"
	}

	var builder strings.Builder
	c.lineMapping = make(map[int]int) // Initialize line mapping
	renderedLineNum := 0 // Track which rendered line we're on

	// Helper function to add content and track line mapping
	addLine := func(content string, docLineIdx int) {
		builder.WriteString(content)
		builder.WriteString("\n")
		// Map this rendered line to the document line
		c.lineMapping[renderedLineNum] = docLineIdx
		renderedLineNum++
	}

	// Helper function to count and map multiple lines in content
	addMultilineContent := func(content string, docLineIdx int) {
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			builder.WriteString(line)
			builder.WriteString("\n")
			c.lineMapping[renderedLineNum] = docLineIdx
			renderedLineNum++
		}
	}

	// Define styles
	heading1Style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		MarginTop(1).
		MarginBottom(1)

	heading2Style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("14")).
		MarginTop(1)

	heading3Style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("10"))

	linkStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")).
		Underline(true)

	linkNumStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Bold(true)

	listStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7"))

	quoteStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Italic(true).
		PaddingLeft(2)

	preformatStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("11")).
		Background(lipgloss.Color("235"))

	for i, line := range c.document.Lines {
		switch line.Type {
		case types.LineHeading1:
			rendered := heading1Style.Render("# " + line.Text)
			// Styles with margins produce multiple lines
			addMultilineContent(rendered, i)

		case types.LineHeading2:
			rendered := heading2Style.Render("## " + line.Text)
			// Styles with margins produce multiple lines
			addMultilineContent(rendered, i)

		case types.LineHeading3:
			rendered := heading3Style.Render("### " + line.Text)
			addMultilineContent(rendered, i)

		case types.LineLink:
			linkText := line.Text
			if linkText == "" {
				linkText = line.URL
			}

			// Add link number for keyboard navigation
			numStr := linkNumStyle.Render(fmt.Sprintf("[%d]", line.LinkNum))
			linkStr := linkStyle.Render(linkText)

			addLine(numStr + " " + linkStr, i)

		case types.LineList:
			addLine(listStyle.Render("  • " + line.Text), i)

		case types.LineQuote:
			rendered := quoteStyle.Render(line.Text)
			addMultilineContent(rendered, i)

		case types.LinePreformatStart:
			// Optionally show alt text
			if line.Text != "" {
				addLine(preformatStyle.Render("``` " + line.Text), i)
			}
			// Note: If text is empty, we don't render anything but the mapping continues

		case types.LinePreformatText:
			addLine(preformatStyle.Render(line.Text), i)

		case types.LinePreformatEnd:
			addLine(preformatStyle.Render("```"), i)

		case types.LineText:
			// Word wrap for long lines
			if len(line.Text) == 0 {
				addLine("", i)
			} else {
				wrapped := wordWrap(line.Text, c.width)
				// wordWrap may produce multiple lines
				addMultilineContent(wrapped, i)
			}
		}
	}

	return builder.String()
}

// wordWrap wraps text to a specified width
func wordWrap(text string, width int) string {
	if width <= 0 {
		width = 80
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	var currentLine string

	for _, word := range words {
		if len(currentLine) == 0 {
			currentLine = word
		} else if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}

	if len(currentLine) > 0 {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}

// GetScrollPercent returns the scroll percentage
func (c *ContentViewport) GetScrollPercent() float64 {
	return c.viewport.ScrollPercent()
}

// ScrollUp scrolls the viewport up
func (c *ContentViewport) ScrollUp() {
	c.viewport.LineUp(1)
}

// ScrollDown scrolls the viewport down
func (c *ContentViewport) ScrollDown() {
	c.viewport.LineDown(1)
}

// PageUp scrolls up one page
func (c *ContentViewport) PageUp() {
	c.viewport.ViewUp()
}

// PageDown scrolls down one page
func (c *ContentViewport) PageDown() {
	c.viewport.ViewDown()
}

// SelectNextLink selects the next link
func (c *ContentViewport) SelectNextLink() tea.Cmd {
	if c.document == nil || len(c.document.Links) == 0 {
		return nil
	}

	c.selectedLink++
	if c.selectedLink >= len(c.document.Links) {
		c.selectedLink = 0
	}

	return nil
}

// SelectPrevLink selects the previous link
func (c *ContentViewport) SelectPrevLink() tea.Cmd {
	if c.document == nil || len(c.document.Links) == 0 {
		return nil
	}

	c.selectedLink--
	if c.selectedLink < 0 {
		c.selectedLink = len(c.document.Links) - 1
	}

	return nil
}

// ActivateSelectedLink activates the currently selected link
func (c *ContentViewport) ActivateSelectedLink() tea.Cmd {
	if c.document == nil || c.selectedLink < 0 || c.selectedLink >= len(c.document.Links) {
		return nil
	}

	link := c.document.Links[c.selectedLink]
	return func() tea.Msg { return NavigateMsg{URL: link.URL} }
}

// SelectLinkByNumber selects a link by its number
func (c *ContentViewport) SelectLinkByNumber(num int) tea.Cmd {
	if c.document == nil || len(c.document.Links) == 0 {
		return nil
	}

	// Find link with this number
	for _, link := range c.document.Links {
		if link.LinkNum == num {
			return func() tea.Msg { return NavigateMsg{URL: link.URL} }
		}
	}

	// Number not found
	return nil
}
