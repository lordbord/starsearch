package themes

import "starsearch/internal/types"

// GetTheme returns a color configuration for the given theme name
func GetTheme(themeName string) *types.ColorConfig {
	switch themeName {
	case "dark":
		return &types.ColorConfig{
			Theme:             "dark",
			LinkColor:         "12",  // Bright blue
			VisitedLinkColor:  "13",  // Bright magenta
			Heading1Color:     "11",  // Bright yellow
			Heading2Color:     "14",  // Bright cyan
			Heading3Color:     "10",  // Bright green
			TextColor:         "15",  // White
			QuoteColor:        "8",   // Gray
			PreformatColor:    "7",   // Silver
			BackgroundColor:   "0",   // Black
		}
	case "light":
		return &types.ColorConfig{
			Theme:             "light",
			LinkColor:         "4",   // Blue
			VisitedLinkColor:  "5",   // Magenta
			Heading1Color:     "3",   // Yellow
			Heading2Color:     "6",   // Cyan
			Heading3Color:     "2",   // Green
			TextColor:         "0",   // Black
			QuoteColor:        "8",   // Gray
			PreformatColor:    "7",   // Silver
			BackgroundColor:   "15",  // White
		}
	case "solarized-dark":
		return &types.ColorConfig{
			Theme:             "solarized-dark",
			LinkColor:         "33",  // Blue
			VisitedLinkColor:  "35",  // Magenta
			Heading1Color:     "136", // Yellow
			Heading2Color:     "37",  // Cyan
			Heading3Color:     "64",  // Green
			TextColor:         "254", // Base0
			QuoteColor:        "66",  // Base01
			PreformatColor:    "244", // Base01
			BackgroundColor:   "235", // Base03
		}
	case "solarized-light":
		return &types.ColorConfig{
			Theme:             "solarized-light",
			LinkColor:         "33",  // Blue
			VisitedLinkColor:  "35",  // Magenta
			Heading1Color:     "136", // Yellow
			Heading2Color:     "37",  // Cyan
			Heading3Color:     "64",  // Green
			TextColor:         "235", // Base03
			QuoteColor:        "244", // Base01
			PreformatColor:    "66",  // Base01
			BackgroundColor:   "254", // Base0
		}
	case "monochrome":
		return &types.ColorConfig{
			Theme:             "monochrome",
			LinkColor:         "15",  // White
			VisitedLinkColor:  "7",   // Silver
			Heading1Color:     "15",  // White
			Heading2Color:     "15",  // White
			Heading3Color:     "15",  // White
			TextColor:         "15",  // White
			QuoteColor:        "7",   // Silver
			PreformatColor:    "7",   // Silver
			BackgroundColor:   "0",   // Black
		}
	case "nord":
		return &types.ColorConfig{
			Theme:             "nord",
			LinkColor:         "75",  // Nord Blue
			VisitedLinkColor:  "141", // Nord Purple
			Heading1Color:     "143", // Nord Green
			Heading2Color:     "109", // Nord Teal
			Heading3Color:     "180", // Nord Cyan
			TextColor:         "255", // Nord Snow Storm
			QuoteColor:        "243", // Nord Polar Night 3
			PreformatColor:    "244", // Nord Polar Night 4
			BackgroundColor:   "235", // Nord Polar Night 0
		}
	case "dracula":
		return &types.ColorConfig{
			Theme:             "dracula",
			LinkColor:         "63",  // Dracula Purple
			VisitedLinkColor:  "141", // Dracula Pink
			Heading1Color:     "11",  // Dracula Yellow
			Heading2Color:     "81",  // Dracula Cyan
			Heading3Color:     "84",  // Dracula Green
			TextColor:         "255", // Dracula Foreground
			QuoteColor:        "241", // Dracula Comment
			PreformatColor:    "241", // Dracula Comment
			BackgroundColor:   "235", // Dracula Background
		}
	default:
		// Default theme
		return &types.ColorConfig{
			Theme:             "default",
			LinkColor:         "12",  // Blue
			VisitedLinkColor:  "13",  // Magenta
			Heading1Color:     "11",  // Yellow
			Heading2Color:     "14",  // Cyan
			Heading3Color:     "10",  // Green
			TextColor:         "15",  // White
			QuoteColor:        "8",   // Gray
			PreformatColor:    "7",   // Silver
			BackgroundColor:   "0",   // Black
		}
	}
}

// ApplyTheme applies theme colors to a config, replacing all colors with theme values
func ApplyTheme(config *types.ColorConfig, themeName string) {
	theme := GetTheme(themeName)
	
	// Apply all theme colors
	config.LinkColor = theme.LinkColor
	config.VisitedLinkColor = theme.VisitedLinkColor
	config.Heading1Color = theme.Heading1Color
	config.Heading2Color = theme.Heading2Color
	config.Heading3Color = theme.Heading3Color
	config.TextColor = theme.TextColor
	config.QuoteColor = theme.QuoteColor
	config.PreformatColor = theme.PreformatColor
	config.BackgroundColor = theme.BackgroundColor
	config.Theme = themeName
}

// getDefaultConfig returns the default color configuration
func getDefaultConfig() *types.ColorConfig {
	return &types.ColorConfig{
		Theme:             "default",
		LinkColor:         "12",
		VisitedLinkColor:  "13",
		Heading1Color:     "11",
		Heading2Color:     "14",
		Heading3Color:     "10",
		TextColor:         "15",
		QuoteColor:        "8",
		PreformatColor:    "7",
		BackgroundColor:   "0",
	}
}

// GetAvailableThemes returns a list of available theme names
func GetAvailableThemes() []string {
	return []string{
		"default",
		"dark",
		"light",
		"solarized-dark",
		"solarized-light",
		"monochrome",
		"nord",
		"dracula",
	}
}

