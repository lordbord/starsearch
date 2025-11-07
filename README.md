# starsearch - A Gemini Browser for the Terminal

A modern, feature-rich Gemini protocol browser built with Go and Bubble Tea TUI framework. Browse Geminispace with full mouse and keyboard support, TOFU certificate handling, and a beautiful terminal interface.

## Features

- **Full Gemini Protocol Support**: Implements the complete Gemini protocol specification
- **Mouse & Keyboard Support**: Click links with your mouse or use keyboard shortcuts
- **TOFU Security**: Trust On First Use certificate management keeps you safe
- **Beautiful TUI**: Clean, styled interface with syntax highlighting for Gemini documents
- **Fast & Lightweight**: Native Go performance with minimal resource usage
- **History Navigation**: Full back/forward navigation with persistent history
- **Bookmarks**: Save and manage your favorite Gemini capsules
- **Tab Support**: Browse multiple capsules simultaneously (coming soon)
- **Download Support**: Save binary files from Gemini capsules (coming soon)

## Installation

### Prerequisites

- Go 1.21 or higher

### Build from Source

```bash
# Clone or navigate to the repository
cd starsearch

# Build the binary
go build -o starsearch ./cmd/starsearch

# Run the browser
./starsearch
```

### Install Globally

```bash
go install ./cmd/starsearch
```

## Usage

### Starting the Browser

```bash
./starsearch
```

The browser will start with an empty page. Use `Ctrl+L` to focus the address bar and enter a Gemini URL.

### Keyboard Shortcuts

#### Navigation
- `Ctrl+L` - Focus the address bar to enter a URL
- `Enter` - Navigate to the URL in the address bar
- `R` - Reload the current page
- `H` / `←` / `Alt+←` - Go back in history
- `L` / `→` / `Alt+→` - Go forward in history
- `Esc` - Cancel current input/action

#### Scrolling
- `↑` / `K` - Scroll up one line
- `↓` / `J` - Scroll down one line
- `PgUp` - Scroll up one page
- `PgDn` / `Space` - Scroll down one page

#### Link Selection
- `G` - Enter link number mode
- `0-9` - Type link number
- `Enter` - Navigate to the selected link
- Click links with your mouse!

#### Bookmarks
- `D` - Add current page to bookmarks (or remove if already bookmarked)

#### Application
- `?` - Show help screen with all keyboard shortcuts
- `Q` / `Ctrl+C` - Quit the browser (when not in input mode)

### Browsing Geminispace

1. Press `Ctrl+L` to focus the address bar
2. Type a Gemini URL (e.g., `gemini://gemini.circumlunar.space/`)
3. Press `Enter` to navigate
4. Click links with your mouse or press `G` and type a link number
5. Enjoy browsing Geminispace!

### Example URLs to Try

- `gemini://gemini.circumlunar.space/` - Project Gemini homepage
- `gemini://geminispace.info/` - Geminispace search and directory
- `gemini://gus.guru/` - Gemini Universal Search
- `gemini://warmedal.se/~antenna/` - Antenna: Gemini feed aggregator
- `gemini://spacewalk.fedi.buzz/` - Spacewalk: Mastodon/Fediverse gateway

## Text/Gemini Format

The browser fully supports the text/gemini format with styled rendering:

- **Headings**: Three levels of headings with distinct styling
- **Links**: Numbered links with colors (click or use 'g' + number)
- **Lists**: Bulleted list items
- **Quotes**: Italic quoted text with indentation
- **Preformatted Text**: Code blocks and ASCII art with monospace styling

## Certificate Management (TOFU)

starsearch uses Trust On First Use (TOFU) for certificate management:

- Certificates are automatically trusted on first visit
- Certificate fingerprints are stored in `~/.config/starsearch/known_hosts.json`
- Changed certificates will trigger a warning (currently auto-accepted)
- Manual certificate management coming soon

## Configuration

Configuration files are stored in:
- Linux/BSD: `~/.config/starsearch/`
- macOS: `~/Library/Application Support/starsearch/`
- Windows: `%APPDATA%\starsearch\`

### Files

- `known_hosts.json` - TOFU certificate store
- `bookmarks.json` - Saved bookmarks
- `history.json` - Browsing history
- `config.toml` - User configuration (coming soon)

## Roadmap

### Phase 3: Interactive Features ✅
- [x] Link highlighting and selection
- [x] Mouse click support for links
- [x] Keyboard-based link navigation (g + number)
- [x] Status code handling (redirects, errors)
- [x] Input prompts (status 10/11) for search

### Phase 4: History & Bookmarks ✅
- [x] Back/forward navigation
- [x] Persistent history storage
- [x] Bookmark manager
- [x] Add/remove bookmarks with 'D' key
- [ ] Bookmarks UI (sidebar or modal to view all bookmarks)

### Phase 5: Tabs
- [ ] Multiple tab support
- [ ] Tab bar UI
- [ ] Tab switching shortcuts
- [ ] Per-tab state management

### Phase 6: Downloads
- [ ] Binary file detection
- [ ] Download prompts
- [ ] Progress indicators
- [ ] Configurable download directory

### Phase 7: Polish
- [x] Help screen
- [ ] Configuration system
- [ ] Custom themes and colors
- [ ] Search in page (Ctrl+F)
- [ ] Certificate manager UI

## Technical Details

### Architecture

- **Language**: Go 1.21+
- **TUI Framework**: Bubble Tea (The Elm Architecture)
- **Styling**: Lipgloss
- **Components**: Bubbles (textinput, viewport)
- **Gemini Client**: go-gemini

### Project Structure

```
starsearch/
├── cmd/starsearch/       # Main entry point
├── internal/
│   ├── app/             # Main application model
│   ├── gemini/          # Gemini client, parser, TOFU
│   ├── ui/              # UI components (viewport, addressbar, statusbar)
│   ├── storage/         # History, bookmarks, config (planned)
│   └── types/           # Shared types
├── go.mod
├── go.sum
└── README.md
```

## Contributing

Contributions are welcome! This is a personal project but feel free to:

- Report bugs
- Suggest features
- Submit pull requests
- Share feedback

## License

[Add your chosen license here]

## Acknowledgments

- [Project Gemini](https://gemini.circumlunar.space/) for the protocol
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the excellent TUI framework
- [go-gemini](https://git.sr.ht/~adnano/go-gemini) for the Gemini client library
- The Gemini community for creating an amazing corner of the internet

## About Gemini

Gemini is a new internet protocol which:
- Is heavier than Gopher
- Is lighter than the Web
- Will not replace either
- Strives for maximum power-to-weight ratio
- Takes user privacy seriously

Learn more at: `gemini://gemini.circumlunar.space/`

---

**Happy browsing! 🚀**
