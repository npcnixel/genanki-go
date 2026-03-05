package tests

import (
	"archive/zip"
	"database/sql"
	"io"
	"os"
	"path/filepath"
	"testing"

	genanki "github.com/npcnixel/genanki-go"

	_ "github.com/mattn/go-sqlite3"
)

func TestExportedNoteTagsUseAnkiFormat(t *testing.T) {
	model := genanki.NewBasicModel(1234567890, "Tag Format Model")
	deck := genanki.NewDeck(9876543210, "Tag Format Deck", "")
	note := genanki.NewNote(model.Model.ID, []string{"Front", "Back"}, []string{"one", "two words", "one"})
	deck.AddNote(note)

	pkg := genanki.NewPackage([]*genanki.Deck{deck}).AddModel(model.Model)

	tmpFile, err := os.CreateTemp("", "tags-format-*.apkg")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	if err := pkg.WriteToFile(tmpPath); err != nil {
		t.Fatalf("write package: %v", err)
	}

	noteTags := firstNoteTagsFromAPKG(t, tmpPath)
	if noteTags != " one two_words one " {
		t.Fatalf("expected Anki tag format, got %q", noteTags)
	}
}

func firstNoteTagsFromAPKG(t *testing.T, apkgPath string) string {
	t.Helper()

	collectionPath := extractCollectionDBFromAPKG(t, apkgPath)

	db, err := sql.Open("sqlite3", collectionPath)
	if err != nil {
		t.Fatalf("open extracted sqlite db: %v", err)
	}
	defer db.Close()

	var tags string
	if err := db.QueryRow("SELECT tags FROM notes LIMIT 1").Scan(&tags); err != nil {
		t.Fatalf("query note tags: %v", err)
	}

	return tags
}

func extractCollectionDBFromAPKG(t *testing.T, apkgPath string) string {
	t.Helper()

	archive, err := zip.OpenReader(apkgPath)
	if err != nil {
		t.Fatalf("open apkg: %v", err)
	}
	defer archive.Close()

	entry := findCollectionEntry(archive.File)
	if entry == nil {
		t.Fatalf("collection db not found in %s", apkgPath)
	}

	rc, err := entry.Open()
	if err != nil {
		t.Fatalf("open collection from zip: %v", err)
	}
	defer rc.Close()

	payload, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("read collection from zip: %v", err)
	}

	collectionPath := filepath.Join(t.TempDir(), filepath.Base(entry.Name))
	if err := os.WriteFile(collectionPath, payload, 0o600); err != nil {
		t.Fatalf("write extracted collection db: %v", err)
	}

	return collectionPath
}

func findCollectionEntry(files []*zip.File) *zip.File {
	for _, file := range files {
		if file.Name == "collection.anki21" || file.Name == "collection.anki2" {
			return file
		}
	}
	return nil
}
