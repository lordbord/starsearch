package app

import (
	"crypto/x509"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"starsearch/internal/gemini"
	"starsearch/internal/gopher"
	"starsearch/internal/storage"
	"starsearch/internal/types"
	"starsearch/internal/ui"
)

// Model is the main application model
type Model struct {
	client         *gemini.Client
	gopherClient   *gopher.Client
	tofuStore      *gemini.TOFUStore
	history        *storage.History
	bookmarks      *storage.Bookmarks
	addressBar     *ui.AddressBar
	viewport       *ui.ContentViewport
	statusBar      *ui.StatusBar
	helpModal      *ui.HelpModal
	inputModal     *ui.InputModal
	width          int
	height         int
	currentURL     string
	currentDoc     *types.Document
	linkNumbers    bool   // Whether we're in link number input mode
	linkInput      string
	showHelp       bool   // Whether to show the help modal
	showInput      bool   // Whether to show the input modal
	pendingInputURL string // URL that triggered input request
	quitting       bool
	isNavigating   bool   // Whether currently navigating (to avoid adding to history during back/forward)
}

// NewModel creates a new application model
func NewModel() (*Model, error) {
	// Get config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = os.TempDir()
	}

	starsearchDir := filepath.Join(configDir, "starsearch")
	tofuPath := filepath.Join(starsearchDir, "known_hosts.json")
	historyPath := filepath.Join(starsearchDir, "history.json")
	bookmarksPath := filepath.Join(starsearchDir, "bookmarks.json")

	// Create TOFU store
	tofuStore, err := gemini.NewTOFUStore(tofuPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create TOFU store: %w", err)
	}

	// Setup TOFU callbacks (for now, auto-accept all)
	tofuStore.OnNewCert = func(host string, cert *x509.Certificate) bool {
		return true // Auto-accept new certificates
	}
	tofuStore.OnCertChange = func(host string, old, new *x509.Certificate) bool {
		return true // Auto-accept changed certificates (user will see warning)
	}

	// Create clients
	client := gemini.NewClient(tofuStore)
	gopherClient := gopher.NewClient()

	// Create history and bookmarks
	history := storage.NewHistory(historyPath, 1000)
	bookmarks := storage.NewBookmarks(bookmarksPath)

	// Create UI components
	addressBar := ui.NewAddressBar()
	viewport := ui.NewContentViewport(80, 20)
	statusBar := ui.NewStatusBar(80)
	helpModal := ui.NewHelpModal()
	inputModal := ui.NewInputModal()

	return &Model{
		client:       client,
		gopherClient: gopherClient,
		tofuStore:    tofuStore,
		history:      history,
		bookmarks:    bookmarks,
		addressBar:   addressBar,
		viewport:     viewport,
		statusBar:    statusBar,
		helpModal:    helpModal,
		inputModal:   inputModal,
		width:        80,
		height:       24,
	}, nil
}

