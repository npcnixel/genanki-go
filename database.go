package genanki

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	pragmas := []string{
		"PRAGMA foreign_keys = ON",
		"PRAGMA page_size = 4096",
		"PRAGMA encoding = 'UTF-8'",
		"PRAGMA legacy_file_format = OFF",
		"PRAGMA journal_mode = DELETE",
		"PRAGMA synchronous = OFF",
		"PRAGMA temp_store = MEMORY",
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to set pragma: %v", err)
		}
	}

	d := &Database{db: db}
	if err := d.initialize(); err != nil {
		db.Close()
		return nil, err
	}

	return d, nil
}

func (d *Database) initialize() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY,
			guid TEXT NOT NULL,
			mid INTEGER NOT NULL,
			mod INTEGER NOT NULL,
			usn INTEGER NOT NULL,
			tags TEXT NOT NULL,
			flds TEXT NOT NULL,
			sfld TEXT NOT NULL,
			csum INTEGER NOT NULL,
			flags INTEGER NOT NULL,
			data TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS cards (
			id INTEGER PRIMARY KEY,
			nid INTEGER NOT NULL,
			did INTEGER NOT NULL,
			ord INTEGER NOT NULL,
			mod INTEGER NOT NULL,
			usn INTEGER NOT NULL,
			type INTEGER NOT NULL,
			queue INTEGER NOT NULL,
			due INTEGER NOT NULL,
			ivl INTEGER NOT NULL,
			factor INTEGER NOT NULL,
			reps INTEGER NOT NULL,
			lapses INTEGER NOT NULL,
			left INTEGER NOT NULL,
			odue INTEGER NOT NULL,
			odid INTEGER NOT NULL,
			flags INTEGER NOT NULL,
			data TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS col (
			id INTEGER PRIMARY KEY,
			crt INTEGER NOT NULL,
			mod INTEGER NOT NULL,
			scm INTEGER NOT NULL,
			ver INTEGER NOT NULL,
			dty INTEGER NOT NULL,
			usn INTEGER NOT NULL,
			ls INTEGER NOT NULL,
			conf TEXT NOT NULL,
			models TEXT NOT NULL,
			decks TEXT NOT NULL,
			dconf TEXT NOT NULL,
			tags TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS graves (
			usn INTEGER NOT NULL,
			oid INTEGER NOT NULL,
			type INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS revlog (
			id INTEGER PRIMARY KEY,
			cid INTEGER NOT NULL,
			usn INTEGER NOT NULL,
			ease INTEGER NOT NULL,
			ivl INTEGER NOT NULL,
			lastIvl INTEGER NOT NULL,
			factor INTEGER NOT NULL,
			time INTEGER NOT NULL,
			type INTEGER NOT NULL
		)`,
	}

	for _, query := range queries {
		if _, err := d.db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %v", err)
		}
	}

	return nil
}

func (d *Database) AddModel(model *Model) error {
	conf := map[string]interface{}{
		"nextPos":              1,
		"estTimes":             true,
		"activeDecks":          []int64{1},
		"sortType":             "noteFld",
		"timeLim":              0,
		"sortBackwards":        false,
		"addToCur":             true,
		"curDeck":              1,
		"newSpread":            0,
		"dueCounts":            true,
		"curModel":             model.ID,
		"collapseTime":         1200,
		"schedVer":             1,
		"newBury":              true,
		"_nt_":                 0,
		"_deck_1_lastNotetype": model.ID,
	}

	confJSON, err := json.Marshal(conf)
	if err != nil {
		return fmt.Errorf("failed to marshal conf: %v", err)
	}

	modelConfig := map[string]interface{}{
		"id":        model.ID,
		"name":      model.Name,
		"type":      0, // Basic type
		"mod":       0,
		"usn":       0,
		"sortf":     0, // Sort by first field
		"did":       1, // Default deck ID
		"vers":      []interface{}{},
		"tags":      []interface{}{},
		"css":       ".card {\n font-family: arial;\n font-size: 20px;\n text-align: center;\n color: black;\n background-color: white;\n}\n",
		"latexPre":  "\\documentclass[12pt]{article}\n\\special{papersize=3in,5in}\n\\usepackage[utf8]{inputenc}\n\\usepackage{amssymb,amsmath}\n\\pagestyle{empty}\n\\setlength{\\parindent}{0in}\n\\begin{document}\n",
		"latexPost": "\\end{document}",
		"latexsvg":  false,
		"req":       []interface{}{[]interface{}{0, "all", []interface{}{0}}},
		"flds": []map[string]interface{}{
			{
				"name":              "Front",
				"ord":               0,
				"sticky":            false,
				"rtl":               false,
				"font":              "Arial",
				"size":              20,
				"media":             []interface{}{},
				"description":       "",
				"plainText":         false,
				"collapsed":         false,
				"excludeFromSearch": false,
				"preventDeletion":   false,
			},
			{
				"name":              "Back",
				"ord":               1,
				"sticky":            false,
				"rtl":               false,
				"font":              "Arial",
				"size":              20,
				"media":             []interface{}{},
				"description":       "",
				"plainText":         false,
				"collapsed":         false,
				"excludeFromSearch": false,
				"preventDeletion":   false,
			},
		},
		"tmpls": []map[string]interface{}{
			{
				"name":  "Card 1",
				"ord":   0,
				"qfmt":  "{{Front}}",
				"afmt":  "{{FrontSide}}\n\n<hr id=answer>\n\n{{Back}}",
				"bqfmt": "",
				"bafmt": "",
				"did":   nil,
				"bfont": "",
				"bsize": 0,
			},
		},
		"originalStockKind": 1,
	}

	modelJSON, err := json.Marshal(modelConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal model config: %v", err)
	}

	deckConf := map[string]interface{}{
		"1": map[string]interface{}{
			"id":       1,
			"mod":      0,
			"name":     "Default",
			"usn":      0,
			"maxTaken": 60,
			"autoplay": true,
			"timer":    0,
			"replayq":  true,
			"new": map[string]interface{}{
				"delays":        []float64{1.0, 10.0},
				"ints":          []int{1, 4, 0},
				"initialFactor": 2500,
				"separate":      true,
				"order":         1,
				"perDay":        20,
				"bury":          false,
			},
			"rev": map[string]interface{}{
				"perDay":     200,
				"ease4":      1.3,
				"fuzz":       0.05,
				"ivlFct":     1.0,
				"maxIvl":     36500,
				"bury":       false,
				"minSpace":   1,
				"hardFactor": 1.2,
			},
			"lapse": map[string]interface{}{
				"delays":      []float64{10.0},
				"mult":        0.0,
				"minInt":      1,
				"leechFails":  8,
				"leechAction": 1,
			},
			"dyn":                     false,
			"newMix":                  0,
			"newPerDayMinimum":        0,
			"interdayLearningMix":     0,
			"reviewOrder":             0,
			"newSortOrder":            0,
			"newGatherPriority":       0,
			"buryInterdayLearning":    false,
			"fsrsWeights":             []interface{}{},
			"fsrsParams5":             []interface{}{},
			"desiredRetention":        0.9,
			"ignoreRevlogsBeforeDate": "",
			"easyDaysPercentages":     []float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0},
			"stopTimerOnAnswer":       false,
			"secondsToShowQuestion":   0.0,
			"secondsToShowAnswer":     0.0,
			"questionAction":          0,
			"answerAction":            0,
			"waitForAudio":            true,
			"sm2Retention":            0.9,
			"weightSearch":            "",
		},
	}

	deckConfJSON, err := json.Marshal(deckConf)
	if err != nil {
		return fmt.Errorf("failed to marshal deck config: %v", err)
	}

	_, err = d.db.Exec(`
		INSERT OR REPLACE INTO col (id, crt, mod, scm, ver, dty, usn, ls, conf, models, decks, dconf, tags)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		1,                 // id                 
		time.Now().Unix(), // crt 
		time.Now().Unix(), // mod 
		time.Now().Unix(), // scm 
		11,                // ver (schema version)                
		0,                 // dty                 
		0,                 // usn                 
		0,                 // ls                 
		string(confJSON),  // conf  
		fmt.Sprintf(`{"%d": %s}`, model.ID, modelJSON), // models 
		"{}",                 // decks (will be updated by AddDeck)                 
		string(deckConfJSON), // dconf 
		"{}",                 // tags                 
	)
	if err != nil {
		return fmt.Errorf("failed to insert model: %v", err)
	}

	return nil
}

func (d *Database) AddDeck(deck *Deck) error {
	deckConfig := map[string]interface{}{
		"id":               deck.ID,
		"mod":              time.Now().Unix(),
		"name":             deck.Name,
		"usn":              -1,
		"lrnToday":         []int{0, 0},
		"revToday":         []int{0, 0},
		"newToday":         []int{0, 0},
		"timeToday":        []int{0, 0},
		"collapsed":        false,
		"browserCollapsed": false,
		"desc":             deck.Desc,
		"dyn":              0,
		"conf":             1,
		"extendNew":        10,
		"extendRev":        50,
	}

	var decksJSON string
	err := d.db.QueryRow("SELECT decks FROM col WHERE id = 1").Scan(&decksJSON)
	if err != nil {
		return fmt.Errorf("failed to read decks: %v", err)
	}

	var decks map[string]interface{}
	if err := json.Unmarshal([]byte(decksJSON), &decks); err != nil {
		decks = make(map[string]interface{})
	}

	decks[fmt.Sprintf("%d", deck.ID)] = deckConfig

	newDecksJSON, err := json.Marshal(decks)
	if err != nil {
		return fmt.Errorf("failed to marshal decks: %v", err)
	}

	_, err = d.db.Exec("UPDATE col SET decks = ? WHERE id = 1", string(newDecksJSON))
	return err
}

func (d *Database) AddNote(note *Note) error {
	tagsJSON, err := json.Marshal(note.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %v", err)
	}

	fields := make([]byte, 0)
	for i, field := range note.Fields {
		if i > 0 {
			fields = append(fields, 0x1f) // Unit separator
		}
		fields = append(fields, []byte(field)...)
	}
	fieldsStr := string(fields)

	csum := int64(0)
	if len(note.Fields) > 0 {
		for _, c := range note.Fields[0] {
			csum = (csum + int64(c)) % 0xffff
		}
	}

	noteData := map[string]interface{}{
		"tags": note.Tags,
	}
	noteDataJSON, err := json.Marshal(noteData)
	if err != nil {
		return fmt.Errorf("failed to marshal note data: %v", err)
	}

	log.Printf("Note fields string: %q", fieldsStr)

	_, err = d.db.Exec(`
		INSERT INTO notes (id, guid, mid, mod, usn, tags, flds, sfld, csum, flags, data)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		note.ID,
		fmt.Sprintf("%x", note.ID),
		note.ModelID,
		note.Modified.Unix(),
		-1,
		string(tagsJSON),
		fieldsStr,
		note.Fields[0],
		csum,
		0,
		string(noteDataJSON),
	)

	if err != nil {
		return fmt.Errorf("failed to insert note: %v", err)
	}

	var verifyFields string
	err = d.db.QueryRow("SELECT flds FROM notes WHERE id = ?", note.ID).Scan(&verifyFields)
	if err != nil {
		return fmt.Errorf("failed to verify note: %v", err)
	}

	log.Printf("Verified note fields: %q", verifyFields)
	return nil
}

func (d *Database) AddCard(noteID, deckID int64, templateOrd int) error {
	now := time.Now().Unix()
	_, err := d.db.Exec(`
		INSERT INTO cards (id, nid, did, ord, mod, usn, type, queue, due, ivl, factor, reps, lapses, left, odue, odid, flags, data)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		GenerateIntID(),
		noteID,
		deckID,
		templateOrd,
		now,
		-1,
		0,    // new card
		0,    // new queue
		0,    // due today
		0,    // initial interval
		2500, // initial factor
		0,    // no reps yet
		0,    // no lapses yet
		0,    // no cards left
		0,    // no original due date
		0,    // no original deck
		0,    // no flags
		"{}",
	)

	return err
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) VerifyContent() error {
	var noteCount int
	err := d.db.QueryRow("SELECT COUNT(*) FROM notes").Scan(&noteCount)
	if err != nil {
		return fmt.Errorf("failed to count notes: %v", err)
	}
	log.Printf("Number of notes in database: %d", noteCount)

	rows, err := d.db.Query("SELECT id, flds FROM notes")
	if err != nil {
		return fmt.Errorf("failed to query notes: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var fields string
		if err := rows.Scan(&id, &fields); err != nil {
			return fmt.Errorf("failed to scan note: %v", err)
		}
		log.Printf("Note %d fields: %q", id, fields)
	}

	var cardCount int
	err = d.db.QueryRow("SELECT COUNT(*) FROM cards").Scan(&cardCount)
	if err != nil {
		return fmt.Errorf("failed to count cards: %v", err)
	}
	log.Printf("Number of cards in database: %d", cardCount)

	var modelsJSON, decksJSON string
	err = d.db.QueryRow("SELECT models, decks FROM col WHERE id = 1").Scan(&modelsJSON, &decksJSON)
	if err != nil {
		return fmt.Errorf("failed to read collection: %v", err)
	}

	var models, decks map[int64]interface{}
	if err := json.Unmarshal([]byte(modelsJSON), &models); err != nil {
		return fmt.Errorf("failed to unmarshal models: %v", err)
	}
	if err := json.Unmarshal([]byte(decksJSON), &decks); err != nil {
		return fmt.Errorf("failed to unmarshal decks: %v", err)
	}

	log.Printf("Number of models in collection: %d", len(models))
	log.Printf("Number of decks in collection: %d", len(decks))

	return nil
}

func (d *Database) GetFilePath() (string, error) {
	tmpFile, err := os.CreateTemp("", "anki-*.db")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}

	if err := d.backupToFile(tmpFile.Name()); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to backup database: %v", err)
	}

	return tmpFile.Name(), nil
}

