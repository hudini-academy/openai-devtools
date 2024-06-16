package mysql

import (
	"database/sql"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"

	"OpenAIDevTools/pkg/models"
)

// UserModel represents the model for interacting with the users table in the database.
type UserModel struct {
	DB *sql.DB
}

// Insert inserts a new user into the users table.
func (m *UserModel) Insert(name, email string, password []byte) error {
	stmt := "INSERT INTO users (name, email, password) VALUES (?, ?, ?)"
	_, err := m.DB.Exec(stmt, name, email, password)
	if err != nil {
		return err
	}
	return nil
}

// Authenticate verifies the user credentials (email and password).
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	row := m.DB.QueryRow("SELECT id, password FROM users WHERE email = ?", email)
	err := row.Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return 0, models.ErrInvalidCredentials
	}

	return id, nil
}

// GetUsernameByID retrieves the username with a given user ID.
func (m *UserModel) GetUsernameByID(id int) (string, error) {
	var username string
	err := m.DB.QueryRow("SELECT name FROM users WHERE id = ?", id).Scan(&username)
	if err == sql.ErrNoRows {
		return "", errors.New("user not found")
	} else if err != nil {
		return "", err
	}
	return username, nil
}

// UserExist checks if a user with the given email already exists in the database.
func (m *UserModel) UserExist(email string) bool {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE email = ?)"
	err := m.DB.QueryRow(query, email).Scan(&exists)
	if err != nil {
		log.Println("Error checking email existence:", err)
		return false
	}
	return exists
}
