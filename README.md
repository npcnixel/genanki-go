# genanki-go

A Go library for generating Anki decks programmatically.

## Features

- Create Anki decks with notes and cards
- Support for basic and cloze models
- Add media files (images, audio, video)
- Generate `.apkg` files for Anki import
- Simple and intuitive API
- **Method chaining for more fluent API usage**
- Optional debug logging for troubleshooting

## Installation

```bash
go get github.com/npcnixel/genanki-go
```

## Quick Start

Here's a simple example of creating a basic Anki deck:

```go
package main

import (
	"fmt"
	"os"

	"github.com/npcnixel/genanki-go"
)

func main() {
	// Create a basic model and deck
	model := genanki.StandardBasicModel("My Model")
	deck := genanki.StandardDeck("My Deck", "A deck for testing")

	// Add notes to the deck using method chaining
	deck.
		AddNote(
			genanki.NewNote(
				model.ID,
				[]string{"What is the capital of France?", "Paris"},[]string{"geography"},
			),
		).
		AddNote(
			genanki.NewNote(
				model.ID,
				[]string{"What is 2+2?", "4"}, []string{"math"},
			),
		)

	// Create and write package
	pkg := genanki.NewPackage([]*genanki.Deck{deck}).AddModel(model.Model)
	if err := pkg.WriteToFile("output.apkg"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully created Anki deck: output.apkg")
}
```

### Debug Logging

To enable debug logging, you can set the `DEBUG` environment variable to `true`:

```bash
DEBUG=true go run main.go
```

Or programmatically:

```go
pkg := genanki.NewPackage([]*genanki.Deck{deck}).AddModel(model.Model)
pkg.SetDebug(true)  // Enable debug logging
```

Debug logging will show detailed information about the database operations, including:
- Number of notes and cards
- Note field contents
- Number of models and decks
- Database verification steps

<details>
<summary><h2 style="display: inline-block">Advanced Usage</h2></summary>

### Creating Different Types of Models

```go
// Create a basic model
basicModel := genanki.StandardBasicModel("Basic Model")

// Create a cloze model
clozeModel := genanki.StandardClozeModel("Cloze Model")

// Create a deck
deck := genanki.StandardDeck("My Deck", "A deck for testing")
```

### Customizing Models

```go
// Create a model
model := genanki.StandardBasicModel("Custom Model")

// Customize the model (each method call separately)
model.Model.SetCSS(`
    .card { 
        font-family: Arial; 
        font-size: 20px;
        text-align: center;
        color: #333;
        background-color: #f5f5f5;
    }
    .question { 
        font-weight: bold; 
        color: navy;
    }
`)

model.Model.AddField(genanki.Field{
    Name: "Extra Info", 
    Font: "Arial",
    Size: 16,
})

basicModel.Model.SetCSS(`...`).AddField(genanki.Field{
    Name:  "Source",
    Ord:   2,
    Font:  "Arial",
    Size:  14,
    Color: "#666666",
})
```

### Adding Media Files

```go
// Create a package
pkg := genanki.NewPackage([]*genanki.Deck{deck})

// Add models to the package
pkg.AddModel(basicModel.Model)
pkg.AddModel(clozeModel.Model)

// Add media files
imageData, _ := os.ReadFile("image.jpg")
pkg.AddMedia("image.jpg", imageData)

audioData, _ := os.ReadFile("audio.mp3")
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
</details>

## Windows Support

For Windows platforms, you can use this branch
* dev/feature-enable-windows-compatibility

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
