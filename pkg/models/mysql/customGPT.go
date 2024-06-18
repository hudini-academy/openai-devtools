package mysql

import (
	"OpenAIDevTools/pkg/models"
	"database/sql"
	"fmt"
)

type CustomGPTModel struct {
	DB *sql.DB
}

// InsertFunction inserts a new record into the CustomGPT table in the database.
// It takes two parameters: buttonName and systemMessage, which represent the SystemName and SystemPrompt columns respectively.
// The function returns an error if any occurs during the insertion process.
func (m *CustomGPTModel) InsertFunction(buttonName, systemMessage string) error {
	stmt := "INSERT INTO CustomGPT (SystemName,SystemPrompt) VALUES (?,?)"
	_, err := m.DB.Exec(stmt, buttonName, systemMessage)
	if err != nil {
		return err
	}
	return nil
}

// GetFunction retrieves all CustomGPT table records, returns []*models.CustomGPT and error,
// querying ID, SystemName, and SystemPrompt columns.
func (m *CustomGPTModel) GetFunction() ([]*models.CustomGPT, error) {
	stmt := "SELECT ID,SystemName,SystemPrompt FROM CustomGPT"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		fmt.Println(err, "error fetching data")
		return nil, err
	}
	defer rows.Close()
	Buttons := []*models.CustomGPT{}
	for rows.Next() {
		s := &models.CustomGPT{}
		err := rows.Scan(&s.ID, &s.SystemName, &s.SystemPrompt)
		if err != nil {
			return nil, err
		}
		Buttons = append(Buttons, s)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err, "error running though rows data")
		return nil, err
	}

	return Buttons, nil
}

// GetIndividualFunction retrieves single CustomGPT table record by ID,
// returns *models.CustomGPT and error; populates SystemName and SystemPrompt,
// returns nil and sql.ErrNoRows if not found, or nil and error on other issues.
func (m *CustomGPTModel) GetIndividualFunction(id int) (*models.CustomGPT, error) {
	stmt := "SELECT SystemName,SystemPrompt FROM CustomGPT WHERE ID = ?"
	row, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	s := &models.CustomGPT{}

	// Scan the values from the row into variables
	if row.Next() {
		err = row.Scan(&s.SystemName, &s.SystemPrompt)
		if err != nil {
			return nil, err
		}

		// Return the retrieved values and nil error
		return s, nil
	}

	// If no rows were returned, return a nil pointer and a specific error
	return nil, sql.ErrNoRows
}

// DeleteFunction deletes CustomGPT table record by ID, returns error; returns nil if record doesn't exist.
func (m *CustomGPTModel) DeleteFunction(id int) error {
	stmt := `DELETE FROM CustomGPT WHERE ID = ?`
	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
}
