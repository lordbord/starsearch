package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"starsearch/internal/app"
)

func main() {
	// Get initial URL from command-line arguments if provided
	var initialURL string
	if len(os.Args) > 1 {
		initialURL = os.Args[1]
	}

	// Create the application model
	model, err := app.NewModel(initialURL)
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
