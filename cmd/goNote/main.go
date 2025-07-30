package main

import (
	"errors"
	"goNote/internal/models"
	"log"
	"log/slog"
	"os"
	"path/filepath"
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
	List   int  `short:"l" long:"list" description:"List out n number of notes. If n is not passed list all" optional:"true" optional-value:"-1"`
	Delete int  `short:"d" long:"delete" description:"Delete specific note based on ID field"`
	Global bool `short:"g" long:"global" description:"Decide what notes to load, global or current directory notes." optional:"true" optional-value:"false"`
	Edit   int  `short:"e" long:"edit" description:"Edit a specific note based on the ID passed." optional:"true" optional-value:"0"`
}

func getArgs() (args []string, opts *Opts, err error) {
	opts = &Opts{}                // Create pointer to new EMPTY struct. Otherwise, the var declared in the sig is just nil (pointer to struct == nil)
	args, err = flags.Parse(opts) // Pass the pointer directly, cannot pass nil here so we need to create the struct above.
	if err != nil {
		return nil, nil, err
	}
	return
}

func EnvFileExists(dir string) (bool, error) {
	fullPath := filepath.Join(dir, ".env")
	_, err := os.Stat(fullPath)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func main() {

	args, opts, err := getArgs()
	if err != nil {
		log.Println("Unable to parse args.", err)
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		slog.Error("Unable to get cwd.", "error", err)
	}

	exists, err := EnvFileExists(cwd)
	if err != nil {
		slog.Error("Unable to check current directory for a .env file.")
	}

	if exists {
		err = godotenv.Load()
		if err != nil {
			log.Printf("Unable to load env file.")
		}
	}

	dsn := os.Getenv("dsn")
	if dsn == "" {
		slog.Error("DSN value cannot be empty, set env var. Or create a .env file in the current dir.", "dsn", dsn)
		os.Exit(1)
	}

	db, err := openDB(dsn)
	if err != nil {
		slog.Error("Unable to open connection to DB.", "error", err, "dsn", dsn)
		os.Exit(1)
	}
	defer db.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := application{
		logger:    logger,
		NoteModel: &models.NoteModel{DB: db},
	}

	// Get note from command-line arguments
	myNote := strings.TrimSpace(strings.Join(args, " "))
	if myNote != "" {
		wd, err := os.Getwd()
		if err != nil {
			app.logger.Error("Unable to get working dir.", "error", err)
			os.Exit(1)
		}
		app.NoteModel.Insert(myNote, wd)
	}

	if opts.Delete != 0 {
		app.NoteModel.Delete(opts.Delete)
	}

	if opts.Edit != 0 {
		err = app.NoteModel.Edit(opts.Edit)
		if err != nil {
			app.logger.Error("Unable to edit note", "error", err, "id", opts.Edit)
		}
	}

	// list at the very end, any new commands need to be before this point!
	var notes []*models.Note
	// list notes if -l or if we just run goNote
	if opts.List != 0 || myNote == "" {
		notes, _ = app.NoteModel.List(opts.List, opts.Global)
	}

	if notes != nil {
		models.DisplayNote(notes)
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
