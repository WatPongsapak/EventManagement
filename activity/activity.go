package activity

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

type Activity struct {
	ID          int
	Name        string
	Location    string
	Speaker     string
	Description string
	Maxjoin     int
	StartDate   time.Time
	EndDate     time.Time
	StartTime   time.Time
	EndTime     time.Time
	Round       int
	Amount      int
	DateRange   string
	TimeRange   string
	MaxRound    int
}

type Manager struct {
	DB *sql.DB
}

func (m *Manager) Insert(a *Activity) error {
	row := m.DB.QueryRow("SELECT round FROM activities WHERE name = $1 ORDER BY round DESC", a.Name)
	var i int
	err := row.Scan(&i)
	if err == nil {
		a.Round = i + 1
	} else {
		a.Round = 1
	}
	stmt := "INSERT INTO activities(name, location, speaker, description, max_joinable, start_date, end_date, start_time, end_time, round) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id"
	r := m.DB.QueryRow(stmt, a.Name, a.Location, a.Speaker, a.Description, a.Maxjoin, a.StartDate, a.EndDate, a.StartTime, a.EndTime, a.Round)
	err = r.Scan(&a.ID)
	return err
}

func (m *Manager) Update(a *Activity) error {
	row := m.DB.QueryRow("SELECT name FROM activities WHERE id = $1", a.ID)
	var oldname string
	row.Scan(&oldname)
	stmt := "UPDATE activities SET name = $2 WHERE name = $1"
	_, err := m.DB.Exec(stmt, oldname, a.Name)
	stmt = "UPDATE activities SET name = $2, location = $3, speaker = $4, description = $5, max_joinable = $6, start_date = $7, end_date = $8, start_time = $9, end_time = $10, round= $11 WHERE id = $1"
	_, err = m.DB.Exec(stmt, a.ID, a.Name, a.Location, a.Speaker, a.Description, a.Maxjoin, a.StartDate, a.EndDate, a.StartTime, a.EndTime, a.Round)
	return err
}

func (m *Manager) FindByID(id int) (*Activity, error) {
	row := m.DB.QueryRow("SELECT id, name, location, speaker, description, max_joinable, start_date, end_date, start_time, end_time, round, (SELECT COUNT(*) FROM pinactivities WHERE pinactivities.activities_id = activities.id) FROM activities WHERE id = $1", id)
	var a Activity
	err := row.Scan(&a.ID, &a.Name, &a.Location, &a.Speaker, &a.Description, &a.Maxjoin, &a.StartDate, &a.EndDate, &a.StartTime, &a.EndTime, &a.Round, &a.Amount)
	if err != nil {
		return nil, err
	}
	a.DateRange = a.StartDate.Format("01/02/2006") + " - " + a.EndDate.Format("01/02/2006")
	a.TimeRange = a.StartTime.Format("03:04 pm") + " - " + a.EndTime.Format("03:04 pm")
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
	countname := map[string]int{}
	rows, _ := m.DB.Query("SELECT name, count(name) FROM activities GROUP BY name")
	for rows.Next() {
		var name string
		var i int
		rows.Scan(&name, &i)
		countname[name] = i
	}

	rows, err := m.DB.Query("SELECT id, name, location, speaker, description, max_joinable, start_date, end_date, start_time, end_time, round, (SELECT COUNT(*) FROM pinactivities WHERE pinactivities.activities_id = activities.id) AS amount FROM activities")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var a Activity
		err := rows.Scan(&a.ID, &a.Name, &a.Location, &a.Speaker, &a.Description, &a.Maxjoin, &a.StartDate, &a.EndDate, &a.StartTime, &a.EndTime, &a.Round, &a.Amount)
		if err != nil {
			return nil, err
		}
		a.DateRange = a.StartDate.Format("01/02/2006") + " - " + a.EndDate.Format("01/02/2006")
		a.TimeRange = a.StartTime.Format("03:04 pm") + " - " + a.EndTime.Format("03:04 pm")
		a.MaxRound = countname[a.Name]
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
