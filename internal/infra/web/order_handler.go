package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mobenaus/fc-pos-go-desafio-cleancode/internal/usecase"
)

type WebOrderHandler struct {
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrdersUseCase  usecase.ListOrdersUseCase
}

func NewWebOrderHandler(
	CreateOrderUseCase usecase.CreateOrderUseCase,
	ListOrdersUseCase usecase.ListOrdersUseCase,
) *WebOrderHandler {
	return &WebOrderHandler{
		CreateOrderUseCase: CreateOrderUseCase,
		ListOrdersUseCase:  ListOrdersUseCase,
	}
}

func (h *WebOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := h.CreateOrderUseCase.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebOrderHandler) List(w http.ResponseWriter, r *http.Request) {
	var page int64 = 1
	var limit int64 = 10
	params := r.URL.Query()
	spage := params.Get("page")
	slimit := params.Get("limit")
	if spage != "" {
		page, _ = strconv.ParseInt(spage, 10, 0)
	}
	if slimit != "" {
		limit, _ = strconv.ParseInt(slimit, 10, 0)
	}
	var dto = usecase.OrderListInputDTO{
		Page:  int(page),
		Limit: int(limit),
	}

	output, err := h.ListOrdersUseCase.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
