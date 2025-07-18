package main

import (
	"flag"
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
	var listFlag bool

	flag.BoolVar(&listFlag, "list", true, "Disable listing out notes by passing false to this flag.")
	flag.BoolVar(&listFlag, "l", true, "Disable listing out notes by passing false to this flag.")
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
	if len(os.Args) > 1 {
		note := strings.Join(os.Args[1:], " ")
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
