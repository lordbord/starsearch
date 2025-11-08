package storage

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"starsearch/internal/types"
)

// Config manages application configuration
type Config struct {
	config     *types.Config
	configPath string
}

// NewConfig creates a new configuration manager
func NewConfig(configPath string) *Config {
	c := &Config{
		config:     getDefaultConfig(),
		configPath: configPath,
	}

	// Try to load existing config
	_ = c.Load() // Ignore errors, use defaults if file doesn't exist

	return c
}

// getDefaultConfig returns the default configuration
func getDefaultConfig() *types.Config {
	return &types.Config{
		General: types.GeneralConfig{
			HomeURL:         "gemini://gemini.circumlunar.space/",
			SearchEngine:    "gemini://gus.guru/",
			MaxHistory:      1000,
			AutoSaveHistory: true,
		},
		UI: types.UIConfig{
			ShowLineNumbers: false,
			ShowLinkNumbers: true,
			EnableMouse:     true,
			ScrollSpeed:     3,
		},
		Colors: types.ColorConfig{
			Theme:            "default",
			LinkColor:        "12", // Blue
			VisitedLinkColor: "13", // Magenta
			Heading1Color:    "11", // Yellow
			Heading2Color:    "14", // Cyan
			Heading3Color:    "10", // Green
			TextColor:        "15", // White
			QuoteColor:       "8",  // Gray
			PreformatColor:   "7",  // Silver
			BackgroundColor:  "0",  // Black
		},
		Downloads: types.DownloadConfig{
			Directory:         "~/Downloads",
			AskBeforeDownload: true,
			MaxConcurrent:     3,
			Timeout:           30,
		},
	}
}

// Get returns the current configuration
func (c *Config) Get() *types.Config {
	return c.config
}

// Set updates the configuration
func (c *Config) Set(config *types.Config) {
	c.config = config
}

// Load loads configuration from disk
func (c *Config) Load() error {
	data, err := os.ReadFile(c.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Config file doesn't exist, create it with defaults
			return c.Save()
		}
		return err
	}

	var config types.Config
	if err := toml.Unmarshal(data, &config); err != nil {
		return err
	}

	// Merge with defaults to ensure all fields are present
	c.config = c.mergeWithDefaults(&config)

	return nil
}

// Save saves configuration to disk
func (c *Config) Save() error {
	// Ensure directory exists
	dir := filepath.Dir(c.configPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := toml.Marshal(c.config)
	if err != nil {
		return err
	}

	return os.WriteFile(c.configPath, data, 0600)
}

// mergeWithDefaults merges loaded config with defaults
func (c *Config) mergeWithDefaults(loaded *types.Config) *types.Config {
	defaults := getDefaultConfig()

	// General settings
	if loaded.General.HomeURL != "" {
		defaults.General.HomeURL = loaded.General.HomeURL
	}
	if loaded.General.SearchEngine != "" {
		defaults.General.SearchEngine = loaded.General.SearchEngine
	}
	if loaded.General.MaxHistory > 0 {
		defaults.General.MaxHistory = loaded.General.MaxHistory
	}
	defaults.General.AutoSaveHistory = loaded.General.AutoSaveHistory

	// UI settings
	defaults.UI.ShowLineNumbers = loaded.UI.ShowLineNumbers
	defaults.UI.ShowLinkNumbers = loaded.UI.ShowLinkNumbers
	defaults.UI.EnableMouse = loaded.UI.EnableMouse
	if loaded.UI.ScrollSpeed > 0 {
		defaults.UI.ScrollSpeed = loaded.UI.ScrollSpeed
	}

	// Color settings
	if loaded.Colors.Theme != "" {
		defaults.Colors.Theme = loaded.Colors.Theme
	}
	if loaded.Colors.LinkColor != "" {
		defaults.Colors.LinkColor = loaded.Colors.LinkColor
	}
	if loaded.Colors.VisitedLinkColor != "" {
		defaults.Colors.VisitedLinkColor = loaded.Colors.VisitedLinkColor
	}
	if loaded.Colors.Heading1Color != "" {
		defaults.Colors.Heading1Color = loaded.Colors.Heading1Color
	}
	if loaded.Colors.Heading2Color != "" {
		defaults.Colors.Heading2Color = loaded.Colors.Heading2Color
	}
	if loaded.Colors.Heading3Color != "" {
		defaults.Colors.Heading3Color = loaded.Colors.Heading3Color
	}
	if loaded.Colors.TextColor != "" {
		defaults.Colors.TextColor = loaded.Colors.TextColor
	}
	if loaded.Colors.QuoteColor != "" {
		defaults.Colors.QuoteColor = loaded.Colors.QuoteColor
	}
	if loaded.Colors.PreformatColor != "" {
		defaults.Colors.PreformatColor = loaded.Colors.PreformatColor
	}
	if loaded.Colors.BackgroundColor != "" {
		defaults.Colors.BackgroundColor = loaded.Colors.BackgroundColor
	}

	// Download settings
	if loaded.Downloads.Directory != "" {
		defaults.Downloads.Directory = loaded.Downloads.Directory
	}
	defaults.Downloads.AskBeforeDownload = loaded.Downloads.AskBeforeDownload
	if loaded.Downloads.MaxConcurrent > 0 {
		defaults.Downloads.MaxConcurrent = loaded.Downloads.MaxConcurrent
	}
	if loaded.Downloads.Timeout > 0 {
		defaults.Downloads.Timeout = loaded.Downloads.Timeout
	}

	return defaults
}

// GetDownloadDirectory returns the expanded download directory path
func (c *Config) GetDownloadDirectory() string {
	dir := c.config.Downloads.Directory
	if dir == "" {
		dir = "~/Downloads"
	}

	// Expand ~ to home directory
	if len(dir) >= 2 && dir[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, dir[2:])
		}
	}

	return dir
}
