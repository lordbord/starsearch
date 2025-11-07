package gopher

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
	"time"

	"starsearch/internal/types"
)

// Client handles Gopher protocol requests
type Client struct {
	timeout time.Duration
}

// NewClient creates a new Gopher client
func NewClient() *Client {
	return &Client{
		timeout: 30 * time.Second,
	}
}

// Fetch retrieves a Gopher URL and returns a response
func (c *Client) Fetch(urlStr string) (*types.Response, error) {
	// Parse URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Ensure scheme is gopher
	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "gopher"
		urlStr = parsedURL.String()
	} else if parsedURL.Scheme != "gopher" {
		return nil, fmt.Errorf("unsupported scheme: %s (only gopher:// is supported)", parsedURL.Scheme)
	}

	// Get host and port
	host := parsedURL.Hostname()
	port := parsedURL.Port()
	if port == "" {
		port = "70" // Default gopher port
	}

	// Get selector (path)
	// Gopher URL format: gopher://host:port/[type][selector]
	// where type is a single character (0, 1, 7, etc.)
	path := parsedURL.Path
	itemType := "1" // Default to directory
	selector := ""

	if path == "" || path == "/" {
		// Root directory
		itemType = "1"
		selector = ""
	} else if len(path) > 1 {
		// Extract item type (character after first /)
		itemType = string(path[1])
		if len(path) > 2 {
			// Selector is everything after the type character
			selector = path[2:]
		}
	}

	// Connect to server
	address := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", address, c.timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	// Set read deadline
	conn.SetDeadline(time.Now().Add(c.timeout))

	// Send selector followed by CRLF
	_, err = conn.Write([]byte(selector + "\r\n"))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Read response until connection closes
	body, err := io.ReadAll(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Determine MIME type based on item type
	mimeType := GetMIMEType(itemType)

	// Create response
	// Gopher doesn't have status codes, so we use 20 (success) for Gemini compatibility
	response := &types.Response{
		Status: 20, // Success
		Meta:   mimeType,
		Body:   body,
		URL:    urlStr,
	}

	return response, nil
}

// GetMIMEType returns the MIME type for a Gopher item type
func GetMIMEType(itemType string) string {
	switch itemType {
	case "0":
		return "text/plain"
	case "1":
		return "text/gopher" // Gopher menu
	case "g", "I":
		return "image/gif"
	case "h":
		return "text/html"
	case "s":
		return "audio/basic"
	case "9", "5":
		return "application/octet-stream"
	default:
		return "text/gopher" // Default to menu format
	}
}

// IsGopherMenu checks if the response is a Gopher menu
func IsGopherMenu(mimeType string) bool {
	return strings.HasPrefix(mimeType, "text/gopher")
}
