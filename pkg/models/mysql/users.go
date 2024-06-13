package mysql

import (
	"OpenAIDevTools/pkg/models"
	"database/sql"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	stmt := "INSERT INTO users (name, email, password) VALUES (?,?,?)"
	_, err := m.DB.Exec(stmt, name, email, password)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (bool, error) {
	stmt := "SELECT id FROM users WHERE email = ? AND password = ?"

	rows,err := m.DB.Query(stmt, email, password)
	if err != nil {
		return false, nil
	}
	defer rows.Close()
	return rows.Next(), nil

	// var id int

	// err := row.Scan(&id)
	// if err == sql.ErrNoRows {
	// 	return 0, nil
	// } else if err != nil {
	// 	return id, err
	// }

	// Any other error
	//return 0, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}

