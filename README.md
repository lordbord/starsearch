# starsearch - A Gemini Browser for the Terminal

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24%2B-blue.svg)](https://go.dev/)

A modern, feature-rich Gemini protocol browser built with Go and Bubble Tea TUI framework. Browse Geminispace with full mouse and keyboard support, TOFU certificate handling, and a beautiful terminal interface.

**Available on Linux, macOS, and Windows** through multiple package managers including AUR, Homebrew, and Chocolatey.

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

Starsearch is available for **Linux**, **macOS**, and **Windows** through multiple package managers.

### Linux

#### Arch Linux (AUR)

```bash
# Using yay
yay -S starsearch

# Using paru
paru -S starsearch
```

#### Other Linux Distributions (via Homebrew)

```bash
brew tap lordbord/starsearch
brew install starsearch
```

### macOS

#### Using Homebrew

```bash
brew tap lordbord/starsearch
brew install starsearch
```

### Windows

#### Using Chocolatey

```powershell
choco install starsearch
```

### Pre-built Binaries

Pre-built binaries for all platforms are available on the [Releases page](https://github.com/lordbord/starsearch/releases).

1. Download the appropriate archive for your system:
   - **Linux (x86_64)**: `starsearch-VERSION-linux-amd64.tar.gz`
   - **Linux (ARM64)**: `starsearch-VERSION-linux-arm64.tar.gz`
   - **macOS (Intel)**: `starsearch-VERSION-darwin-amd64.tar.gz`
   - **macOS (Apple Silicon)**: `starsearch-VERSION-darwin-arm64.tar.gz`
   - **Windows (x86_64)**: `starsearch-VERSION-windows-amd64.zip`
2. Extract the binary
3. Move it to a directory in your PATH
4. (Optional) Verify the checksum from `checksums.txt`

### Build from Source

Requires Go 1.24 or higher:

```bash
# Clone the repository
git clone https://github.com/lordbord/starsearch.git
cd starsearch

# Build the binary
go build -o starsearch ./cmd/starsearch

# Run the browser
./starsearch

# Or install globally
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

### ✅ v0.1.0 - Ready for Release!

All major features are complete and the project is ready for its first stable release:

**Core Features:**
- ✅ Full Gemini Protocol Support
- ✅ Interactive TUI with mouse and keyboard navigation
- ✅ TOFU certificate management
- ✅ History navigation and bookmarks
- ✅ Multi-tab browsing
- ✅ Download support with progress tracking
- ✅ Search in page functionality
- ✅ Configuration system with TOML

**Cross-Platform Distribution:**
- ✅ Pre-built binaries for Linux (x86_64, ARM64), macOS (Intel, Apple Silicon), Windows
- ✅ AUR package for Arch Linux
- ✅ Homebrew formula ready (tap setup required)
- ✅ Chocolatey package ready (submission pending)
- ✅ GitHub Actions workflow for automated releases

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
├── cmd/starsearch/              # Main entry point
├── internal/
│   ├── app/                    # Main application model
│   ├── gemini/                 # Gemini client, parser, TOFU
│   ├── ui/                     # UI components (viewport, addressbar, statusbar, modals)
│   ├── storage/                # History, bookmarks, config, downloads
│   └── types/                  # Shared types
├── homebrew/                    # Homebrew formula
├── chocolatey/                  # Chocolatey package
├── .github/workflows/           # CI/CD automation
├── go.mod
├── go.sum
├── PKGBUILD                     # AUR package definition
├── DISTRIBUTION.md              # Distribution guide
└── README.md
```

## Distribution

For detailed information about packaging and distribution across platforms, see [DISTRIBUTION.md](DISTRIBUTION.md).

### Package Maintainers

If you'd like to package starsearch for additional platforms or distributions:

1. Pre-built binaries are available in [GitHub Releases](https://github.com/lordbord/starsearch/releases)
2. See [DISTRIBUTION.md](DISTRIBUTION.md) for checksums and platform-specific notes
3. Open an issue to let us know about your package!

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