// Init initializes the application
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// If input modal is showing, handle it first
		if m.showInput {
			var cmd tea.Cmd
			m.inputModal, cmd = m.inputModal.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		}

		// Global key handlers
		switch msg.String() {
		case "ctrl+c", "q":
			if m.showHelp {
				m.showHelp = false
				return m, nil
			}
			if !m.addressBar.IsFocused() && !m.linkNumbers {
				m.quitting = true
				return m, tea.Quit
			}

		case "ctrl+l":
			// Focus address bar
			m.linkNumbers = false
			m.linkInput = ""
			m.addressBar.SetValue(m.currentURL)
			return m, m.addressBar.Focus()

		case "g":
			// Enter link number mode
			if !m.addressBar.IsFocused() && !m.linkNumbers {
				m.linkNumbers = true
				m.linkInput = ""
				m.statusBar.SetMessage("Enter link number: ")
				// Viewport moves down by 1 line due to help text
				m.viewport.SetYPosition(4)
				return m, nil
			}

		case "esc":
			// Exit help modal
			if m.showHelp {
				m.showHelp = false
				return m, nil
			}
			// Exit link number mode
			if m.linkNumbers {
				m.linkNumbers = false
				m.linkInput = ""
				m.statusBar.SetMessage("Ready")
				// Viewport moves back up when help text disappears
				m.viewport.SetYPosition(3)
				return m, nil
			}

		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			// Handle link number input
			if m.linkNumbers {
				m.linkInput += msg.String()
				m.statusBar.SetMessage("Enter link number: " + m.linkInput)
				return m, nil
			}

		case "enter":
			// Activate link number
			if m.linkNumbers {
				num, err := strconv.Atoi(m.linkInput)
				if err == nil {
					m.linkNumbers = false
					m.linkInput = ""
					m.statusBar.SetMessage("Ready")
					// Viewport moves back up when help text disappears
					m.viewport.SetYPosition(3)
					return m, m.viewport.SelectLinkByNumber(num)
				}
				m.linkNumbers = false
				m.linkInput = ""
				m.statusBar.SetMessage("Invalid link number")
				// Viewport moves back up when help text disappears
				m.viewport.SetYPosition(3)
				return m, nil
			}

		case "r":
			// Reload current page
			if !m.addressBar.IsFocused() && !m.linkNumbers && m.currentURL != "" {
				m.isNavigating = true
				return m, m.navigate(m.currentURL)
			}

		case "d":
			// Add/remove bookmark
			if !m.addressBar.IsFocused() && !m.linkNumbers && m.currentURL != "" {
				if m.bookmarks.HasBookmark(m.currentURL) {
					// Remove bookmark
					if err := m.bookmarks.Remove(m.currentURL); err == nil {
						m.statusBar.SetMessage("Bookmark removed")
					} else {
						m.statusBar.SetError("Failed to remove bookmark")
					}
				} else {
					// Add bookmark
					title := "Untitled"
					if m.currentDoc != nil {
						title = gemini.GetTitle(m.currentDoc)
					}
					if err := m.bookmarks.Add(m.currentURL, title, nil); err == nil {
						m.statusBar.SetMessage("Bookmark added")
					} else {
						m.statusBar.SetError("Failed to add bookmark")
					}
				}
				return m, nil
			}

		case "h", "left", "alt+left":
			// Go back in history
			if !m.addressBar.IsFocused() && !m.linkNumbers {
				if m.history.CanGoBack() {
					url := m.history.Back()
					if url != "" {
						m.isNavigating = true
						m.statusBar.SetMessage("Going back...")
						return m, m.navigate(url)
					}
				} else {
					m.statusBar.SetMessage("No more history to go back")
				}
			}

		case "l", "right", "alt+right":
			// Go forward in history
			if !m.addressBar.IsFocused() && !m.linkNumbers {
				if m.history.CanGoForward() {
					url := m.history.Forward()
					if url != "" {
						m.isNavigating = true
						m.statusBar.SetMessage("Going forward...")
						return m, m.navigate(url)
					}
				} else {
					m.statusBar.SetMessage("No more history to go forward")
				}
			}

		case "j", "down":
			if !m.addressBar.IsFocused() && !m.linkNumbers {
				m.viewport.ScrollDown()
			}

		case "k", "up":
			if !m.addressBar.IsFocused() && !m.linkNumbers {
				m.viewport.ScrollUp()
			}

		case "pgdown", " ":
			if !m.addressBar.IsFocused() && !m.linkNumbers {
				m.viewport.PageDown()
			}

		case "pgup":
			if !m.addressBar.IsFocused() && !m.linkNumbers {
				m.viewport.PageUp()
			}

		case "?":
			// Toggle help modal
			if !m.addressBar.IsFocused() && !m.linkNumbers {
				m.showHelp = !m.showHelp
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update component sizes
		m.addressBar.SetWidth(m.width)

		viewportHeight := m.height - 4 // Leave space for address bar and status bar
		if viewportHeight < 1 {
			viewportHeight = 1
		}
		m.viewport.SetSize(m.width, viewportHeight)

		// Set viewport Y position (address bar with border takes 3 lines)
		m.viewport.SetYPosition(3)

		m.statusBar.SetWidth(m.width)
		m.helpModal.SetSize(m.width, m.height)
		m.inputModal.SetSize(m.width, m.height)

		return m, nil

	case ui.InputSubmitMsg:
		// User submitted input
		m.showInput = false
		if m.pendingInputURL != "" && msg.Input != "" {
			// Append input as URL-encoded query parameter
			inputURL := m.pendingInputURL + "?" + url.QueryEscape(msg.Input)
			m.pendingInputURL = ""
			return m, m.navigate(inputURL)
		}
		m.pendingInputURL = ""
		return m, nil

	case ui.InputCancelMsg:
		// User cancelled input
		m.showInput = false
		m.pendingInputURL = ""
		m.statusBar.SetMessage("Input cancelled")
		return m, nil

	case ui.NavigateMsg:
		// Handle navigation
		return m, m.navigate(msg.URL)

	case fetchCompleteMsg:
		// Handle fetch completion
		m.statusBar.SetLoading(false)

		if msg.err != nil {
			m.statusBar.SetError(msg.err.Error())
			return m, nil
		}

		// Handle Gopher protocol
		if msg.protocol == "gopher" {
			// Parse the document using Gopher parser
			parser := gopher.NewParser(msg.resp.URL)
			doc, err := parser.Parse(msg.resp)
			if err != nil {
				m.statusBar.SetError(fmt.Sprintf("Failed to parse Gopher document: %v", err))
				return m, nil
			}

			m.currentDoc = doc
			m.currentURL = msg.resp.URL
			m.viewport.SetDocument(doc)
			m.statusBar.SetURL(m.currentURL)

			// Get title from URL for Gopher
			title := msg.resp.URL
			m.statusBar.SetMessage(fmt.Sprintf("Loaded: %s", title))

			// Add to history (unless we're navigating back/forward)
			if !m.isNavigating {
				m.history.Add(m.currentURL, title)
			}
			m.isNavigating = false

			return m, nil
		}

		// Handle Gemini protocol (default)
		// Handle different status codes
		if gemini.IsSuccessStatus(msg.resp.Status) {
			// Parse the document
			parser := gemini.NewParser(msg.resp.URL)
			doc, err := parser.Parse(msg.resp)
			if err != nil {
				m.statusBar.SetError(fmt.Sprintf("Failed to parse document: %v", err))
				return m, nil
			}

			m.currentDoc = doc
			m.currentURL = msg.resp.URL
			m.viewport.SetDocument(doc)
			m.statusBar.SetURL(m.currentURL)

			// Get title for status
			title := gemini.GetTitle(doc)
			m.statusBar.SetMessage(fmt.Sprintf("Loaded: %s", title))

			// Add to history (unless we're navigating back/forward)
			if !m.isNavigating {
				m.history.Add(m.currentURL, title)
			}
			m.isNavigating = false

		} else if gemini.IsRedirectStatus(msg.resp.Status) {
			// Handle redirect
			newURL := msg.resp.Meta
			m.statusBar.SetMessage(fmt.Sprintf("Redirecting to: %s", newURL))
			return m, m.navigate(newURL)

		} else if gemini.IsInputStatus(msg.resp.Status) {
			// Handle input request (status 10 or 11)
			prompt := msg.resp.Meta
			if prompt == "" {
				prompt = "Input required"
			}

			// Determine if this is sensitive input (status 11)
			sensitive := (msg.resp.Status == 11)

			// Store the URL that triggered input request
			m.pendingInputURL = msg.resp.URL

			// Show input modal
			m.showInput = true
			return m, m.inputModal.Show(prompt, sensitive)

		} else {
			// Handle error status
			statusMsg := gemini.GetStatusMessage(msg.resp.Status)
			m.statusBar.SetError(fmt.Sprintf("%s: %s", statusMsg, msg.resp.Meta))
		}

		return m, nil

	case externalLinkOpenedMsg:
		// External link was opened successfully
		m.statusBar.SetMessage(fmt.Sprintf("Opened external link: %s", msg.url))
		return m, nil

	case tea.MouseMsg:
		// Check if click is on address bar (first 3 lines)
		if msg.Type == tea.MouseLeft && msg.Y <= 2 && !m.addressBar.IsFocused() {
			// Focus address bar, same as Ctrl+L
			if m.linkNumbers {
				m.linkNumbers = false
				m.linkInput = ""
				m.statusBar.SetMessage("Ready")
				// Viewport moves back up when help text disappears
				m.viewport.SetYPosition(3)
			}
			m.addressBar.SetValue(m.currentURL)
			return m, m.addressBar.Focus()
		}

		// Pass mouse events to viewport
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// Update address bar
	if m.addressBar.IsFocused() {
		var cmd tea.Cmd
		m.addressBar, cmd = m.addressBar.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// Update scroll percentage in status bar
	m.statusBar.SetScrollPercent(m.viewport.GetScrollPercent())

	return m, tea.Batch(cmds...)
}

// View renders the application
func (m *Model) View() string {
	if m.quitting {
		return "Thanks for using starsearch!\n"
	}

	// Show input modal if active
	if m.showInput {
		return m.inputModal.View()
	}

	// Show help modal if active
	if m.showHelp {
		return m.helpModal.View()
	}

	// Layout components vertically
	components := []string{
		m.addressBar.View(),
		m.viewport.View(),
		m.statusBar.View(),
	}

	// Add help text if in link mode
	if m.linkNumbers {
		helpStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")).
			Background(lipgloss.Color("235")).
			Padding(0, 1)
		helpText := helpStyle.Render(" Type link number and press Enter (ESC to cancel) ")
		components = append([]string{helpText}, components...)
	}

	return lipgloss.JoinVertical(lipgloss.Left, components...)
}

// navigate fetches and displays a URL
func (m *Model) navigate(urlStr string) tea.Cmd {
	// Parse URL to detect protocol
	parsedURL, err := url.Parse(urlStr)
	if err == nil && parsedURL.Scheme != "" {
		switch parsedURL.Scheme {
		case "gopher":
			// Handle Gopher protocol
			m.statusBar.SetLoading(true)
			m.statusBar.SetMessage("Fetching " + urlStr + "...")

			return func() tea.Msg {
				resp, err := m.gopherClient.Fetch(urlStr)
				return fetchCompleteMsg{resp: resp, err: err, protocol: "gopher"}
			}

		case "gemini":
			// Handle Gemini protocol (continue below)

		default:
			// Handle other external protocols (http, https, etc.)
			return m.openExternalURL(urlStr)
		}
	}

	// Normalize URL for Gemini protocol
	if !strings.HasPrefix(urlStr, "gemini://") {
		urlStr = "gemini://" + urlStr
	}

	m.statusBar.SetLoading(true)
	m.statusBar.SetMessage("Fetching " + urlStr + "...")

	return func() tea.Msg {
		resp, err := m.client.Fetch(urlStr)
		return fetchCompleteMsg{resp: resp, err: err, protocol: "gemini"}
	}
}

// openExternalURL opens a URL in the system's default browser
func (m *Model) openExternalURL(urlStr string) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd

		switch runtime.GOOS {
		case "linux", "freebsd", "openbsd", "netbsd":
			cmd = exec.Command("xdg-open", urlStr)
		case "darwin":
			cmd = exec.Command("open", urlStr)
		case "windows":
			cmd = exec.Command("cmd", "/c", "start", urlStr)
		default:
			return fetchCompleteMsg{
				resp: nil,
				err:  fmt.Errorf("unsupported platform for opening external links: %s", runtime.GOOS),
			}
		}

		err := cmd.Start()
		if err != nil {
			return fetchCompleteMsg{
				resp: nil,
				err:  fmt.Errorf("failed to open external link: %w", err),
			}
		}

		// Return a message indicating the link was opened externally
		m.statusBar.SetMessage(fmt.Sprintf("Opened external link in browser: %s", urlStr))
		return externalLinkOpenedMsg{url: urlStr}
	}
}

// externalLinkOpenedMsg is sent when an external link is opened
type externalLinkOpenedMsg struct {
	url string
}

// fetchCompleteMsg is sent when a fetch completes
type fetchCompleteMsg struct {
	resp     *types.Response
	err      error
	protocol string // "gemini" or "gopher"
}
