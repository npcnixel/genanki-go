# genanki-go

A Go library for generating Anki decks programmatically.

## Features

- Create Anki decks with notes and cards
- Support for basic and cloze models
- Add media files (images, audio, video)
- Generate `.apkg` files for Anki import
- Simple and intuitive API

## Installation

```bash
go get github.com/npcnixel/genanki-go
```

## Quick Start

Here's a simple example of creating a basic Anki deck:

```go
package main

import (
    "github.com/npcnixel/genanki-go"
)

func main() {
    // Create a basic model
    model := genanki.StandardBasicModel("My Model")

    // Create a deck
    deck := genanki.StandardDeck("My Deck", "A deck for testing")

    // Create a note
    note := genanki.NewNote(model.ID, []string{"What is the capital of France?", "Paris"}, []string{"geography"})
    
    // Add note to deck
    deck.AddNote(note)
    
    // Create a package with the deck
    pkg := genanki.NewPackage([]*genanki.Deck{deck})
    
    // Add the model to the package
    pkg.AddModel(model.Model)
    
    // Write the package to a file
    pkg.WriteToFile("output.apkg")
}
```

## Advanced Usage

### Creating Different Types of Models

```go
// Create a basic model
basicModel := genanki.StandardBasicModel("Basic Model")

// Create a cloze model
clozeModel := genanki.StandardClozeModel("Cloze Model")

// Create a deck
deck := genanki.StandardDeck("My Deck", "A deck for testing")
```

### Adding Media Files

```go
// Create a package
pkg := genanki.NewPackage([]*genanki.Deck{deck})

// Add models to the package
pkg.AddModel(basicModel.Model)
pkg.AddModel(clozeModel.Model)

// Add an image file
imageData, _ := ioutil.ReadFile("image.jpg")
pkg.AddMedia("image.jpg", imageData)

// Add an audio file
audioData, _ := ioutil.ReadFile("audio.mp3")
pkg.AddMedia("audio.mp3", audioData)
```

### Creating Cloze Notes

```go
// Create a cloze note
note := genanki.NewNote(clozeModel.ID, []string{
    "The capital of France is {{c1::Paris}}.",
    "The capital of Spain is {{c1::Madrid}}.",
}, []string{"geography"})
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.