func (d *Database) backupToFile(path string) error {
	destDB, err := sql.Open("sqlite3", path)
	if err != nil {
		return fmt.Errorf("failed to open destination database: %v", err)
	}
	defer destDB.Close()

	pragmas := []string{
		"PRAGMA foreign_keys = ON",
		"PRAGMA page_size = 4096",
		"PRAGMA encoding = 'UTF-8'",
		"PRAGMA legacy_file_format = OFF",
		"PRAGMA journal_mode = DELETE",
		"PRAGMA synchronous = OFF",
	}

	for _, pragma := range pragmas {
		if _, err := destDB.Exec(pragma); err != nil {
			return fmt.Errorf("failed to set pragma: %v", err)
		}
	}

	if err := (&Database{db: destDB}).initialize(); err != nil {
		return fmt.Errorf("failed to initialize destination database: %v", err)
	}

	tx, err := destDB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	tables := []string{"notes", "cards", "col", "graves", "revlog"}
	for _, tableName := range tables {
		data, err := d.db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
		if err != nil {
			return fmt.Errorf("failed to get data from %s: %v", tableName, err)
		}
		defer data.Close()

		cols, err := data.Columns()
		if err != nil {
			return fmt.Errorf("failed to get columns for %s: %v", tableName, err)
		}

		placeholders := make([]string, len(cols))
		for i := range placeholders {
			placeholders[i] = "?"
		}
		insertSQL := fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES (%s)",
			tableName,
			strings.Join(cols, ", "),
			strings.Join(placeholders, ", "),
		)

		stmt, err := tx.Prepare(insertSQL)
		if err != nil {
			return fmt.Errorf("failed to prepare statement: %v", err)
		}
		defer stmt.Close()

		for data.Next() {
			values := make([]interface{}, len(cols))
			valuePtrs := make([]interface{}, len(cols))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			if err := data.Scan(valuePtrs...); err != nil {
				return fmt.Errorf("failed to scan row: %v", err)
			}

			if _, err := stmt.Exec(values...); err != nil {
				return fmt.Errorf("failed to insert row: %v", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	if _, err := destDB.Exec("PRAGMA wal_checkpoint(TRUNCATE)"); err != nil {
		return fmt.Errorf("failed to checkpoint WAL: %v", err)
	}

	if _, err := destDB.Exec("VACUUM"); err != nil {
		return fmt.Errorf("failed to vacuum database: %v", err)
	}

	return nil
}
