package tests

import (
	"testing"
	"time"

	genanki "github.com/npcnixel/genanki-go"
)

func TestNewNote(t *testing.T) {
	modelID := int64(1234567890)
	fields := []string{"Question", "Answer"}
	tags := []string{"test", "basic"}

	note := genanki.NewNote(modelID, fields, tags)

	if note.ModelID != modelID {
		t.Errorf("Expected note model ID %d, got %d", modelID, note.ModelID)
	}

	if len(note.Fields) != len(fields) {
		t.Errorf("Expected %d fields, got %d", len(fields), len(note.Fields))
	} else {
		for i, field := range fields {
			if note.Fields[i] != field {
				t.Errorf("Expected field %d to be '%s', got '%s'", i, field, note.Fields[i])
			}
		}
	}

	if len(note.Tags) != len(tags) {
		t.Errorf("Expected %d tags, got %d", len(tags), len(note.Tags))
	} else {
		for i, tag := range tags {
			if note.Tags[i] != tag {
				t.Errorf("Expected tag %d to be '%s', got '%s'", i, tag, note.Tags[i])
			}
		}
	}

	if note.ID == 0 {
		t.Error("Expected note ID to be non-zero")
	}

	// Check that the checksum is calculated correctly
	expectedChecksum := int64(0)
	if len(fields) > 0 {
		for _, c := range fields[0] {
			expectedChecksum = (expectedChecksum + int64(c)) % 0xffff
		}
	}

	if note.CheckSum != expectedChecksum {
		t.Errorf("Expected checksum %d, got %d", expectedChecksum, note.CheckSum)
	}

	// Check that the modified time is set
	zeroTime := time.Time{}
	if note.Modified == zeroTime {
		t.Error("Expected modified time to be set")
	}

	// Check that SortField is set to the first field
	if len(fields) > 0 && note.SortField != fields[0] {
		t.Errorf("Expected sort field to be '%s', got '%s'", fields[0], note.SortField)
	}
}
