package crud

import (
	"database/sql"
	"errors"
)

type repository struct {
	db1 *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db1: db}
	// return &repository{}
}

func (repo *repository) select1() ([]dual, error) {
	query := "select 1 as id"
	var d dual
	err := repo.db1.QueryRow(query).Scan(&d.ID)
	if err == sql.ErrNoRows {
		return nil, errors.New(err.Error())
	}
	if err != nil {
		return nil, err
	}
	return []dual{d}, nil
}

func (repo *repository) selectAll() ([]Category, error) {
	query := "SELECT id, nama, description FROM \"Category\""
	rows, err := repo.db1.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]Category, 0)
	for rows.Next() {
		var c Category
		err := rows.Scan(&c.ID, &c.Nama, &c.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (repo *repository) insert(category *Category) error {
	query := "INSERT INTO \"Category\" (nama, description) VALUES ($1, $2) RETURNING id"
	err := repo.db1.QueryRow(query, category.Nama, category.Description).Scan(&category.ID)
	return err
}

func (repo *repository) selectByID(id int) (*Category, error) {
	query := "SELECT id, nama, description FROM \"Category\" WHERE id = $1"

	var c Category
	err := repo.db1.QueryRow(query, id).Scan(&c.ID, &c.Nama, &c.Description)
	if err == sql.ErrNoRows {
		return nil, errors.New(err.Error())
	}
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (repo *repository) update(category *Category) error {
	query := "UPDATE \"Category\" SET nama = $1, description = $2 WHERE id = $3"
	result, err := repo.db1.Exec(query, category.Nama, category.Description, category.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("update error")
	}

	return nil
}
func (repo *repository) delete(id int) error {
	query := "DELETE FROM \"Category\" WHERE id = $1"
	result, err := repo.db1.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("delete error")
	}

	return err
}