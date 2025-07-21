package main

import (
	"flag"
	"fmt"
	"goNote/internal/models"
	"log"
	"log/slog"
	"os"
	"strings"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type application struct {
	logger    *slog.Logger
	NoteModel *models.NoteModel
}

func main() {
	var listFlag int

	flag.IntVar(&listFlag, "list", 0, "List out n number of notes, if no number is passed, list all.")
	flag.IntVar(&listFlag, "l", 0, "List out n number of notes, if 0 is passed, list all.")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Printf("Unable to load env file.")
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	dsn := os.Getenv("dsn")

	db, err := openDB(dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := application{
		logger:    logger,
		NoteModel: &models.NoteModel{DB: db},
	}

	// Get note from command-line arguments
	args := flag.Args()
	if len(args) > 0 {
		note := strings.Join(args, " ")
		// Use note variable as needed, e.g.:
		// app.NoteModel.Insert(note, "default", time.Now())

		wd, err := os.Getwd()
		if err != nil {
			logger.Error("Unable to get working dir.", "error", err)
			os.Exit(1)
		}

		if note != "" {
			app.NoteModel.Insert(note, wd)
		}
	}
	// TODO: make this better dont care rn.
	// need to use something else since default for list flag is only
	// used whenever the flag isn't passed at all.
	// cleanup the if else...
	if listFlag == 0 && flag.NFlag() == 1 {
		// --list or -l was passed without a value
		notes, err := app.NoteModel.List(0)
		if err != nil {
			fmt.Println(err)
		}
		for _, v := range notes {
			fmt.Println(v)
		}

	} else if listFlag >= 1 {
		notes, err := app.NoteModel.List(listFlag)
		if err != nil {
			fmt.Println(err)
		}
		for _, v := range notes {
			fmt.Println(v)
		}
	}

}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
