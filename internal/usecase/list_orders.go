package usecase

import (
	"github.com/mobenaus/fc-pos-go-desafio-cleancode/internal/entity"
)

type OrderListInputDTO struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type OrderListResultOutputDTO struct {
	Total  int
	Orders []OrderListOutputDTO
}

type OrderListOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(
	OrderRepository entity.OrderRepositoryInterface,
) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: OrderRepository,
	}
}

func (c *ListOrdersUseCase) Execute(input OrderListInputDTO) (OrderListResultOutputDTO, error) {
	orderList, err := c.OrderRepository.List(input.Page, input.Limit)
	if err != nil {
		return OrderListResultOutputDTO{}, err
	}

	total, err := c.OrderRepository.GetTotal()
	if err != nil {
		return OrderListResultOutputDTO{}, err
	}

	var orders []OrderListOutputDTO

	for _, order := range orderList {
		orders = append(orders, OrderListOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		})
	}

	return OrderListResultOutputDTO{
		Total:  total,
		Orders: orders,
	}, nil
}
