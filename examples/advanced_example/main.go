package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"time"

	genanki "github.com/npcnixel/genanki-go"
)

// generateSampleImage creates a simple PNG image with text
func generateSampleImage() ([]byte, error) {
	// Create a 200x100 image
	img := image.NewRGBA(image.Rect(0, 0, 200, 100))

	// Fill the image with a light blue background
	for y := 0; y < 100; y++ {
		for x := 0; x < 200; x++ {
			img.Set(x, y, color.RGBA{135, 206, 250, 255}) // Light sky blue
		}
	}

	// Draw a simple border
	for x := 0; x < 200; x++ {
		img.Set(x, 0, color.RGBA{0, 0, 0, 255})
		img.Set(x, 99, color.RGBA{0, 0, 0, 255})
	}
	for y := 0; y < 100; y++ {
		img.Set(0, y, color.RGBA{0, 0, 0, 255})
		img.Set(199, y, color.RGBA{0, 0, 0, 255})
	}

	// Encode the image to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func main() {
	db, err := genanki.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	// Create a basic model using the built-in function
	// Using Anki's standard Basic model ID
	basicModel := genanki.NewBasicModel(
		1, // Anki's standard Basic model ID
		"Basic",
	)

	// Customize the basic model
	basicModel.Model.Fields = []genanki.Field{
		{Name: "Front", Ord: 0, Font: "Arial", Size: 20},
		{Name: "Back", Ord: 1, Font: "Arial", Size: 20},
	}

	basicModel.Model.Templates = []genanki.Template{
		{
			Name: "Card 1",
			Ord:  0,
			Qfmt: "{{Front}}",
			Afmt: "{{FrontSide}}<hr id=answer>{{Back}}",
		},
	}

	basicModel.Model.CSS = `.card {
		font-family: Arial, sans-serif;
		font-size: 20px;
		text-align: center;
		color: #333;
		background-color: #f8f8f8;
		padding: 20px;
	}
	.card img {
		max-width: 90%;
		max-height: 400px;
	}
	hr#answer {
		border: 1px solid #ccc;
		margin: 20px 0;
	}
	small {
		font-size: 14px;
		color: #666;
	}
	.cloze {
		font-weight: bold;
		color: blue;
	}`

	// Add model to database
	if err := db.AddModel(basicModel.Model); err != nil {
		log.Fatalf("Failed to add basic model: %v", err)
	}
	log.Println("Added basic model to database")

	// Create a deck
	deck := genanki.NewDeck(
		1, // Using a simple deck ID
		"Advanced Example Deck",
		"A deck demonstrating advanced features of genanki-go",
	)

	// Add deck to database
	if err := db.AddDeck(deck); err != nil {
		log.Fatalf("Failed to add deck: %v", err)
	}
	log.Println("Added deck to database")

	// Create notes with unique IDs
	basicNote1 := genanki.NewNote(
		basicModel.Model.ID,
		[]string{
			"What is the capital of France?",
			"Paris",
		},
		[]string{"geography", "europe", "capitals"},
	)

	basicNote2 := genanki.NewNote(
		basicModel.Model.ID,
		[]string{
			"What is the largest planet in our solar system?",
			"Jupiter",
		},
		[]string{"astronomy", "planets", "science"},
	)

	// Add a note with the user's image
	basicNote3 := genanki.NewNote(
		basicModel.Model.ID,
		[]string{
			"What does this image show?<br><img src='istockphoto-1263636227-612x612.jpg'>",
			"A stock photo",
		},
		[]string{"images", "examples"},
	)

	// Add a note with a generated image
	basicNote4 := genanki.NewNote(
		basicModel.Model.ID,
		[]string{
			"What does this represent?<br><img src='sample_image.png'>",
			"A generated image",
		},
		[]string{"images", "examples", "generated"},
	)

	// Add cloze-like notes (using basic model)
	clozeNote1 := genanki.NewNote(
		basicModel.Model.ID,
		[]string{
			"What is the capital of France?",
			"The capital of France is Paris.",
		},
		[]string{"geography", "europe"},
	)

	clozeNote2 := genanki.NewNote(
		basicModel.Model.ID,
		[]string{
			"What is the largest planet in our solar system?",
			"Jupiter is the largest planet in our solar system.",
		},
		[]string{"astronomy", "planets"},
	)

	// Add notes to database and create cards
	notes := []*genanki.Note{basicNote1, basicNote2, basicNote3, basicNote4, clozeNote1, clozeNote2}
	for i, note := range notes {
		if err := db.AddNote(note); err != nil {
			log.Fatalf("Failed to add note %d: %v", i+1, err)
		}
		log.Printf("Added note %d to database", i+1)

		// Create a card for each note
		if err := db.AddCard(note.ID, deck.ID, 0); err != nil {
			log.Fatalf("Failed to add card for note %d: %v", i+1, err)
		}
		log.Printf("Added card for note %d", i+1)
	}

	// Create package
	pkg := genanki.NewPackage(db)

	// Add the user's image file
	userImagePath := "istockphoto-1263636227-612x612.jpg"
	userImageData, err := os.ReadFile(userImagePath)
	if err != nil {
		log.Printf("Warning: Failed to read user image file: %v", err)
		log.Println("Continuing without user image...")
	} else {
		pkg.AddMedia("istockphoto-1263636227-612x612.jpg", userImageData)
		log.Println("Added user image to package")
	}

	// Generate and add sample image
	imageData, err := generateSampleImage()
	if err != nil {
		log.Fatalf("Failed to generate sample image: %v", err)
	}
	pkg.AddMedia("sample_image.png", imageData)
	log.Println("Added generated image to package")

	// Ensure output directory exists at same level as example directories
	outputDir := filepath.Join("..", "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Write package to file in the output directory
	outputPath := filepath.Join(outputDir, "advanced_deck.apkg")
	if err := pkg.WriteToFile(outputPath); err != nil {
		log.Fatalf("Failed to write package: %v", err)
	}

	fmt.Printf("Successfully created Anki deck: %s\n", outputPath)
	fmt.Printf("Number of notes: %d\n", len(notes))
	fmt.Printf("Created at: %s\n", time.Now().Format(time.RFC1123))
}
