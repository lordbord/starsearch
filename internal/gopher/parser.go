package gopher

import (
	"bufio"
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"starsearch/internal/types"
)

// Parser parses Gopher menu format documents
type Parser struct {
	baseURL *url.URL
}

// NewParser creates a new Gopher document parser
func NewParser(baseURL string) *Parser {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return &Parser{
			baseURL: nil,
		}
	}
	return &Parser{
		baseURL: parsed,
	}
}

// Parse parses a Gopher response into a structured document
func (p *Parser) Parse(resp *types.Response) (*types.Document, error) {
	doc := &types.Document{
		URL:      resp.URL,
		RawBody:  resp.Body,
		Lines:    make([]types.Line, 0),
		Links:    make([]types.Line, 0),
		MIMEType: resp.Meta,
	}

	// Only parse gopher menu format
	if !IsGopherMenu(doc.MIMEType) {
		// For non-menu content (text files), treat as plain text
		if strings.HasPrefix(doc.MIMEType, "text/plain") {
			scanner := bufio.NewScanner(bytes.NewReader(resp.Body))
			for scanner.Scan() {
				line := types.Line{
					Type: types.LineText,
					Raw:  scanner.Text(),
					Text: scanner.Text(),
				}
				doc.Lines = append(doc.Lines, line)
			}
			return doc, scanner.Err()
		}
		// For binary content, just store the body
		return doc, nil
	}

	// Parse Gopher menu line by line
	scanner := bufio.NewScanner(bytes.NewReader(resp.Body))
	linkNum := 1

	for scanner.Scan() {
		rawLine := scanner.Text()
		line := p.parseGopherLine(rawLine, &linkNum)
		doc.Lines = append(doc.Lines, line)

		// Track links separately for easy access
		if line.Type == types.LineLink {
			doc.Links = append(doc.Links, line)
		}
	}

	return doc, scanner.Err()
}

// parseGopherLine parses a single line of a Gopher menu
// Gopher format: TypeDisplayString\tSelector\tHost\tPort\r\n
func (p *Parser) parseGopherLine(rawLine string, linkNum *int) types.Line {
	line := types.Line{
		Raw: rawLine,
	}

	// Remove any trailing CRLF or LF
	rawLine = strings.TrimRight(rawLine, "\r\n")

	// Empty line
	if len(rawLine) == 0 {
		line.Type = types.LineText
		line.Text = ""
		return line
	}

	// Check for end-of-menu marker
	if rawLine == "." {
		line.Type = types.LineText
		line.Text = ""
		return line
	}

	// Extract item type (first character)
	itemType := rawLine[0:1]
	remaining := ""
	if len(rawLine) > 1 {
		remaining = rawLine[1:]
	}

	// Split by tabs
	parts := strings.Split(remaining, "\t")

	// Get display string
	displayString := ""
	if len(parts) > 0 {
		displayString = parts[0]
	}

	// Check if this is a selectable item (has selector, host, port)
	if len(parts) >= 3 {
		selector := parts[0]
		if len(parts) >= 2 {
			selector = parts[1]
		}
		host := ""
		port := "70"

		if len(parts) >= 3 {
			host = parts[2]
		}
		if len(parts) >= 4 && parts[3] != "" {
			port = parts[3]
		}

		// Build URL based on item type
		isLink := false
		var gopherURL string

		switch itemType {
		case "i", "3":
			// Informational text or error - not a link
			line.Type = types.LineText
			line.Text = displayString
			return line

		case "h":
			// HTML link - check if it's an external URL
			if strings.HasPrefix(selector, "URL:") {
				gopherURL = strings.TrimPrefix(selector, "URL:")
			} else {
				gopherURL = fmt.Sprintf("gopher://%s:%s/h%s", host, port, selector)
			}
			isLink = true

		case "0":
			// Text file
			gopherURL = fmt.Sprintf("gopher://%s:%s/0%s", host, port, selector)
			isLink = true

		case "1":
			// Directory/Menu
			gopherURL = fmt.Sprintf("gopher://%s:%s/1%s", host, port, selector)
			isLink = true

		case "7":
			// Search
			gopherURL = fmt.Sprintf("gopher://%s:%s/7%s", host, port, selector)
			isLink = true

		case "9":
			// Binary file
			gopherURL = fmt.Sprintf("gopher://%s:%s/9%s", host, port, selector)
			isLink = true

		case "g", "I":
			// Image
			gopherURL = fmt.Sprintf("gopher://%s:%s/%s%s", host, port, itemType, selector)
			isLink = true

		default:
			// Unknown type - treat as link
			gopherURL = fmt.Sprintf("gopher://%s:%s/%s%s", host, port, itemType, selector)
			isLink = true
		}

		if isLink {
			line.Type = types.LineLink
			line.Text = displayString
			line.URL = gopherURL
			line.LinkNum = *linkNum
			*linkNum++
			return line
		}
	} else {
		// No selector/host/port - treat as informational text
		line.Type = types.LineText
		line.Text = displayString
		return line
	}

	// Default to text
	line.Type = types.LineText
	line.Text = displayString
	return line
}

// GetItemTypeDescription returns a human-readable description of a Gopher item type
func GetItemTypeDescription(itemType string) string {
	switch itemType {
	case "0":
		return "Text file"
	case "1":
		return "Directory"
	case "2":
		return "CSO phone book"
	case "3":
		return "Error"
	case "4":
		return "BinHex file"
	case "5":
		return "DOS archive"
	case "6":
		return "UUEncoded file"
	case "7":
		return "Search"
	case "8":
		return "Telnet session"
	case "9":
		return "Binary file"
	case "+":
		return "Redundant server"
	case "g":
		return "GIF image"
	case "I":
		return "Image"
	case "T":
		return "TN3270 session"
	case "h":
		return "HTML"
	case "i":
		return "Info"
	case "s":
		return "Sound"
	default:
		return "Unknown"
	}
}
