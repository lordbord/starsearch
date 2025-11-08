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
- **Tab Support**: Browse multiple capsules simultaneously with full tab management
- **Download Support**: Save binary files with progress tracking and queue management
- **Search in Page**: Find text within documents with highlighting and navigation
- **Configuration System**: Customizable settings via TOML configuration file
- **Certificate Manager**: View and manage TOFU certificates with manual trust control

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
- `B` - Open bookmarks manager

#### Search
- `Ctrl+F` - Open search in page
- `n` - Next search result
- `N` - Previous search result
- `Esc` - Close search

#### Tabs
- `Ctrl+T` - New tab
- `Ctrl+W` - Close current tab
- `Ctrl+Tab` - Next tab
- `Ctrl+Shift+Tab` - Previous tab
- `1-9` - Switch to specific tab

#### Application
- `?` - Show help screen with all keyboard shortcuts
- `Q` / `Ctrl+C` - Quit the browser (when not in input mode)

### Browsing Geminispace

1. Press `Ctrl+L` to focus the address bar
2. Type a Gemini URL (e.g., `gemini://gemini.circumlunar.space/`)
3. Press `Enter` to navigate
4. Click links with your mouse or press `G` and type a link number
5. Use `Ctrl+T` to open new tabs for parallel browsing
6. Press `Ctrl+F` to search within pages
7. Enjoy browsing Geminispace!

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
- **Search**: Text highlighting with current match emphasis
- **Images**: Automatic rendering with terminal-compatible display

## Certificate Management (TOFU)

starsearch uses Trust On First Use (TOFU) for certificate management:

- Certificates are automatically trusted on first visit
- Certificate fingerprints are stored in `~/.config/starsearch/known_hosts.json`
- Manual certificate management with trust/untrust controls
- View certificate details including issuer, subject, and validity periods
- Changed certificates trigger warnings with manual review options

## Configuration

Configuration files are stored in:
- Linux/BSD: `~/.config/starsearch/`
- macOS: `~/Library/Application Support/starsearch/`
- Windows: `%APPDATA%\starsearch\`

### Files

- `config.toml` - User configuration (colors, UI settings, downloads)
- `known_hosts.json` - TOFU certificate store
- `bookmarks.json` - Saved bookmarks
- `history.json` - Browsing history
- `downloads.json` - Active and completed downloads

### Configuration Options

The `config.toml` file supports the following sections:

```toml
[general]
home_url = "gemini://gemini.circumlunar.space/"
search_engine = "gemini://gus.guru/"
max_history = 1000
auto_save_history = true

[ui]
show_line_numbers = false
show_link_numbers = true
enable_mouse = true
scroll_speed = 3

[colors]
theme = "default"
link_color = "12"
visited_link_color = "13"
heading1_color = "11"
heading2_color = "14"
heading3_color = "10"
text_color = "15"
quote_color = "8"
preformat_color = "7"
background_color = "0"

[downloads]
directory = "~/Downloads"
ask_before_download = true
max_concurrent = 3
timeout = 30
```

## Development Status

🎉 **All Major Features Complete!**

The browser now includes all planned functionality:

### ✅ Completed Features
- **Interactive Features**: Link highlighting, mouse support, keyboard navigation
- **History & Bookmarks**: Full navigation, persistent storage, bookmark management UI
- **Tab Support**: Multiple tabs, tab bar, keyboard shortcuts, per-tab state
- **Downloads**: Binary file detection, progress tracking, queue management
- **Polish**: Configuration system, themes, search, certificate management

### 🔮 Future Enhancements
Potential areas for future development:
- Plugin system for custom protocols
- Session persistence and restoration
- Advanced bookmark organization (folders, tags)
- Custom themes and color schemes
- Integration with external editors
- Gemini-to-HTML export functionality

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
│   ├── ui/              # UI components (viewport, addressbar, statusbar, modals)
│   ├── storage/         # History, bookmarks, config, downloads
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

MIT License - see the [LICENSE](LICENSE) file for details.

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
