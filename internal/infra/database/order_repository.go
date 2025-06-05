package database

import (
	"database/sql"

	"github.com/mobenaus/fc-pos-go-desafio-cleancode/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) GetTotal() (int, error) {
	var total int
	err := r.Db.QueryRow("Select count(*) from orders").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *OrderRepository) List(page int, limit int) ([]entity.Order, error) {
	var selectpage = page - 1
	var selectlimit = limit
	if page < 1 {
		selectpage = 0
	}
	if limit < 1 {
		limit = 1
	}
	println(selectlimit)
	println(selectpage)
	rows, err := r.Db.Query("select id, price, tax, final_price from orders limit ? offset ?", selectlimit, selectpage)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []entity.Order
	for rows.Next() {
		var o entity.Order
		err = rows.Scan(&o.ID, &o.Price, &o.Tax, &o.FinalPrice)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil

}
