package types

// LineType represents the type of a line in a Gemini document
type LineType int

const (
	LineText LineType = iota
	LineLink
	LineHeading1
	LineHeading2
	LineHeading3
	LineList
	LineQuote
	LinePreformatStart
	LinePreformatEnd
	LinePreformatText
)

// Line represents a single line in a Gemini document
type Line struct {
	Type    LineType
	Raw     string // Raw line content
	Text    string // Display text
	URL     string // For links only
	LinkNum int    // Link number for keyboard selection
}

// Document represents a parsed Gemini document
type Document struct {
	URL      string
	RawBody  []byte
	Lines    []Line
	Links    []Line // All links for easy access
	MIMEType string
}

// Response represents a Gemini protocol response
type Response struct {
	Status     int
	Meta       string
	Body       []byte
	RemoteAddr string
	URL        string
}

// Tab represents a browser tab
type Tab struct {
	ID       int
	Title    string
	URL      string
	Document *Document
	Scroll   int // Scroll position
}

// Bookmark represents a saved bookmark
type Bookmark struct {
	Title string
	URL   string
	Tags  []string
}

// HistoryEntry represents a visited page
type HistoryEntry struct {
	URL       string
	Timestamp int64
	Title     string
}
