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

	// Create a basic model
	basicModel := genanki.NewBasicModel(0, "Geography Model")
	if err := db.AddModel(basicModel.Model); err != nil {
		t.Fatalf("Failed to add model: %v", err)
	}

	// Create two decks
	geoDeck := genanki.NewDeck(0, "Geography Deck", "A deck for geography flashcards")
	mathDeck := genanki.NewDeck(2059400110, "Math Deck", "A deck for math flashcards")

	if err := db.AddDeck(geoDeck); err != nil {
		t.Fatalf("Failed to add geography deck: %v", err)
	}
	if err := db.AddDeck(mathDeck); err != nil {
		t.Fatalf("Failed to add math deck: %v", err)
	}

	// Create notes for geography
	geoNote1 := genanki.NewNote(basicModel.ID, []string{"What is the capital of France?", "Paris"}, []string{"geography", "europe"})
	geoNote2 := genanki.NewNote(basicModel.ID, []string{"What is the capital of Italy?", "Rome"}, []string{"geography", "europe"})

	// Create notes for math
	mathNote1 := genanki.NewNote(basicModel.ID, []string{"What is 2+2?", "4"}, []string{"math", "addition"})
	mathNote2 := genanki.NewNote(basicModel.ID, []string{"What is 3×3?", "9"}, []string{"math", "multiplication"})

	// Add notes to the database
	for _, note := range []*genanki.Note{geoNote1, geoNote2, mathNote1, mathNote2} {
		if err := db.AddNote(note); err != nil {
			t.Fatalf("Failed to add note: %v", err)
		}
	}

	// Add cards for geography deck
	if err := db.AddCard(geoNote1.ID, geoDeck.ID, 0); err != nil {
		t.Fatalf("Failed to add card: %v", err)
	}
	if err := db.AddCard(geoNote2.ID, geoDeck.ID, 0); err != nil {
		t.Fatalf("Failed to add card: %v", err)
	}

	// Add cards for math deck
	if err := db.AddCard(mathNote1.ID, mathDeck.ID, 0); err != nil {
		t.Fatalf("Failed to add card: %v", err)
	}
	if err := db.AddCard(mathNote2.ID, mathDeck.ID, 0); err != nil {
		t.Fatalf("Failed to add card: %v", err)
	}

	// Create a package from the database
	pkg := genanki.NewPackage(db)

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
	// Create a new database
	db, err := genanki.NewDatabase()
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Use convenience functions to create model and deck
	basicModel := genanki.StandardBasicModel("Test Basic Model")
	deck := genanki.StandardDeck("Test Standard Deck", "A test deck using standard IDs")

	// Add model to the database
	if err := db.AddModel(basicModel.Model); err != nil {
		t.Fatalf("Failed to add model: %v", err)
	}

	// Add deck to the database
	if err := db.AddDeck(deck); err != nil {
		t.Fatalf("Failed to add deck: %v", err)
	}

	// Create and add notes
	note1 := genanki.NewNote(basicModel.ID, []string{"Question 1", "Answer 1"}, []string{"test"})
	note2 := genanki.NewNote(basicModel.ID, []string{"Question 2", "Answer 2"}, []string{"test"})

	// Add notes to the database
	for _, note := range []*genanki.Note{note1, note2} {
		if err := db.AddNote(note); err != nil {
			t.Fatalf("Failed to add note: %v", err)
		}
	}

	// Add cards to the deck
	if err := db.AddCard(note1.ID, deck.ID, 0); err != nil {
		t.Fatalf("Failed to add card: %v", err)
	}
	if err := db.AddCard(note2.ID, deck.ID, 0); err != nil {
		t.Fatalf("Failed to add card: %v", err)
	}

	// Create a package from the database
	pkg := genanki.NewPackage(db)

	// Write the package to a temporary file
	tmpFile, err := os.CreateTemp("", "anki-standard-*.apkg")
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
