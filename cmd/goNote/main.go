package main

import (
	"fmt"
	"goNote/internal/models"
	"log"
	"log/slog"
	"os"
	"strings"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
)

type application struct {
	logger    *slog.Logger
	NoteModel *models.NoteModel
}

type Opts struct {
	List   int `short:"l" long:"list" description:"List out n number of notes. If n is not passed list all" optional:"true" optional-value:"-1"`
	Delete int `short:"d" long:"delete" description:"Delete specific note based on ID field"`
}

func main() {
	var opts Opts

	args, err := flags.Parse(&opts)
	if err != nil {
		fmt.Println("unable to parse flags.")
		os.Exit(1)
	}

	err = godotenv.Load()
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
	myNote := strings.TrimSpace(strings.Join(args, " "))
	if myNote != "" {
		wd, err := os.Getwd()
		if err != nil {
			logger.Error("Unable to get working dir.", "error", err)
			os.Exit(1)
		}
		app.NoteModel.Insert(myNote, wd)
	}

	var notes []models.Note
	// list notes if -l or if we just run goNote
	if opts.List != 0 || myNote == "" {
		notes, _ = app.NoteModel.List(opts.List)
	}

	if notes != nil {
		models.DisplayNote(notes)
	}

	if opts.Delete != 0 {
		app.NoteModel.Delete(opts.Delete)
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
