package date

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

type Date struct {
	ID            int
	ActivitiesID  int
	StartDateTime time.Time
	EndDateTime   time.Time
}

type Manager struct {
	DB *sql.DB
}

func (m *Manager) Insert(d *Date) error {
	stmt := "INSERT INTO date(activities_id, start_datetime, end_datetime) VALUES ($1,$2,$3)"
	_, err := m.DB.Exec(stmt, d.ActivitiesID, d.StartDateTime, d.EndDateTime)
	return err
}

func (m *Manager) Update(d *Date) error {
	stmt := "UPDATE date SET start_datetime = $2, end_datetime = $3 WHERE id = $1"
	_, err := m.DB.Exec(stmt, d.ID, d.StartDateTime, d.EndDateTime)
	return err
}

func (m *Manager) FindByID(id int) (*Date, error) {
	row := m.DB.QueryRow("SELECT id, activities_id, start_datetime, end_datetime FROM date WHERE id = $1", id)
	var d Date
	err := row.Scan(&d.ID, &d.ActivitiesID, &d.StartDateTime, &d.EndDateTime)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (m *Manager) Delete(d *Date) error {
	stmt := "DELETE FROM date WHERE id = $1"
	r, err := m.DB.Exec(stmt, d.ID)
	effect, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if effect != 1 {
		return errors.New("date: delete not have effected row")
	}
	return nil
}

func (m *Manager) All() ([]Date, error) {
	users := []Date{}
	rows, err := m.DB.Query("SELECT id, activities_id, start_datetime, end_datetime FROM date")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var d Date
		err := rows.Scan(&d.ID, &d.ActivitiesID, &d.StartDateTime, &d.EndDateTime)
		if err != nil {
			return nil, err
		}
		users = append(users, d)
	}
	return users, nil
}

func (m *Manager) ResetStorage() {
	_, err := m.DB.Exec("TRUNCATE TABLE activities RESTART IDENTITY;")
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Manager) Last() Date {
	var d Date
	stmt := "SELECT id, name, location, speaker, description, Max_joinable FROM activities ORDER BY id DESC LIMIT 1"
	row := m.DB.QueryRow(stmt)
	err := row.Scan(&d.ID, &d.ActivitiesID, &d.StartDateTime, &d.EndDateTime)
	if err != nil {
		log.Fatal(err)
	}
	return d
}
