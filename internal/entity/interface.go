package entity

type OrderRepositoryInterface interface {
	Save(order *Order) error
	List(page int, limit int) ([]Order, error)
	GetTotal() (int, error)
}
