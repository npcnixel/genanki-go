package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/npcnixel/genanki-go"
)

func main() {
	// Create a basic model with auto-generated ID using convenience function
	basicModel := genanki.StandardBasicModel("Basic")

	// Create a new deck with auto-generated ID using convenience function
	deck := genanki.StandardDeck("Test Deck", "A test deck")

	// Print the generated IDs for reference
	fmt.Printf("Generated Basic Model ID: %d\n", basicModel.ID)
	fmt.Printf("Generated Deck ID: %d\n", deck.ID)

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

	// Ensure output directory exists at same level as example directories
	outputDir := filepath.Join("..", "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Write package to file in the output directory
	outputPath := filepath.Join(outputDir, "basic_deck.apkg")
	if err := pkg.WriteToFile(outputPath); err != nil {
		log.Fatalf("Failed to write package: %v", err)
	}

	// Print summary information
	fmt.Printf("Successfully created Anki deck: %s\n", outputPath)
	fmt.Printf("Number of notes: %d\n", len(deck.Notes))
}
