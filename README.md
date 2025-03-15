# genanki-go

A Go library for generating Anki decks programmatically. This is a port of the Python [genanki](https://github.com/kerrickstaley/genanki) library.

## Installation

```bash
go get github.com/npcnixel/genanki-go
```

## Features

- Create Anki decks programmatically
- Support for basic card types
- Support for cloze deletion cards
- Include media files (images, audio) (TODO)
- Generate Anki-compatible .apkg files
- Compatible with Anki 25.02 and newer versions

## Requirements

- Go 1.16 or later
- [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) for SQLite support

## Basic Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/npcnixel/genanki-go/pkg/genanki"
)

func main() {
	// Create a basic model (front/back)
	basicModel := genanki.NewBasicModel(
		genanki.GenerateIntID(),
		"Basic Model",
	)

	// Create a new deck
	deck := genanki.NewDeck(
		genanki.GenerateDeckID(),
		"My First Deck",
		"A simple example deck",
	)

	// Add a note to the deck
	note := genanki.NewNote(
		basicModel.ID,
		[]string{"What is the capital of France?", "Paris"},
		[]string{"geography", "europe"},
	)
	deck.AddNote(note)

	// Write the deck to an .apkg file
	err := deck.WriteToApkg("my_first_deck.apkg", []*genanki.Model{basicModel})
	if err != nil {
		log.Fatalf("Error creating Anki deck: %v", err)
	}

	fmt.Println("Successfully created Anki deck: my_first_deck.apkg")
}
```

## Advanced Features

### Using Images

```go
// Add a note with an image
imageNote := genanki.NewNote(
	basicModel.ID,
	[]string{
		"What famous landmark is this?",
		"The Eiffel Tower<br><img src=\"eiffel_tower.jpg\">",
	},
	[]string{"landmarks", "europe"},
)
deck.AddNote(imageNote)

// Add the image file to the deck's media
deck.AddMedia("eiffel_tower.jpg", "/path/to/eiffel_tower.jpg")
```

### Creating Cloze Deletion Cards

```go
// Create a cloze model
clozeModel := genanki.NewClozeModel(
	genanki.GenerateIntID(),
	"Cloze Model",
)

// Add a cloze note
clozeNote := genanki.NewNote(
	clozeModel.ID,
	[]string{
		"The capital of France is {{c1::Paris}}.",
		"Additional information about Paris.",
	},
	[]string{"geography", "cloze"},
)
deck.AddNote(clozeNote)
```

## Examples

See the `examples` directory for complete examples:

- `examples/basic_example`: Simple example creating a basic deck
- `examples/advanced_example`: Advanced example demonstrating cloze deletions and media

## License

[MIT License](LICENSE)

## Acknowledgements

This project is a Go port of the Python [genanki](https://github.com/kerrickstaley/genanki) library by Kerrick Staley.
