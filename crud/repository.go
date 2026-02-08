package crud

import (
	"database/sql"
	"errors"
	"time"
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
	//
	rows, err := repo.db1.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	//
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

func (repo *repository) selectAllProducts(name string) ([]Product, error) {
	query := "SELECT id, name, price, stock FROM \"Product\""
	var args []interface{}
	if name != "" {
		query += " WHERE name ILIKE $1"
		args = append(args, "%"+name+"%")
	}
	//
	rows, err := repo.db1.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	//
	products := make([]Product, 0)
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (repo *repository) insert(category *Category) error {
	query := "INSERT INTO \"Category\" (nama, description) VALUES ($1, $2) RETURNING id"
	err := repo.db1.QueryRow(query, category.Nama, category.Description).Scan(&category.ID)
	return err
}

func (repo *repository) selectByID(id int) (*Category, error) {
	var c Category
	//
	query := "SELECT id, nama, description FROM \"Category\" WHERE id = $1"
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
		return errors.New("Update error")
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
		return errors.New("Delete error")
	}

	return err
}

func (repo *repository) insertShoppingCart(co *CheckoutItem1) error {
	query := "INSERT INTO \"CheckoutItem\" (product_id, quantity) VALUES ($1, $2) RETURNING product_id"
	err := repo.db1.QueryRow(query, co.ProductID, co.Quantity).Scan(&co.ProductID)
	return err
}

func (repo *repository) createTransaction(items []CheckoutItem1) (*Transaction, error) {
	var res *Transaction
	tx, err := repo.db1.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	details := make([]TransactionDetail, 0)
	totalAmount := 0
	//
	for _, item := range items {
		var (
			Name1            string
			ID, Price, Stock int
		)
		query := "SELECT id, name, price, stock FROM \"Product\" WHERE id = $1"
		err := tx.QueryRow(query, item.ProductID).Scan(&ID, &Name1, &Price, &Stock)
		if err == sql.ErrNoRows {
			return nil, errors.New(err.Error())
		}
		if err != nil {
			return nil, err
		}
		//
		subtotal := item.Quantity * Price
		totalAmount += subtotal
		_, err = tx.Exec("UPDATE \"Product\" SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}
		// Isi struktur data
		details = append(details, TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: Name1,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}
	for i := range details {
		details[i].ID = i
		details[i].TransactionID = transactionID
		_, err := tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2,$3,$4)", transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		if err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	res = &Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}
	return res, nil
}

type ReportRepository interface {
	GetTodaySummary() (totalRevenue int, totalTransaksi int, bestProduct BestProduct, err error)
}

func (r *repository) GetTodaySummary(startDate, endDate time.Time) (int, int, BestProduct, error) {
	var (
		totalRevenue   int
		totalTransaksi int
		best           BestProduct
	)
	query := `
select
  COALESCE(SUM(total_amount), 0) as total_revenue,
  COALESCE(COUNT(*), 0) as total_transaksi
from
  transactions
where
  -- t.created_at::date = CURRENT_DATE 
  created_at >= $1
  and created_at <= $2
	`

	// var args []interface{}
	// if !startDate.IsZero() {
	// 	query += " where created_at::date >= $1"
	// 	args = append(args, startDate)
	// }
	// if !endDate.IsZero() {
	// 	query += " AND created_at::date <= $2"
	// 	args = append(args, endDate.Add(24*time.Hour))
	// }
	err := r.db1.QueryRow(query, startDate, endDate.Add(24*time.Hour)).Scan(&totalRevenue, &totalTransaksi)
	// err := r.db1.QueryRow(query, args...).Scan(&totalRevenue, &totalTransaksi)
	if err != nil {
		return 0, 0, BestProduct{}, err
	}
	err = r.db1.QueryRow(`
select
  p.name,
  COALESCE(SUM(td.quantity), 0) as qty
from
  transaction_details td
  join "Product" p on p.id = td.product_id
  join transactions t on t.id = td.transaction_id
where
  -- t.created_at::date = CURRENT_DATE 
  t.created_at >= $1
  and t.created_at <= $2
group by
  p.name
order by
  qty desc
limit
  1
	`, startDate, endDate.Add(24*time.Hour)).Scan(&best.Nama, &best.QtyTerjual)
	if err == sql.ErrNoRows {
		return totalRevenue, totalTransaksi, BestProduct{}, nil
	}
	if err != nil {
		return 0, 0, BestProduct{}, err
	}
	return totalRevenue, totalTransaksi, best, nil
}
