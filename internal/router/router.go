package router

import (
	"net/http"
	"order-api/internal/handler"

	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(orderHandler *handler.OrderHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("GET /orders", orderHandler.GetOrders)
	mux.HandleFunc("GET /orders/{id}", orderHandler.GetOrderByID)
	mux.HandleFunc("POST /orders", orderHandler.CreateOrder)
	mux.HandleFunc("POST /checkout", orderHandler.Checkout)
	mux.HandleFunc("PUT /orders/{id}", orderHandler.UpdateOrder)
	mux.HandleFunc("DELETE /orders/{id}", orderHandler.DeleteOrder)
	return mux
}
