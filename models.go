package genanki

import (
	"crypto/rand"
	"encoding/binary"
	"math"
	"time"
)

type Model struct {
	ID        int64
	Name      string
	Fields    []Field
	Templates []Template
	CSS       string
}

type Field struct {
	Name   string
	Ord    int
	Sticky bool
	RTF    bool
	Font   string
	Size   int
	Color  string
	Align  string
}

type Template struct {
	Name  string
	Ord   int
	Qfmt  string
	Afmt  string
	Bqfmt string
	Bafmt string
}

type Note struct {
	ID        int64
	ModelID   int64
	Fields    []string
	Tags      []string
	Modified  time.Time
	SortField string
	CheckSum  int64
}

type Deck struct {
	ID       int64
	Name     string
	Desc     string
	Notes    []*Note
	Media    map[string][]byte
	Created  time.Time
	Modified time.Time
}

func GenerateIntID() int64 {
	var b [8]byte
	rand.Read(b[:])
	return int64(binary.LittleEndian.Uint64(b[:])) % math.MaxInt64
}

func NewModel(id int64, name string) *Model {
	return &Model{
		ID:        id,
		Name:      name,
		Fields:    make([]Field, 0),
		Templates: make([]Template, 0),
		CSS:       defaultCSS,
	}
}

func NewNote(modelID int64, fields []string, tags []string) *Note {
	now := time.Now()

	csum := int64(0)
	if len(fields) > 0 {
		for _, c := range fields[0] {
			csum = (csum + int64(c)) % 0xffff
		}
	}

	return &Note{
		ID:        GenerateIntID(),
		ModelID:   modelID,
		Fields:    fields,
		Tags:      tags,
		Modified:  now,
		SortField: fields[0],
		CheckSum:  csum,
	}
}

func NewDeck(id int64, name string, desc string) *Deck {
	// Auto-generate ID if not provided (i.e., if id is 0)
	if id == 0 {
		// Use a standard Anki deck ID for better compatibility
		id = 1347639657110
	}

	now := time.Now()
	return &Deck{
		ID:       id,
		Name:     name,
		Desc:     desc,
		Notes:    make([]*Note, 0),
		Media:    make(map[string][]byte),
		Created:  now,
		Modified: now,
	}
}

func (d *Deck) AddNote(note *Note) *Deck {
	d.Notes = append(d.Notes, note)
	d.Modified = time.Now()
	return d
}

func (d *Deck) AddMedia(filename string, data []byte) *Deck {
	d.Media[filename] = data
	d.Modified = time.Now()
	return d
}

const defaultCSS = `
.card {
    font-family: arial;
    font-size: 20px;
    text-align: center;
    color: black;
    background-color: white;
    line-height: 1.2;
}
`

type BasicModel struct {
	*Model
}

type ClozeModel struct {
	*Model
}

func NewBasicModel(id int64, name string) *BasicModel {
	// Auto-generate ID if not provided (i.e., if id is 0)
	if id == 0 {
		// Use Anki's standard Basic model ID for better compatibility
		id = 1607392319
	}

	model := NewModel(id, name)

	model.Fields = []Field{
		{Name: "Front", Ord: 0, Font: "Arial", Size: 20, Color: "#000000", Align: "left"},
		{Name: "Back", Ord: 1, Font: "Arial", Size: 20, Color: "#000000", Align: "left"},
	}

	model.Templates = []Template{
		{
			Name: "Card 1",
			Ord:  0,
			Qfmt: "{{Front}}",
			Afmt: "{{FrontSide}}\n\n<hr id=answer>\n\n{{Back}}",
		},
	}

	return &BasicModel{Model: model}
}

func NewClozeModel(id int64, name string) *ClozeModel {
	// Auto-generate ID if not provided (i.e., if id is 0)
	if id == 0 {
		// Use Anki's standard Cloze model ID for better compatibility
		id = 1122334455
	}

	model := NewModel(id, name)

	model.Fields = []Field{
		{Name: "Text", Ord: 0, Font: "Arial", Size: 20, Color: "#000000", Align: "left"},
		{Name: "Extra", Ord: 1, Font: "Arial", Size: 20, Color: "#000000", Align: "left"},
	}

	model.Templates = []Template{
		{
			Name: "Cloze",
			Ord:  0,
			Qfmt: "{{cloze:Text}}\n\n{{Extra}}",
			Afmt: "{{cloze:Text}}\n\n{{Extra}}",
		},
	}

	return &ClozeModel{Model: model}
}

func (m *Model) AddField(field Field) *Model {
	field.Ord = len(m.Fields)
	m.Fields = append(m.Fields, field)
	return m
}

func (m *Model) AddTemplate(template Template) *Model {
	template.Ord = len(m.Templates)
	m.Templates = append(m.Templates, template)
	return m
}

func (m *Model) SetCSS(css string) *Model {
	m.CSS = css
	return m
}

// Convenience functions that don't require passing IDs

// StandardBasicModel creates a new basic model with Anki's standard model ID
func StandardBasicModel(name string) *BasicModel {
	return NewBasicModel(0, name)
}

// StandardClozeModel creates a new cloze model with Anki's standard model ID
func StandardClozeModel(name string) *ClozeModel {
	return NewClozeModel(0, name)
}

// StandardDeck creates a new deck with Anki's standard deck ID
func StandardDeck(name string, desc string) *Deck {
	return NewDeck(0, name, desc)
}
