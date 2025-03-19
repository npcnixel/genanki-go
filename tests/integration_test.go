package tests

import (
	"os"
	"testing"

	genanki "github.com/npcnixel/genanki-go"
)

func TestBasicIntegration(t *testing.T) {
	// Create a model
	basicModel := genanki.NewBasicModel(1234567890, "Test Basic Model")

	// Create a deck
	deck := genanki.NewDeck(9876543210, "Test Deck", "Test Description")

	// Create a note
	note := genanki.NewNote(
		basicModel.Model.ID,
		[]string{"Test Question", "Test Answer"},
		[]string{"test", "basic"},
	)

	// Add note to deck
	deck.AddNote(note)

	// Create a package
	pkg := genanki.NewPackage([]*genanki.Deck{deck})
	pkg.AddModel(basicModel.Model)

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
	// Create a basic model
	basicModel := genanki.NewBasicModel(1234567890, "Test Basic Model")

	// Create a deck
	deck := genanki.NewDeck(9876543210, "Test Deck", "Test Description")

	// Create multiple notes and add them to the deck
	deck.AddNote(genanki.NewNote(basicModel.Model.ID, []string{"Question 1", "Answer 1"}, []string{"test"}))
	deck.AddNote(genanki.NewNote(basicModel.Model.ID, []string{"Question 2", "Answer 2"}, []string{"test"}))
	deck.AddNote(genanki.NewNote(basicModel.Model.ID, []string{"Question 3", "Answer 3"}, []string{"test"}))

	// Create package
	pkg := genanki.NewPackage([]*genanki.Deck{deck})
	pkg.AddModel(basicModel.Model)

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
	// Create a basic model
	basicModel := genanki.NewBasicModel(0, "Geography Model")

	// Create two decks
	geoDeck := genanki.NewDeck(0, "Geography Deck", "A deck for geography flashcards")
	mathDeck := genanki.NewDeck(2059400110, "Math Deck", "A deck for math flashcards")

	// Create notes for geography
	geoNote1 := genanki.NewNote(basicModel.ID, []string{"What is the capital of France?", "Paris"}, []string{"geography", "europe"})
	geoNote2 := genanki.NewNote(basicModel.ID, []string{"What is the capital of Italy?", "Rome"}, []string{"geography", "europe"})

	// Create notes for math
	mathNote1 := genanki.NewNote(basicModel.ID, []string{"What is 2+2?", "4"}, []string{"math", "addition"})
	mathNote2 := genanki.NewNote(basicModel.ID, []string{"What is 3×3?", "9"}, []string{"math", "multiplication"})

	// Add notes to decks
	geoDeck.AddNote(geoNote1).AddNote(geoNote2)
	mathDeck.AddNote(mathNote1).AddNote(mathNote2)

	// Create a package with multiple decks
	pkg := genanki.NewPackage([]*genanki.Deck{geoDeck, mathDeck})
	pkg.AddModel(basicModel.Model)

	// Write the package to a temporary file
	tmpFile, err := os.CreateTemp("", "anki-*.apkg")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if err := pkg.WriteToFile(tmpFile.Name()); err != nil {
		t.Fatalf("Failed to write package: %v", err)
	}

	// Verify the file exists and has content
	stat, err := os.Stat(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}
	if stat.Size() == 0 {
		t.Fatalf("File has zero size")
	}
}

// Test using the new convenience functions
func TestNewAPIPattern(t *testing.T) {
	// Use convenience functions to create model and deck
	basicModel := genanki.StandardBasicModel("Test Basic Model")
	deck := genanki.StandardDeck("Test Standard Deck", "A test deck using standard IDs")

	// Create and add notes to deck with method chaining
	deck.
		AddNote(genanki.NewNote(basicModel.ID, []string{"Question 1", "Answer 1"}, []string{"test"})).
		AddNote(genanki.NewNote(basicModel.ID, []string{"Question 2", "Answer 2"}, []string{"test"}))

	// Verify notes were added to deck
	if len(deck.Notes) != 2 {
		t.Errorf("Expected 2 notes in deck, got %d", len(deck.Notes))
	}

	// Create package and add model with method chaining
	pkg := genanki.NewPackage([]*genanki.Deck{deck}).
		AddModel(basicModel.Model).
		AddMedia("test_media.txt", []byte("Test content"))

	// Write to a temporary file
	tmpFile, err := os.CreateTemp("", "api-pattern-*.apkg")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if err := pkg.WriteToFile(tmpFile.Name()); err != nil {
		t.Fatalf("Failed to write package: %v", err)
	}

	// Verify the file exists and has content
	stat, err := os.Stat(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}
	if stat.Size() == 0 {
		t.Fatalf("File has zero size")
	}

	t.Logf("Successfully created Anki deck with new API: %s (%d bytes)",
		tmpFile.Name(), stat.Size())
}
