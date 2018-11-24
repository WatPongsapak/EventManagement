package activity

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

type Activity struct {
	ID            int
	Name          string
	Location      string
	Speaker       string
	Description   string
	Maxjoin       int
	StartDatetime time.Time
	EndDatetime   time.Time
	Round         int
}

type Manager struct {
	DB *sql.DB
}

func (m *Manager) Insert(a *Activity) error {
	row := m.DB.QueryRow("SELECT round FROM activities WHERE name = $1", a.Name)
	var i int
	err := row.Scan(&i)
	if err == nil {
		a.Round = i + 1
	} else {
		a.Round = 1
	}
	stmt := "INSERT INTO activities(name, location, speaker, description, max_joinable, start_datetime, end_datetime, round) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id"
	r := m.DB.QueryRow(stmt, a.Name, a.Location, a.Speaker, a.Description, a.Maxjoin, a.StartDatetime, a.EndDatetime, a.Round)
	err = r.Scan(&a.ID)
	return err
}

func (m *Manager) Update(a *Activity) error {
	stmt := "UPDATE activities SET name = $2, location = $3, speaker = $4, description = $5, max_joinable = $6, start_datetime = $7, end_datetime = $8, round= $9 WHERE id = $1"
	_, err := m.DB.Exec(stmt, a.ID, a.Name, a.Location, a.Speaker, a.Description, a.Maxjoin, a.StartDatetime, a.EndDatetime, a.Round)
	return err
}

func (m *Manager) FindByID(id int) (*Activity, error) {
	row := m.DB.QueryRow("SELECT id, name, location, speaker, description, max_joinable, start_datetime, end_datetime, round FROM activities WHERE id = $1", id)
	var a Activity
	err := row.Scan(&a.ID, &a.Name, &a.Location, &a.Speaker, &a.Description, &a.Maxjoin, &a.StartDatetime, &a.EndDatetime, &a.Round)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (m *Manager) Delete(a *Activity) error {
	stmt := "DELETE FROM activities WHERE id = $1"
	r, err := m.DB.Exec(stmt, a.ID)
	effect, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if effect != 1 {
		return errors.New("user: delete not have effected row")
	}
	return nil
}

func (m *Manager) All() ([]Activity, error) {
	users := []Activity{}
	rows, err := m.DB.Query("SELECT id, name, location, speaker, description, max_joinable, start_datetime, end_datetime, round FROM activities")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var a Activity
		err := rows.Scan(&a.ID, &a.Name, &a.Location, &a.Speaker, &a.Description, &a.Maxjoin, &a.StartDatetime, &a.EndDatetime, &a.Round)
		if err != nil {
			return nil, err
		}
		users = append(users, a)
	}
	return users, nil
}

func (m *Manager) ResetStorage() {
	_, err := m.DB.Exec("TRUNCATE TABLE activities RESTART IDENTITY;")
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Manager) Last() Activity {
	var a Activity
	stmt := "SELECT id, name, location, speaker, description, max_joinable, start_datetime, end_datetime, round FROM activities ORDER BY id DESC LIMIT 1"
	row := m.DB.QueryRow(stmt)
	err := row.Scan(&a.ID, &a.Name, &a.Location, &a.Speaker, &a.Description, &a.Maxjoin)
	if err != nil {
		log.Fatal(err)
	}
	return a
}