package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"starsearch/internal/app"
)

const version = "0.1.2"

func main() {
	// Handle version flag
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("starsearch v%s\n", version)
		os.Exit(0)
	}

	// Get initial URL from command-line arguments if provided
	var initialURL string
	if len(os.Args) > 1 {
		initialURL = os.Args[1]
	}

	// Create the application model with version
	model, err := app.NewModel(initialURL, version)
	if err != nil {
		log.Fatal(err)
	}

	// Create the Bubble Tea program with alternate screen buffer
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
