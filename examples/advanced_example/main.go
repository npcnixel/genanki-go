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
	"strings"
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
	// Create models - using auto-generated IDs
	basicModel := genanki.NewBasicModel(
		0, // Auto-generate ID
		"Basic",
	)

	clozeModel := genanki.NewClozeModel(
		0, // Auto-generate ID
		"Cloze",
	)

	// Create a deck - using auto-generated ID
	deck := genanki.NewDeck(
		0, // Auto-generate ID
		"Advanced Example Deck",
		"An advanced example deck with multiple note types and media",
	)

	// Print the generated IDs for reference
	fmt.Printf("Generated Basic Model ID: %d\n", basicModel.ID)
	fmt.Printf("Generated Cloze Model ID: %d\n", clozeModel.ID)
	fmt.Printf("Generated Deck ID: %d\n", deck.ID)

	// Create basic notes
	basicNote1 := genanki.NewNote(
		basicModel.ID,
		[]string{
			"What is the capital of France?",
			"Paris",
		},
		[]string{"geography", "europe"},
	)

	basicNote2 := genanki.NewNote(
		basicModel.ID,
		[]string{
			"What is the largest planet in our solar system?",
			"Jupiter",
		},
		[]string{"astronomy", "planets"},
	)

	basicNote3 := genanki.NewNote(
		basicModel.ID,
		[]string{
			"What does this image show?<br><img src='istockphoto-1263636227-612x612.jpg'>",
			"A stock photo",
		},
		[]string{"images", "examples"},
	)

	// Add a note with a generated image
	basicNote4 := genanki.NewNote(
		basicModel.ID,
		[]string{
			"What does this represent?<br><img src='sample_image.png'>",
			"A generated image",
		},
		[]string{"images", "examples", "generated"},
	)

	// Add cloze-like notes
	clozeNote1 := genanki.NewNote(
		clozeModel.ID,
		[]string{
			"The capital of France is {{c1::Paris}}.",
			"",
		},
		[]string{"geography", "europe"},
	)

	clozeNote2 := genanki.NewNote(
		clozeModel.ID,
		[]string{
			"{{c1::Jupiter}} is the largest planet in our solar system.",
			"",
		},
		[]string{"astronomy", "planets"},
	)

	// Add notes to deck
	notes := []*genanki.Note{basicNote1, basicNote2, basicNote3, basicNote4, clozeNote1, clozeNote2}
	for _, note := range notes {
		deck.AddNote(note)
		log.Printf("Added note with fields: %q", strings.Join(note.Fields, "\u001f"))
	}

	// Create package
	pkg := genanki.NewPackage([]*genanki.Deck{deck})

	// Add models to package
	pkg.AddModel(basicModel.Model)
	pkg.AddModel(clozeModel.Model)

	// Add the user's image file
	userImagePath := "istockphoto-1263636227-612x612.jpg"
	userImageData, err := os.ReadFile(userImagePath)
	if err != nil {
		log.Printf("Warning: Failed to read user image file: %v", err)
		log.Println("Continuing without user image...")
	} else {
		pkg.AddMedia(userImagePath, userImageData)
		log.Println("Added user image to package")
	}

	// Generate and add sample image
	imageData, err := generateSampleImage()
	if err != nil {
		log.Printf("Warning: Failed to generate sample image: %v", err)
		log.Println("Continuing without sample image...")
	} else {
		pkg.AddMedia("sample_image.png", imageData)
		log.Println("Added generated image to package")
	}

	// Create output directory
	outputDir := filepath.Join("..", "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Write package to file
	outputPath := filepath.Join(outputDir, "advanced_deck.apkg")
	if err := pkg.WriteToFile(outputPath); err != nil {
		log.Fatalf("Failed to write package: %v", err)
	}

	// Summary
	fmt.Printf("Successfully created Anki deck: %s\n", outputPath)
	fmt.Printf("Number of notes: %d\n", len(deck.Notes))
	fmt.Printf("Created: %s\n", time.Now().Format(time.RFC1123))
}
