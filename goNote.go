package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
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
	if len(nm.Notes) < 1 {
		slog.Info("No notes currently exist in this directory.")
		return
	}

	fmt.Println("id", "description") // TODO: do this better... lol
	for i, note := range nm.Notes {
		fmt.Println(i+1, note.Body)
	}
}

func main() {
	nm := &NoteManager{}

	var createFlag string
	var listFlag bool

	flag.StringVar(&createFlag, "create", "", "If defined, we will create a note and add it to our list.")
	flag.StringVar(&createFlag, "c", "", "If defined, we will create a note and add it to our list.")

	flag.BoolVar(&listFlag, "list", true, "Disable listing out notes by passing false to this flag.")
	flag.BoolVar(&listFlag, "l", true, "Disable listing out notes by passing false to this flag.")

	flag.Parse()
	if createFlag != "" {
		nm.createNote(createFlag)
	}

	// nm.createNote(*createFlag)

	// fmt.Println(*nm)
	if listFlag {
		nm.listNotes()
	}
}
