package genanki

import (
	"archive/zip"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Package struct {
	decks  []*Deck
	models []*Model
	media  map[string][]byte
	db     *Database
}

// NewPackage creates a new package from decks or a database
func NewPackage(data interface{}) *Package {
	switch v := data.(type) {
	case []*Deck:
		return &Package{
			decks:  v,
			models: make([]*Model, 0),
			media:  make(map[string][]byte),
		}
	case *Database:
		return &Package{
			db:     v,
			decks:  make([]*Deck, 0),
			models: make([]*Model, 0),
			media:  make(map[string][]byte),
		}
	default:
		panic("NewPackage: unsupported type")
	}
}

// AddModel adds a model to the package
func (p *Package) AddModel(model *Model) {
	p.models = append(p.models, model)
}

func (p *Package) AddMedia(filename string, data []byte) {
	p.media[filename] = data
}

func (p *Package) WriteToFile(path string) error {
	var dbToUse *Database
	var err error

	// Determine which database to use
	if p.db != nil {
		// Using existing database
		dbToUse = p.db
	} else {
		// Create a new database for the package
		dbToUse, err = NewDatabase()
		if err != nil {
			return fmt.Errorf("failed to create database: %v", err)
		}
		defer dbToUse.Close()

		// Add all models
		for _, model := range p.models {
			if err := dbToUse.AddModel(model); err != nil {
				return fmt.Errorf("failed to add model to database: %v", err)
			}
		}

		// Add all decks
		for _, deck := range p.decks {
			if err := dbToUse.AddDeck(deck); err != nil {
				return fmt.Errorf("failed to add deck to database: %v", err)
			}

			// Add all notes from this deck
			for _, note := range deck.Notes {
				if err := dbToUse.AddNote(note); err != nil {
					return fmt.Errorf("failed to add note to database: %v", err)
				}

				// Add a card for each note
				if err := dbToUse.AddCard(note.ID, deck.ID, 0); err != nil {
					return fmt.Errorf("failed to add card to database: %v", err)
				}
			}

			// Add deck media to package media
			for filename, data := range deck.Media {
				p.media[filename] = data
			}
		}
	}

	// Create the output file
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	w1, err := zw.Create("collection.anki2")
	if err != nil {
		return fmt.Errorf("failed to create collection.anki2: %v", err)
	}

	dbFile, err := dbToUse.GetFilePath()
	if err != nil {
		return fmt.Errorf("failed to get database file: %v", err)
	}
	defer os.Remove(dbFile)

	dbContent, err := os.ReadFile(dbFile)
	if err != nil {
		return fmt.Errorf("failed to read database: %v", err)
	}
	if _, err := w1.Write(dbContent); err != nil {
		return fmt.Errorf("failed to write collection.anki2: %v", err)
	}

	mediaMap := make(map[string]string)
	for filename, data := range p.media {
		mediaFilename := fmt.Sprintf("%d", len(mediaMap))
		mediaMap[mediaFilename] = filename

		w, err := zw.Create(mediaFilename)
		if err != nil {
			return fmt.Errorf("failed to create media file: %v", err)
		}
		if _, err := w.Write(data); err != nil {
			return fmt.Errorf("failed to write media file: %v", err)
		}
	}

	w3, err := zw.Create("media")
	if err != nil {
		return fmt.Errorf("failed to create media: %v", err)
	}
	mediaJSON, err := json.Marshal(mediaMap)
	if err != nil {
		return fmt.Errorf("failed to marshal media: %v", err)
	}
	if _, err := w3.Write(mediaJSON); err != nil {
		return fmt.Errorf("failed to write media: %v", err)
	}

	return nil
}

func GenerateMediaHash(data []byte) string {
	hash := sha1.Sum(data)
	return hex.EncodeToString(hash[:])
}

func SanitizeFilename(filename string) string {
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	sanitized := filename
	for _, char := range invalid {
		sanitized = strings.ReplaceAll(sanitized, char, "_")
	}
	return sanitized
}
