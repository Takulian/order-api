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

// GetOrders godoc
//
//	@Summary		Get all orders
//	@Description	Get all orders
//	@Tags			Orders
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	model.Order
//	@Router			/orders [get]
func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	orders := h.service.GetAll()
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetOrderByID godoc
//
//	@Summary	Get order
//	@Tags		Orders
//	@Produce	json
//	@Param		id	path		int	true	"Order ID"
//	@Success	200	{object}	model.Order
//	@Router		/orders/{id} [get]
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

// CreateOrder godoc
//
//	@Summary		Create order
//	@Description	Create new order
//	@Tags			Orders
//	@Accept			json
//	@Produce		json
//	@Param			order	body		dto.CreateOrderRequest	true	"Create Order"
//	@Success		201		{object}	model.Order
//	@Router			/orders [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	order, err := h.service.Create(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// UpdateOrder godoc
//
//	@Summary	Update order
//	@Tags		Orders
//	@Accept		json
//	@Produce	json
//	@Param		id		path	int						true	"Order ID"
//	@Param		order	body	dto.UpdateOrderRequest	true	"Order"
//	@Success	200		{object}	model.Order
//	@Router		/orders/{id} [put]
func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	order, err := h.service.Update(id, req)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// DeleteOrder godoc
//
//	@Summary	Delete order
//	@Tags		Orders
//	@Produce	json
//	@Param		id	path	int	true	"Order ID"
//	@Success	204
//	@Router		/orders/{id} [delete]
func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(id); err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
