package mysql

import (
	"OpenAIDevTools/pkg/models"
	"database/sql"
)

type CustomGPTModel struct {
	DB *sql.DB
}

func (m *CustomGPTModel) InsertFunction(buttonName, systemMessage string) error {
	stmt := "INSERT INTO customgpt (title,message) VALUES (?,?)"
	_, err := m.DB.Exec(stmt, buttonName, systemMessage)
	if err != nil {
		return err
	}
	return nil
}
func (m *CustomGPTModel) GetFunction() ([]*models.CustomGPT, error) {
	stmt := "SELECT id,title,message FROM customgpt"
	rows, err := m.DB.Query(stmt)
	if err != nil {
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
		return nil, err
	}

	return Buttons, nil
}
func (m *CustomGPTModel) GetIndividualFunction(id int) (*models.CustomGPT, error) {
	stmt := "SELECT title, message FROM customgpt WHERE id = ?"
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

func (m *CustomGPTModel) DeleteFunction(id int) error {
	stmt:= `DELETE FROM customgpt WHERE id = ?`
	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
}