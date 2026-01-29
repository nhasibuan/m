package crud

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Category struct {
	ID          int    `json:"id"`
	Nama        string `json:"nama"`
	Description string `json:"description"`
}

var Category1 = []Category{
	{ID: 1, Nama: " ", Description: " "},
}

func Connect(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(5)
	return db, nil
}

//
