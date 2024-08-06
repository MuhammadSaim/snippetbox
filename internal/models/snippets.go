package models

import (
	"database/sql"
	"time"
)

// Define a Snippet type to hold teh data for an individual snippet. Notice how
// the fields of the struct correspond to the fields in our MySQL snippets table
type Snippet struct {
	ID int
	Title string
	Content string
	ExpiredAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Define a  SnippetModel type which wraps a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the DB
func (m *SnippetModel) Insert(title string, conetnt string, expires int) (int, error)  {
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
func (m *SnippetModel) Get(id int) (Snippet, error)  {
	return Snippet{}, nil
}

// This will return teh 10 most recently created snippets
func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
