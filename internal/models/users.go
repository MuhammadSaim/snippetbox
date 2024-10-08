package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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
	// Create a bcrypt hash of the plain-text password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	// SQL statement
	stmt := `INSERT INTO users (name, email, password, created_at, updated_at)
    VALUES (?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	// Use the Exec method to insert the user details
	_, err = m.DB.Exec(stmt, name, email, hashedPassword)
	if err != nil {
		// If this returns an error we use the errors.As func
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
	}
	return nil
}

// We'll use the Authenticate method to verify whether a user exists
// with the provided details
func (m *UserModel) Authenticate(email string, password string) (int, error) {
	// Retrive the id hashed password associated with the given email
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, password FROM users WHERE email = ?"

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Check wether the hashed password and plain password provided match.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// if everything goes well we have our user and its id
	return id, nil
}

// We'll use the Exists method to check if a user exists or not
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
