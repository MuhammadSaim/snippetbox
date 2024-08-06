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
func (m *SnippetModel) Insert(title string, conetnt string, expiredAt int, createdAt int, updatedAt int) (int, error)  {
	return 0, nil
}

// This will return a specific snipped based on its id
func (m *SnippetModel) Get(id int) (Snippet, error)  {
	return Snippet{}, nil
}

// This will return teh 10 most recently created snippets
func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
