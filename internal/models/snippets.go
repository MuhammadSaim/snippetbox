package models

import (
	"database/sql"
	"errors"
	"time"
)

// Define a Snippet type to hold teh data for an individual snippet. Notice how
// the fields of the struct correspond to the fields in our MySQL snippets table
type Snippet struct {
	ID        int
	Title     string
	Content   string
	ExpiredAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Define a  SnippetModel type which wraps a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the DB
func (m *SnippetModel) Insert(title string, conetnt string, expires int) (int, error) {
	// Write the SQL statement we want to execute. I have split it over Three lines
	// for readability that's why I used backquotes
	stmt := `INSERT INTO snippets (title, content, expired_at, created_at, updated_at)
			VALUES
			(?, ?, DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY), UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	// Use the Exec() method on the embeded connection pool to execute the
	// statement . The first parameter is the SQL statement, followed by the
	// values for the placeholder parameters
	result, err := m.DB.Exec(stmt, title, conetnt, expires)
	if err != nil {
		return 0, err
	}

	// Use the LastInsertId method on the result to get teh ID
	// of our newly inserted record
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID returned has the type int64, so we convert it to an int type
	// before returning
	return int(id), nil
}

// This will return a specific snipped based on its id
func (m *SnippetModel) Get(id int) (Snippet, error) {
	// SQL statement to get the snippet
	stmt := `SELECT * FROM snippets WHERE expired_at > UTC_TIMESTAMP() AND id = ?`

	// Use the QueryRow func on the connection pool to execute our
	// SQL statement, passing in the untrusted id as the value for the placeholder
	// and return the sql.Row
	row := m.DB.QueryRow(stmt, id)

	// initialize the Snippet struct
	var snippet Snippet

	// Use row.Scan to copy the values from each field to crosponding structs
	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.ExpiredAt, &snippet.CreatedAt, &snippet.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	// If everything is okay then return the snippet
	return snippet, nil
}

// This will return teh 10 most recently created snippets
func (m *SnippetModel) Latest() ([]Snippet, error) {
	// Write a SQL statement we want to execute
	stmt := `SELECT * FROM snippets
			WHERE expired_at > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	// Use the Query method on the connection pool to execute error
	// SQL statement This returns a sql.Rows resultset containing the result of
	// our query
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// We defer rows.Close() to ensure that resultset is properly closed
	defer rows.Close()

	// Initialize an empty slice to hold the Snippet struct
	var snippets []Snippet

	// loop through the resultset and append it to the snippets 1 by 1
	for rows.Next() {
		// create a pointer to new snippet
		var snippet Snippet
		err = rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.ExpiredAt, &snippet.CreatedAt, &snippet.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Append it to the slice
		snippets = append(snippets, snippet)
	}

	// After the loops end so we have to retrive any errors
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
