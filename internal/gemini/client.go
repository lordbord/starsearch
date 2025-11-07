package gemini

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"git.sr.ht/~adnano/go-gemini"
	"starsearch/internal/types"
)

// Client wraps the go-gemini client with additional functionality
type Client struct {
	client     *gemini.Client
	tofuStore  *TOFUStore
	userAgent  string
	timeout    time.Duration
}

// NewClient creates a new Gemini client with TOFU support
func NewClient(tofuStore *TOFUStore) *Client {
	return &Client{
		client:    &gemini.Client{},
		tofuStore: tofuStore,
		userAgent: "starsearch/1.0",
		timeout:   30 * time.Second,
	}
}

// Fetch retrieves a Gemini URL and returns a parsed response
func (c *Client) Fetch(urlStr string) (*types.Response, error) {
	// Parse and validate URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Ensure scheme is gemini
	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "gemini"
		urlStr = parsedURL.String()
	} else if parsedURL.Scheme != "gemini" {
		return nil, fmt.Errorf("unsupported scheme: %s (only gemini:// is supported)", parsedURL.Scheme)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	// Fetch the URL
	resp, err := c.client.Do(ctx, &gemini.Request{
		URL: parsedURL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch: %w", err)
	}
	defer resp.Body.Close()

	// Verify certificate using TOFU
	tlsState := resp.TLS()
	if tlsState != nil && len(tlsState.PeerCertificates) > 0 {
		cert := tlsState.PeerCertificates[0]
		host := parsedURL.Hostname()

		if err := c.tofuStore.Verify(host, cert); err != nil {
			return nil, fmt.Errorf("certificate verification failed: %w", err)
		}
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Create response
	response := &types.Response{
		Status: int(resp.Status),
		Meta:   resp.Meta,
		Body:   body,
		URL:    urlStr,
	}

	return response, nil
}

// IsSuccessStatus checks if a status code indicates success
func IsSuccessStatus(status int) bool {
	return status >= 20 && status < 30
}

// IsRedirectStatus checks if a status code indicates a redirect
func IsRedirectStatus(status int) bool {
	return status >= 30 && status < 40
}

// IsInputStatus checks if a status code requests input
func IsInputStatus(status int) bool {
	return status >= 10 && status < 20
}

// IsSensitiveInput checks if input should be hidden (password)
func IsSensitiveInput(status int) bool {
	return status == 11
}

// IsTemporaryFailure checks if a status code indicates a temporary failure
func IsTemporaryFailure(status int) bool {
	return status >= 40 && status < 50
}

// IsPermanentFailure checks if a status code indicates a permanent failure
func IsPermanentFailure(status int) bool {
	return status >= 50 && status < 60
}

// IsCertificateRequired checks if a status code requires a client certificate
func IsCertificateRequired(status int) bool {
	return status >= 60 && status < 70
}

// GetStatusMessage returns a human-readable message for a status code
func GetStatusMessage(status int) string {
	switch {
	case status == 10:
		return "Input required"
	case status == 11:
		return "Sensitive input required"
	case status == 20:
		return "Success"
	case status == 30:
		return "Temporary redirect"
	case status == 31:
		return "Permanent redirect"
	case status == 40:
		return "Temporary failure"
	case status == 41:
		return "Server unavailable"
	case status == 42:
		return "CGI error"
	case status == 43:
		return "Proxy error"
	case status == 44:
		return "Slow down (rate limited)"
	case status == 50:
		return "Permanent failure"
	case status == 51:
		return "Not found"
	case status == 52:
		return "Gone"
	case status == 53:
		return "Proxy request refused"
	case status == 59:
		return "Bad request"
	case status == 60:
		return "Client certificate required"
	case status == 61:
		return "Certificate not authorized"
	case status == 62:
		return "Certificate not valid"
	default:
		return fmt.Sprintf("Unknown status: %d", status)
	}
}

// GetMIMEType returns the MIME type from the meta field for successful responses
func GetMIMEType(resp *types.Response) string {
	if IsSuccessStatus(resp.Status) {
		return resp.Meta
	}
	return ""
}

// IsTextGemini checks if the response is text/gemini
func IsTextGemini(mimeType string) bool {
	return mimeType == "text/gemini" ||
	       mimeType == "text/gemini; charset=utf-8" ||
	       mimeType == "text/gemini;charset=utf-8"
}

// IsTextPlain checks if the response is plain text
func IsTextPlain(mimeType string) bool {
	return mimeType == "text/plain" ||
	       mimeType == "text/plain; charset=utf-8" ||
	       mimeType == "text/plain;charset=utf-8"
}
