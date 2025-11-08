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

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"starsearch/internal/gemini"
	"starsearch/internal/gopher"
	"starsearch/internal/renderer"
	"starsearch/internal/storage"
	"starsearch/internal/types"
	"starsearch/internal/ui"
)

// Model is the main application model
type Model struct {
	client          *gemini.Client
	gopherClient    *gopher.Client
	tofuStore       *gemini.TOFUStore
	history         *storage.History
	bookmarks       *storage.Bookmarks
	config          *storage.Config
	addressBar      *ui.AddressBar
	viewport        *ui.ContentViewport
	statusBar       *ui.StatusBar
	tabBar          *ui.TabBar
	helpModal       *ui.HelpModal
	inputModal      *ui.InputModal
	bookmarksModal  *ui.BookmarksModal
	searchModal     *ui.SearchModal
	width           int
	height          int
	currentURL      string
	currentDoc      *types.Document
	linkNumbers     bool // Whether we're in link number input mode
	linkInput       string
	showHelp        bool   // Whether to show the help modal
	showInput       bool   // Whether to show the input modal
	showBookmarks   bool   // Whether to show the bookmarks modal
	showSearch      bool   // Whether to show the search modal
	pendingInputURL string // URL that triggered input request
	quitting        bool
	isNavigating    bool // Whether currently navigating (to avoid adding to history during back/forward)
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
	configPath := filepath.Join(starsearchDir, "config.toml")

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

	// Create config, history and bookmarks
	config := storage.NewConfig(configPath)
	history := storage.NewHistory(historyPath, config.Get().General.MaxHistory)
	bookmarks := storage.NewBookmarks(bookmarksPath)

	// Create UI components
	addressBar := ui.NewAddressBar()
	viewport := ui.NewContentViewport(80, 20)
	statusBar := ui.NewStatusBar(80)
	tabBar := ui.NewTabBar()
	helpModal := ui.NewHelpModal()
	inputModal := ui.NewInputModal()
	bookmarksModal := ui.NewBookmarksModal()
	searchModal := ui.NewSearchModal()

	// Create initial tab
	tabBar.AddTab("", "New Tab")

	return &Model{
		client:         client,
		gopherClient:   gopherClient,
		tofuStore:      tofuStore,
		history:        history,
		bookmarks:      bookmarks,
		config:         config,
		addressBar:     addressBar,
		viewport:       viewport,
		statusBar:      statusBar,
		tabBar:         tabBar,
		helpModal:      helpModal,
		inputModal:     inputModal,
		bookmarksModal: bookmarksModal,
		searchModal:    searchModal,
		width:          80,
		height:         24,
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
		// If bookmarks modal is showing, handle it first
		if m.showBookmarks {
			var cmd tea.Cmd
			m.bookmarksModal, cmd = m.bookmarksModal.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			// Check if modal was closed
			if !m.bookmarksModal.IsVisible() {
				m.showBookmarks = false
			}
			return m, tea.Batch(cmds...)
		}

		// If search modal is showing, handle it
		if m.showSearch {
			var cmd tea.Cmd
			m.searchModal, cmd = m.searchModal.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			// Check if modal was closed
			if !m.searchModal.IsVisible() {
				m.showSearch = false
			}
			return m, tea.Batch(cmds...)
		}

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
		case "ctrl+t":
			// New tab
			if !m.addressBar.IsFocused() && !m.linkNumbers {
				m.saveCurrentTabState()
				m.tabBar.AddTab("", "New Tab")
				m.loadTabState()
				return m, nil
			}

		case "ctrl+w":
			// Close current tab
			if !m.addressBar.IsFocused() && !m.linkNumbers {
				if len(m.tabBar.GetTabs()) > 1 {
					currentIdx := m.tabBar.GetActiveIndex()
					m.tabBar.CloseTab(currentIdx)
					m.loadTabState()
				} else {
					// Last tab - quit application
					m.quitting = true
					return m, tea.Quit
				}
				return m, nil
			}

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

		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			// Handle link number input
			if m.linkNumbers {
				m.linkInput += msg.String()
				m.statusBar.SetMessage("Enter link number: " + m.linkInput)
				return m, nil
			}
			// Tab switching (1-9)
			if !m.addressBar.IsFocused() && !m.linkNumbers {
				num, _ := strconv.Atoi(msg.String())
				tabIdx := num - 1
				if tabIdx >= 0 && tabIdx < len(m.tabBar.GetTabs()) {
					m.saveCurrentTabState()
					m.tabBar.SwitchTab(tabIdx)
					m.loadTabState()
				}
				return m, nil
			}

		case "0":
			// Handle link number input only
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

		case "ctrl+f":
			// Open search modal
			if !m.addressBar.IsFocused() && !m.linkNumbers && m.currentDoc != nil {
				m.showSearch = true
				return m, m.searchModal.Show(m.currentDoc)
			}

		case "ctrl+y":
			// Copy page content to clipboard
			if !m.addressBar.IsFocused() && !m.linkNumbers && m.currentDoc != nil {
				return m, m.copyPageContent()
			}

		case "b":
			// Toggle bookmarks modal
			if !m.addressBar.IsFocused() && !m.linkNumbers {
				// Close help modal if open
				m.showHelp = false
				m.showBookmarks = true
				m.bookmarksModal.Show(m.bookmarks.GetAll())
				return m, nil
			}

		case "ctrl+tab":
			// Next tab
			if !m.addressBar.IsFocused() && !m.linkNumbers {
				tabs := m.tabBar.GetTabs()
				if len(tabs) > 1 {
					m.saveCurrentTabState()
					nextIdx := (m.tabBar.GetActiveIndex() + 1) % len(tabs)
					m.tabBar.SwitchTab(nextIdx)
					m.loadTabState()
				}
				return m, nil
			}

		case "ctrl+shift+tab":
			// Previous tab
			if !m.addressBar.IsFocused() && !m.linkNumbers {
				tabs := m.tabBar.GetTabs()
				if len(tabs) > 1 {
					m.saveCurrentTabState()
					prevIdx := m.tabBar.GetActiveIndex() - 1
					if prevIdx < 0 {
						prevIdx = len(tabs) - 1
					}
					m.tabBar.SwitchTab(prevIdx)
					m.loadTabState()
				}
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update component sizes (subtract 2 to account for terminal edges)
		m.addressBar.SetWidth(m.width - 2)

		// Calculate viewport height: total - tab bar (1) - address bar (3) - status bar (1)
		viewportHeight := m.height - 5
		if viewportHeight < 1 {
			viewportHeight = 1
		}
		m.viewport.SetSize(m.width, viewportHeight)

		// Set viewport Y position (tab bar (1) + address bar with border (3) = 4 lines)
		m.viewport.SetYPosition(4)

		m.statusBar.SetWidth(m.width)
		m.tabBar.SetSize(m.width, 1)
		m.helpModal.SetSize(m.width, m.height)
		m.inputModal.SetSize(m.width, m.height)
		m.bookmarksModal.SetSize(m.width, m.height)
		m.searchModal.SetSize(m.width, m.height)

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

	case ui.BookmarkSelectedMsg:
		// User selected a bookmark to navigate to
		m.showBookmarks = false
		m.statusBar.SetMessage("Navigating to bookmark...")
		return m, m.navigate(msg.URL)

	case ui.BookmarkDeleteMsg:
		// User deleted a bookmark
		if err := m.bookmarks.Remove(msg.URL); err == nil {
			m.statusBar.SetMessage("Bookmark deleted")
			// Refresh the bookmarks modal with updated list
			m.bookmarksModal.Show(m.bookmarks.GetAll())
		} else {
			m.statusBar.SetError("Failed to delete bookmark")
		}
		return m, nil

	case ui.SearchSubmitMsg:
		// User submitted a search
		m.viewport.SetSearch(msg.Query, m.searchModal.GetResults(), msg.CaseSensitive)
		return m, nil

	case ui.SearchNavigateMsg:
		// User is navigating search results
		if msg.Direction == "next" || msg.Direction == "prev" {
			// Navigation is handled by the search modal
			result := m.searchModal.GetCurrentResult()
			if result != nil {
				m.viewport.GoToSearchResult(result)
			}
		} else if msg.Direction == "goto" {
			// Go to selected result
			result := m.searchModal.GetCurrentResult()
			if result != nil {
				m.viewport.GoToSearchResult(result)
			}
		}
		return m, nil

	case ui.SearchCloseMsg:
		// User closed search modal
		m.showSearch = false
		m.viewport.ClearSearch()
		return m, nil

	case ui.NavigateMsg:
		// Handle navigation
		return m, m.navigate(msg.URL)

	case fetchCompleteMsg:
		// Handle fetch completion
		m.statusBar.SetLoading(false)

		if msg.err != nil {
			m.statusBar.SetError(msg.err.Error())
			m.saveCurrentTabState()
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

			// Save tab state
			m.saveCurrentTabState()

			return m, nil
		}

		// Handle Gemini protocol (default)
		// Handle different status codes
		if gemini.IsSuccessStatus(msg.resp.Status) {
			mimeType := gemini.GetMIMEType(msg.resp)

			// Check if this is an image
			if renderer.IsImageMIME(mimeType) {
				// Render image
				imgRenderer := renderer.NewImageRenderer(m.width-4, m.height-8)
				renderedImage, err := imgRenderer.RenderImage(msg.resp.Body)
				if err != nil {
					m.statusBar.SetError(fmt.Sprintf("Failed to render image: %v", err))
					return m, nil
				}

				// Create a document with the rendered image as preformatted text
				doc := &types.Document{
					URL:      msg.resp.URL,
					RawBody:  msg.resp.Body,
					MIMEType: mimeType,
					Lines:    []types.Line{},
					Links:    []types.Line{},
				}

				// Split rendered image into lines
				for _, line := range strings.Split(renderedImage, "\n") {
					doc.Lines = append(doc.Lines, types.Line{
						Type: types.LineText,
						Text: line,
						Raw:  line,
					})
				}

				m.currentDoc = doc
				m.currentURL = msg.resp.URL
				m.viewport.SetDocument(doc)
				m.statusBar.SetURL(m.currentURL)

				// Use filename or URL as title
				title := msg.resp.URL
				m.statusBar.SetMessage(fmt.Sprintf("Image loaded: %s", mimeType))

				// Add to history
				if !m.isNavigating {
					m.history.Add(m.currentURL, title)
				}
				m.isNavigating = false

				// Save tab state
				m.saveCurrentTabState()
			} else {
				// Parse text document
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

				// Save tab state
				m.saveCurrentTabState()
			}

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
		// If bookmarks modal is showing, handle mouse events there
		if m.showBookmarks {
			var cmd tea.Cmd
			m.bookmarksModal, cmd = m.bookmarksModal.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			// Check if modal was closed
			if !m.bookmarksModal.IsVisible() {
				m.showBookmarks = false
			}
			return m, tea.Batch(cmds...)
		}

		// If search modal is showing, handle mouse events there
		if m.showSearch {
			var cmd tea.Cmd
			m.searchModal, cmd = m.searchModal.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			// Check if modal was closed
			if !m.searchModal.IsVisible() {
				m.showSearch = false
			}
			return m, tea.Batch(cmds...)
		}

		if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
			// Check if click is on tab bar (line 0)
			if msg.Y == 0 {
				// Pass to tab bar for handling
				var cmd tea.Cmd
				m.tabBar, cmd = m.tabBar.Update(msg)
				if cmd != nil {
					// Check if this is a tab switch message
					if switchMsg, ok := cmd().(ui.TabSwitchMsg); ok {
						// Save current tab state and load new tab state
						m.saveCurrentTabState()
						m.tabBar.SwitchTab(switchMsg.Index)
						m.loadTabState()
					}
				}
				return m, nil
			}

			// Check if click is on address bar (lines 1-3)
			if msg.Y >= 1 && msg.Y <= 3 {
				if !m.addressBar.IsFocused() {
					// Focus address bar, same as Ctrl+L
					if m.linkNumbers {
						m.linkNumbers = false
						m.linkInput = ""
						m.statusBar.SetMessage("Ready")
						// Viewport moves back up when help text disappears
						m.viewport.SetYPosition(4)
					}
					m.addressBar.SetValue(m.currentURL)
					focusCmd := m.addressBar.Focus()
					cmds = append(cmds, focusCmd)
				}
				// Pass mouse event to address bar for cursor positioning
				var cmd tea.Cmd
				m.addressBar, cmd = m.addressBar.Update(msg)
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
				return m, tea.Batch(cmds...)
			}

			// Click anywhere else - blur address bar if focused
			if m.addressBar.IsFocused() {
				m.addressBar.Blur()
				return m, nil
			}
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

	// Show bookmarks modal if active (highest priority for overlay)
	if m.showBookmarks {
		return m.bookmarksModal.View()
	}

	// Show search modal if active
	if m.showSearch {
		return m.searchModal.View()
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
		m.tabBar.View(),
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

// copyPageContent copies the current page content to the clipboard
func (m *Model) copyPageContent() tea.Cmd {
	if m.currentDoc == nil {
		return nil
	}

	_ = clipboard.WriteAll(string(m.currentDoc.RawBody))
	return nil
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

// saveCurrentTabState saves the current browsing state to the active tab
func (m *Model) saveCurrentTabState() {
	if m.tabBar.GetActiveTab() != nil {
		url := m.currentURL
		doc := m.currentDoc
		scroll := m.viewport.GetScrollOffset()
		title := ""
		if doc != nil {
			title = gemini.GetTitle(doc)
		} else if url != "" {
			title = url
		}
		idx := m.tabBar.GetActiveIndex()
		m.tabBar.UpdateTab(idx, url, title, doc, scroll)
	}
}

// loadTabState loads the state from the active tab
func (m *Model) loadTabState() {
	tab := m.tabBar.GetActiveTab()
	if tab != nil {
		m.currentURL = tab.URL
		m.currentDoc = tab.Document
		if tab.Document != nil {
			m.viewport.SetDocument(tab.Document)
			m.viewport.SetScrollOffset(tab.Scroll)
		} else {
			// Clear viewport if tab has no document
			m.viewport.SetDocument(nil)
		}
		m.statusBar.SetURL(m.currentURL)
		m.addressBar.SetValue(m.currentURL)
	}
}
