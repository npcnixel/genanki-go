package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/npcnixel/genanki-go"
)

func main() {
	// Parse command line flags
	debugFlag := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	// Enable debug logging if either flag or environment variable is set
	debug := *debugFlag || strings.ToLower(os.Getenv("DEBUG")) == "true"
	if debug {
		log.Printf("Debug mode enabled")
	}

	// Create a basic model with auto-generated ID using convenience function
	basicModel := genanki.StandardBasicModel("Basic")

	// Create a new deck with auto-generated ID using convenience function
	deck := genanki.StandardDeck("Test Deck", "A test deck")

	// Create a note
	note := genanki.NewNote(
		basicModel.ID,
		[]string{"What is 2+2?", "4"},
		[]string{"math", "basic"},
	)

	// Add note to the deck using chaining
	deck.AddNote(note)

	// Create a package with the deck using chaining
	pkg := genanki.NewPackage([]*genanki.Deck{deck}).AddModel(basicModel.Model)

	// Enable debug logging
	if debug {
		pkg.SetDebug(true)
	}

	// Check if output path is specified in environment
	outputPath := os.Getenv("OUTPUT_PATH")
	if outputPath == "" {
		// Use default path if no environment variable is set
		outputDir := filepath.Join("..", "output")
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}
		outputPath = filepath.Join(outputDir, "basic_deck.apkg")
	} else {
		// Ensure the directory for the specified path exists
		outputDir := filepath.Dir(outputPath)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}
	}

	// Write package to file
	if err := pkg.WriteToFile(outputPath); err != nil {
		log.Fatalf("Failed to write package: %v", err)
	}
}
