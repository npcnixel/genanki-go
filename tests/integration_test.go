package tests

import (
	"os"
	"testing"

	genanki "github.com/npcnixel/genanki-go"
)

func TestBasicIntegration(t *testing.T) {
	// Create a new database
	db, err := genanki.NewDatabase()
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Create a model
	basicModel := genanki.NewBasicModel(1234567890, "Test Basic Model")
	if err := db.AddModel(basicModel.Model); err != nil {
		t.Fatalf("Failed to add model: %v", err)
	}

	// Create a deck
	deck := genanki.NewDeck(9876543210, "Test Deck", "Test Description")
	if err := db.AddDeck(deck); err != nil {
		t.Fatalf("Failed to add deck: %v", err)
	}

	// Create a note
	note := genanki.NewNote(
		basicModel.Model.ID,
		[]string{"Test Question", "Test Answer"},
		[]string{"test", "basic"},
	)
	if err := db.AddNote(note); err != nil {
		t.Fatalf("Failed to add note: %v", err)
	}

	// Create a card
	if err := db.AddCard(note.ID, deck.ID, 0); err != nil {
		t.Fatalf("Failed to add card: %v", err)
	}

	// Create a package
	pkg := genanki.NewPackage(db)

	// Add media
	mediaData := []byte("test media data")
	pkg.AddMedia("test_media.txt", mediaData)

	// Write to file
	tmpFile, err := os.CreateTemp("", "test_deck_*.apkg")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	if err := pkg.WriteToFile(tmpFile.Name()); err != nil {
		t.Fatalf("Failed to write package: %v", err)
	}

	// Verify file exists and has content
	fileInfo, err := os.Stat(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	if fileInfo.Size() == 0 {
		t.Error("Expected file to have content")
	}
}

func TestMultipleNotesIntegration(t *testing.T) {
	// Create a new database
	db, err := genanki.NewDatabase()
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Create a basic model
	basicModel := genanki.NewBasicModel(1234567890, "Test Basic Model")
	if err := db.AddModel(basicModel.Model); err != nil {
		t.Fatalf("Failed to add basic model: %v", err)
	}

	// Create a deck
	deck := genanki.NewDeck(9876543210, "Test Deck", "Test Description")
	if err := db.AddDeck(deck); err != nil {
		t.Fatalf("Failed to add deck: %v", err)
	}

	// Create multiple notes
	notes := []*genanki.Note{
		genanki.NewNote(basicModel.Model.ID, []string{"Question 1", "Answer 1"}, []string{"test"}),
		genanki.NewNote(basicModel.Model.ID, []string{"Question 2", "Answer 2"}, []string{"test"}),
		genanki.NewNote(basicModel.Model.ID, []string{"Question 3", "Answer 3"}, []string{"test"}),
	}

	// Add notes and cards to database
	for _, note := range notes {
		if err := db.AddNote(note); err != nil {
			t.Fatalf("Failed to add note: %v", err)
		}

		if err := db.AddCard(note.ID, deck.ID, 0); err != nil {
			t.Fatalf("Failed to add card: %v", err)
		}
	}

	// Create package
	pkg := genanki.NewPackage(db)

	// Write to file
	tmpFile, err := os.CreateTemp("", "test_multi_*.apkg")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	if err := pkg.WriteToFile(tmpFile.Name()); err != nil {
		t.Fatalf("Failed to write package: %v", err)
	}

	// Verify file exists and has content
	fileInfo, err := os.Stat(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	if fileInfo.Size() == 0 {
		t.Error("Expected file to have content")
	}
}

func TestMultipleModelsDeckIntegration(t *testing.T) {
	// Create a new database
	db, err := genanki.NewDatabase()
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Create models
	basicModel := genanki.NewBasicModel(1234567890, "Basic Model")
	if err := db.AddModel(basicModel.Model); err != nil {
		t.Fatalf("Failed to add basic model: %v", err)
	}

	// Create decks
	deck1 := genanki.NewDeck(1111111111, "Geography Deck", "A deck about geography")
	if err := db.AddDeck(deck1); err != nil {
		t.Fatalf("Failed to add deck1: %v", err)
	}

	deck2 := genanki.NewDeck(2222222222, "Math Deck", "A deck about math")
	if err := db.AddDeck(deck2); err != nil {
		t.Fatalf("Failed to add deck2: %v", err)
	}

	// Create notes for geography deck
	geoNotes := []*genanki.Note{
		genanki.NewNote(basicModel.Model.ID, []string{"What is the capital of France?", "Paris"}, []string{"geography"}),
		genanki.NewNote(basicModel.Model.ID, []string{"What is the capital of Italy?", "Rome"}, []string{"geography"}),
	}

	// Create notes for math deck
	mathNotes := []*genanki.Note{
		genanki.NewNote(basicModel.Model.ID, []string{"What is 2+2?", "4"}, []string{"math"}),
		genanki.NewNote(basicModel.Model.ID, []string{"What is 3×3?", "9"}, []string{"math"}),
	}

	// Add geography notes and cards
	for _, note := range geoNotes {
		if err := db.AddNote(note); err != nil {
			t.Fatalf("Failed to add geo note: %v", err)
		}

		if err := db.AddCard(note.ID, deck1.ID, 0); err != nil {
			t.Fatalf("Failed to add card for geo note: %v", err)
		}
	}

	// Add math notes and cards
	for _, note := range mathNotes {
		if err := db.AddNote(note); err != nil {
			t.Fatalf("Failed to add math note: %v", err)
		}

		if err := db.AddCard(note.ID, deck2.ID, 0); err != nil {
			t.Fatalf("Failed to add card for math note: %v", err)
		}
	}

	// Create package
	pkg := genanki.NewPackage(db)

	// Write to file
	tmpFile, err := os.CreateTemp("", "test_multi_models_decks_*.apkg")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	if err := pkg.WriteToFile(tmpFile.Name()); err != nil {
		t.Fatalf("Failed to write package: %v", err)
	}

	// Verify file exists and has content
	fileInfo, err := os.Stat(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	if fileInfo.Size() == 0 {
		t.Error("Expected file to have content")
	}
}
