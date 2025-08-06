package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Opts struct {
	List   int  `short:"l" long:"list" description:"List out n number of notes. If n is not passed list all" optional:"true" optional-value:"-1"`
	Delete int  `short:"d" long:"delete" description:"Delete specific note based on ID field"`
	Global bool `short:"g" long:"global" description:"Decide what notes to load, global or current directory notes." optional:"true" optional-value:"false"`
	Edit   int  `short:"e" long:"edit" description:"Edit a specific note based on the ID passed." optional:"true" optional-value:"0"`
}

// func getArgs() (args []string, opts *Opts, err error) {
// 	opts = &Opts{}                // Create pointer to new EMPTY struct. Otherwise, the var declared in the sig is just nil (pointer to struct == nil)
// 	args, err = flags.Parse(opts) // Pass the pointer directly, cannot pass nil here so we need to create the struct above.
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	return
// }

// insert a note into the sqlite db, return the last id inserted into db
func insertNote(db *sql.DB, body, dir string) (int64, error) {
	insertSQL := `INSERT INTO notes (body, dir) VALUES (?, ?)`
	result, err := db.Exec(insertSQL, body, dir)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func main() {

	// args, opts, err := getArgs()
	// if err != nil {
	// 	log.Println("Unable to parse args.", err)
	// }

	db, err := sql.Open("sqlite3", "./note.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
		CREATE TABLE IF NOT EXISTS notes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        body TEXT,
        dir TEXT,
        time DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Print("Unable to get cwd.")
	}

	_, err = insertNote(db, "test", cwd)
	if err != nil {
		log.Printf("Unable to insert note into db. Raw err: %s", err)
	}

	rows, err := db.Query("select * from notes")
	if err != nil {
		log.Printf("Unable to query rows from notes", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var body string
		var dir string
		var time string

		err = rows.Scan(&id, &body, &dir, &time)
		if err != nil {
			fmt.Println("unable to scan rows")
		}

		fmt.Println(id, body, dir, time)
	}
}
