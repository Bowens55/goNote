package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// func listNote(){}

const file string = "./.notesdb/db"

const create string = `
	CREATE TABLE IF NOT EXISTS notes (
	id INTEGER NOT NULL PRIMARY KEY,
	time DATETIME NOT NULL,
	note TEXT,
	directory TEXT
	);`

type Notes struct {
	mu sync.Mutex
	db *sql.DB
}

// create the table if not exists
// return the db, well technically the notes struct that contains db and sync mutex
func NewNotes() (*Notes, error) {
	os.MkdirAll("./.notesdb", 0755)
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(create); err != nil {
		db.Close()
		return nil, err
	}
	return &Notes{
		db: db,
	}, nil
}

func createNote(note string, db *sql.DB) error {
	currentDir, err := os.Getwd()
	stmt := `INSERT INTO notes(note, directory, time) VALUES(?, ?, datetime('now'))`
	_, err = db.Exec(stmt, note, currentDir)
	if err != nil {
		return fmt.Errorf("failed to insert note: %v", err)
	}
	return nil
}

// func listNote() {

// }

// delete note by ID
func deleteNote(id int, db *sql.DB) (int64, error) {
	sql := `DELETE FROM notes WHERE id = ?`
	result, err := db.Exec(sql, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func main() {

	createPtr := flag.String("create", "", "The note to create.")
	flag.Parse()

	notes, err := NewNotes()
	if err != nil {
		fmt.Printf("Failed to initialize notes DB: %v\n", err)
		return
	}
	defer notes.db.Close()

	if *createPtr != "" {
		if err := createNote(*createPtr, notes.db); err != nil {
			fmt.Printf("Failed to create note: %v\n", err)
			return
		}
	}

	rows, err := notes.db.Query("SELECT * FROM notes")
	if err != nil {
		fmt.Printf("Query failed: %v\n", err)
		return
	}
	defer rows.Close()

	// _, err = deleteNote(7, notes.db)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	for rows.Next() {
		var id int
		var ts string
		var note string
		var dir string

		if err := rows.Scan(&id, &ts, &note, &dir); err != nil {
			fmt.Printf("Scan failed: %v\n", err)
			continue
		}

		fmt.Printf("ID=%d Time=%s Note=%q Dir=%s\n", id, ts, note, dir)
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("Rows error: %v\n", err)
	}
}
