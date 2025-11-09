package types

import "time"

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

// Config represents the application configuration
type Config struct {
	General     GeneralConfig     `toml:"general"`
	UI          UIConfig          `toml:"ui"`
	Colors      ColorConfig       `toml:"colors"`
	Downloads   DownloadConfig    `toml:"downloads"`
	Performance PerformanceConfig `toml:"performance"`
}

// GeneralConfig contains general application settings
type GeneralConfig struct {
	HomeURL         string `toml:"home_url"`
	SearchEngine    string `toml:"search_engine"`
	MaxHistory      int    `toml:"max_history"`
	AutoSaveHistory bool   `toml:"auto_save_history"`
	RestoreSession  bool   `toml:"restore_session"`
}

// UIConfig contains user interface settings
type UIConfig struct {
	ShowLineNumbers bool `toml:"show_line_numbers"`
	ShowLinkNumbers bool `toml:"show_link_numbers"`
	EnableMouse     bool `toml:"enable_mouse"`
	ScrollSpeed     int  `toml:"scroll_speed"`
}

// ColorConfig contains color theme settings
type ColorConfig struct {
	Theme           string `toml:"theme"`
	LinkColor       string `toml:"link_color"`
	VisitedLinkColor string `toml:"visited_link_color"`
	Heading1Color   string `toml:"heading1_color"`
	Heading2Color   string `toml:"heading2_color"`
	Heading3Color   string `toml:"heading3_color"`
	TextColor       string `toml:"text_color"`
	QuoteColor      string `toml:"quote_color"`
	PreformatColor  string `toml:"preformat_color"`
	BackgroundColor string `toml:"background_color"`
}

// DownloadConfig contains download settings
type DownloadConfig struct {
	Directory       string `toml:"directory"`
	AskBeforeDownload bool `toml:"ask_before_download"`
	MaxConcurrent   int    `toml:"max_concurrent"`
	Timeout         int    `toml:"timeout"`
}

// PerformanceConfig contains performance settings
type PerformanceConfig struct {
	EnableCache      bool `toml:"enable_cache"`
	CacheTTL         int  `toml:"cache_ttl"`
	CacheSizeMB      int  `toml:"cache_size_mb"`
	EnablePrefetch   bool `toml:"enable_prefetch"`
	PrefetchIdleDelay int `toml:"prefetch_idle_delay"`
	ConnectionPoolSize int `toml:"connection_pool_size"`
}

// DownloadStatus represents the status of a download
type DownloadStatus int

const (
	DownloadPending DownloadStatus = iota
	Downloading
	DownloadCompleted
	DownloadFailed
	DownloadCancelled
)

// Download represents a file download
type Download struct {
	ID          string         `json:"id"`
	URL         string         `json:"url"`
	Filename    string         `json:"filename"`
	Size        int64          `json:"size"`
	Downloaded  int64          `json:"downloaded"`
	Status      DownloadStatus `json:"status"`
	Error       string         `json:"error"`
	StartTime   int64          `json:"start_time"`
	FinishTime  int64          `json:"finish_time"`
}

// SearchResult represents a search match in a document
type SearchResult struct {
	Line     int    `json:"line"`
	Start    int    `json:"start"`
	End      int    `json:"end"`
	Text     string `json:"text"`
	Selected bool   `json:"selected"`
}

// CertificateInfo represents certificate information for display
type CertificateInfo struct {
	Host         string    `json:"host"`
	Fingerprint  string    `json:"fingerprint"`
	NotBefore    time.Time `json:"not_before"`
	NotAfter     time.Time `json:"not_after"`
	Issuer       string    `json:"issuer"`
	Subject      string    `json:"subject"`
	Trusted      bool      `json:"trusted"`
	FirstSeen    time.Time `json:"first_seen"`
	LastSeen     time.Time `json:"last_seen"`
}

// SessionTab represents a tab in a saved session
type SessionTab struct {
	URL    string `json:"url"`
	Title  string `json:"title"`
	Scroll int    `json:"scroll"`
}

// Session represents a saved browser session
type Session struct {
	Tabs        []SessionTab `json:"tabs"`
	ActiveIndex int          `json:"active_index"`
	Timestamp   int64        `json:"timestamp"`
}
