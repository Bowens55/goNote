package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type NoteManager struct {
	Notes []Note
}

type Note struct {
	// title string dont think I want a title, I just want quick notes
	Body      string
	Directory string
	SavedAt   time.Time
}

func (nm *NoteManager) createNote(body string) {
	wd, err := os.Getwd()
	if err != nil {
		log.Println("Failed to get working directory.", err)
	}

	for _, note := range nm.Notes {
		// if body is the same as existing just return
		if note.Body == body {
			return
		}
	}
	note := Note{Body: body, Directory: wd, SavedAt: time.Now()}
	nm.Notes = append(nm.Notes, note)
}

func (nm *NoteManager) listNotes() {
	for _, n := range nm.Notes {
		fmt.Println(n.Body)
	}
}

func main() {
	nm := &NoteManager{}
	nm.createNote("test note.")
	nm.createNote("another")
	nm.createNote("another")
	fmt.Println(*nm)
	nm.listNotes()
}
