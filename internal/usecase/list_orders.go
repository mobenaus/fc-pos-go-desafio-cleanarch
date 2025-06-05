package usecase

import (
	"github.com/mobenaus/fc-pos-go-desafio-cleancode/internal/entity"
)

type OrderListInputDTO struct {
	Page  int `json:"page"`
	limit int `json:"limit"`
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

func (c *ListOrdersUseCase) Execute(input OrderListInputDTO) ([]OrderListOutputDTO, error) {
	orderList, err := c.OrderRepository.List(input.Page, input.limit)
	if err != nil {
		return []OrderListOutputDTO{}, err
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

	return orders, nil
}
