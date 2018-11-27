package pinactivity

import (
	"database/sql"
	"errors"
	"log"
)

type Pinactivity struct {
	ActivitiesID int
	EmployeeCode string
	Name         string
	Phone        string
}

type Manager struct {
	DB *sql.DB
}

func (m *Manager) Insert(p *Pinactivity) error {
	stmt := "INSERT INTO pinactivities(activities_id, employee_code, name, phone) VALUES ($1,$2,$3,$4)"
	_, err := m.DB.Exec(stmt, p.ActivitiesID, p.EmployeeCode, p.Name, p.Phone)
	return err
}

func (m *Manager) Delete(employeecode, id string) error {
	stmt := "DELETE FROM pinactivities WHERE employee_code = $1 AND activities_id=$2"
	r, err := m.DB.Exec(stmt, employeecode, id)
	effect, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if effect != 1 {
		return errors.New("pinactivity: delete not have effected row")
	}
	return nil
}

func (m *Manager) All(id int) ([]Pinactivity, error) {
	users := []Pinactivity{}
	rows, err := m.DB.Query("SELECT activities_id, employee_code, name, phone FROM pinactivities WHERE activities_id = $1 ", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var p Pinactivity
		err := rows.Scan(&p.ActivitiesID, &p.EmployeeCode, &p.Name, &p.Phone)
		if err != nil {
			return nil, err
		}
		users = append(users, p)
	}
	return users, nil
}

func (m *Manager) ResetStorage() {
	_, err := m.DB.Exec("TRUNCATE TABLE activities RESTART IDENTITY;")
	if err != nil {
		log.Fatal(err)
	}
}
