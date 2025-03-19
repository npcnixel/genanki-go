package tests

import (
	"os"
	"testing"

	"github.com/npcnixel/genanki-go"
)

func TestChainedAPI(t *testing.T) {
	// Create a basic model and deck
	basicModel := genanki.StandardBasicModel("Chained API Test Model")
	deck := genanki.StandardDeck("Chained API Test Deck", "Testing the chained API")

	// Create notes
	note1 := genanki.NewNote(basicModel.ID, []string{"Question 1", "Answer 1"}, []string{"test"})
	note2 := genanki.NewNote(basicModel.ID, []string{"Question 2", "Answer 2"}, []string{"test"})

	// Add notes to deck using chaining
	deck.AddNote(note1).AddNote(note2)

	// Verify notes were added to deck
	if len(deck.Notes) != 2 {
		t.Errorf("Expected 2 notes in deck, got %d", len(deck.Notes))
	}

	// Create a package with direct deck approach
	pkg := genanki.NewPackage([]*genanki.Deck{deck})
	pkg.AddModel(basicModel.Model)
	pkg.AddMedia("test_file.txt", []byte("Test content"))

	// Write the package to a temporary file
	tmpFile, err := os.CreateTemp("", "anki-chained-*.apkg")
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

	t.Logf("Successfully created Anki deck with chained API: %s", tmpFile.Name())
	t.Logf("File size: %d bytes", stat.Size())
}

func TestModelCustomizationChaining(t *testing.T) {
	// Create a model
	model := genanki.StandardBasicModel("Custom Chained Model")

	// Customize the model using method chaining
	model.Model.SetCSS(`
		.card { 
			font-family: Arial; 
			font-size: 20px;
		}
		.question { 
			font-weight: bold; 
		}
	`).AddField(genanki.Field{
		Name: "Extra Info",
		Font: "Arial",
		Size: 16,
	}).AddTemplate(genanki.Template{
		Name: "Extra Template",
		Qfmt: "{{Extra Info}}?",
		Afmt: "{{Front}}<hr>{{Back}}<br>{{Extra Info}}",
	})

	// Verify customizations
	if len(model.Model.Fields) != 3 { // Front, Back, Extra Info
		t.Errorf("Expected 3 fields, got %d", len(model.Model.Fields))
	}

	if len(model.Model.Templates) != 2 { // Card 1 and Extra Template
		t.Errorf("Expected 2 templates, got %d", len(model.Model.Templates))
	}

	if model.Model.Fields[2].Name != "Extra Info" {
		t.Errorf("Expected field name 'Extra Info', got '%s'", model.Model.Fields[2].Name)
	}

	if model.Model.Templates[1].Name != "Extra Template" {
		t.Errorf("Expected template name 'Extra Template', got '%s'", model.Model.Templates[1].Name)
	}
}

func TestPackageChaining(t *testing.T) {
	// Create models and deck
	basicModel := genanki.StandardBasicModel("Package Chaining Model")
	clozeModel := genanki.StandardClozeModel("Package Chaining Cloze")
	deck1 := genanki.StandardDeck("Deck 1", "First test deck")
	deck2 := genanki.StandardDeck("Deck 2", "Second test deck")

	// Add notes to decks using chaining
	note1 := genanki.NewNote(basicModel.ID, []string{"Q1", "A1"}, nil)
	note2 := genanki.NewNote(basicModel.ID, []string{"Q2", "A2"}, nil)
	note3 := genanki.NewNote(clozeModel.ID, []string{"This is a {{c1::cloze}} test", ""}, nil)

	deck1.AddNote(note1).AddNote(note2)
	deck2.AddNote(note3)

	// Create package with chaining
	pkg := genanki.NewPackage([]*genanki.Deck{deck1, deck2}).
		AddModel(basicModel.Model).
		AddModel(clozeModel.Model).
		AddMedia("test1.txt", []byte("Test file 1")).
		AddMedia("test2.txt", []byte("Test file 2"))

	// Write to temp file and verify
	tmpFile, err := os.CreateTemp("", "anki-package-chain-*.apkg")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if err := pkg.WriteToFile(tmpFile.Name()); err != nil {
		t.Fatalf("Failed to write package: %v", err)
	}

	// Verify file was created successfully
	stat, err := os.Stat(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}
	if stat.Size() == 0 {
		t.Fatalf("File has zero size")
	}

	t.Logf("Successfully created package with chained API: %s (%d bytes)",
		tmpFile.Name(), stat.Size())
}
