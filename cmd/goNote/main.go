package main

import (
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
}

func getArgs() (args []string, opts *Opts, err error) {
	opts = &Opts{}                // Create pointer to new EMPTY struct. Otherwise, the var declared in the sig is just nil (pointer to struct == nil)
	args, err = flags.Parse(opts) // Pass the pointer directly, cannot pass nil here so we need to create the struct above.
	if err != nil {
		return nil, nil, err
	}
	return
}

func main() {

	args, opts, err := getArgs()
	if err != nil {
		log.Println("Unable to parse args.", err)
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Println("Unable to get cwd.", err)
	}

	fullPath := filepath.Join(cwd, ".env")

	_, err = os.Stat(fullPath)
	if !os.IsNotExist(err) {
		err = godotenv.Load()
		if err != nil {
			log.Printf("Unable to load env file.")
		}
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

	if opts.Delete != 0 {
		app.NoteModel.Delete(opts.Delete)
	}

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
