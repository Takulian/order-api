package router

import(
	"net/http"
	"order-api/internal/handler"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /orders", handler.GetOrders)
	mux.HandleFunc("POST /orders", handler.CreateOrder)
	return mux
}