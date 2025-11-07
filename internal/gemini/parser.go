package gemini

import (
	"bufio"
	"bytes"
	"net/url"
	"strings"

	"starsearch/internal/types"
)

// Parser parses text/gemini format documents
type Parser struct {
	baseURL *url.URL
}

// NewParser creates a new Gemini document parser
func NewParser(baseURL string) *Parser {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		// If URL parsing fails, create parser with nil baseURL
		// Relative URLs won't be resolved, but absolute URLs will still work
		return &Parser{
			baseURL: nil,
		}
	}
	return &Parser{
		baseURL: parsed,
	}
}

// Parse parses a Gemini response into a structured document
func (p *Parser) Parse(resp *types.Response) (*types.Document, error) {
	doc := &types.Document{
		URL:      resp.URL,
		RawBody:  resp.Body,
		Lines:    make([]types.Line, 0),
		Links:    make([]types.Line, 0),
		MIMEType: GetMIMEType(resp),
	}

	// Only parse text/gemini documents
	if !IsTextGemini(doc.MIMEType) && !IsTextPlain(doc.MIMEType) {
		// For non-text content, just store the body
		return doc, nil
	}

	// Parse line by line
	scanner := bufio.NewScanner(bytes.NewReader(resp.Body))
	inPreformat := false
	linkNum := 1

	for scanner.Scan() {
		rawLine := scanner.Text()
		line := p.parseLine(rawLine, &inPreformat, &linkNum)
		doc.Lines = append(doc.Lines, line)

		// Track links separately for easy access
		if line.Type == types.LineLink {
			doc.Links = append(doc.Links, line)
		}
	}

	return doc, scanner.Err()
}

// parseLine parses a single line of Gemini text
func (p *Parser) parseLine(rawLine string, inPreformat *bool, linkNum *int) types.Line {
	line := types.Line{
		Raw: rawLine,
	}

	// Check for preformat toggle
	if strings.HasPrefix(rawLine, "```") {
		*inPreformat = !*inPreformat
		if *inPreformat {
			line.Type = types.LinePreformatStart
			// Alt text is everything after the ```
			line.Text = strings.TrimSpace(strings.TrimPrefix(rawLine, "```"))
		} else {
			line.Type = types.LinePreformatEnd
		}
		return line
	}

	// If we're in preformat mode, return as-is
	if *inPreformat {
		line.Type = types.LinePreformatText
		line.Text = rawLine
		return line
	}

	// Link line: => URL [optional text]
	if strings.HasPrefix(rawLine, "=>") {
		line.Type = types.LineLink
		content := strings.TrimSpace(strings.TrimPrefix(rawLine, "=>"))
		parts := strings.Fields(content)

		if len(parts) > 0 {
			linkURL := parts[0]

			// Parse and resolve URL
			parsed, err := url.Parse(linkURL)
			if err == nil {
				// Resolve relative URLs against base
				if p.baseURL != nil {
					parsed = p.baseURL.ResolveReference(parsed)
				}
				line.URL = parsed.String()
			} else {
				line.URL = linkURL
			}

			// Link text is everything after the URL
			if len(parts) > 1 {
				line.Text = strings.Join(parts[1:], " ")
			} else {
				// No text provided, use URL as text
				line.Text = linkURL
			}

			line.LinkNum = *linkNum
			*linkNum++
		}
		return line
	}

	// Heading line: # ## ###
	if strings.HasPrefix(rawLine, "#") {
		if strings.HasPrefix(rawLine, "###") {
			line.Type = types.LineHeading3
			line.Text = strings.TrimSpace(strings.TrimPrefix(rawLine, "###"))
		} else if strings.HasPrefix(rawLine, "##") {
			line.Type = types.LineHeading2
			line.Text = strings.TrimSpace(strings.TrimPrefix(rawLine, "##"))
		} else {
			line.Type = types.LineHeading1
			line.Text = strings.TrimSpace(strings.TrimPrefix(rawLine, "#"))
		}
		return line
	}

	// List item: * item
	if strings.HasPrefix(rawLine, "*") && len(rawLine) > 1 && rawLine[1] == ' ' {
		line.Type = types.LineList
		line.Text = strings.TrimSpace(strings.TrimPrefix(rawLine, "*"))
		return line
	}

	// Quote line: > quote
	if strings.HasPrefix(rawLine, ">") {
		line.Type = types.LineQuote
		line.Text = strings.TrimSpace(strings.TrimPrefix(rawLine, ">"))
		return line
	}

	// Default: regular text
	line.Type = types.LineText
	line.Text = rawLine
	return line
}

// GetTitle attempts to extract a title from the document (first heading)
func GetTitle(doc *types.Document) string {
	for _, line := range doc.Lines {
		if line.Type == types.LineHeading1 ||
		   line.Type == types.LineHeading2 ||
		   line.Type == types.LineHeading3 {
			if line.Text != "" {
				return line.Text
			}
		}
	}

	// No heading found, try to use URL
	parsed, err := url.Parse(doc.URL)
	if err == nil {
		return parsed.Host
	}

	return doc.URL
}
