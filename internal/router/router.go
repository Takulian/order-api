package router

import(
	"net/http"
	"order-api/internal/handler"
)

func NewRouter(orderHandler *handler.OrderHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /orders", orderHandler.GetOrders)
	mux.HandleFunc("GET /orders/{id}", orderHandler.GetOrderByID)
	mux.HandleFunc("POST /orders", orderHandler.CreateOrder)
	return mux
}