package tests

import (
	"testing"
	"time"

	genanki "github.com/npcnixel/genanki-go"
)

func TestNewDeck(t *testing.T) {
	deckID := int64(1234567890)
	deckName := "Test Deck"
	deckDesc := "A test deck"

	deck := genanki.NewDeck(deckID, deckName, deckDesc)

	if deck.ID != deckID {
		t.Errorf("Expected deck ID %d, got %d", deckID, deck.ID)
	}

	if deck.Name != deckName {
		t.Errorf("Expected deck name %s, got %s", deckName, deck.Name)
	}

	if deck.Desc != deckDesc {
		t.Errorf("Expected deck description %s, got %s", deckDesc, deck.Desc)
	}

	if len(deck.Notes) != 0 {
		t.Errorf("Expected 0 notes, got %d", len(deck.Notes))
	}

	zeroTime := time.Time{}
	if deck.Created == zeroTime {
		t.Error("Expected created time to be set")
	}

	if deck.Modified == zeroTime {
		t.Error("Expected modified time to be set")
	}
}

func TestAddNote(t *testing.T) {
	deck := genanki.NewDeck(1234567890, "Test Deck", "A test deck")
	note := genanki.NewNote(9876543210, []string{"Question", "Answer"}, []string{"test", "basic"})

	initialModified := deck.Modified
	time.Sleep(1 * time.Millisecond) // Ensure time difference

	deck.AddNote(note)

	if len(deck.Notes) != 1 {
		t.Errorf("Expected 1 note, got %d", len(deck.Notes))
	}

	if deck.Notes[0] != note {
		t.Error("Expected note to be added to deck")
	}

	if deck.Modified.Equal(initialModified) {
		t.Error("Expected modified time to be updated")
	}
}

func TestAddMedia(t *testing.T) {
	deck := genanki.NewDeck(1234567890, "Test Deck", "A test deck")

	filename := "test.jpg"
	data := []byte("test image data")

	initialModified := deck.Modified
	time.Sleep(1 * time.Millisecond) // Ensure time difference

	deck.AddMedia(filename, data)

	if len(deck.Media) != 1 {
		t.Errorf("Expected 1 media item, got %d", len(deck.Media))
	}

	if string(deck.Media[filename]) != string(data) {
		t.Error("Expected media data to match provided data")
	}

	if deck.Modified.Equal(initialModified) {
		t.Error("Expected modified time to be updated")
	}
}
