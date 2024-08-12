package models

import (
	"database/sql"
	"time"
)

// Define a new User struct
type User struct {
	ID        int64
	Name      string
	Email     string
	password  []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Define a new UserModel struct which wraps DB conn
type UserModel struct {
	DB *sql.DB
}

// We'll use the Insert method to add a new record
func (m *UserModel) Insert(name string, email string, password string) error {
	return nil
}

// We'll use the Authenticate method to verify whether a user exists
// with the provided details
func (m *UserModel) Authenticate(email string, password string) (int, error) {
	return 0, nil
}

// We'll use the Exists method to check if a user exists or not
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
