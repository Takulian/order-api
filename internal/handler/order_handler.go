package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	
	"order-api/internal/dto"
	"order-api/internal/service"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	orders := h.service.GetAll()
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	order := h.service.Create(req)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
