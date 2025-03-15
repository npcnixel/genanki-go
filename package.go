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
	db    *Database
	media map[string][]byte
}

func NewPackage(db *Database) *Package {
	return &Package{
		db:    db,
		media: make(map[string][]byte),
	}
}

func (p *Package) AddMedia(filename string, data []byte) {
	p.media[filename] = data
}

func (p *Package) WriteToFile(path string) error {
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

	dbFile, err := p.db.GetFilePath()
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
