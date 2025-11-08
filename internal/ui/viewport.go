package ui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"starsearch/internal/types"
)

// ContentViewport displays Gemini document content
type ContentViewport struct {
	viewport       viewport.Model
	document       *types.Document
	width          int
	height         int
	yPosition      int // Y position of viewport in screen layout
	selectedLink   int // Currently selected link for keyboard navigation
	lineMapping    map[int]int // Maps rendered line number to document line index
	linkBounds     map[int][]linkBound // Maps rendered line to clickable link regions
	searchResults  []types.SearchResult
	currentSearch  string
	searchHighlight bool
	caseSensitive  bool
}

// linkBound represents the clickable region of a link on a rendered line
type linkBound struct {
	startX int
	endX   int
	url    string
}

// NewContentViewport creates a new content viewport
func NewContentViewport(width, height int) *ContentViewport {
	vp := viewport.New(width, height)
	vp.MouseWheelEnabled = true

	return &ContentViewport{
		viewport:       vp,
		width:          width,
		height:         height,
		selectedLink:   -1,
		searchResults:  []types.SearchResult{},
		searchHighlight: false,
		caseSensitive:  false,
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

				// Check if there are link bounds for this rendered line
				if bounds, ok := c.linkBounds[renderedLineNum]; ok {
					// Check if click X position is within any link bound
					for _, bound := range bounds {
						if msg.X >= bound.startX && msg.X < bound.endX {
							return c, func() tea.Msg { return NavigateMsg{URL: bound.url} }
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
	c.searchResults = []types.SearchResult{}
	c.currentSearch = ""
	c.searchHighlight = false
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

// SetSearch sets search results and highlights them
func (c *ContentViewport) SetSearch(query string, results []types.SearchResult, caseSensitive bool) {
	c.currentSearch = query
	c.searchResults = results
	c.searchHighlight = len(results) > 0
	c.caseSensitive = caseSensitive

	// Re-render document with highlights
	content := c.renderDocument()
	c.viewport.SetContent(content)
}

// ClearSearch clears search highlighting
func (c *ContentViewport) ClearSearch() {
	c.currentSearch = ""
	c.searchResults = []types.SearchResult{}
	c.searchHighlight = false

	// Re-render document without highlights
	content := c.renderDocument()
	c.viewport.SetContent(content)
}

// GoToSearchResult navigates to a specific search result
func (c *ContentViewport) GoToSearchResult(result *types.SearchResult) {
	if result == nil || c.document == nil {
		return
	}

	// Find rendered line number for this document line
	targetLine := -1
	for renderedLine, docLine := range c.lineMapping {
		if docLine == result.Line {
			targetLine = renderedLine
			break
		}
	}

	if targetLine >= 0 {
		// Scroll to make line visible
		c.viewport.YOffset = targetLine
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
	c.linkBounds = make(map[int][]linkBound) // Initialize link bounds
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
			// Wrap heading text before styling
			wrapped := wordWrap("# "+line.Text, c.width)
			rendered := heading1Style.Render(wrapped)
			// Styles with margins produce multiple lines
			addMultilineContent(rendered, i)

		case types.LineHeading2:
			// Wrap heading text before styling
			wrapped := wordWrap("## "+line.Text, c.width)
			rendered := heading2Style.Render(wrapped)
			// Styles with margins produce multiple lines
			addMultilineContent(rendered, i)

		case types.LineHeading3:
			// Wrap heading text before styling
			wrapped := wordWrap("### "+line.Text, c.width)
			rendered := heading3Style.Render(wrapped)
			addMultilineContent(rendered, i)

		case types.LineLink:
			linkText := line.Text
			if linkText == "" {
				linkText = line.URL
			}

			// Apply search highlighting if enabled
			if c.searchHighlight && c.currentSearch != "" {
				linkText = c.highlightSearchText(linkText, i)
			}

			// Add link number for keyboard navigation
			numStrPlain := fmt.Sprintf("[%d] ", line.LinkNum)
			linkPrefix := len(numStrPlain)

			// Wrap link text to fit viewport width (accounting for the link number prefix)
			availableWidth := c.width - linkPrefix
			if availableWidth < 20 {
				availableWidth = 20 // Minimum width for readability
			}
			wrappedLink := wordWrap(linkText, availableWidth)
			wrappedLines := strings.Split(wrappedLink, "\n")

			// Render each wrapped line
			for lineIdx, wrappedLine := range wrappedLines {
				var displayLine string
				if lineIdx == 0 {
					// First line includes the link number
					numStr := linkNumStyle.Render(fmt.Sprintf("[%d]", line.LinkNum))
					linkStr := linkStyle.Render(wrappedLine)
					displayLine = numStr + " " + linkStr

					// Calculate clickable bounds for first line
					startX := linkPrefix
					endX := startX + len(stripANSI(wrappedLine))
					c.linkBounds[renderedLineNum] = []linkBound{
						{startX: startX, endX: endX, url: line.URL},
					}
				} else {
					// Continuation lines are indented to align with first line
					indent := strings.Repeat(" ", linkPrefix)
					linkStr := linkStyle.Render(wrappedLine)
					displayLine = indent + linkStr

					// Calculate clickable bounds for continuation line
					startX := linkPrefix
					endX := startX + len(stripANSI(wrappedLine))
					c.linkBounds[renderedLineNum] = []linkBound{
						{startX: startX, endX: endX, url: line.URL},
					}
				}

				addLine(displayLine, i)
			}
			// Skip the default addLine call since we handled it above
			continue

		case types.LineList:
			// Wrap list text (accounting for bullet point)
			listPrefix := "  • "
			availableWidth := c.width - len(listPrefix)
			if availableWidth < 20 {
				availableWidth = 20
			}
			wrapped := wordWrap(line.Text, availableWidth)
			wrappedLines := strings.Split(wrapped, "\n")

			for lineIdx, wrappedLine := range wrappedLines {
				if lineIdx == 0 {
					// First line with bullet
					addLine(listStyle.Render(listPrefix+wrappedLine), i)
				} else {
					// Continuation lines indented
					indent := strings.Repeat(" ", len(listPrefix))
					addLine(listStyle.Render(indent+wrappedLine), i)
				}
			}
			continue

		case types.LineQuote:
			// Wrap quote text (accounting for padding)
			quotePadding := 2 // PaddingLeft(2) from quoteStyle
			availableWidth := c.width - quotePadding
			if availableWidth < 20 {
				availableWidth = 20
			}
			wrapped := wordWrap(line.Text, availableWidth)
			rendered := quoteStyle.Render(wrapped)
			addMultilineContent(rendered, i)

		case types.LinePreformatStart:
			// Optionally show alt text, hard-wrap if needed
			if line.Text != "" {
				wrapped := hardWrap("``` "+line.Text, c.width)
				addMultilineContent(preformatStyle.Render(wrapped), i)
			}
			// Note: If text is empty, we don't render anything but the mapping continues

		case types.LinePreformatText:
			// Hard-wrap preformatted text to prevent overflow
			wrapped := hardWrap(line.Text, c.width)
			addMultilineContent(preformatStyle.Render(wrapped), i)

		case types.LinePreformatEnd:
			addLine(preformatStyle.Render("```"), i)

		case types.LineText:
			// Word wrap for long lines
			if len(line.Text) == 0 {
				addLine("", i)
			} else {
				text := line.Text
				// Apply search highlighting if enabled
				if c.searchHighlight && c.currentSearch != "" {
					text = c.highlightSearchText(text, i)
				}
				wrapped := wordWrap(text, c.width)
				// wordWrap may produce multiple lines
				addMultilineContent(wrapped, i)
			}
		}
	}

	return builder.String()
}

// highlightSearchText applies highlighting to search matches in text
func (c *ContentViewport) highlightSearchText(text string, lineIdx int) string {
	if !c.searchHighlight || c.currentSearch == "" {
		return text
	}

	// Find all search results for this line
	var lineResults []types.SearchResult
	for _, result := range c.searchResults {
		if result.Line == lineIdx {
			lineResults = append(lineResults, result)
		}
	}

	if len(lineResults) == 0 {
		return text
	}

	// Sort results by start position
	for i := 0; i < len(lineResults)-1; i++ {
		for j := i + 1; j < len(lineResults); j++ {
			if lineResults[i].Start > lineResults[j].Start {
				lineResults[i], lineResults[j] = lineResults[j], lineResults[i]
			}
		}
	}

	// Apply highlighting
	result := ""
	lastEnd := 0
	
	searchHighlightStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("11")).
		Bold(true)

	searchCurrentStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("3")).
		Bold(true)

	for _, searchResult := range lineResults {
		// Add text before match
		result += text[lastEnd:searchResult.Start]
		
		// Add highlighted match
		matchText := text[searchResult.Start:searchResult.End]
		
		// Check if this is the current match
		isCurrent := false
		for _, currentResult := range c.searchResults {
			if currentResult.Line == lineIdx && 
			   currentResult.Start == searchResult.Start && 
			   currentResult.End == searchResult.End {
				isCurrent = true
				break
			}
		}
		
		if isCurrent {
			result += searchCurrentStyle.Render(matchText)
		} else {
			result += searchHighlightStyle.Render(matchText)
		}
		
		lastEnd = searchResult.End
	}
	
	// Add remaining text
	result += text[lastEnd:]
	
	return result
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

// hardWrap wraps text by breaking at exact width (for preformatted text)
func hardWrap(text string, width int) string {
	if width <= 0 {
		width = 80
	}

	if len(text) <= width {
		return text
	}

	var lines []string
	for len(text) > width {
		lines = append(lines, text[:width])
		text = text[width:]
	}
	if len(text) > 0 {
		lines = append(lines, text)
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

// GetScrollOffset returns the current scroll offset
func (c *ContentViewport) GetScrollOffset() int {
	return c.viewport.YOffset
}

// SetScrollOffset sets the scroll offset
func (c *ContentViewport) SetScrollOffset(offset int) {
	c.viewport.YOffset = offset
}

// stripANSI removes ANSI escape codes from a string
func stripANSI(str string) string {
	// Regex to match ANSI escape sequences
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansiRegex.ReplaceAllString(str, "")
}
