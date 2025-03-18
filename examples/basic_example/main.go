package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/npcnixel/genanki-go"
)

func main() {

	db, err := genanki.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	basicModel := genanki.NewBasicModel(
		1234567890,
		"Basic",
	)

	if err := db.AddModel(basicModel.Model); err != nil {
		log.Fatalf("Failed to add model: %v", err)
	}
	log.Printf("Added model to database")

	deck := genanki.NewDeck(
		1122334455,
		"Test Deck",
		"A test deck",
	)

	if err := db.AddDeck(deck); err != nil {
		log.Fatalf("Failed to add deck: %v", err)
	}
	log.Printf("Added deck to database")

	note := genanki.NewNote(
		basicModel.ID,
		[]string{
			"What is 2+2?",
			"4",
		},
		[]string{"math", "basic"},
	)

	deck.AddNote(note)

	if err := db.AddNote(note); err != nil {
		log.Fatalf("Failed to add note: %v", err)
	}
	log.Printf("Added note to database")

	if err := db.AddCard(note.ID, deck.ID, 0); err != nil {
		log.Fatalf("Failed to add card: %v", err)
	}
	log.Printf("Added card to database")

	pkg := genanki.NewPackage(db)

	log.Println("Verifying database content...")
	if err := db.VerifyContent(); err != nil {
		log.Printf("Warning: Failed to verify database content: %v", err)
	}

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

	fmt.Printf("Successfully created Anki deck: %s\n", outputPath)
	fmt.Printf("Number of notes: %d\n", len(deck.Notes))
}
