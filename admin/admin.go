package admin

import (
	"database/sql"
)

type Admin struct {
	Username string
}

type Manager struct {
	DB *sql.DB
}

func (m *Manager) Login(username, password string) (*Admin, error) {
	row := m.DB.QueryRow("SELECT username FROM admins WHERE username = $1 AND password = $2", username, password)
	var a Admin
	err := row.Scan(&a.Username)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
