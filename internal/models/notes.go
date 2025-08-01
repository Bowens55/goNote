package models

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"text/tabwriter"
	"time"
)

type NoteModel struct {
	DB *sql.DB
}

type Note struct {
	// title string dont think I want a title, I just want quick notes
	ID        int
	Body      string
	Directory string
	SavedAt   time.Time
}

// This will insert a new snippet into the database.
func (m *NoteModel) Insert(body, directory string) (int, error) {
	// Write the SQL statement we want to execute.
	stmt := `INSERT INTO notes (body, directory, saved_at) VALUES (?, ?, ?)`

	// Use the Exec() method on the embedded connection pool to execute the
	// statement. The first parameter is the SQL statement, followed by the
	// values for the placeholder parameters: title, content and expiry in
	// that order. This method returns a sql.Result type, which contains some
	// basic information about what happened when the statement was executed.
	result, err := m.DB.Exec(stmt, body, directory, time.Now())
	if err != nil {
		return 0, err
	}

	// Use the LastInsertId() method on the result to get the ID of our
	// newly inserted record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID returned has the type int64, so we convert it to an int type
	// before returning.
	return int(id), nil
}

// func (m *NoteModel) getID(n int) {
// 	pass
// }

func (m *NoteModel) List(n int, global bool) (Notes []*Note, err error) {
	// Build base query
	stmt := "SELECT * FROM notes"
	args := []any{} // Empty list that can hold values of any type (strings, numbers, booleans, etc.)

	if !global {
		stmt += " WHERE directory = ?" // SQL becomes: "SELECT * FROM notes WHERE directory = ?"
		dir, err := os.Getwd()         // Get current folder path (like "/home/user/projects")
		if err != nil {
			return nil, err
		}
		args = append(args, dir) // Add that path to our list - Now args = ["/home/user/projects"]
	}

	if n > 0 {
		stmt += " LIMIT ?"     // SQL becomes: "SELECT * FROM notes WHERE directory = ? LIMIT ?"
		args = append(args, n) // Add the number to our list - Now args = ["/home/user/projects", 5]
	}

	// The args... spreads out the list, so Query("SELECT * FROM notes WHERE directory = ? LIMIT ?", "/home/user/projects", 5)
	rows, err := m.DB.Query(stmt, args...)
	if err != nil {
		return nil, fmt.Errorf("unable to pull data from sql: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		note := &Note{} // Create NEW instance each iteration to avoid same pointer issue
		err = rows.Scan(&note.ID, &note.Body, &note.Directory, &note.SavedAt)
		if err != nil {
			fmt.Println("Unable to pull data from sql into note struct.", err)
			return nil, err
		}
		Notes = append(Notes, note)
	}
	return
}

func (m *NoteModel) Delete(id int) error {
	stmt := "DELETE FROM notes WHERE id = ?;"
	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
}

// UPDATE your_table_name
// SET body = 'New content', directory = '/new/path'
// WHERE id = 2;

func (m *NoteModel) Edit(id int) error {
	fmt.Print("New note content: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return fmt.Errorf("failed to read input")
	}
	userInput := scanner.Text()

	stmt := "UPDATE notes SET body = ?, saved_at = ? WHERE ID = ?"
	_, err := m.DB.Exec(stmt, userInput, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update note: %w", err)
	}

	return nil
}

// this might become formatNote depending on how we integrate fuzzy finder.
func DisplayNote(notes []*Note) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "ID\tNote\tDirectory\tDate")
	for _, note := range notes {
		localTime := note.SavedAt.In(time.Local)
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
			note.ID,
			note.Body,
			note.Directory,
			localTime.Format("Jan 2, 2006 3:04 PM MST"))
	}
	w.Flush()
}

// // This will return a specific snippet based on its id.
// func (m *SnippetModel) Get(id int) (Snippet, error) {
// 	// Write the SQL statement we want to execute. Again, I've split it over two
// 	// lines for readability.
// 	stmt := `SELECT id, title, content, created, expires FROM snippets
//     WHERE expires > UTC_TIMESTAMP() AND id = ?`

// 	// Use the QueryRow() method on the connection pool to execute our
// 	// SQL statement, passing in the untrusted id variable as the value for the
// 	// placeholder parameter. This returns a pointer to a sql.Row object which
// 	// holds the result from the database.
// 	row := m.DB.QueryRow(stmt, id)

// 	// Initialize a new zeroed Snippet struct.
// 	var s Snippet

// 	// Use row.Scan() to copy the values from each field in sql.Row to the
// 	// corresponding field in the Snippet struct. Notice that the arguments
// 	// to row.Scan are *pointers* to the place you want to copy the data into,
// 	// and the number of arguments must be exactly the same as the number of
// 	// columns returned by your statement.
// 	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
// 	if err != nil {
// 		// If the query returns no rows, then row.Scan() will return a
// 		// sql.ErrNoRows error. We use the errors.Is() function check for that
// 		// error specifically, and return our own ErrNoRecord error
// 		// instead (we'll create this in a moment).
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return Snippet{}, ErrNoRecord
// 		} else {
// 			return Snippet{}, err
// 		}
// 	}

// 	// If everything went OK, then return the filled Snippet struct.
// 	return s, nil
// }

// func (m *SnippetModel) Latest() ([]Snippet, error) {
// 	// Write the SQL statement we want to execute.
// 	stmt := `SELECT id, title, content, created, expires FROM snippets
//     WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

// 	// Use the Query() method on the connection pool to execute our
// 	// SQL statement. This returns a sql.Rows resultset containing the result of
// 	// our query.
// 	rows, err := m.DB.Query(stmt)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// We defer rows.Close() to ensure the sql.Rows resultset is
// 	// always properly closed before the Latest() method returns. This defer
// 	// statement should come *after* you check for an error from the Query()
// 	// method. Otherwise, if Query() returns an error, you'll get a panic
// 	// trying to close a nil resultset.
// 	defer rows.Close()

// 	// Initialize an empty slice to hold the Snippet structs.
// 	var snippets []Snippet

// 	// Use rows.Next to iterate through the rows in the resultset. This
// 	// prepares the first (and then each subsequent) row to be acted on by the
// 	// rows.Scan() method. If iteration over all the rows completes then the
// 	// resultset automatically closes itself and frees-up the underlying
// 	// database connection.
// 	for rows.Next() {
// 		// Create a new zeroed Snippet struct.
// 		var s Snippet
// 		// Use rows.Scan() to copy the values from each field in the row to the
// 		// new Snippet object that we created. Again, the arguments to row.Scan()
// 		// must be pointers to the place you want to copy the data into, and the
// 		// number of arguments must be exactly the same as the number of
// 		// columns returned by your statement.
// 		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
// 		if err != nil {
// 			return nil, err
// 		}
// 		// Append it to the slice of snippets.
// 		snippets = append(snippets, s)
// 	}

// 	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
// 	// error that was encountered during the iteration. It's important to
// 	// call this - don't assume that a successful iteration was completed
// 	// over the whole resultset.
// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	// If everything went OK then return the Snippets slice.
// 	return snippets, nil
// }
