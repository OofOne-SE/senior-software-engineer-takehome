package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"senior-software-engineer-takehome/models"

	"github.com/go-openapi/strfmt"
	_ "github.com/lib/pq"
)

type Postgres struct {
	log func(string, ...any)
	db  *sql.DB
}

func New(l func(string, ...any), dbName, user, pwd, addr string) Postgres {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, pwd, addr, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := db.Query(`CREATE TABLE IF NOT EXISTS weatherunits (date date, humidity float8, temperature float8)`); err != nil {
		log.Fatal(err)
	}
	return Postgres{
		log: l,
		db:  db,
	}
}

func (p Postgres) AddUnit(u *models.WeatherUnit) error {
	rows, err := p.db.Query(`INSERT INTO weatherunits(date, humidity, temperature) VALUES($1, $2, $3)`, u.Date.String(), u.Humidity, u.Temperature)
	defer rows.Close()
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	} else {
		return nil
	}
}

func (p Postgres) GetUnits(from, to *strfmt.Date) ([]*models.WeatherUnit, error) {
	rows, err := p.db.Query("SELECT * FROM weatherunits WHERE date BETWEEN to_date($1,'YYYY-MM-DD') AND to_date($2,'YYYY-MM-DD')", from.String(), to.String())
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	} else {
		var res []*models.WeatherUnit

		for rows.Next() {
			wu := &models.WeatherUnit{}
			if err := rows.Scan(&wu.Date, &wu.Humidity, &wu.Temperature); err != nil {
				return res, err
			}
			res = append(res, wu)
		}
		if err = rows.Err(); err != nil {
			return res, err
		}
		return res, nil
	}
}
